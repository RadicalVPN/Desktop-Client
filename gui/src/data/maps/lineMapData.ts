import { computed, ComputedRef, Ref } from '@vue/reactivity'
import { useColors } from 'vuestic-ui'

type GeoBounds = {
  bottom: number
  left: number
  right: number
  top: number
}

export type PointGeoCoord = {
  latitude: number
  longitude: number
}

export type DataGeometry = {
  geometry: {
    type: string
    coordinates: [number, number][]
  }
}

export type CityItem = {
  id?: string
  title: string
  country: string
  latitude: number
  longitude: number
  svgPath: string
  color: string
  flights?: PointGeoCoord[]
}

const targetSVG =
  'M9,0C4.029,0,0,4.029,0,9s4.029,9,9,9s9-4.029,9-9S13.971,0,9,0z M9,15.93 c-3.83,0-6.93-3.1-6.93-6.93S5.17,2.07,9,2.07s6.93,3.1,6.93,6.93S12.83,15.93,9,15.93 M12.5,9c0,1.933-1.567,3.5-3.5,3.5S5.5,10.933,5.5,9S7.067,5.5,9,5.5 S12.5,7.067,12.5,9z'

const london = {
  id: 'london',
  color: 'info',
  svgPath: targetSVG,
  title: 'London',
  country: 'United Kingdom',
  latitude: 51.5002,
  longitude: -0.1262,
  flights: [
    {
      latitude: 50.4422,
      longitude: 30.5367,
    },
    {
      latitude: 46.948,
      longitude: 7.4481,
    },
    {
      latitude: 59.3328,
      longitude: 18.0645,
    },
    {
      latitude: 40.4167,
      longitude: -3.7033,
    },
    {
      latitude: 46.0514,
      longitude: 14.506,
    },
    {
      latitude: 48.2116,
      longitude: 17.1547,
    },
    {
      latitude: 44.8048,
      longitude: 20.4781,
    },
    {
      latitude: 55.7558,
      longitude: 37.6176,
    },
    {
      latitude: 38.7072,
      longitude: -9.1355,
    },
    {
      latitude: 54.6896,
      longitude: 25.2799,
    },
    {
      latitude: 64.1353,
      longitude: -21.8952,
    },
    {
      latitude: 40.43,
      longitude: -74.0,
    },
  ],
}

const vilnius = {
  id: 'vilnius',
  color: 'info',
  svgPath: targetSVG,
  title: 'Vilnius',
  country: 'Lithuania',
  latitude: 54.6896,
  longitude: 25.2799,
  flights: [
    {
      latitude: 50.8371,
      longitude: 4.3676,
    },
    {
      latitude: 59.9138,
      longitude: 10.7387,
    },
    {
      latitude: 40.4167,
      longitude: -3.7033,
    },
    {
      latitude: 50.0878,
      longitude: 14.4205,
    },
    {
      latitude: 48.2116,
      longitude: 17.1547,
    },
    {
      latitude: 44.8048,
      longitude: 20.4781,
    },
    {
      latitude: 55.7558,
      longitude: 37.6176,
    },
    {
      latitude: 37.9792,
      longitude: 23.7166,
    },
    {
      latitude: 51.5002,
      longitude: -0.1262,
    },
    {
      latitude: 53.3441,
      longitude: -6.2675,
    },
  ],
}

const cities: CityItem[] = [
  london,
  vilnius,
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Brussels',
    country: 'Belgium',
    latitude: 50.8371,
    longitude: 4.3676,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Prague',
    country: 'Czech Republic',
    latitude: 50.0878,
    longitude: 14.4205,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Athens',
    country: 'Greece',
    latitude: 37.9792,
    longitude: 23.7166,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Reykjavik',
    country: 'Iceland',
    latitude: 64.1353,
    longitude: -21.8952,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Dublin',
    country: 'Ireland',
    latitude: 53.3441,
    longitude: -6.2675,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Oslo',
    country: 'Norway',
    latitude: 59.9138,
    longitude: 10.7387,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Lisbon',
    country: 'Portugal',
    latitude: 38.7072,
    longitude: -9.1355,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Moscow',
    country: 'Russia',
    latitude: 55.7558,
    longitude: 37.6176,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Belgrade',
    country: 'Serbia',
    latitude: 44.8048,
    longitude: 20.4781,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Bratislava',
    country: 'Slovakia',
    latitude: 48.2116,
    longitude: 17.1547,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Ljubljana',
    country: 'Slovenia',
    latitude: 46.0514,
    longitude: 14.506,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Madrid',
    country: 'Spain',
    latitude: 40.4167,
    longitude: -3.7033,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Stockholm',
    country: 'Sweden',
    latitude: 59.3328,
    longitude: 18.0645,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Groß Godems',
    country: 'Germany',
    latitude: 53.371943,
    longitude: 11.813868,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Bern',
    country: 'Switzerland',
    latitude: 46.948,
    longitude: 7.4481,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Kiev',
    country: 'Ukraine',
    latitude: 50.4422,
    longitude: 30.5367,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'Paris',
    country: 'France',
    latitude: 48.8567,
    longitude: 2.351,
  },
  {
    svgPath: targetSVG,
    color: 'info',
    title: 'New York',
    country: 'United States of America',
    latitude: 40.43,
    longitude: -74,
  },
]

export const lineMapData = {
  cities,
  mainCity: london.title,
  homeCity: london.title,
}

export const useMapData = (data: Ref<CityItem[]>): ComputedRef<CityItem[]> => {
  const { getColor } = useColors()

  return computed(() =>
    data.value.map((item) => ({
      ...item,
      color: getColor(item.color),
    })),
  )
}

export const getGeoBounds = (item?: CityItem): GeoBounds | undefined => {
  if (!item || !item.flights || !item.flights.length) {
    return
  }

  const latitudes = [...item.flights.map(({ latitude }) => latitude), item.latitude]
  const longitudes = [...item.flights.map(({ longitude }) => longitude), item.longitude]

  return {
    bottom: Math.min(...latitudes),
    left: Math.min(...longitudes),
    right: Math.max(...longitudes),
    top: Math.max(...latitudes),
  }
}

export const compareStrings = (first: string, second: string) => first.toLowerCase() === second.toLowerCase()
