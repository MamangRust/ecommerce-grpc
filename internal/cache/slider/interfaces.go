package slider_cache

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
)

type SliderQueryCache interface {
	GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, *int, bool)
	SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersRow, total *int)

	GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, *int, bool)
	SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersActiveRow, total *int)

	GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, *int, bool)
	SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersTrashedRow, total *int)
}

type SliderCommandCache interface {
	DeleteSliderCache(ctx context.Context, slider_id int)
}
