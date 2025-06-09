<template>
  <form @submit="onSubmit" class="space-y-6 w-full">
    <FormField v-slot="{ field }" name="site_title">
      <FormItem>
        <FormLabel>{{ t('admin.general.siteTitle') }}</FormLabel>
        <FormControl>
          <Input type="text" placeholder="" v-bind="field" />
        </FormControl>
        <FormDescription>
          {{ t('admin.general.siteTitle.description') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ field }" name="site_subtitle">
      <FormItem>
        <FormLabel>{{ t('admin.general.siteSubtitle') }}</FormLabel>
        <FormControl>
          <Input type="text" placeholder="" v-bind="field" />
        </FormControl>
        <FormDescription>
          {{ t('admin.general.siteSubtitle.description') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ field }" name="site_description">
      <FormItem>
        <FormLabel>{{ t('admin.general.siteDescription') }}</FormLabel>
        <FormControl>
          <Textarea 
            placeholder="" 
            v-bind="field" 
            class="min-h-[100px]"
          />
        </FormControl>
        <FormDescription>
          {{ t('admin.general.siteDescription.description') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="primary_font">
      <FormItem>
        <FormLabel>{{ t('admin.general.primaryFont') }}</FormLabel>
        <FormControl>
          <Select v-bind="componentField" :modelValue="componentField.modelValue">
            <SelectTrigger>
              <SelectValue :placeholder="t('admin.general.primaryFont.placeholder')" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="inter">Inter</SelectItem>
                <SelectItem value="roboto">Roboto</SelectItem>
                <SelectItem value="open-sans">Open Sans</SelectItem>
                <SelectItem value="lato">Lato</SelectItem>
                <SelectItem value="montserrat">Montserrat</SelectItem>
                <SelectItem value="poppins">Poppins</SelectItem>
                <SelectItem value="raleway">Raleway</SelectItem>
                <SelectItem value="playfair-display">Playfair Display</SelectItem>
                <SelectItem value="merriweather">Merriweather</SelectItem>
                <SelectItem value="ubuntu">Ubuntu</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          {{ t('admin.general.primaryFont.description') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ field }" name="logo_url" :value="props.initialValues.logo_url">
      <FormItem>
        <FormLabel>{{ t('admin.general.logoURL') }}</FormLabel>
        <FormControl>
          <Input type="text" placeholder="" v-bind="field" />
        </FormControl>
        <FormDescription>{{ t('admin.general.logoURL.description') }}</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <Button type="submit" :isLoading="formLoading"> {{ submitLabel }} </Button>
  </form>
</template>

<script setup>
import { watch, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { useI18n } from 'vue-i18n'

const emitter = useEmitter()
const { t } = useI18n()
const formLoading = ref(false)

const props = defineProps({
  initialValues: {
    type: Object,
    required: false
  },
  submitForm: {
    type: Function,
    required: true
  },
  submitLabel: {
    type: String,
    required: false,
    default: ''
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const submitLabel = props.submitLabel || t('globals.messages.save')
const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t))
})

const onSubmit = form.handleSubmit(async (values) => {
  try {
    formLoading.value = true
    await props.submitForm(values)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.updatedSuccessfully', {
        name: t('globals.terms.setting', 2)
      })
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    formLoading.value = false
  }
})

// Watch for changes in initialValues and update the form.
watch(
  () => props.initialValues,
  (newValues) => {
    if (Object.keys(newValues).length === 0) {
      return
    }
    form.setValues(newValues)
  },
  { deep: true }
)
</script>