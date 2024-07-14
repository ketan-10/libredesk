package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abhinavxd/artemis/internal/envelope"
	umodels "github.com/abhinavxd/artemis/internal/user/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

func handleGetAllConversations(r *fastglue.Request) error {
	var (
		app         = r.Context.(*App)
		order       = string(r.RequestCtx.QueryArgs().Peek("order"))
		orderBy     = string(r.RequestCtx.QueryArgs().Peek("order_by"))
		filter      = string(r.RequestCtx.QueryArgs().Peek("filter"))
		page, _     = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page")))
		pageSize, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page_size")))
	)
	c, err := app.conversationManager.GetAll(order, orderBy, filter, page, pageSize)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(c)
}

func handleGetAssignedConversations(r *fastglue.Request) error {
	var (
		app         = r.Context.(*App)
		user        = r.RequestCtx.UserValue("user").(umodels.User)
		order       = string(r.RequestCtx.QueryArgs().Peek("order"))
		orderBy     = string(r.RequestCtx.QueryArgs().Peek("order_by"))
		filter      = string(r.RequestCtx.QueryArgs().Peek("filter"))
		page, _     = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page")))
		pageSize, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page_size")))
	)
	c, err := app.conversationManager.GetAssigned(user.ID, order, orderBy, filter, page, pageSize)
	if err != nil {
		return r.SendErrorEnvelope(http.StatusInternalServerError, err.Error(), nil, "")
	}
	return r.SendEnvelope(c)
}

func handleGetTeamConversations(r *fastglue.Request) error {
	var (
		app         = r.Context.(*App)
		user        = r.RequestCtx.UserValue("user").(umodels.User)
		order       = string(r.RequestCtx.QueryArgs().Peek("order"))
		orderBy     = string(r.RequestCtx.QueryArgs().Peek("order_by"))
		filter      = string(r.RequestCtx.QueryArgs().Peek("filter"))
		page, _     = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page")))
		pageSize, _ = strconv.Atoi(string(r.RequestCtx.QueryArgs().Peek("page_size")))
	)
	c, err := app.conversationManager.GetTeamConversations(user.ID, order, orderBy, filter, page, pageSize)
	if err != nil {
		return r.SendErrorEnvelope(http.StatusInternalServerError, err.Error(), nil, "")
	}
	return r.SendEnvelope(c)
}

func handleGetConversation(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		uuid = r.RequestCtx.UserValue("uuid").(string)
	)
	c, err := app.conversationManager.Get(uuid)
	if err != nil {
		return r.SendErrorEnvelope(http.StatusInternalServerError, err.Error(), nil, "")
	}
	return r.SendEnvelope(c)
}

func handleUpdateAssigneeLastSeen(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		uuid = r.RequestCtx.UserValue("uuid").(string)
	)
	err := app.conversationManager.UpdateAssigneeLastSeen(uuid)
	if err != nil {
		return r.SendErrorEnvelope(http.StatusInternalServerError, err.Error(), nil, "")
	}
	return r.SendEnvelope(true)
}

func handleGetConversationParticipants(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		uuid = r.RequestCtx.UserValue("uuid").(string)
	)
	p, err := app.conversationManager.GetParticipants(uuid)
	if err != nil {
		return r.SendErrorEnvelope(http.StatusInternalServerError, err.Error(), nil, "")
	}
	return r.SendEnvelope(p)
}

func handleUpdateUserAssignee(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		uuid = r.RequestCtx.UserValue("uuid").(string)
		user = r.RequestCtx.UserValue("user").(umodels.User)
	)
	assigneeID, err := r.RequestCtx.PostArgs().GetUint("assignee_id")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid assignee `id`.", nil, envelope.InputError)
	}
	if err := app.conversationManager.UpdateUserAssignee(uuid, assigneeID, user); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUpdateTeamAssignee(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		uuid = r.RequestCtx.UserValue("uuid").(string)
		user = r.RequestCtx.UserValue("user").(umodels.User)
	)
	assigneeID, err := r.RequestCtx.PostArgs().GetUint("assignee_id")
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid assignee `id`.", nil, envelope.InputError)
	}
	if err := app.conversationManager.UpdateTeamAssignee(uuid, assigneeID, user); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUpdatePriority(r *fastglue.Request) error {
	var (
		app      = r.Context.(*App)
		p        = r.RequestCtx.PostArgs()
		priority = p.Peek("priority")
		uuid     = r.RequestCtx.UserValue("uuid").(string)
		user     = r.RequestCtx.UserValue("user").(umodels.User)
	)
	if err := app.conversationManager.UpdatePriority(uuid, priority, user); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUpdateStatus(r *fastglue.Request) error {
	var (
		app    = r.Context.(*App)
		p      = r.RequestCtx.PostArgs()
		status = p.Peek("status")
		uuid   = r.RequestCtx.UserValue("uuid").(string)
		user   = r.RequestCtx.UserValue("user").(umodels.User)
	)
	if err := app.conversationManager.UpdateStatus(uuid, status, user); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleAddConversationTags(r *fastglue.Request) error {
	var (
		app     = r.Context.(*App)
		p       = r.RequestCtx.PostArgs()
		tagIDs  = []int{}
		tagJSON = p.Peek("tag_ids")
		uuid    = r.RequestCtx.UserValue("uuid").(string)
	)
	err := json.Unmarshal(tagJSON, &tagIDs)
	if err != nil {
		app.lo.Error("unmarshalling tag ids", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "error adding tags", nil, "")
	}

	if err := app.conversationManager.UpsertTags(uuid, tagIDs); err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUserDashboardCounts(r *fastglue.Request) error {
	var (
		app  = r.Context.(*App)
		user = r.RequestCtx.UserValue("user").(umodels.User)
	)

	stats, err := app.conversationManager.GetAssigneeStats(user.ID)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(stats)
}

func handleUserDashboardCharts(r *fastglue.Request) error {
	var app = r.Context.(*App)
	stats, err := app.conversationManager.GetNewConversationsStats()
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(stats)
}
