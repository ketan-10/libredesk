// Package autoassigner continuously assigns conversations at regular intervals to users.
package autoassigner

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/abhinavxd/libredesk/internal/conversation/models"
	tmodels "github.com/abhinavxd/libredesk/internal/team/models"
	umodels "github.com/abhinavxd/libredesk/internal/user/models"
	"github.com/mr-karan/balance"
	"github.com/zerodha/logf"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

const (
	AssignmentTypeRoundRobin = "Round robin"
)

type conversationStore interface {
	GetUnassignedConversations() ([]models.Conversation, error)
	UpdateConversationUserAssignee(conversationUUID string, userID int, user umodels.User) error
	ActiveUserConversationsCount(userID int) (int, error)
}

type teamStore interface {
	GetAll() ([]tmodels.Team, error)
	GetMembers(teamID int) ([]umodels.User, error)
}

// Engine represents a manager for assigning unassigned conversations
// to team agents in a round-robin pattern.
type Engine struct {
	roundRobinBalancer map[int]*balance.Balance
	// Mutex to protect the balancer map
	balanceMu              sync.Mutex
	teamMaxAutoAssignments map[int]int

	systemUser        umodels.User
	conversationStore conversationStore
	teamStore         teamStore
	lo                *logf.Logger
	closed            bool
	closedMu          sync.Mutex
	wg                sync.WaitGroup
}

// New initializes a new Engine instance, set up with the provided team manager,
// conversation manager, and logger.
func New(teamStore teamStore, conversationStore conversationStore, systemUser umodels.User, lo *logf.Logger) (*Engine, error) {
	var e = Engine{
		conversationStore:      conversationStore,
		teamStore:              teamStore,
		systemUser:             systemUser,
		lo:                     lo,
		teamMaxAutoAssignments: make(map[int]int),
	}
	balancer, err := e.populateTeamBalancer()
	if err != nil {
		return nil, err
	}
	e.roundRobinBalancer = balancer
	return &e, nil
}

// Run initiates the conversation assignment process and is to be invoked as a goroutine.
// This function continuously assigns unassigned conversations to agents at regular intervals.
func (e *Engine) Run(ctx context.Context, autoAssignInterval time.Duration) {
	ticker := time.NewTicker(autoAssignInterval)
	defer ticker.Stop()

	e.wg.Add(1)
	defer e.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.closedMu.Lock()
			closed := e.closed
			e.closedMu.Unlock()
			if closed {
				return
			}
			if err := e.reloadBalancer(); err != nil {
				e.lo.Error("error reloading balancer", "error", err)
			}
			if err := e.assignConversations(); err != nil {
				e.lo.Error("error assigning conversations", "error", err)
			}
		}
	}
}

// Close signals the Engine to stop its auto-assignment process.
// It sets the closed flag, which will cause the Run loop to exit.
func (e *Engine) Close() {
	e.closedMu.Lock()
	defer e.closedMu.Unlock()
	if e.closed {
		return
	}
	e.closed = true
	e.wg.Wait()
}

// reloadBalancer updates the round-robin balancer with the latest user and team data.
func (e *Engine) reloadBalancer() error {
	e.balanceMu.Lock()
	defer e.balanceMu.Unlock()

	balancer, err := e.populateTeamBalancer()
	if err != nil {
		e.lo.Error("error updating team balancer pool", "error", err)
		return err
	}
	e.roundRobinBalancer = balancer
	return nil
}

// populateTeamBalancer populates the team balancer pool with the team members.
func (e *Engine) populateTeamBalancer() (map[int]*balance.Balance, error) {
	var (
		balancer   = make(map[int]*balance.Balance)
		teams, err = e.teamStore.GetAll()
	)
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		if team.ConversationAssignmentType != AssignmentTypeRoundRobin {
			continue
		}

		users, err := e.teamStore.GetMembers(team.ID)
		if err != nil {
			return nil, err
		}

		// Add users to team's balancer pool.
		for _, user := range users {
			if _, ok := balancer[team.ID]; !ok {
				balancer[team.ID] = balance.NewBalance()
			}
			balancer[team.ID].Add(strconv.Itoa(user.ID), 1)
		}

		// Set max auto assigned conversations for the team.
		e.teamMaxAutoAssignments[team.ID] = team.MaxAutoAssignedConversations
	}
	return balancer, nil
}

// assignConversations function fetches conversations that have been assigned to teams but not to any individual user,
// and then proceeds to assign them to team members based on a round-robin strategy.
func (e *Engine) assignConversations() error {
	unassignedConversations, err := e.conversationStore.GetUnassignedConversations()
	if err != nil {
		return fmt.Errorf("fetching unassigned conversations: %w", err)
	}

	if len(unassignedConversations) > 0 {
		e.lo.Debug("found unassigned conversations", "count", len(unassignedConversations))
	}

	for _, conversation := range unassignedConversations {
		// Get user from the pool.
		userIDStr, err := e.getUserFromPool(conversation.AssignedTeamID.Int)
		if err != nil {
			if err != ErrTeamNotFound {
				e.lo.Error("error fetching user from balancer pool", "conversation_uuid", conversation.UUID, "error", err)
			}
			continue
		}

		// Convert to int.
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			e.lo.Error("error converting user id from string to int", "user_id", userIDStr, "error", err)
			continue
		}

		// Get active conversations count for the user.
		activeConversationsCount, err := e.conversationStore.ActiveUserConversationsCount(userID)
		if err != nil {
			e.lo.Error("error fetching active conversations count for user", "user_id", userID, "error", err)
			continue
		}

		// Check if user has reached the max auto assigned conversations limit.
		if activeConversationsCount >= e.teamMaxAutoAssignments[conversation.AssignedTeamID.Int] {
			e.lo.Debug("user has reached max auto assigned conversations limit, skipping auto assignment", "user_id", userID, "user_active_conversations_count", activeConversationsCount, "max_auto_assigned_conversations", e.teamMaxAutoAssignments[conversation.AssignedTeamID.Int])
			continue
		}

		// Assign conversation to user.
		if err := e.conversationStore.UpdateConversationUserAssignee(conversation.UUID, userID, e.systemUser); err != nil {
			e.lo.Error("error assigning conversation", "conversation_uuid", conversation.UUID, "error", err)
			continue
		}
	}
	return nil
}

// getUserFromPool returns user ID from the team balancer pool.
func (e *Engine) getUserFromPool(assignedTeamID int) (string, error) {
	e.balanceMu.Lock()
	defer e.balanceMu.Unlock()

	pool, ok := e.roundRobinBalancer[assignedTeamID]
	if !ok {
		return "", ErrTeamNotFound
	}
	return pool.Get(), nil
}
