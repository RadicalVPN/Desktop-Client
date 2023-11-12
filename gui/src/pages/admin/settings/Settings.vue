<template>
  <h1 style="color: var(--va-text-primary)" class="text-xl mb-6">Settings</h1>

  <settings-personalization class="mb-8" />
  <settings-versions />
</template>

<script lang="ts" setup>
  import SettingsPersonalization from './SettingsPersonalization.vue'
  import SettingsVersions from './SettingsVersions.vue'
</script>
