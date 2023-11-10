<template>
  <line-map v-model="mainCity" :map-data="cities" />

  <div class="absolute top-1/4 ml-4 transform -translate-y-1/4 w-50 h-80 pt-6">
    <va-card>
      <va-card-content>
        <div class="flex items-center pb-6">
          <va-icon
            :color="store.vpnConnected ? 'success' : 'danger'"
            :name="store.vpnConnected ? 'fa-lock' : 'fa-lock-open'"
          />
          <p class="pl-2 text--secondary font-bold" :style="{ color: store.vpnConnected ? 'success' : 'danger' }">
            {{ store.vpnConnected ? 'Connected' : 'Disconnected' }}
          </p>
        </div>

        <div v-if="mainCity != 'N/A'">
          <p class="pb-4">Selected Server: {{ mainCity }}</p>

          <va-button v-if="!store.vpnConnected" :loading="isConnectionStateSwitching" @click="connect()"
            >Connect</va-button
          >
          <va-button v-if="store.vpnConnected" :loading="isConnectionStateSwitching" @click="disconnect()"
            >Disconnect</va-button
          >
        </div>
        <div v-else>
          <va-button v-if="!store.vpnConnected" :loading="isConnectionStateSwitching" @click="fastConnect()"
            >Fast connect (Fastest Server)</va-button
          >
        </div>

        <va-divider class="pt-4 pb-4" />

        <va-input class="pb-4" placeholder="Server Name" />
        <div
          v-for="(option, id) in store.serverList"
          :key="id"
          class="server__item flex flex-1 flex-wrap items-center pt-1 pb-1 mt-2 mb-2"
        >
          <va-icon :name="`flag-icon-${option.country} small`" />
          <span class="dropdown-item__text pl-4">
            {{ option.city }}
          </span>

          <div class="ml-auto text-sm flex items-center">
            <a class="ml-4">{{ option.latency + ' ms' }}</a>
          </div>

          <va-divider />
        </div>
      </va-card-content>
    </va-card>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import LineMap from '../../../components/maps/LineMap.vue'
  import { useGlobalStore } from '../../../stores/global-store'
  import { targetSVG } from '../../../data/maps/lineMapData'
  import { DaemonHelper } from '../../../helper/daemon'
  import { useModal } from 'vuestic-ui'
  import { useI18n } from 'vue-i18n'

  async function fastConnect() {
    isConnectionStateSwitching.value = true

    const fastestServer = store.serverList.sort((a, b) => a.latency - b.latency)[0]
    mainCity.value = `${fastestServer.country_name} - ${fastestServer.city}`

    await connect()
  }

  async function connect() {
    isConnectionStateSwitching.value = true

    //parse the the server from the selection
    const split = mainCity.value.split(' - ')
    const countryName = split[0]
    const cityName = split[1]

    //get the server from the store
    const server = store.serverList.find((server) => server.city === cityName && server.country_name === countryName)

    if (!server) {
      console.error('server not found')
      return
    }

    const res = await new DaemonHelper().connectToServer(server.id)

    if (res.status === true) {
      store.vpnConnected = true
    }

    if (res.status === false && res.data.error === 'vpn connection limit') {
      await confirm({
        title: t('vpn.connectionLimit.title'),
        message: t('vpn.connectionLimit.description'),

        hideDefaultActions: true,
        closeButton: true,
        blur: true,
      })
    } else {
      if (!store.disableNotifications) {
        new Notification('RadicalVPN', {
          body: `${t('notifications.vpn.connect')} ${mainCity.value}`,
        })
      }
    }

    isConnectionStateSwitching.value = false
  }

  async function disconnect() {
    isConnectionStateSwitching.value = true

    const res = await new DaemonHelper().disconnectFromServer()

    if (res) {
      if (!store.disableNotifications) {
        new Notification('RadicalVPN', {
          body: `${t('notifications.vpn.disconnect')}`,
        })
      }
      store.vpnConnected = false
    }

    isConnectionStateSwitching.value = false
  }

  async function syncConnectionState() {
    const res = await new DaemonHelper().getConnectionState()

    if (res) {
      store.vpnConnected = true
      mainCity.value = 'Unknown'
    } else {
      store.vpnConnected = false
    }
  }

  const store = useGlobalStore()
  const mainCity = ref('N/A')
  const isConnectionStateSwitching = ref(false)
  const cities = ref(
    store.serverList.map((server) => ({
      color: 'info',
      title: `${server.country_name} - ${server.city}`,
      country: server.country,
      latitude: parseInt(server.latitude),
      longitude: parseInt(server.longitude),
      svgPath: targetSVG,
    })),
  )

  const { confirm } = useModal()
  const { t } = useI18n()

  onMounted(async () => {
    await syncConnectionState()
  })
</script>

<style lang="scss" scoped>
  @import 'flag-icons/css/flag-icons.css';

  .server {
    &__item {
      cursor: pointer;
      flex-wrap: nowrap;

      &:last-of-type {
        padding-bottom: 0 !important;
      }

      &:hover {
        color: var(--va-primary);
      }
    }
  }
</style>
