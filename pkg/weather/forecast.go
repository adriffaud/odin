package weather

import (
	"math"
	"time"

	"driffaud.fr/odin/pkg/util"
)

// ForecastHour represents a single hour of forecast data
type ForecastHour struct {
	DateTime                 time.Time
	Hour                     int
	Clouds                   int
	CloudsLow                int
	CloudsMid                int
	CloudsHigh               int
	Temperature              float64
	WindSpeed                float64
	WindDirection            float64
	Humidity                 int
	DewPoint                 float64
	PrecipitationProbability int
	Rating                   int
	Seeing                   int
}

// BestObservationInfo represents the best time range for astronomical observation
type BestObservationInfo struct {
	TimeRange        *TimeRange
	LowestCloudCover int
}

// TimeRange represents a start and end hour for observation
type TimeRange struct {
	Start int
	End   int
}

// NightForecast contains forecast and analysis for an astronomical night
type NightForecast struct {
	BestObservation      BestObservationInfo
	ExtremeCloudCover    int
	DisplayCloudCover    int
	NightlyTemperature   int
	NightlyHumidity      int
	NightlyWindSpeed     int
	NightlyDewPoint      int
	MaxPrecipProbability int
	NightlyWindDirection int
	WindDirectionText    string
	SeeingIndex          int
}

// GenerateForecastData converts Open-Meteo weather data into a slice of hourly forecast data
func GenerateForecastData(data WeatherData) []ForecastHour {
	if len(data.Hourly.Time) == 0 {
		return []ForecastHour{}
	}

	forecast := make([]ForecastHour, len(data.Hourly.Time))

	for i, timeStr := range data.Hourly.Time {
		dateTime, _ := time.Parse(util.ISO8601Format, timeStr)

		var clouds, cloudsLow, cloudsMid, cloudsHigh int
		var temp, windSpeed, windDir, dewPoint float64
		var humidity, precipProb int

		if i < len(data.Hourly.CloudCover) {
			clouds = data.Hourly.CloudCover[i]
		}
		if i < len(data.Hourly.CloudCoverLow) {
			cloudsLow = data.Hourly.CloudCoverLow[i]
		}
		if i < len(data.Hourly.CloudCoverMid) {
			cloudsMid = data.Hourly.CloudCoverMid[i]
		}
		if i < len(data.Hourly.CloudCoverHigh) {
			cloudsHigh = data.Hourly.CloudCoverHigh[i]
		}
		if i < len(data.Hourly.Temperature) {
			temp = data.Hourly.Temperature[i]
		}
		if i < len(data.Hourly.WindSpeed) {
			windSpeed = data.Hourly.WindSpeed[i]
		}
		if i < len(data.Hourly.WindDirection) {
			windDir = data.Hourly.WindDirection[i]
		}
		if i < len(data.Hourly.RelativeHumidity) {
			humidity = data.Hourly.RelativeHumidity[i]
		}
		if i < len(data.Hourly.DewPoint) {
			dewPoint = data.Hourly.DewPoint[i]
		}
		if i < len(data.Hourly.PrecipitationProbability) {
			precipProb = data.Hourly.PrecipitationProbability[i]
		}

		seeingIndex := calculateSeeingIndex(temp, dewPoint, windSpeed, humidity)
		ratingIndex := calculateSkyQualityIndex(clouds, humidity, windSpeed, temp, dewPoint, seeingIndex)

		forecast[i] = ForecastHour{
			DateTime:                 dateTime,
			Hour:                     dateTime.Hour(),
			Clouds:                   clouds,
			CloudsLow:                cloudsLow,
			CloudsMid:                cloudsMid,
			CloudsHigh:               cloudsHigh,
			Temperature:              temp,
			WindSpeed:                windSpeed,
			WindDirection:            windDir,
			Humidity:                 humidity,
			DewPoint:                 dewPoint,
			PrecipitationProbability: precipProb,
			Rating:                   ratingIndex,
			Seeing:                   seeingIndex,
		}
	}

	return forecast
}

