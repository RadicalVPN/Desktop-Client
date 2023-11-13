<template>
  <h1 style="color: var(--va-text-primary)" class="text-xl mb-6">{{ t('settings.settings') }}</h1>

  <settings-personalization class="mb-8" />
  <settings-versions />
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import SettingsPersonalization from './SettingsPersonalization.vue'
  import SettingsVersions from './SettingsVersions.vue'

  const { t } = useI18n()
</script>
