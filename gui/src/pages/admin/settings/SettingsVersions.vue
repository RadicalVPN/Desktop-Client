<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.versions') }}</va-card-title>

    <va-card-content class="my-3 grid grid-cols-12 gap-6">
      <div class="col-span-10">
        <div class="mb-4">
          <p class="text-sm font-semibold">{{ t('settings.frontendVersion') }}</p>
          <p class="text-lg">{{ 'v' + frontendVersion }}</p>
        </div>

        <div class="mb-4">
          <p class="text-sm font-semibold">{{ t('settings.daemonVersion') }}</p>
          <p class="text-lg">{{ 'v' + daemonVersion }}</p>
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
