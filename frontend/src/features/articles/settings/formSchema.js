import * as z from 'zod'

export const createFormSchema = (t) => {
  return z.object({
    site_title: z
      .string()
      .min(1, { message: t('validation.required', { field: t('admin.general.siteTitle') }) })
      .max(100, { message: t('validation.maxLength', { field: t('admin.general.siteTitle'), max: 100 }) }),
    
    site_subtitle: z
      .string()
      .max(200, { message: t('validation.maxLength', { field: t('admin.general.siteSubtitle'), max: 200 }) })
      .optional(),
    
    site_description: z
      .string()
      .max(500, { message: t('validation.maxLength', { field: t('admin.general.siteDescription'), max: 500 }) })
      .optional(),
    
    primary_font: z
      .string()
      .min(1, { message: t('validation.required', { field: t('admin.general.primaryFont') }) }),
    
    logo_url: z
      .string()
      .url({ message: t('validation.invalidUrl', { field: t('admin.general.logoURL') }) })
      .optional()
      .or(z.literal(''))
  })
}