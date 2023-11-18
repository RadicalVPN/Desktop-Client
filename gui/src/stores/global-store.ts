import { defineStore } from 'pinia'
import { useColors } from 'vuestic-ui'
import { useStorage } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import languages from '../i18n/languages'
import { DaemonHelper } from '../helper/daemon'

export interface Server {
  id: string
  hostname: string
  country: string
  city: string
  internal_ip: string
  external_ip: string
  public_key: string
  online: boolean
  latency: number
  country_name: string
  latitude: string
  longitude: string
  location: string
}

export interface Location {
  id: string
  name: string
  country: string
  country_name: string
  city: string
  latitude: string
  longitude: string
  latency: number
}

export const useGlobalStore = defineStore('global', {
  state: () => {
    const { applyPreset } = useColors()
    const { locale } = useI18n()

    const theme = localStorage.getItem('theme') || 'light'
    applyPreset(theme)

    //load i18n on start
    let i18n = localStorage.getItem('language')
    if (i18n) {
      if (i18n === '_system') {
        const sysLang = navigator.language

        //valid locale?
        const isValidLocale = languages.some((language) => language.text === sysLang)

        if (isValidLocale) {
          i18n = sysLang
        } else {
          console.warn('no translation for system language', sysLang)
          i18n = 'gb'
        }
      }

      locale.value = i18n
    }

    return {
      isSidebarMinimized: false,
      userName: 'Vasili S',
      theme: theme,
      serverList: [] as Server[],
      locationList: [] as Location[],
      isDaemonConfirmed: false,
      vpnConnected: false,
      animatedMap: useStorage('animatedMap', true),
      disableNotifications: useStorage('disableNotifications', false),
      privacyFirewallLevel: useStorage('privacyFirewall', 'basic'),
      language: useStorage('language', 'gb'),
      showServerOnMap: useStorage('showCountryOnMap', false),
      mainCity: 'N/A',
      auth: {
        isAuthChecking: true,
      },
    }
  },

  actions: {
    toggleSidebar() {
      this.isSidebarMinimized = !this.isSidebarMinimized
    },
    async loadServerList() {
      try {
        this.serverList = (await new DaemonHelper().getServerList()).filter((server: Server) => server.online)
        this.computeLocationList()
      } catch (e) {
        console.log('failed to load server list', e)
      }
    },
    computeLocationList() {
      this.locationList = Object.values(
        this.serverList.reduce((acc: any, server: Server) => {
          if (!acc[server.location]) {
            acc[server.location] = {
              id: server.location,
              name: server.hostname,
              country: server.country,
              country_name: server.country_name,
              city: server.city,
              latitude: server.latitude,
              longitude: server.longitude,
              latency: server.latency,
            }
          }

          return acc
        }, {}),
      ).map((location: any) => {
        return {
          ...location,
        }
      })
    },
    changeUserName(userName: string) {
      this.userName = userName
    },
  },
})
