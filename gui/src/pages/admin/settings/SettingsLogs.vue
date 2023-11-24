<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.logs') }}</va-card-title>

    <va-card-content>
      <va-switch v-model="logsActivated" size="small" :label="t('settings.activeLogs')" />

      <va-virtual-scroller
        v-if="logsActivated"
        v-slot="{ item }"
        v-el:scroller
        :items="logs.reverse()"
        :wrapper-size="200"
        :bench="5"
      >
        <va-badge class="pb-2" :color="item.color" :text="item.message" />
      </va-virtual-scroller>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { DaemonHelper } from '../../../helper/daemon'

  const { t } = useI18n()

  const logs = ref<any[]>([])
  const logsActivated = ref()

  async function updateLogs() {
    logs.value = await new DaemonHelper().getLogs()
    console.log(logs.value)
  }

  onMounted(async () => {
    await updateLogs()

    setInterval(async () => {
      await updateLogs()
    }, 10_000)
  })
</script>
