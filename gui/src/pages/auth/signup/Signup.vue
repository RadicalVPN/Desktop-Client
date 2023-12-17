<template>
  <div class="flex flex-col items-center">
    <va-button @click="openPage">Signup now</va-button>
  </div>
</template>

<script setup lang="ts">
  import { shell } from 'electron'

  function openPage() {
    shell.openExternal('https://radicalvpn.com/portal/auth/signup')
  }
</script>
