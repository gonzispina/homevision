package repository

import (
	"fmt"
	"github.com/gonzispina/gokit/context"
	"github.com/gonzispina/gokit/logs"
	"homevision/internal/houses"
	"io"
	"os"
	"path"
	"path/filepath"
)

type entry struct {
	House    *houses.House
	FilePath string
}

type HouseStorageConfig struct {
	FilesPath string
}

func DefaultHouseStorageConfig() *HouseStorageConfig {
	wd, _ := os.Getwd()
	return &HouseStorageConfig{FilesPath: path.Join(wd, "tmp")}
}

func NewHousesStorage(config *HouseStorageConfig, logger logs.Logger) *HousesStorage {
	if config == nil {
		panic("config must be initialized")
	}
	if logger == nil {
		panic("logger must be initialized")
	}
	return &HousesStorage{
		config: config,
		logger: logger,
	}
}

// HousesStorage In memory storage
type HousesStorage struct {
	config *HouseStorageConfig
	logger logs.Logger
}

func (s *HousesStorage) Setup(ctx context.Context) error {
	_, err := os.Stat(s.config.FilesPath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		s.logger.Error(ctx, "Unexpected err", logs.Error(err))
		return err
	}

	if err := os.Mkdir(s.config.FilesPath, os.ModePerm); err != nil {
		s.logger.Error(ctx, "Couldn't create folder", logs.Error(err))
		return err
	}

	return nil
}

func (s *HousesStorage) SaveFile(ctx context.Context, h *houses.House, content io.ReadCloser) error {
	defer content.Close()

	ext := filepath.Ext(h.PhotoURL)
	filepath := path.Join(s.config.FilesPath, fmt.Sprintf("%v-%s.%s", h.ID, h.Address, ext))
	file, err := os.Create(filepath)
	if err != nil {
		s.logger.Error(ctx, "Couldn't open file", logs.Error(err))
		return err
	}

	defer file.Close()
	_, err = io.Copy(file, content)
	if err != nil {
		s.logger.Info(ctx, "Couldn't write file", logs.Error(err))
		return err
	}

	return nil
}
