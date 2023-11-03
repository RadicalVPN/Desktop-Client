<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.personalization') }}</va-card-title>

    <va-card-content>
      <va-select v-model="selectedTheme" :options="themeOptions" class="mb-8" :label="t('menu.theme')"></va-select>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { ref, watchEffect } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useColors } from 'vuestic-ui'
  import { useGlobalStore } from '../../../stores/global-store'

  const { t } = useI18n()
  const { applyPreset } = useColors()
  const store = useGlobalStore()

  const selectedTheme = ref(localStorage.getItem('theme') || 'Light')
  const themeOptions = ['Light', 'Dark']

  watchEffect(() => {
    setTheme(selectedTheme.value)
  })

  function setTheme(theme: string) {
    theme = theme.toLocaleLowerCase()

    localStorage.setItem('theme', theme)
    store.theme = theme
    applyPreset(theme)
  }
</script>
