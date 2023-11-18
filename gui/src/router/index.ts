import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

import AuthLayout from '../layouts/AuthLayout.vue'
import AppLayout from '../layouts/AppLayout.vue'
import UIRoute from '../pages/admin/ui/route'
import DaemonConnect from '../pages/daemon/ConnectDaemon.vue'
import DaemonInstall from '../pages/daemon/DaemonInstall.vue'
import { DaemonHelper } from '../helper/daemon'
import { Server, useGlobalStore } from '../stores/global-store'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/:catchAll(.*)',
    redirect: { name: 'daemon-install' },
  },
  {
    name: 'admin',
    path: '/admin',
    component: AppLayout,
    beforeEnter: async (to, from, next) => {
      const store = useGlobalStore()

      await store.loadServerList()

      //reload server list every minute
      await setInterval(store.loadServerList, 60 * 1000)

      if (!store.isDaemonConfirmed) {
        next('/daemon-install')
        return
      }

      next()
    },
    children: [
      {
        name: 'dashboard',
        path: 'dashboard',
        component: () => import('../pages/admin/serverpage/ServerMap.vue'),
      },
      {
        name: 'privacy-firewall',
        path: 'privacy-firewall',
        component: () => import('../pages/admin/privacyfirewall/PrivacyFirewall.vue'),
      },
      {
        name: 'settings',
        path: 'settings',
        component: () => import('../pages/admin/settings/Settings.vue'),
      },
      UIRoute,
    ],
  },
  {
    name: 'daemon',
    path: '/daemon',
    component: DaemonConnect,
  },
  {
    name: 'daemon-install',
    path: '/daemon-install',
    component: DaemonInstall,
  },
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        name: 'login',
        path: 'login',
        component: () => import('../pages/auth/login/Login.vue'),
      },
      {
        name: 'signup',
        path: 'signup',
        component: () => import('../pages/auth/signup/Signup.vue'),
      },
      {
        name: 'recover-password',
        path: 'recover-password',
        component: () => import('../pages/auth/recover-password/RecoverPassword.vue'),
      },
      {
        path: '',
        redirect: { name: 'login' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    // For some reason using documentation example doesn't scroll on page navigation.
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    } else {
      document.querySelector('.app-layout__page')?.scrollTo(0, 0)
    }
  },
  routes,
})

export default router
