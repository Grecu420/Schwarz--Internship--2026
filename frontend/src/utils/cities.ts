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
 * Escapes special characters for ODSQL double-quoted string literals.
 */
function escapeODSQL(value: string): string {
  return value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')
}


async function fetchCitiesFromAPI(
  whereClause: string,
  limit: number,
): Promise<City[]> {
  const url = `https://public.opendatasoft.com/api/explore/v2.1/catalog/datasets/geonames-all-cities-with-a-population-1000/records?where=${encodeURIComponent(whereClause)}&limit=${limit}`
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

export async function getCitiesByCountry(
  countryCode: string,
  minPopulation: number = 40000,
  limit: number = 100,
  name?: string,
): Promise<City[]> {
  const conditions: string[] = [
    `country_code="${escapeODSQL(countryCode.toUpperCase())}"`,
    `population > ${minPopulation}`,
  ]

  if (name && name.trim() !== '') {
    conditions.push(`search(name, "${escapeODSQL(name.trim())}")`)
  }

  return fetchCitiesFromAPI(conditions.join(' AND '), limit)
}