<template>
  <h1 style="color: var(--va-text-primary)" class="text-xl mb-6">{{ t('privacyFirewall.title') }}</h1>

  <privacy-firewall-information class="mb-8 max-w-3xl" />
  <privacy-firewall-configuration class="max-w-3xl" />
  <!--  -->
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n'
  import PrivacyFirewallInformation from './Information.vue'
  import PrivacyFirewallConfiguration from './Configuration.vue'

  const { t } = useI18n()
</script>
