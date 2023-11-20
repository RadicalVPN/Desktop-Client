<template>
  <va-card>
    <va-card-title style="font-size: 0.875rem">{{ t('privacyFirewall.settings.title') }}</va-card-title>

    <va-card-content>
      <va-select
        v-model="store.privacyFirewallLevel"
        class="mb-6 max-w-3xl"
        :label="t('privacyFirewall.settings.dropDownLabel')"
        :options="levels"
        text-by="textBy"
        value-by="text"
        :error="isAggresive"
        :error-messages="isAggresive ? [t('errors.privacyFirewallAggresive')] : []"
      />
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import { useGlobalStore } from '../../../stores/global-store'
  import { SelectableOption } from 'vuestic-ui/dist/types/composables'
  import { computed } from 'vue'

  const { t } = useI18n()
  const store = useGlobalStore()

  const isAggresive = computed(() => store.privacyFirewallLevel === 'aggresive')
  const levels = [
    {
      text: 'basic',
      textBy: t('privacyFirewall.level.basic'),
    },
    {
      text: 'recommended',
      textBy: t('privacyFirewall.level.recommended'),
    },
    {
      text: 'comprehensive',
      textBy: t('privacyFirewall.level.comprehensive'),
    },
    {
      text: 'aggresive',
      textBy: t('privacyFirewall.level.aggresive'),
    },
  ] as SelectableOption[]
</script>
