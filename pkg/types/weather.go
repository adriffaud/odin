package types

type WeatherData struct {
	Current        CurrentWeather `json:"current"`
	CurrentUnits   CurrentUnits   `json:"current_units"`
	Hourly         HourlyWeather  `json:"hourly"`
	HourlyUnits    HourlyUnits    `json:"hourly_units"`
	Daily          DailyWeather   `json:"daily"`
	DailyUnits     DailyUnits     `json:"daily_units"`
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	Elevation      float64        `json:"elevation"`
	GenerationTime float64        `json:"generationtime_ms"`
	Timezone       string         `json:"timezone"`
	TimezoneAbbr   string         `json:"timezone_abbreviation"`
}

type CurrentWeather struct {
	Time                     string  `json:"time"`
	Temperature              float64 `json:"temperature_2m"`
	RelativeHumidity         int     `json:"relative_humidity_2m"`
	CloudCover               int     `json:"cloud_cover"`
	WindSpeed                float64 `json:"wind_speed_10m"`
	WindDirection            float64 `json:"wind_direction_10m"`
	PrecipitationProbability int     `json:"precipitation_probability"`
	DewPoint                 float64 `json:"dew_point_2m"`
}

type CurrentUnits struct {
	Temperature              string `json:"temperature_2m"`
	RelativeHumidity         string `json:"relative_humidity_2m"`
	CloudCover               string `json:"cloud_cover"`
	WindSpeed                string `json:"wind_speed_10m"`
	WindDirection            string `json:"wind_direction_10m"`
	PrecipitationProbability string `json:"precipitation_probability"`
	DewPoint                 string `json:"dew_point_2m"`
}

type HourlyWeather struct {
	Time                     []string  `json:"time"`
	Temperature              []float64 `json:"temperature_2m"`
	RelativeHumidity         []int     `json:"relative_humidity_2m"`
	CloudCover               []int     `json:"cloud_cover"`
	CloudCoverLow            []int     `json:"cloud_cover_low"`
	CloudCoverMid            []int     `json:"cloud_cover_mid"`
	CloudCoverHigh           []int     `json:"cloud_cover_high"`
	WindSpeed                []float64 `json:"wind_speed_10m"`
	WindDirection            []float64 `json:"wind_direction_10m"`
	PrecipitationProbability []int     `json:"precipitation_probability"`
	DewPoint                 []float64 `json:"dew_point_2m"`
}

type HourlyUnits struct {
	Temperature              string `json:"temperature_2m"`
	RelativeHumidity         string `json:"relative_humidity_2m"`
	CloudCover               string `json:"cloud_cover"`
	CloudCoverLow            string `json:"cloud_cover_low"`
	CloudCoverMid            string `json:"cloud_cover_mid"`
	CloudCoverHigh           string `json:"cloud_cover_high"`
	WindSpeed                string `json:"wind_speed_10m"`
	WindDirection            string `json:"wind_direction_10m"`
	PrecipitationProbability string `json:"precipitation_probability"`
	DewPoint                 string `json:"dew_point_2m"`
}

type DailyWeather struct {
	Time    []string `json:"time"`
	Sunrise []string `json:"sunrise"`
	Sunset  []string `json:"sunset"`
}

type DailyUnits struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type WeatherResultMsg struct {
	Data WeatherData
}
