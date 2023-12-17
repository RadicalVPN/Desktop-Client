<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.personalization') }}</va-card-title>

    <va-card-content>
      <div class="flex items-center justify-between">
        <p>{{ t('settings.language') }}</p>
        <div class="w-40">
          <va-select v-model="store.language" :options="languages" text-by="textBy" value-by="text" class="mb-6" />
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.theme') }}</p>
        <div class="w-40">
          <va-select v-model="store.theme" :options="themeOptions" class="mb-6" />
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.animatedMap') }}</p>
        <div class="w-40">
          <va-switch v-model="store.animatedMap" size="small" class="mb-6" />
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.disableNotifications') }}</p>
        <div class="w-40">
          <va-switch v-model="store.disableNotifications" size="small" class="mb-6" />
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.showServerTitle') }}</p>
        <div class="w-40">
          <va-switch v-model="store.showServerOnMap" size="small" class="mb-6" />
        </div>
      </div>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { watchEffect } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useColors } from 'vuestic-ui'
  import { useGlobalStore } from '../../../stores/global-store'
  import rawLanguages from '../../../i18n/languages'
  import { computed } from 'vue'

  const { t, locale } = useI18n()
  const { applyPreset } = useColors()
  const store = useGlobalStore()

  const themeOptions = ['light', 'dark']
  const languages = computed(() =>
    rawLanguages.map((lang) => ({
      ...lang,
      textBy: t(lang.textBy),
    })),
  )

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
    if (newLocale === '_system') {
      const sysLang = navigator.language

      //valid locale?
      const isValidLocale = languages.value.some((language) => language.text === sysLang)

      if (isValidLocale) {
        newLocale = sysLang
      } else {
        console.warn('no translation for system language', sysLang)
        newLocale = 'gb'
      }
    }

    console.log('updating locale', newLocale)
    locale.value = newLocale
  }
</script>
