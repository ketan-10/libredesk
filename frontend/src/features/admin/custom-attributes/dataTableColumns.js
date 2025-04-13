import { h } from 'vue'
import dataTableDropdown from '@/features/admin/custom-attributes/dataTableDropdown.vue'
import { format } from 'date-fns'

export const createColumns = (t) => [
    {
        accessorKey: 'name',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.name'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center font-medium' }, row.getValue('name'))
        }
    },
    {
        accessorKey: 'key',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.key'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center font-medium' }, row.getValue('key'))
        }
    },
    {
        accessorKey: 'data_type',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.type'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center font-medium' }, row.getValue('data_type'))
        }
    },
    {
        accessorKey: 'applies_to',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.appliesTo'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center font-medium' }, row.getValue('applies_to'))
        }
    },
    {
        accessorKey: 'created_at',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.createdAt'))
        },
        cell: function ({ row }) {
            return h(
                'div',
                { class: 'text-center font-medium' },
                format(row.getValue('created_at'), 'PPpp')
            )
        }
    },
    {
        accessorKey: 'updated_at',
        header: function () {
            return h('div', { class: 'text-center' }, t('form.field.updatedAt'))
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
            const customAttribute = row.original
            return h(
                'div',
                { class: 'relative' },
                h(dataTableDropdown, {
                    customAttribute
                })
            )
        }
    }
]
