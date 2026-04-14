package db

import (
	"context"
)

func (q *Queries) GetShippingAddressByOrderIDTrashed(ctx context.Context, orderID int32) (*ShippingAddress, error) {
	row := q.db.QueryRow(ctx, "SELECT * FROM shipping_addresses WHERE order_id = $1 AND deleted_at IS NOT NULL", orderID)
	var i ShippingAddress
	err := row.Scan(
		&i.ShippingAddressID,
		&i.OrderID,
		&i.Alamat,
		&i.Provinsi,
		&i.Negara,
		&i.Kota,
		&i.Courier,
		&i.ShippingMethod,
		&i.ShippingCost,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}