// AnalyzeNightForecast generates a complete night forecast analysis for astronomical observation
func AnalyzeNightForecast(forecastData []ForecastHour, sunsetTime, sunriseTime time.Time) NightForecast {
	const goodCloudCoverThreshold = 30
	const consecutiveGoodHoursRequired = 2

	nightForecastData := filterNightForecastData(forecastData, sunsetTime, sunriseTime)

	bestObservationInfo := getBestObservationTimeRange(
		nightForecastData,
		goodCloudCoverThreshold,
		consecutiveGoodHoursRequired,
	)
	extremeCloudCover := calculateExtremeCloudCover(nightForecastData)
	displayCloudCover := extremeCloudCover
	if bestObservationInfo.TimeRange != nil {
		displayCloudCover = bestObservationInfo.LowestCloudCover
	}
	nightlyTemperature := int(math.Floor(calculateNightlyAverage(nightForecastData, "temperature")))
	nightlyHumidity := int(math.Floor(calculateNightlyAverage(nightForecastData, "humidity")))
	nightlyWindSpeed := int(math.Floor(calculateNightlyAverage(nightForecastData, "windSpeed")))
	nightlyDewPoint := int(math.Floor(calculateNightlyAverage(nightForecastData, "dewPoint")))
	maxPrecipProbability := maxPrecipitationProbability(nightForecastData)
	seeingIndex := generateSeeingIndexForNight(nightForecastData)
	nightlyWindDirection := calculateWindDirectionAverage(nightForecastData)
	windDirectionText := convertWindDirectionToNSEW(nightlyWindDirection)

	return NightForecast{
		BestObservation:      bestObservationInfo,
		ExtremeCloudCover:    extremeCloudCover,
		DisplayCloudCover:    displayCloudCover,
		NightlyTemperature:   nightlyTemperature,
		NightlyHumidity:      nightlyHumidity,
		NightlyWindSpeed:     nightlyWindSpeed,
		NightlyDewPoint:      nightlyDewPoint,
		MaxPrecipProbability: maxPrecipProbability,
		NightlyWindDirection: nightlyWindDirection,
		WindDirectionText:    windDirectionText,
		SeeingIndex:          seeingIndex,
	}
}

// filterNightForecastData filters forecast data for the astronomical night
func filterNightForecastData(forecastData []ForecastHour, sunsetTime, sunriseTime time.Time) []ForecastHour {
	var nightForecast []ForecastHour

	for _, hour := range forecastData {
		if (hour.DateTime.Equal(sunsetTime) || hour.DateTime.After(sunsetTime)) &&
			(hour.DateTime.Equal(sunriseTime) || hour.DateTime.Before(sunriseTime)) {
			nightForecast = append(nightForecast, hour)
		}
	}

	return nightForecast
}

// getBestObservationTimeRange finds the best observation time range within tonight's astronomical night
func getBestObservationTimeRange(data []ForecastHour, goodCloudCoverThreshold int, consecutiveGoodHoursRequired int) BestObservationInfo {
	var observationStartTime *int
	var observationEndTime *int
	consecutiveGoodHours := 0
	lowestCloudCover := math.MaxInt32
	var extendedEndTime *int

	for i := range data {
		hourlyData := data[i]
		hour := hourlyData.Hour
		clouds := hourlyData.Clouds

		if clouds <= goodCloudCoverThreshold {
			if consecutiveGoodHours == 0 {
				observationStartTime = &hour
			}
			consecutiveGoodHours++
			if clouds < lowestCloudCover {
				lowestCloudCover = clouds
			}
			extendedEndTime = &hour

			if consecutiveGoodHours >= consecutiveGoodHoursRequired {
				observationEndTime = extendedEndTime
			}
		} else {
			if observationEndTime != nil {
				break // Stop extending if minimum is met and conditions change
			}
			consecutiveGoodHours = 0
			observationStartTime = nil
			observationEndTime = nil
			lowestCloudCover = math.MaxInt32
		}
	}

	var timeRange *TimeRange
	if observationStartTime != nil && observationEndTime != nil {
		timeRange = &TimeRange{
			Start: *observationStartTime,
			End:   *observationEndTime,
		}
	}

	if lowestCloudCover == math.MaxInt32 {
		lowestCloudCover = 0
	}

	return BestObservationInfo{
		TimeRange:        timeRange,
		LowestCloudCover: lowestCloudCover,
	}
}

