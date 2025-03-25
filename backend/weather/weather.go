package weather

import "time"

type WeatherInfo struct {
	Location      string    `json:"location"`
	TemperatureC  float64   `json:"temperature_c"`
	Conditions    string    `json:"conditions"`
	ForecastSunHr float64   `json:"forecast_sun_hours"`
	Timestamp     time.Time `json:"timestamp"`
}

func GetMockWeather() WeatherInfo {
	return WeatherInfo{
		Location:      "Chester",
		TemperatureC:  21.2,
		Conditions:    "Partly Cloudy",
		ForecastSunHr: 8,
		Timestamp:     time.Now(),
	}
}
