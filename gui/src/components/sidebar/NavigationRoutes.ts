export interface INavigationRoute {
  name: string
  displayName: string
  meta: { icon: string }
  children?: INavigationRoute[]
}

export default {
  root: {
    name: '/',
    displayName: 'navigationRoutes.home',
  },
  routes: [
    {
      name: 'dashboard',
      displayName: 'menu.dashboard',
      meta: {
        icon: 'vuestic-iconset-maps',
      },
    },
    {
      name: 'privacy-firewall',
      displayName: 'menu.privacyFirewall',
      meta: {
        icon: 'vuestic-iconset-tables',
      },
    },
    {
      name: 'settings',
      displayName: 'menu.settings',
      meta: {
        icon: 'vuestic-iconset-settings',
      },
    },
  ] as INavigationRoute[],
}