// calculateNightlyAverage calculates average for a specific meteorological parameter
func calculateNightlyAverage(data []ForecastHour, parameter string) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64
	count := 0

	for _, hour := range data {
		var value float64
		switch parameter {
		case "temperature":
			value = hour.Temperature
		case "humidity":
			value = float64(hour.Humidity)
		case "windSpeed":
			value = hour.WindSpeed
		case "dewPoint":
			value = hour.DewPoint
		default:
			continue
		}
		sum += value
		count++
	}

	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// maxPrecipitationProbability finds the maximum precipitation probability in the forecast period
func maxPrecipitationProbability(data []ForecastHour) int {
	if len(data) == 0 {
		return 0
	}

	max := 0
	for _, hour := range data {
		if hour.PrecipitationProbability > max {
			max = hour.PrecipitationProbability
		}
	}
	return max
}

// calculateExtremeCloudCover finds the maximum cloud cover in the forecast period
func calculateExtremeCloudCover(data []ForecastHour) int {
	if len(data) == 0 {
		return 0
	}

	max := 0
	for _, hour := range data {
		if hour.Clouds > max {
			max = hour.Clouds
		}
	}
	return max
}

// calculateSeeingIndex calculates seeing conditions for astronomical observation
func calculateSeeingIndex(temperature, dewPoint, windSpeed float64, humidity int) int {
	tempWeight := 0.25
	windWeight := 0.4
	humidityWeight := 0.15
	dewPointWeight := 0.2

	tempDiff := math.Abs(temperature - dewPoint)
	tempFactor := math.Max(0.1, math.Min(1, (15-tempDiff)/15))
	windFactor := math.Max(0.1, math.Min(1, 1-windSpeed/25))
	humidityFactor := math.Max(0.1, math.Min(1, 1-float64(humidity)/100))
	dewPointFactor := math.Max(0.1, math.Min(1, (10-tempDiff)/10))

	weightedIndex := tempWeight*tempFactor +
		windWeight*windFactor +
		humidityWeight*humidityFactor +
		dewPointWeight*dewPointFactor

	return int(math.Round(math.Max(1, weightedIndex*5)))
}

func calculateSkyQualityIndex(clouds, humidity int, windSpeed, temp, dewPoint float64, seeing int) int {
	tempDiff := math.Abs(temp - 15)
	dewPointDiff := math.Abs(temp - dewPoint)

	cloudsFactor := 1 - clouds/15
	humidityFactor := 5 - humidity/20
	windFactor := 5 - windSpeed/10
	tempFactor := 5 - tempDiff/10
	dewPointFactor := 5 - dewPointDiff/5
	seeingFactor := 5 - seeing

	skyQualityIndex := 0.5*float64(cloudsFactor) +
		0.2*float64(humidityFactor) +
		0.2*windFactor +
		0.1*tempFactor +
		0.15*dewPointFactor +
		0.5*float64(seeingFactor)

	return int(math.Max(0, math.Min(5, skyQualityIndex)))
}

// generateSeeingIndexForNight calculates the average seeing index for a night
func generateSeeingIndexForNight(nightForecastData []ForecastHour) int {
	if len(nightForecastData) == 0 {
		return 0
	}

	var totalIndex float64
	count := 0

	for _, hour := range nightForecastData {
		seeingIndex := calculateSeeingIndex(
			hour.Temperature,
			hour.DewPoint,
			hour.WindSpeed,
			hour.Humidity,
		)
		totalIndex += float64(seeingIndex)
		count++
	}

	if count == 0 {
		return 0
	}
	return int(math.Round(totalIndex / float64(count)))
}

// calculateWindDirectionAverage calculates the average wind direction using vector averaging
func calculateWindDirectionAverage(data []ForecastHour) int {
	var x, y float64
	count := 0

	for _, hour := range data {
		if hour.WindDirection >= 0 {
			radians := hour.WindDirection * math.Pi / 180
			x += math.Cos(radians)
			y += math.Sin(radians)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	averageDirection := math.Atan2(y, x) * (180 / math.Pi)
	return int(math.Round(math.Mod(averageDirection+360, 360)))
}

// convertWindDirectionToNSEW converts wind direction degrees to cardinal directions
func convertWindDirectionToNSEW(degrees int) string {
	if degrees < 0 {
		return "N/A"
	}

	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int(math.Floor(float64(degrees)+22.5/45)) % 8
	return directions[index]
}
