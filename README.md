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
- **🌐 Internationalization**: Supports English and French languages

## 🛠️ Installation

### Prerequisites

- Go 1.19 or later

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/odin.git
cd odin

# Build the application
go build -o odin cmd/odin/main.go

# Run the application
./odin
```

## 📖 Usage

### 🎮 Navigation

- **Tab**: Switch between search field and favorites list
- **Enter**: Confirm selection
- **Esc**: Go back
- **CTRL+C**: Exit application
- **F2**: Add current location to favorites
- **F3**: Remove location from favorites

### 🚀 Workflow

1. Search for a location or select from your favorites
2. View weather data optimized for astronomical observation
3. Check the best time period for tonight's viewing conditions
4. See detailed weather parameters that affect observation quality

## ⚙️ Configuration

Odin stores favorites in the user configuration directory:
- **Linux**: `~/.config/odin/favorites.json`
- **macOS**: `~/Library/Application Support/odin/favorites.json`
- **Windows**: `%APPDATA%\odin\favorites.json`

## 🔧 Technical Details

Odin is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles): UI components for Bubble Tea
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Style definitions for terminal UI
- [SunCalc](https://github.com/sixdouglas/suncalc): Astronomical calculations
- [go-i18n](https://github.com/nicksnyder/go-i18n): Internationalization support
- [go-locale](https://github.com/Xuanwo/go-locale): Locale detection

### Supported Languages

- English
- French

### Data Sources

- Weather data provided by [Open-Meteo API](https://open-meteo.com/)
- Geocoding provided by [Photon API](https://photon.komoot.io/)

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- Weather data provided by [Open-Meteo](https://open-meteo.com/)
- Geocoding powered by [Photon](https://photon.komoot.io/)
- Astronomical calculations using [SunCalc](https://github.com/sixdouglas/suncalc)
