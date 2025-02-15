package recordmapper

import (
	"ecommerce/internal/domain/record"
	db "ecommerce/pkg/database/schema"
)

type cartRecordMapper struct {
}

func NewCartRecordMapper() *cartRecordMapper {
	return &cartRecordMapper{}
}

func (s *cartRecordMapper) ToCartRecord(cart *db.Cart) *record.CartRecord {
	var deletedAt *string

	if cart.DeletedAt.Valid {
		formattedDeletedAt := cart.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.CartRecord{
		ID:        int(cart.CartID),
		UserID:    int(cart.UserID),
		ProductID: int(cart.ProductID),
		Name:      cart.Name,
		Price:     int(cart.Price),
		Image:     cart.Image,
		Quantity:  int(cart.Quantity),
		Weight:    int(cart.Weight),
		CreatedAt: cart.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: cart.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}

func (s *cartRecordMapper) ToCartRecordPagination(cart *db.GetCartsRow) *record.CartRecord {
	var deletedAt *string

	if cart.DeletedAt.Valid {
		formattedDeletedAt := cart.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.CartRecord{
		ID:        int(cart.CartID),
		UserID:    int(cart.UserID),
		ProductID: int(cart.ProductID),
		Name:      cart.Name,
		Price:     int(cart.Price),
		Image:     cart.Image,
		Quantity:  int(cart.Quantity),
		Weight:    int(cart.Weight),
		CreatedAt: cart.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: cart.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}

func (s *cartRecordMapper) ToCartsRecordPagination(carts []*db.GetCartsRow) []*record.CartRecord {
	var cartRecords []*record.CartRecord

	for _, cart := range carts {
		cartRecords = append(cartRecords, s.ToCartRecordPagination(cart))
	}

	return cartRecords
}
