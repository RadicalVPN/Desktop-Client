<template>
  <div class="auth-layout grid grid-cols-12 content-center">
    <div class="flex col-span-12 p-4 justify-center">
      <radical-logo height="80" />
    </div>

    <div v-if="!connectionFailed" class="flex col-span-12 p-4 justify-center">
      <spring-spinner :animation-duration="3000" :size="60" />
    </div>

    <div class="flex justify-center col-span-12 p-4">
      <va-card class="auth-layout__card">
        <va-card-content class="text-center fold-bold">
          <div v-if="!connectionFailed">
            <p class="font-bold pb-6">Connecting to Daemon..</p>

            <p class="pb-6">Attempt: {{ attempts }}</p>
            <p>Delay: {{ backOff }} ms</p>
          </div>

          <div v-if="connectionFailed">
            <p class="text-xl">❗️❗️ No connection to Daemon ❗️❗️</p>
            <br />
            <p>You may restart your PC and try again. Or contact our customer support.<br /><br /></p>
            <a class="content-center font-bold underline" href="mailto:support@radicalvpn.com">
              support@radicalvpn.com
            </a>

            <div class="mt-8">
              <va-button @click="tryConnect()">Try again</va-button>
            </div>
          </div>
        </va-card-content>
      </va-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import RadicalLogo from '../../components/RadicalLogo.vue'
  import { SpringSpinner } from 'epic-spinners'
  import { setTimeout } from 'timers/promises'
  import { DaemonHelper } from '../../helper/daemon'
  import { useRouter } from 'vue-router'
  import { useGlobalStore } from '../../stores/global-store'
  import { setMapStoreSuffix } from 'pinia'

  const daemonHelper = new DaemonHelper()
  const router = useRouter()
  const store = useGlobalStore()

  const attempts = ref(0)
  const connectionFailed = ref(false)
  const backOff = ref(0)

  function redirectLogin() {
    store.isDaemonConfirmed = true
    router.push('/auth/login')
  }

  async function tryConnect() {
    connectionFailed.value = false
    backOff.value = 0

    for (let i = 0; i <= 10; i++) {
      attempts.value = i

      if (await daemonHelper.daemonIsStarted()) {
        console.log('connection success')
        redirectLogin()
        return
      }

      await setTimeout(backOff.value)
      backOff.value += 500
    }

    connectionFailed.value = true
  }

  onMounted(async () => {
    await tryConnect()
  })
</script>

<style lang="scss">
  .auth-layout {
    min-height: 100vh;
    background-image: linear-gradient(to right, var(--va-background-primary), var(--va-white));

    &__card {
      width: 100%;
      max-width: 600px;
    }
  }
</style>
