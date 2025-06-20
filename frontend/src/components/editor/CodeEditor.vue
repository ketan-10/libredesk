<template>
    <div ref="codeEditor" @click="editorView?.focus()" class="w-full h-[28rem] border overflow-y-scroll" />
</template>

<script setup>
import { ref, onMounted, watch, nextTick, useTemplateRef } from 'vue'
import { EditorView, basicSetup } from 'codemirror'
import { html } from '@codemirror/lang-html'
import { oneDark } from '@codemirror/theme-one-dark'

const props = defineProps({
    modelValue: { type: String, default: '' },
    language: { type: String, default: 'html' },
    disabled: Boolean
})

const emit = defineEmits(['update:modelValue'])
const data = ref('')
let editorView = null 
const codeEditor = useTemplateRef('codeEditor')

const initCodeEditor = (body) => {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches

    editorView = new EditorView({
        doc: body,
        extensions: [
            basicSetup,
            html(),
            ...(isDark ? [oneDark] : []),
            EditorView.editable.of(!props.disabled),
            EditorView.updateListener.of((update) => {
                if (!update.docChanged) return
                const v = update.state.doc.toString()
                emit('update:modelValue', v)
                data.value = v
                
            })
        ],
        parent: codeEditor.value
    })

    nextTick(() => {
        editorView?.focus()
    })
}

onMounted(() => {
    initCodeEditor(props.modelValue || '')
})

watch(() => props.modelValue, (newVal) => {
    if (newVal !== data.value) {
        editorView?.dispatch({
            changes: { from: 0, to: editorView.state.doc.length, insert: newVal }
        })
    }
})
</script>