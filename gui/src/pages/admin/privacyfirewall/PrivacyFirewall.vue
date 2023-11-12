<template>
  <va-card class="mb-8">
    <va-card-title>Information</va-card-title>

    <va-card-content>
      <va-collapse v-model="test" header="Collapse header"> Collapse content </va-collapse></va-card-content
    >
  </va-card>

  <va-card class="mb-8">
    <va-card-title>Configuration</va-card-title>

    <va-card-content>
      <va-select
        v-model="store.privacyFirewallLevel"
        class="mb-6"
        label="asdasd"
        :options="levels"
        text-by="textBy"
        value-by="text"
      />
    </va-card-content>

    {{ store.privacyFirewallLevel }}
  </va-card>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import { useGlobalStore } from '../../../stores/global-store'
  import { SelectableOption } from 'vuestic-ui/dist/types/composables'
  import { ref } from 'vue'

  const { t } = useI18n()
  const store = useGlobalStore()

  const test = ref()

  const levels = [
    {
      text: 'basic',
      textBy: t('privacyFirewall.basic'),
    },
    {
      text: 'recommended',
      textBy: t('privacyFirewall.recommended'),
    },
    {
      text: 'comprehensive',
      textBy: t('privacyFirewall.comprehensive'),
    },
    {
      text: 'agressive',
      textBy: t('privacyFirewall.agressive'),
    },
  ] as SelectableOption[]
</script>
