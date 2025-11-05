package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CatBreedExists(breed string) (bool, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return false, fmt.Errorf("failed to fetch breeds: %w", err)
	}
	defer resp.Body.Close()

	var breeds []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false, fmt.Errorf("invalid response from Cat API: %w", err)
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true, nil
		}
	}

	return false, nil
}
