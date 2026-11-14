package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const regionBaseURL = "https://www.emsifa.com/api-wilayah-indonesia/v2"

var regionHTTPClient = &http.Client{Timeout: 10 * time.Second}

type RegionService struct{}

func (s *RegionService) fetchFromEmsifa(ctx context.Context, endpoint string) (any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, regionBaseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := regionHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return data, nil
}

func (s *RegionService) GetProvinces(ctx context.Context) (any, error) {
	return s.fetchFromEmsifa(ctx, "/provinces.json")
}

func (s *RegionService) GetCitiesByProvince(ctx context.Context, provinceID string) (any, error) {
	if provinceID == "" {
		return nil, fmt.Errorf("province ID cannot be empty")
	}
	return s.fetchFromEmsifa(ctx, fmt.Sprintf("/regencies/%s.json", provinceID))
}

func (s *RegionService) GetProvinceDetail(ctx context.Context, provinceID string) (any, error) {
	if provinceID == "" {
		return nil, fmt.Errorf("province ID cannot be empty")
	}
	return s.fetchFromEmsifa(ctx, fmt.Sprintf("/provinces/%s.json", provinceID))
}

func (s *RegionService) GetCityDetail(ctx context.Context, cityID string) (any, error) {
	if cityID == "" {
		return nil, fmt.Errorf("city ID cannot be empty")
	}
	return s.fetchFromEmsifa(ctx, fmt.Sprintf("/regencies/%s.json", cityID))
}
