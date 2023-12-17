<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.logs') }}</va-card-title>

    <va-card-content>
      <va-switch v-model="logsActivated" class="mb-4" size="small" :label="t('settings.activeLogs')" />

      <va-virtual-scroller v-if="logsActivated" v-slot="{ item }" :items="logs" :wrapper-size="300" :bench="5">
        <va-badge class="pb-0.5" :color="item.color" :text="item.message" />
      </va-virtual-scroller>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { DaemonHelper, ParsedLog } from '../../../helper/daemon'
  import { onBeforeUnmount } from 'vue'

  const { t } = useI18n()

  const logs = ref<ParsedLog[]>([])
  const logsActivated = ref()

  async function updateLogs() {
    logs.value = await new DaemonHelper().getLogs()
  }

  let refreshTimer: any
  onMounted(async () => {
    await updateLogs()

    refreshTimer = setInterval(async () => {
      await updateLogs()
    }, 10_000)
  })

  //make sure to clean the interval before unmounting
  onBeforeUnmount(() => {
    clearInterval(refreshTimer)
  })
</script>
