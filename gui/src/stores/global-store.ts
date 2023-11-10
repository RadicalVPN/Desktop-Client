import { defineStore } from 'pinia'
import { useColors } from 'vuestic-ui'

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

    const theme = localStorage.getItem('theme') || 'light'
    applyPreset(theme)

    const localMapSettings = localStorage.getItem('animatedMap')
    const animatedMap = localMapSettings ? JSON.parse(localMapSettings) : true

    const localNotificationsSettings = localStorage.getItem('disableNotifications')
    const disableNotifications = localNotificationsSettings ? JSON.parse(localNotificationsSettings) : false

    return {
      isSidebarMinimized: false,
      userName: 'Vasili S',
      theme: theme,
      serverList: [] as Server[],
      isDaemonConfirmed: false,
      vpnConnected: false,
      animatedMap: animatedMap,
      disableNotifications: disableNotifications,
    }
  },

  actions: {
    toggleSidebar() {
      this.isSidebarMinimized = !this.isSidebarMinimized
    },
    setMapAnimation() {
      localStorage.setItem('animatedMap', JSON.stringify(this.animatedMap))
    },
    setNotifications() {
      localStorage.setItem('disableNotifications', JSON.stringify(this.disableNotifications))
    },
    changeUserName(userName: string) {
      this.userName = userName
    },
  },
})
