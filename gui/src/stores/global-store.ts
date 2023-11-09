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

    return {
      isSidebarMinimized: false,
      userName: 'Vasili S',
      theme: theme,
      serverList: [] as Server[],
      isDaemonConfirmed: false,
      vpnConnected: false,
      animatedMap: animatedMap,
    }
  },

  actions: {
    toggleSidebar() {
      this.isSidebarMinimized = !this.isSidebarMinimized
    },
    setMapAnimation() {
      localStorage.setItem('animatedMap', JSON.stringify(this.animatedMap))
    },
    changeUserName(userName: string) {
      this.userName = userName
    },
  },
})
