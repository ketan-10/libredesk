import { h } from 'vue'
import UserDataTableDropDown from '@/components/admin/team/users/UserDataTableDropDown.vue'
import { format } from 'date-fns'

export const columns = [
  {
    accessorKey: 'first_name',
    header: function () {
      return h('div', { class: 'text-center' }, 'First name')
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center font-medium' }, row.getValue('first_name'))
    }
  },
  {
    accessorKey: 'last_name',
    header: function () {
      return h('div', { class: 'text-center' }, 'Last name')
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center font-medium' }, row.getValue('last_name'))
    }
  },
  {
    accessorKey: 'email',
    header: function () {
      return h('div', { class: 'text-center' }, 'Email')
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center font-medium' }, row.getValue('email'))
    }
  },
  {
    accessorKey: 'updated_at',
    header: function () {
      return h('div', { class: 'text-center' }, 'Modified at')
    },
    cell: function ({ row }) {
      return h(
        'div',
        { class: 'text-center font-medium' },
        format(row.getValue('updated_at'), 'PPpp')
      )
    }
  },
  {
    id: 'actions',
    enableHiding: false,
    cell: ({ row }) => {
      const user = row.original
      return h(
        'div',
        { class: 'relative' },
        h(UserDataTableDropDown, {
          user
        })
      )
    }
  }
]
