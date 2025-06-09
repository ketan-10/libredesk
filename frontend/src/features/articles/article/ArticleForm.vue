<template>
  <form @submit.prevent="onSubmit" class="space-y-6 w-full">

    <!-- Article Title -->
    <FormField v-slot="{ componentField }" name="title">
      <FormItem>
        <FormLabel>{{ $t('globals.terms.title') }}</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" class="text-lg font-medium" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Section Selection -->
    <FormField v-slot="{ componentField }" name="section_id">
      <FormItem>
        <FormLabel class="flex items-center whitespace-nowrap">
          {{ $t('globals.terms.section') }}
        </FormLabel>
        <FormControl>
          <div class="relative">
            <ComboBox v-bind="componentField" :items="sections"
              :placeholder="$t('globals.messages.select', { name: $t('globals.terms.section') })"
              :disabled="isSectionLoading">
                <template #item="{ item }">
                  <div class="flex items-center gap-2">
                    <span class="text-sm">{{ item.name }}</span>
                  </div>
                </template>

                <template #selected="{ selected }">
                  <div class="flex items-center">
                    <span v-if="selected">{{ selected.name }}</span>
                  </div>
                </template>
              </ComboBox>
            <div v-if="isSectionLoading" class="absolute inset-0 flex items-center justify-center">
              <Loader2 class="h-4 w-4 animate-spin" />
            </div>
          </div>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="content">
      <FormItem>
        <FormControl>
          <CodeMirror v-model="componentField.modelValue" @update:modelValue="handleChange" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
    
    <Button type="submit" :isLoading="isLoading"> {{ hasFormChanges ? 'Submit' : 'Publish' }} </Button>
  </form>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import ComboBox from '@/components/ui/combobox/ComboBox.vue'
import { useI18n } from 'vue-i18n'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import api from '@/api/index.js'
import CodeMirror from '@/components/editor/CodeHtmlPreview.vue'
import Button from '@/components/ui/button/Button.vue'
import { Loader2 } from 'lucide-vue-next'

const props = defineProps({
  initialValues: {
    type: Object,
    required: false,
    default: () => ({})
  },
  submitForm: {
    type: Function,
    required: true
  },
  isNewForm: {
    type: Boolean,
    required: false,
    default: false
  },
  isLoading: {
    type: Boolean,
    required: false
  }
})

const { t } = useI18n()
const emitter = useEmitter()
const sections = ref([])
const isSectionLoading = ref(false)

const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: props.initialValues
})

// Track if form has changes
const hasFormChanges = computed(() => {
  // For new forms, always show submit button
  if (props.isNewForm) return true

  // Check if form is dirty (has changes)
  return form.meta.value.dirty
})

const onSubmit = form.handleSubmit((values) => {

  props.submitForm(values)
})

// Fetch sections
const fetchData = async () => {
  try {
    isSectionLoading.value = true
    const resp = await api.getArticleSections()

    sections.value = resp?.data?.data?.map(sec => ({
      value: sec.id,
      name: sec.name
    })) ?? []

  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: t('globals.messages.errorFetching')
    })
  } finally {
    isSectionLoading.value = false
  }
}

onMounted(() => {
  fetchData()

  // Set initial values if provided
  if (Object.keys(props.initialValues).length > 0) {
    form.setValues(props.initialValues)
  }
})
</script>
