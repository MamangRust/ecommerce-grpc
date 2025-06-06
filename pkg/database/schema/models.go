// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type Banner struct {
	BannerID  int32        `json:"banner_id"`
	Name      string       `json:"name"`
	StartDate time.Time    `json:"start_date"`
	EndDate   time.Time    `json:"end_date"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	IsActive  sql.NullBool `json:"is_active"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Cart struct {
	CartID    int32        `json:"cart_id"`
	UserID    int32        `json:"user_id"`
	ProductID int32        `json:"product_id"`
	Name      string       `json:"name"`
	Price     int32        `json:"price"`
	Image     string       `json:"image"`
	Quantity  int32        `json:"quantity"`
	Weight    int32        `json:"weight"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Category struct {
	CategoryID    int32          `json:"category_id"`
	Name          string         `json:"name"`
	Description   sql.NullString `json:"description"`
	SlugCategory  sql.NullString `json:"slug_category"`
	ImageCategory sql.NullString `json:"image_category"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at"`
	DeletedAt     sql.NullTime   `json:"deleted_at"`
}

type Merchant struct {
	MerchantID   int32          `json:"merchant_id"`
	UserID       int32          `json:"user_id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	Address      sql.NullString `json:"address"`
	ContactEmail sql.NullString `json:"contact_email"`
	ContactPhone sql.NullString `json:"contact_phone"`
	Status       string         `json:"status"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}

type MerchantBusinessInformation struct {
	MerchantBusinessInfoID int32          `json:"merchant_business_info_id"`
	MerchantID             int32          `json:"merchant_id"`
	BusinessType           sql.NullString `json:"business_type"`
	TaxID                  sql.NullString `json:"tax_id"`
	EstablishedYear        sql.NullInt32  `json:"established_year"`
	NumberOfEmployees      sql.NullInt32  `json:"number_of_employees"`
	WebsiteUrl             sql.NullString `json:"website_url"`
	CreatedAt              sql.NullTime   `json:"created_at"`
	UpdatedAt              sql.NullTime   `json:"updated_at"`
	DeletedAt              sql.NullTime   `json:"deleted_at"`
}

type MerchantCertificationsAndAward struct {
	MerchantCertificationID int32          `json:"merchant_certification_id"`
	MerchantID              int32          `json:"merchant_id"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	IssuedBy                sql.NullString `json:"issued_by"`
	IssueDate               sql.NullTime   `json:"issue_date"`
	ExpiryDate              sql.NullTime   `json:"expiry_date"`
	CertificateUrl          sql.NullString `json:"certificate_url"`
	CreatedAt               sql.NullTime   `json:"created_at"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	DeletedAt               sql.NullTime   `json:"deleted_at"`
}

type MerchantDetail struct {
	MerchantDetailID int32          `json:"merchant_detail_id"`
	MerchantID       int32          `json:"merchant_id"`
	DisplayName      sql.NullString `json:"display_name"`
	CoverImageUrl    sql.NullString `json:"cover_image_url"`
	LogoUrl          sql.NullString `json:"logo_url"`
	ShortDescription sql.NullString `json:"short_description"`
	WebsiteUrl       sql.NullString `json:"website_url"`
	CreatedAt        sql.NullTime   `json:"created_at"`
	UpdatedAt        sql.NullTime   `json:"updated_at"`
	DeletedAt        sql.NullTime   `json:"deleted_at"`
}

