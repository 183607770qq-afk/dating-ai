package models

import (
	"time"
)

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PlatformID    string    `gorm:"uniqueIndex:idx_platform_product;not null" json:"platform_id"`
	Platform      string    `gorm:"index;not null" json:"platform"`
	Title         string    `gorm:"type:text;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	ImageURL      string    `gorm:"type:text" json:"image_url"`
	ProductURL    string    `gorm:"type:text;not null" json:"product_url"`
	Price         float64   `gorm:"type:decimal(10,2)" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2)" json:"original_price"`
	Currency      string    `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	SalesCount    int       `gorm:"default:0" json:"sales_count"`
	ReviewCount   int       `gorm:"default:0" json:"review_count"`
	Rating        float64   `gorm:"type:decimal(3,2)" json:"rating"`
	Category      string    `gorm:"index" json:"category"`
	Tags          string    `gorm:"type:text" json:"tags"`
	Badge         string    `json:"badge"`
	ShopName      string    `json:"shop_name"`
	ShopID        string    `json:"shop_id"`
	IsHot         bool      `gorm:"default:false;index" json:"is_hot"`
	TrendingScore float64   `gorm:"type:decimal(10,2);default:0" json:"trending_score"`
	CrawledAt     time.Time `json:"crawled_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	ParentID    *uint     `json:"parent_id"`
	Platform    string    `gorm:"index" json:"platform"`
	Icon        string    `json:"icon"`
	Description string    `gorm:"type:text" json:"description"`
	ProductCount int      `gorm:"default:0" json:"product_count"`
	HotCount     int      `gorm:"default:0" json:"hot_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

type PriceHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"index;not null" json:"product_id"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64 `gorm:"type:decimal(10,2)" json:"original_price"`
	SalesCount  int       `json:"sales_count"`
	CrawledAt   time.Time `json:"crawled_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (PriceHistory) TableName() string {
	return "price_history"
}

type Platform struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	BaseURL     string    `gorm:"type:text" json:"base_url"`
	Logo        string    `json:"logo"`
	Country     string    `json:"country"`
	Language    string    `json:"language"`
	Currency    string    `json:"currency"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	HotURL      string    `gorm:"type:text" json:"hot_url"`
	CategoryURL string    `gorm:"type:text" json:"category_url"`
	ProductURL  string    `gorm:"type:text" json:"product_url"`
	Priority    int       `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Platform) TableName() string {
	return "platforms"
}

type CrawlTask struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PlatformID  uint      `gorm:"index" json:"platform_id"`
	CategoryID  *uint     `json:"category_id"`
	Keyword     string    `json:"keyword"`
	URL         string    `gorm:"type:text;not null" json:"url"`
	Type        string    `gorm:"index" json:"type"`
	Status      string    `gorm:"index;default:'pending'" json:"status"`
	Priority    int       `gorm:"default:0" json:"priority"`
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`
	RetryCount  int       `gorm:"default:0" json:"retry_count"`
	ProductCount int      `gorm:"default:0" json:"product_count"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CrawlTask) TableName() string {
	return "crawl_tasks"
}

type Analytics struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProductID    uint      `gorm:"index;not null" json:"product_id"`
	Date         time.Time `gorm:"index;not null" json:"date"`
	Views        int       `gorm:"default:0" json:"views"`
	UniqueViews  int       `gorm:"default:0" json:"unique_views"`
	CartAdds     int       `gorm:"default:0" json:"cart_adds"`
	Purchases    int       `gorm:"default:0" json:"purchases"`
	Revenue      float64   `gorm:"type:decimal(12,2);default:0" json:"revenue"`
	ConversionRate float64 `gorm:"type:decimal(5,3)" json:"conversion_rate"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Analytics) TableName() string {
	return "analytics"
}
