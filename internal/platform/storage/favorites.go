package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"driffaud.fr/odin/internal/domain"
)

// FavoritesStore manages the storage and retrieval of favorite places
type FavoritesStore struct {
	Favorites []domain.Place
	FilePath  string
}

// NewFavoritesStore creates a new store for favorite places
func NewFavoritesStore() (*FavoritesStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appConfigDir := configDir + "/odin"
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(appConfigDir, "favorites.json")
	store := &FavoritesStore{
		FilePath: filePath,
	}

	if err := store.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return store, nil
}

// Load reads favorites from the JSON file
func (s *FavoritesStore) Load() error {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.Favorites)
}

// Save stores favorites to the JSON file
func (fs *FavoritesStore) Save() error {
	data, err := json.MarshalIndent(fs.Favorites, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.FilePath, data, 0644)
}

// AddFavorite adds a place to favorites
func (fs *FavoritesStore) AddFavorite(place domain.Place) error {
	// Check if the place is already in favorites
	for _, fav := range fs.Favorites {
		if fav.Name == place.Name && fav.Latitude == place.Latitude && fav.Longitude == place.Longitude {
			return nil // Already a favorite
		}
	}

	fs.Favorites = append(fs.Favorites, place)
	return fs.Save()
}

// RemoveFavorite removes a place from favorites
func (fs *FavoritesStore) RemoveFavorite(place domain.Place) error {
	newFavorites := []domain.Place{}
	for _, fav := range fs.Favorites {
		if !(fav.Name == place.Name && fav.Latitude == place.Latitude && fav.Longitude == place.Longitude) {
			newFavorites = append(newFavorites, fav)
		}
	}

	fs.Favorites = newFavorites
	return fs.Save()
}

// IsFavorite checks if a place is in the favorites list
func (fs *FavoritesStore) IsFavorite(place domain.Place) bool {
	for _, fav := range fs.Favorites {
		if fav.Name == place.Name && fav.Latitude == place.Latitude && fav.Longitude == place.Longitude {
			return true
		}
	}
	return false
}
