package record

type ShippingAddressRecord struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	Alamat    string  `json:"alamat"`
	Provinsi  string  `json:"provinsi"`
	Negara    string  `json:"negara"`
	Kota      string  `json:"kota"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
