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
        <va-switch
          v-model="store.disableNotifications"
          size="small"
          class="pb-6"
          :label="t('settings.disableNotifications')"
        />
      </div>

      <div>
        <va-switch v-model="store.showServerOnMap" size="small" :label="t('settings.showServerTitle')" />
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
