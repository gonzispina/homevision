package repository

import (
	"encoding/json"
	"fmt"
	"github.com/gonzispina/gokit/context"
	"github.com/gonzispina/gokit/logs"
	"homevision/internal/houses"
	"io"
	"net/http"
)

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		URL:        "http://app-homevision-staging.herokuapp.com",
		HousesPath: "/api_project/houses",
	}
}

type ClientConfig struct {
	URL        string
	HousesPath string
}

// NewAPIClient constructor
func NewAPIClient(c *ClientConfig, logger logs.Logger) *APIClient {
	if c == nil {
		panic("config must be initialized")
	}
	if logger == nil {
		panic("logger must be initialized")
	}
	return &APIClient{
		config: c,
		client: &http.Client{},
		logger: logger,
	}
}

type APIClient struct {
	config *ClientConfig
	client *http.Client
	logger logs.Logger
}

// GetHousesPaged returns a page of houses for the given page number and offset
func (c *APIClient) GetHousesPaged(ctx context.Context, page, offset int) ([]*houses.House, error) {
	url := fmt.Sprintf("%s%s?page=%v&per_page=%v",
		c.config.URL,
		c.config.HousesPath,
		page,
		offset,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error(ctx, "Couldn't create request", logs.Error(err))
		return nil, err
	}

	var res *http.Response
	for i := 0; i < 3; i++ {
		req = req.WithContext(ctx)
		res, err = c.client.Do(req)
		if err != nil {
			c.logger.Error(ctx, "Couldn't perform request", logs.Error(err))
			return nil, err
		}

		if res.StatusCode == 200 {
			break
		}

		_ = res.Body.Close()
		c.logger.Info(ctx, "Received Invalid response")
	}

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var r struct {
		Houses []*houses.House `json:"houses"`
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error(ctx, "Couldn't read body", logs.Error(err))
		return nil, err
	}

	if err = json.Unmarshal(body, &r); err != nil {
		c.logger.Error(ctx, "Couldn't unmarshal response", logs.Error(err))
		return nil, err
	}

	return r.Houses, nil
}

// GetHousePhoto by its URL
func (c *APIClient) GetHousePhoto(ctx context.Context, photoURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, photoURL, nil)
	if err != nil {
		c.logger.Error(ctx, "Couldn't create request", logs.Error(err))
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(ctx, "Couldn't perform request", logs.Error(err))
		return nil, err
	}

	if res.StatusCode != 200 {
		c.logger.Info(ctx, "WTF2")
		return nil, nil
	}

	return res.Body, nil
}
