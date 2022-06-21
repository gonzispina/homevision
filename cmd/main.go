package main

import (
	"github.com/gonzispina/gokit/context"
	"github.com/gonzispina/gokit/logs"
	"homevision/internal/houses"
	"homevision/internal/houses/repository"
)

func main() {
	ctx := context.Background()
	logger := logs.InitDefault()

	client := repository.NewAPIClient(repository.DefaultClientConfig(), logger)

	storage := repository.NewHousesStorage(repository.DefaultHouseStorageConfig(), logger)
	err := storage.Setup(ctx)
	if err != nil {
		panic(err)
	}

	manager := houses.NewHouseManager(client, storage)

	for {
		err := manager.GetAllHouses(ctx)
		if err == nil {
			break
		}
	}

}
