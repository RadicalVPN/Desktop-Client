import { defineStore } from 'pinia'
import { useColors } from 'vuestic-ui'
import { useStorage } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import languages from '../i18n/languages'

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
      isDaemonConfirmed: false,
      vpnConnected: false,
      animatedMap: useStorage('animatedMap', true),
      disableNotifications: useStorage('disableNotifications', false),
      privacyFirewallLevel: useStorage('privacyFirewall', 'basic'),
      language: useStorage('language', 'gb'),
    }
  },

  actions: {
    toggleSidebar() {
      this.isSidebarMinimized = !this.isSidebarMinimized
    },
    changeUserName(userName: string) {
      this.userName = userName
    },
  },
})
