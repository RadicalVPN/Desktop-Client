<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.personalization') }}</va-card-title>

    <va-card-content>
      <va-select v-model="store.theme" :options="themeOptions" class="mb-8" :label="t('settings.theme')" />
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { watchEffect } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useColors } from 'vuestic-ui'
  import { useGlobalStore } from '../../../stores/global-store'

  const { t } = useI18n()
  const { applyPreset } = useColors()
  const store = useGlobalStore()

  const themeOptions = ['light', 'dark']

  watchEffect(() => {
    setTheme(store.theme)
  })

  function setTheme(theme: string) {
    theme = theme.toLocaleLowerCase()

    localStorage.setItem('theme', theme)
    store.theme = theme
    applyPreset(theme)
  }
</script>
