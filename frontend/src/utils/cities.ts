interface CityCoordinates {
  lon: number
  lat: number
}

interface OpenDataSoftRecord {
  name: string
  country_code: string
  cou_name_en: string
  coordinates: CityCoordinates
  population?: number
  timezone?: string
}

interface OpenDataSoftResponse {
  total_count: number
  results: OpenDataSoftRecord[]
}

export interface City {
  name: string
  countryCode: string
  latitude: number
  longitude: number
  population?: number
}

/**
 * Fetches cities for a given 2-letter ISO country code.
 *
 * @param countryCode Two-letter ISO country code (e.g., "JP", "FR", "US")
 * @param limit Maximum number of records to return (default: 100, max: 100 per page)
 */
export async function getCitiesByCountry(
  countryCode: string,
  minPopulation: number = 40000,
  limit: number = 100,
): Promise<City[]> {
  const where = `country_code="${countryCode.toUpperCase()}" AND population > ${minPopulation}`
  const url = `https://public.opendatasoft.com/api/explore/v2.1/catalog/datasets/geonames-all-cities-with-a-population-1000/records?where=${encodeURIComponent(where)}&limit=${limit}`
  const response = await fetch(url)

  if (!response.ok) {
    throw new Error(`Failed to fetch cities: ${response.status} ${response.statusText}`)
  }

  const data: OpenDataSoftResponse = await response.json()

  return data.results.map((record) => ({
    name: record.name,
    countryCode: record.country_code,
    latitude: record.coordinates.lat,
    longitude: record.coordinates.lon,
    population: record.population,
  }))
}
