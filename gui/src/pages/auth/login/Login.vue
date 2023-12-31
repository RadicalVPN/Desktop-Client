<template>
  <form @submit.prevent="onsubmit">
    <va-input
      v-model="email"
      class="mb-4"
      type="email"
      :label="t('auth.email')"
      :error="!!emailErrors.length"
      :error-messages="emailErrors"
    />

    <va-input
      v-model="password"
      class="mb-4"
      type="password"
      :label="t('auth.password')"
      :error="!!passwordErrors.length"
      :error-messages="passwordErrors"
    />

    <va-input
      v-if="showTotp === true"
      v-model="totp"
      class="mb-4"
      type="text"
      :label="t('auth.totpCode')"
      :error="!!totpErrors.length"
      :error-messages="totpErrors"
    />

    <div class="flex justify-center mt-4">
      <va-button class="my-0" @click="onsubmit">{{ t('auth.login') }}</va-button>
    </div>
  </form>
</template>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useI18n } from 'vue-i18n'
  import axios from 'axios'
  import { DaemonCredentials } from '../../../helper/credentials'

  const { t } = useI18n()

  const email = ref('')
  const password = ref('')
  const emailErrors = ref<string[]>([])
  const passwordErrors = ref<string[]>([])

  const showTotp = ref<boolean>(false)
  const totp = ref('')
  const totpErrors = ref<string[]>([])

  const router = useRouter()

  const formReady = computed(() => !emailErrors.value.length && !passwordErrors.value.length)

  async function onsubmit() {
    emailErrors.value = email.value ? [] : ['Email is required']
    passwordErrors.value = password.value ? [] : ['Password is required']

    if (!formReady.value) return

    const credentials = DaemonCredentials.getCredentials()
    // await axios.get('https://radicalvpn.com/geoip/current')
    const resp = await axios.post(
      `http://localhost:${credentials.port}/login`,
      {
        email: email.value,
        password: password.value,
        turnstileChallenge: '',
        ...(totp.value && {
          totpToken: totp.value,
        }),
      },
      {
        headers: {
          'x-radical-daemon-secret': credentials.secret,
        },
        validateStatus: () => true,
      },
    )

    if (resp.status === 200) {
      console.log('logged in')
      router.push({ name: 'dashboard' })
    } else {
      if (resp.data == 'totp required') {
        showTotp.value = true
        totpErrors.value = ['TOTP is required']
      } else if (resp.data === 'invalid totp token') {
        totpErrors.value = ['Invalid TOTP token']
      } else if (resp.data === 'turnstile challenge failed') {
        const error = ['Captcha (Turnstile) challenge failed']

        emailErrors.value = error
        passwordErrors.value = error
      } else {
        emailErrors.value = ['Invalid email or password']
        passwordErrors.value = ['Invalid email or password']
      }

      console.log('failed to login')
    }
  }
</script>
