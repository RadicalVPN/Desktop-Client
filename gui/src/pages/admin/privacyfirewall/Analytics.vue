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
          <h2 class="va-h2 m-0 text-white">{{ statistic.value }}</h2>
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
  let refreshTimer: any
  const statistics = ref([
    {
      title: 'Total queries',
      value: 'N/A',
      icon: 'fa-globe',
      color: '#33a65a',
    },
    {
      title: 'Queries Blocked',
      value: 'N/A',
      icon: 'fa-hand',
      color: '#41c0ef',
    },
    {
      title: 'Percent Blocked',
      value: 'N/A',
      icon: 'fa-chart-pie',
      color: '#f39c12',
    },
  ])

  async function updateStatistics() {
    const stats = await new DaemonHelper().getPrivacyFirewallStats()
    const formatter = new Intl.NumberFormat()

    console.log('updating privacy firewall stats')

    if (stats.status) {
      const data = stats.data
      const total = data.total
      const blocked = data.blocked

      statistics.value[0].value = formatter.format(total)
      statistics.value[1].value = formatter.format(blocked)
      statistics.value[2].value = `${((blocked / total) * 100).toFixed(2).toString()}%`
    }
  }

  onMounted(async () => {
    await updateStatistics()
    refreshTimer = setInterval(async () => {
      await updateStatistics()
    }, 30_000) //backend caches 60 secs; let's refresh every 30 secs
  })

  onBeforeUnmount(() => {
    clearInterval(refreshTimer)
  })
</script>