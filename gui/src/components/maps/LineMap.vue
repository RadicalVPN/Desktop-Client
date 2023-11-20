<template>
  <div ref="mapRef" class="line-map"></div>
</template>

<script setup lang="ts">
  import { ref, toRef, computed, onMounted, onBeforeUnmount, watch, shallowRef, onUpdated } from 'vue'
  import * as am5 from '@amcharts/amcharts5'
  import * as am5map from '@amcharts/amcharts5/map'
  import am5geodata_worldHigh from '@amcharts/amcharts5-geodata/worldHigh'
  import am5themes_Animated from '@amcharts/amcharts5/themes/Animated'
  import { useColors } from 'vuestic-ui'
  import { useGlobalStore } from '../../stores/global-store'

  import { useMapData, CityItem, getGeoBounds, compareStrings } from '../../data/maps/lineMapData'

  const props = withDefaults(
    defineProps<{
      mapData: CityItem[]
      homeCity: string
      modelValue: string
    }>(),
    {
      homeCity: 'New York',
      modelValue: 'New York',
    },
  )

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
  }>()

  const { colors } = useColors()
  const store = useGlobalStore()
  const mapRef = ref()
  const mapRoot = shallowRef()
  const mapChart = shallowRef()
  const mapPolygonSeries = shallowRef()
  const mapPointSeries = shallowRef()
  const mapLineSeries = shallowRef()
  const mapZoomControl = shallowRef()

  const mainCity = computed({
    get() {
      return props.modelValue
    },
    set(value) {
      emit('update:modelValue', value)
    },
  })

  const mapPointSeriesData = useMapData(toRef(props, 'mapData'))

  const getItemByMainCityTitle = () =>
    mapPointSeriesData.value.find(({ title }) => compareStrings(title, mainCity.value))

  const zoomToGeoBounds = () => {
    const item = getItemByMainCityTitle()
    const geoBounds = getGeoBounds(item)

    if (geoBounds) {
      mapChart.value.zoomToGeoBounds(geoBounds)
    }
  }

  const createMap = () => {
    const root = am5.Root.new(mapRef.value)

    if (store.animatedMap) {
      root.setThemes([am5themes_Animated.new(root)])
    }

    const chart = root.container.children.push(
      am5map.MapChart.new(root, {
        minZoomLevel: 1,
        maxZoomLevel: 20,
        zoomStep: store.animatedMap ? 3 : 2,
      }),
    )

    //click on background in map
    chart.chartContainer.get('background')?.events.on('click', () => {
      if (!store.vpnConnected) {
        mainCity.value = 'N/A' as any
      }
    })

    const zoomControl = chart.set('zoomControl', am5map.ZoomControl.new(root, {}))

    // polygon series
    const polygonSeries = chart.series.push(
      am5map.MapPolygonSeries.new(root, {
        geoJSON: am5geodata_worldHigh,
        exclude: ['AQ'], // i think we will never have servers in antarctica :)
      }),
    )

    polygonSeries.mapPolygons.template.setAll({
      fill: am5.color(colors.secondary),
      fillOpacity: 0.4,
      strokeWidth: 0.5,
      toggleKey: 'active',
      interactive: true,
    })

    let previousPolygon: am5map.MapPolygon | undefined
    polygonSeries.mapPolygons.template.on('active', function (active, target) {
      if (previousPolygon && previousPolygon != target) {
        previousPolygon.set('active', false)
      }
      if (target?.get('active') && target.dataItem) {
        polygonSeries.zoomToDataItem(target.dataItem as any)
      } else {
        chart.goHome()
      }
      previousPolygon = target
    })

    polygonSeries.events.on('datavalidated', zoomToGeoBounds)

    // point series
    const pointSeries = chart.series.push(
      am5map.MapPointSeries.new(root, {
        latitudeField: 'latitude',
        longitudeField: 'longitude',
      }),
    )

    // point series bullets
    const bulletTemplate = am5.Template.new({}) as am5.Template<am5.Graphics>
    bulletTemplate.events.on('click', (ev) => {
      mainCity.value = (ev as any).target._dataItem.dataContext.title
    })

    pointSeries.bullets.push((root, series, dataItem) => {
      const itemData = dataItem.dataContext as CityItem
      const isMainCity = compareStrings(itemData.title, mainCity.value)

      return am5.Bullet.new(root, {
        sprite: am5.Graphics.new(
          root,
          {
            svgPath: itemData.svgPath,
            x: am5.percent(100),
            y: am5.percent(100),
            centerX: am5.percent(100),
            centerY: am5.percent(100),
            fill: am5.color(isMainCity ? itemData.color : colors.primary),
            scale: isMainCity ? 2 : 1.25,
            tooltipText: '{title}',
          },
          bulletTemplate,
        ),
      })
    })

    if (store.showServerOnMap) {
      pointSeries.bullets.push((root, series, dataItem) => {
        const itemData = dataItem.dataContext as CityItem
        const labelText = itemData.title.split(' - ')[1]

        return am5.Bullet.new(root, {
          sprite: am5.Label.new(root, {
            fill: am5.color(colors.textPrimary),
            text: labelText,
            populateText: true,
            fontWeight: 'bold',
            fontSize: 10,
            centerY: am5.p100,
          }),
        })
      })
    }

    // set map data
    pointSeries.data.setAll(mapPointSeriesData.value)

    // assign objects to refs
    mapRoot.value = root
    mapChart.value = chart
    mapZoomControl.value = zoomControl
    mapPointSeries.value = pointSeries
    mapPolygonSeries.value = polygonSeries
  }

  const setPointSeriesData = () => {
    mapPointSeries.value.data.setAll(mapPointSeriesData.value)
  }

  const updateChartDataOnChangeTheme = () => {
    if (mapRoot.value) {
      mapPolygonSeries.value.mapPolygons.template.setAll({
        fill: am5.color(colors.secondary),
      })

      mapLineSeries.value.mapLines.template.setAll({
        stroke: am5.color(colors.primary),
      })

      setPointSeriesData()
    }
  }

  const updateChartDataOnUpdateProps = () => {
    if (mapRoot.value) {
      setPointSeriesData()
      zoomToGeoBounds()
    }
  }

  const disposeMap = () => {
    if (mapRoot.value) {
      mapRoot.value.dispose()
    }
  }

  onMounted(createMap)
  onUpdated(updateChartDataOnUpdateProps)
  watch(colors, updateChartDataOnChangeTheme)
  onBeforeUnmount(disposeMap)
</script>

<style lang="scss" scoped>
  .line-map {
    width: 100%;
    height: 100%;
    border-radius: inherit;

    :deep(div),
    :deep(canvas) {
      border-radius: inherit;
    }
  }
</style>
