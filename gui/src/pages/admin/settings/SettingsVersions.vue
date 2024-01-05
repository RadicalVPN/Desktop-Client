<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.versions') }}</va-card-title>

    <va-card-content>
      <div class="flex items-center justify-between">
        <p>{{ t('settings.version') }}</p>
        <div class="w-40">
          <p class="text-lg mb-6">{{ version }}</p>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <p>{{ t('settings.appType') }}</p>
        <div class="w-40">
          <p class="text-lg mb-6">{{ appType }}</p>
        </div>
      </div>
    </va-card-content>
  </va-card>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import { DaemonHelper } from '../../../helper/daemon'
  import { onMounted, ref } from 'vue'

  const { t } = useI18n()

  const version = ref('N/A')
  const appType = ref<'Production' | 'Nightly' | 'Unknown'>('Unknown')

  onMounted(async () => {
    await loadDaemonVersionInfo()
  })

  async function loadDaemonVersionInfo() {
    const versionInfo = await new DaemonHelper().getDaemonVersionInfo()

    version.value = versionInfo.version

    if (versionInfo.nightly.isNightly) {
      appType.value = 'Nightly'
    } else if (versionInfo.release.isRelease) {
      appType.value = 'Production'
    }
  }
</script>
