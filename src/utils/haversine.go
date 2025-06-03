package haversine

import "math"

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Asin(math.Sqrt(a))

	return R * c * 1000 // haversine distance in meteres
}
