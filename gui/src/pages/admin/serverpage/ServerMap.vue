<template>
  <line-map v-model="mainCity" :map-data="cities" :home-city="homeCity" />

  <div class="absolute top-1/4 ml-4 transform -translate-y-1/4 w-50 h-80 pt-6">
    <va-card>
      <va-card-content>
        <div class="flex items-center pb-6">
          <va-icon :color="isConnected ? 'success' : 'danger'" :name="isConnected ? 'fa-lock' : 'fa-lock-open'" />
          <p class="pl-2 text--secondary" :style="{ color: isConnected ? 'success' : 'danger' }">Not Connected</p>
        </div>
        <va-button>Connect</va-button>

        <va-divider class="pt-4 pb-4" />

        <va-input class="pb-4" placeholder="Server Name" />
        <div
          v-for="(option, id) in store.serverList"
          :key="id"
          class="server__item flex flex-1 flex-wrap items-center pt-1 pb-1 mt-2 mb-2"
          @click="console.log(option.city)"
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
  import { ref } from 'vue'
  import LineMap from '../../../components/maps/LineMap.vue'
  import { lineMapData } from '../../../data/maps/lineMapData'
  import { useGlobalStore } from '../../../stores/global-store'

  const cities = ref(lineMapData.cities)
  const mainCity = ref('Vilnius')
  const homeCity = ref('Vilnius')
  const store = useGlobalStore()
  const isConnected = ref(false)
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
