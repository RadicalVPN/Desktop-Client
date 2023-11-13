<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.personalization') }}</va-card-title>

    <va-card-content>
      <va-select v-model="store.theme" :options="themeOptions" class="mb-4" :label="t('settings.theme')" />

      <va-select
        v-model="store.language"
        class="mb-6"
        :label="t('settings.language')"
        :options="languages"
        text-by="textBy"
        value-by="text"
      />

      <va-switch v-model="store.animatedMap" size="small" class="pb-6" :label="t('settings.animatedMap')" />

      <div>
        <va-switch v-model="store.disableNotifications" size="small" :label="t('settings.disableNotifications')" />
      </div>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { computed, watchEffect } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useColors } from 'vuestic-ui'
  import { useGlobalStore } from '../../../stores/global-store'

  const { t, locale } = useI18n()
  const { applyPreset } = useColors()
  const store = useGlobalStore()

  const themeOptions = ['light', 'dark']
  const languages = [
    {
      text: 'gb',
      textBy: t('language.english'),
    },
    {
      text: 'de',
      textBy: t('language.german'),
    },
  ]

  watchEffect(() => {
    setTheme(store.theme)
    applyLocale(store.language)
  })

  function setTheme(theme: string) {
    theme = theme.toLocaleLowerCase()

    localStorage.setItem('theme', theme)
    store.theme = theme
    applyPreset(theme)
  }

  function applyLocale(newLocale: string) {
    console.log('updating locale', newLocale)
    locale.value = newLocale
  }
</script>
