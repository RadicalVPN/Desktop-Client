<template>
  <h1 style="color: var(--va-text-primary)" class="text-xl mb-6">{{ t('settings.settings') }}</h1>

  <settings-personalization class="mb-8" />
  <settings-versions class="mb-8" />
  <settings-logs />
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import SettingsPersonalization from './SettingsPersonalization.vue'
  import SettingsVersions from './SettingsVersions.vue'
  import SettingsLogs from './SettingsLogs.vue'

  const { t } = useI18n()
</script>
