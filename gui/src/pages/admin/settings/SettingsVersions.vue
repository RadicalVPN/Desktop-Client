<template>
  <va-card class="horizontal-bars">
    <va-card-title style="font-size: 0.875rem">{{ t('settings.versions') }}</va-card-title>

    <va-card-content>
      <div class="flex items-center justify-between">
        <p>{{ t('settings.version') }}</p>
        <div class="w-40">
          <p class="text-lg mb-6">{{ 'v' + version }}</p>
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
  const version = ref('unk')

  onMounted(async () => {
    await loadDaemonVersion()
  })

  async function loadDaemonVersion() {
    version.value = await new DaemonHelper().getDaemonVersion()
  }
</script>
