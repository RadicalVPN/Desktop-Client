<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.versions') }}</va-card-title>

    <va-card-content>
      <div class="flex items-center justify-between">
        <p>{{ t('settings.frontendVersion') }}</p>
        <div class="w-40">
          <p class="text-lg mb-6">{{ 'v' + frontendVersion }}</p>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.daemonVersion') }}</p>
        <div class="w-40">
          <p class="text-lg mb-6">{{ 'v' + daemonVersion }}</p>
        </div>
      </div>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import { version as frontendVersion } from '../../../../package.json'
  import { DaemonHelper } from '../../../helper/daemon'
  import { onMounted, ref } from 'vue'

  const { t } = useI18n()
  const daemonVersion = ref('Unknown')

  onMounted(async () => {
    await loadDaemonVersion()
  })

  async function loadDaemonVersion() {
    daemonVersion.value = await new DaemonHelper().getDaemonVersion()
  }
</script>
