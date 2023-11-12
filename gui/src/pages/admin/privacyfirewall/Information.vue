<template>
  <va-badge overlap>
    <template #text> {{ t('new') }} </template>
    <va-card>
      <va-card-title style="font-size: 0.875rem">{{ t('privacyFirewall.info.title') }}</va-card-title>
      <va-card-content class="space-y-2">
        <p>{{ t('privacyFirewall.info.description.1') }}</p>
        <p>{{ t('privacyFirewall.info.description.2') }}</p>
        <p>{{ t('privacyFirewall.info.description.3') }}</p>
      </va-card-content>
    </va-card>
  </va-badge>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'

  const { t } = useI18n()
</script>
