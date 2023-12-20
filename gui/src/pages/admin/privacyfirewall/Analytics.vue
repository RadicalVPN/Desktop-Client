<template>
  <div class="grid grid-cols-12 gap-6 pb-6">
    <va-card
      v-for="(statistic, index) in statistics"
      :key="index"
      class="col-span-12 sm:col-span-6 md:col-span-4 lg:col-span-2"
      :color="statistic.color"
    >
      <va-card-content class="flex items-center">
        <div class="flex-1">
          <p class="text-white pb-2">{{ statistic.title }}</p>
          <VaSkeleton v-if="isLoading" variant="text" class="va-h2 m-1 w-20" />
          <h2 v-else class="va-h2 m-0 text-white">{{ statistic.value }}</h2>
        </div>
        <div class="absolute top-10 right-3">
          <va-icon :name="statistic.icon" class="opcacity-25" color="secondary" :size="60" />
        </div> </va-card-content
    ></va-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue'
  import { DaemonHelper } from '../../../helper/daemon'
  import { onBeforeUnmount } from 'vue'
  const isLoading = ref(false)
  let refreshTimer: any
  const statistics = ref([
    {
      title: 'Total queries',
      value: '0',
      icon: 'fa-globe',
      color: '#33a65a',
    },
    {
      title: 'Queries Blocked',
      value: '0',
      icon: 'fa-hand',
      color: '#41c0ef',
    },
    {
      title: 'Percent Blocked',
      value: '0',
      icon: 'fa-chart-pie',
      color: '#f39c12',
    },
  ])

  async function updateStatistics() {
    console.log('updating privacy firewall stats')
    isLoading.value = true

    const stats = await new DaemonHelper().getPrivacyFirewallStats()
    if (stats.status) {
      const data = stats.data
      const total = data.total
      const blocked = data.blocked

      statistics.value[0].value = total.toString()
      statistics.value[1].value = blocked.toString()
      statistics.value[2].value = `${((blocked / total) * 100).toFixed(2).toString()}%`
    }

    isLoading.value = false
  }

  onMounted(async () => {
    await updateStatistics()
    refreshTimer = setInterval(async () => {
      await updateStatistics()
    }, 30_000)
  })

  onBeforeUnmount(() => {
    clearInterval(refreshTimer)
  })
</script>
