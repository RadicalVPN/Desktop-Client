import { defineStore } from 'pinia'

export const useGlobalStore = defineStore('global', {
  state: () => {
    return {
      isSidebarMinimized: false,
      userName: 'Vasili S',
      theme: 'light',
      serverList: [] as {
        id: string
        hostname: string
        country: string
        city: string
        internal_ip: string
        external_ip: string
        public_key: string
        online: boolean
        latency: number
      }[],
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
