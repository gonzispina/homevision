package houses

import (
	"github.com/gonzispina/gokit/context"
	"io"
	"sync"
)

type HouseGetter interface {
	GetHousesPaged(ctx context.Context, page, offset int) ([]*House, error)
	GetHousePhoto(ctx context.Context, photoURL string) (io.ReadCloser, error)
}

type HouseSaver interface {
	SaveHouse(ctx context.Context, h *House)
	SaveFile(ctx context.Context, h *House, content io.ReadCloser) error
}

func NewHouseManager(g HouseGetter, s HouseSaver) HouseManager {
	if g == nil {
		panic("getter must be initialized")
	}
	if s == nil {
		panic("saver must be initialized")
	}
	return &houseManager{
		getter: g,
		saver:  s,
	}
}

// HouseManager ...
type HouseManager interface {
	GetAllHouses(ctx context.Context) error
}

type houseManager struct {
	getter HouseGetter
	saver  HouseSaver
}

func (m *houseManager) GetAllHouses(ctx context.Context) error {
	offset := 20
	pages := 5

	wg := sync.WaitGroup{}

	for page := 1; page <= pages; page++ {
		houses, err := m.getter.GetHousesPaged(ctx, page, offset)
		if err != nil {
			return err
		}

		for i, h := range houses {
			wg.Add(1)

			go func(i int, h *House) {
				defer wg.Done()
				// Save to retry later
				m.saver.SaveHouse(ctx, h)

				photo, err := m.getter.GetHousePhoto(ctx, h.PhotoURL)
				if err != nil {
					return
				}

				err = m.saver.SaveFile(ctx, h, photo)
				if err != nil {
					return
				}
			}(i, h)
		}
	}

	wg.Wait()
	return nil
}
