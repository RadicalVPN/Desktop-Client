<template>
  <va-sidebar :width="width" :minimized="minimized" :minimized-width="minimizedWidth" :animated="animated">
    <template v-for="route in items" :key="route.title">
      <VaSidebarItem :active="route.name === useRoute().name" :to="{ name: route.name }">
        <VaSidebarItemContent>
          <va-icon :name="route.meta.icon" class="va-sidebar-item__icon" />
        </VaSidebarItemContent>
      </VaSidebarItem>

      <va-spacer v-if="route.name === items.slice(-2)[0]?.name" />
    </template>
  </va-sidebar>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import NavigationRoutes from './NavigationRoutes'
  import { useRoute } from 'vue-router'

  withDefaults(
    defineProps<{
      width?: string
      color?: string
      animated?: boolean
      minimized?: boolean
      minimizedWidth?: string
    }>(),
    {
      width: '16rem',
      color: 'secondary',
      animated: true,
      minimized: true,
      minimizedWidth: undefined,
    },
  )

  const items = ref(NavigationRoutes.routes)
</script>

<style lang="scss">
  .va-sidebar {
    &__menu {
      padding: 1rem 0 0 0;
    }

    &-item {
      filter: blur(50);
      &__icon {
        width: 1.5rem;
        height: 1.5rem;
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }
  }
</style>
