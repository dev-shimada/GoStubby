package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dev-shimada/gostubby/internal/domain/model"
)

type ConfigRepository struct{}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{}
}

func (c ConfigRepository) Load(path string) ([]model.Endpoint, error) {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access config path: %v", err)
	}

	// If path is a file, load it directly
	if !info.IsDir() {
		if filepath.Ext(path) != ".json" {
			return nil, fmt.Errorf("config file must be a JSON file")
		}
		return c.loadFile(path)
	}

	// If path is a directory, walk through it
	var allEndpoints []model.Endpoint
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		endpoints, err := c.loadFile(path)
		if err != nil {
			return fmt.Errorf("error loading file %s: %v", path, err)
		}
		allEndpoints = append(allEndpoints, endpoints...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return allEndpoints, nil
}

func (c ConfigRepository) loadFile(path string) ([]model.Endpoint, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error(fmt.Sprintf("Failed to close file: %s", err))
		}
	}()

	byteValue, _ := io.ReadAll(file)
	var endpoints []model.Endpoint
	err = json.Unmarshal(byteValue, &endpoints)
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}
