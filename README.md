# 🔭 Odin - Weather for Astronomy

Odin is a command-line application designed for amateur astronomers to check weather conditions for stargazing. It provides detailed weather forecasts optimized for astronomical observation, including cloud cover, seeing conditions, and ideal viewing periods.

## ✨ Features

- **🔍 Location Search**: Find any location worldwide
- **☁️ Astronomical Weather Data**: Get specialized weather data relevant for astronomy
- **🌌 Night Viewing Forecast**: Calculates the best time periods for observation during the night
- **🌓 Sun and Moon Information**: Shows rise/set times and moon phase data
- **⭐ Favorites Management**: Save and quickly access your favorite observation locations
- **📊 Seeing Index**: Numerical rating of overall viewing conditions
- **🕒 Detailed Hourly Forecast**: View temperature, humidity, cloud cover, and more

## 🛠️ Installation

### Prerequisites

- Go 1.19 or later

### Building from Source

1. Clone the repository
2. Build the application:

```bash
go build -o odin main.go
```

3. Run the application:

```bash
./odin
```

## 📖 Usage

### 🎮 Navigation

- **Tab**: Switch between search field and favorites list
- **Enter**: Confirm selection
- **Esc**: Go back or exit
- **F2**: Add current location to favorites
- **F3**: Remove location from favorites

### 🚀 Workflow

1. Search for a location or select from your favorites
2. View weather data optimized for astronomical observation
3. Check the best time period for tonight's viewing conditions
4. See detailed weather parameters that affect observation quality

## 🔧 Technical Details

Odin is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Style definitions for terminal UI
- [Open-Meteo API](https://open-meteo.com/): Weather data
- [Photon API](https://photon.komoot.io/): Geocoding for location search
- [SunCalc](https://github.com/sixdouglas/suncalc): Astronomical calculations

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.
