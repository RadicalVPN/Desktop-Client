<template>
  <div class="auth-layout grid grid-cols-12 content-center">
    <div class="flex col-span-12 p-4 justify-center">
      <radical-logo height="80" />
    </div>

    <div class="flex justify-center col-span-12 p-4">
      <va-card class="auth-layout__card">
        <va-card-content>
          <div class="text-center font-bold">
            <p v-if="installRequired">
              Good news everyone! 🚀<br /><br />RadicalVPN is yearning for a little buddy (Daemon) to make its day.<br />Follow
              your Mac's lead, and you'll be all set to conquer the virtual world!
            </p>

            <p v-if="isConnectingToDaemon">Connecting to RadicalVPN daemon...</p>

            <p class="text-center font-bold text-xl" v-if="daemonInstalledFailed">
              ❗️❗️ The daemon was not installed. ❗️❗️
            </p>
            <br />
            <p v-if="daemonInstalledFailed">
              You may restart the client and try again. Or contact our customer support.<br /><br />
            </p>
            <a class="content-center font-bold underline" href="mailto:support@radicalvpn.com">
              support@radicalvpn.com
            </a>
          </div>
        </va-card-content>
      </va-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import RadicalLogo from '../../components/RadicalLogo.vue'
  import { DaemonHelper } from '../../helper/daemon'
  import { useRouter } from 'vue-router'
  import { useGlobalStore } from '../../stores/global-store'

  const daemonHelper = new DaemonHelper()
  const router = useRouter()
  const store = useGlobalStore()

  const installRequired = ref(false)
  const isConnectingToDaemon = ref(false)
  const daemonInstalledFailed = ref(false)

  function redirectLogin() {
    store.isDaemonConfirmed = true
    router.push('/auth/login')
  }

  async function load() {
    try {
      installRequired.value = await daemonHelper.isDaemonInstallRequired()
    } catch (error) {
      console.log(error)
    }

    console.log('installing daemon')
    if (installRequired.value === true) {
      try {
        if (await daemonHelper.installDaemon()) {
          redirectLogin()
        } else {
          installRequired.value = false
          daemonInstalledFailed.value = true
        }
      } catch (error) {
        installRequired.value = false
        daemonInstalledFailed.value = true
      }
    } else {
      redirectLogin()
    }

    return
  }

  load()
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