type MerchantPolicy struct {
	MerchantPolicyID int32        `json:"merchant_policy_id"`
	MerchantID       int32        `json:"merchant_id"`
	PolicyType       string       `json:"policy_type"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	CreatedAt        sql.NullTime `json:"created_at"`
	UpdatedAt        sql.NullTime `json:"updated_at"`
	DeletedAt        sql.NullTime `json:"deleted_at"`
}

type MerchantSocialMediaLink struct {
	MerchantSocialID int32        `json:"merchant_social_id"`
	MerchantDetailID int32        `json:"merchant_detail_id"`
	Platform         string       `json:"platform"`
	Url              string       `json:"url"`
	CreatedAt        sql.NullTime `json:"created_at"`
	UpdatedAt        sql.NullTime `json:"updated_at"`
}

type Order struct {
	OrderID    int32        `json:"order_id"`
	UserID     int32        `json:"user_id"`
	MerchantID int32        `json:"merchant_id"`
	TotalPrice int32        `json:"total_price"`
	CreatedAt  sql.NullTime `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

type OrderItem struct {
	OrderItemID int32        `json:"order_item_id"`
	OrderID     int32        `json:"order_id"`
	ProductID   int32        `json:"product_id"`
	Quantity    int32        `json:"quantity"`
	Price       int32        `json:"price"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type Product struct {
	ProductID    int32           `json:"product_id"`
	MerchantID   int32           `json:"merchant_id"`
	CategoryID   int32           `json:"category_id"`
	Name         string          `json:"name"`
	Description  sql.NullString  `json:"description"`
	Price        int32           `json:"price"`
	CountInStock int32           `json:"count_in_stock"`
	Brand        sql.NullString  `json:"brand"`
	Weight       sql.NullInt32   `json:"weight"`
	Rating       sql.NullFloat64 `json:"rating"`
	SlugProduct  sql.NullString  `json:"slug_product"`
	ImageProduct sql.NullString  `json:"image_product"`
	CreatedAt    sql.NullTime    `json:"created_at"`
	UpdatedAt    sql.NullTime    `json:"updated_at"`
	DeletedAt    sql.NullTime    `json:"deleted_at"`
}

type RefreshToken struct {
	RefreshTokenID int32        `json:"refresh_token_id"`
	UserID         int32        `json:"user_id"`
	Token          string       `json:"token"`
	Expiration     time.Time    `json:"expiration"`
	CreatedAt      sql.NullTime `json:"created_at"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
}

type Review struct {
	ReviewID  int32        `json:"review_id"`
	UserID    int32        `json:"user_id"`
	ProductID int32        `json:"product_id"`
	Name      string       `json:"name"`
	Comment   string       `json:"comment"`
	Rating    int32        `json:"rating"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type ReviewDetail struct {
	ReviewDetailID int32          `json:"review_detail_id"`
	ReviewID       int32          `json:"review_id"`
	Type           string         `json:"type"`
	Url            string         `json:"url"`
	Caption        sql.NullString `json:"caption"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
}

type Role struct {
	RoleID    int32        `json:"role_id"`
	RoleName  string       `json:"role_name"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type ShippingAddress struct {
	ShippingAddressID int32        `json:"shipping_address_id"`
	OrderID           int32        `json:"order_id"`
	Alamat            string       `json:"alamat"`
	Provinsi          string       `json:"provinsi"`
	Negara            string       `json:"negara"`
	Kota              string       `json:"kota"`
	Courier           string       `json:"courier"`
	ShippingMethod    string       `json:"shipping_method"`
	ShippingCost      float64      `json:"shipping_cost"`
	CreatedAt         sql.NullTime `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
}

type Slider struct {
	SliderID  int32        `json:"slider_id"`
	Name      string       `json:"name"`
	Image     string       `json:"image"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Transaction struct {
	TransactionID int32        `json:"transaction_id"`
	OrderID       int32        `json:"order_id"`
	MerchantID    int32        `json:"merchant_id"`
	PaymentMethod string       `json:"payment_method"`
	Amount        int32        `json:"amount"`
	PaymentStatus string       `json:"payment_status"`
	CreatedAt     sql.NullTime `json:"created_at"`
	UpdatedAt     sql.NullTime `json:"updated_at"`
	DeletedAt     sql.NullTime `json:"deleted_at"`
}

type User struct {
	UserID    int32        `json:"user_id"`
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type UserRole struct {
	UserRoleID int32        `json:"user_role_id"`
	UserID     int32        `json:"user_id"`
	RoleID     int32        `json:"role_id"`
	CreatedAt  sql.NullTime `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}
