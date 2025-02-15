package record

type CategoriesRecord struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	SlugCategory  string  `json:"slug_category"`
	ImageCategory string  `json:"image_category"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}
