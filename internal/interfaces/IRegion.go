package interfaces

import "context"

type IRegionService interface {
	GetProvinces(ctx context.Context) (any, error)
	GetCitiesByProvince(ctx context.Context, provinceID string) (any, error)
	GetProvinceDetail(ctx context.Context, provinceID string) (any, error)
	GetCityDetail(ctx context.Context, cityID string) (any, error)
}
