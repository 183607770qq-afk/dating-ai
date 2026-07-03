package database

import (
	"fmt"
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接管理器
// 使用 GORM ORM 操作 PostgreSQL 数据库
type Database struct {
	DB *gorm.DB // GORM数据库实例
}

// NewDatabase 创建数据库连接实例
// 流程：
// 1. 构建DSN连接字符串
// 2. 初始化GORM连接
// 3. 配置连接池参数
// 4. 自动迁移数据库表结构
// 5. 初始化平台种子数据
//
// 参数：
//   cfg: 数据库配置
// 返回：
//   *Database: 数据库实例
//   error: 连接错误
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	// 构建DSN连接字符串
	dsn := cfg.DSN()

	// 配置GORM日志级别
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 建立数据库连接
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层SQL数据库实例以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(10)    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)   // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 自动迁移数据库表结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 初始化平台种子数据（如果不存在）
	if err := seedPlatforms(db); err != nil {
		log.Printf("Warning: failed to seed platforms: %v", err)
	}

	return &Database{DB: db}, nil
}

// autoMigrate 自动迁移数据库表结构
// 确保所有模型对应的表存在，不会删除已存在的数据
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Platform{},      // 平台表
		&models.Category{},      // 分类表
		&models.Product{},       // 产品表
		&models.PriceHistory{},  // 价格历史表
		&models.CrawlTask{},     // 爬取任务表
		&models.Analytics{},     // 分析数据表
	)
}

// seedPlatforms 初始化平台种子数据
// 如果数据库中没有平台数据，则插入预设的电商平台配置
// 支持的平台包括：Amazon、eBay、淘宝、京东、拼多多、天猫、Shopee等
func seedPlatforms(db *gorm.DB) error {
	platforms := []models.Platform{
		{
			Name:        "Amazon",      // 平台名称
			Code:        "amazon",      // 平台唯一标识
			BaseURL:     "https://www.amazon.com", // 基础URL
			Logo:        "/logos/amazon.png",      // Logo路径
			Country:     "US",          // 国家/地区
			Language:    "en",          // 语言
			Currency:    "USD",         // 货币类型
			IsActive:    true,          // 是否启用
			HotURL:      "https://www.amazon.com/gp/bestsellers",      // 热销页面URL
			CategoryURL: "https://www.amazon.com/gp/site-directory",   // 分类页面URL
			ProductURL:  "https://www.amazon.com/dp/",                 // 产品详情页URL模板
			Priority:    10,           // 爬取优先级（数字越大优先级越高）
		},
		{
			Name:        "eBay",
			Code:        "ebay",
			BaseURL:     "https://www.ebay.com",
			Logo:        "/logos/ebay.png",
			Country:     "US",
			Language:    "en",
			Currency:    "USD",
			IsActive:    true,
			HotURL:      "https://www.ebay.com/b/Best-Selling-Best-Sellers-on-eBay/6000",
			CategoryURL: "https://www.ebay.com/categories",
			ProductURL:  "https://www.ebay.com/itm/",
			Priority:    9,
		},
		{
			Name:        "淘宝",
			Code:        "taobao",
			BaseURL:     "https://www.taobao.com",
			Logo:        "/logos/taobao.png",
			Country:     "CN",
			Language:    "zh",
			Currency:    "CNY",
			IsActive:    true,
			HotURL:      "https://www.taobao.com/tbhome/page/market-list",
			CategoryURL: "https://www.taobao.com/category",
			ProductURL:  "https://item.taobao.com/item.htm?id=",
			Priority:    8,
		},
		{
			Name:        "京东",
			Code:        "jd",
			BaseURL:     "https://www.jd.com",
			Logo:        "/logos/jd.png",
			Country:     "CN",
			Language:    "zh",
			Currency:    "CNY",
			IsActive:    true,
			HotURL:      "https://www.jd.com/renqi.html",
			CategoryURL: "https://www.jd.com/category.html",
			ProductURL:  "https://item.jd.com/",
			Priority:    8,
		},
		{
			Name:        "拼多多",
			Code:        "pdd",
			BaseURL:     "https://www.pinduoduo.com",
			Logo:        "/logos/pdd.png",
			Country:     "CN",
			Language:    "zh",
			Currency:    "CNY",
			IsActive:    true,
			HotURL:      "https://www.pinduoduo.com/home/goods/",
			CategoryURL: "https://www.pinduoduo.com/home/category/",
			ProductURL:  "https://mobile.yangkeduo.com/goods.html?goods_id=",
			Priority:    7,
		},
		{
			Name:        "天猫",
			Code:        "tmall",
			BaseURL:     "https://www.tmall.com",
			Logo:        "/logos/tmall.png",
			Country:     "CN",
			Language:    "zh",
			Currency:    "CNY",
			IsActive:    true,
			HotURL:      "https://www.tmall.com/tbhome/page/market-list",
			CategoryURL: "https://www.tmall.com/category",
			ProductURL:  "https://detail.tmall.com/item.htm?id=",
			Priority:    8,
		},
		{
			Name:        "Shopee",
			Code:        "shopee",
			BaseURL:     "https://www.shopee.com",
			Logo:        "/logos/shopee.png",
			Country:     "SG",
			Language:    "en",
			Currency:    "SGD",
			IsActive:    true,
			HotURL:      "https://shopee.co.id/best-shops",
			CategoryURL: "https://shopee.co.id/category",
			ProductURL:  "https://shopee.co.id/product/",
			Priority:    6,
		},
		{
			Name:        "Lazada",
			Code:        "lazada",
			BaseURL:     "https://www.lazada.com",
			Logo:        "/logos/lazada.png",
			Country:     "SG",
			Language:    "en",
			Currency:    "SGD",
			IsActive:    true,
			HotURL:      "https://www.lazada.com/best-selling/",
			CategoryURL: "https://www.lazada.com/category/",
			ProductURL:  "https://www.lazada.com/products/",
			Priority:    5,
		},
		{
			Name:        "Wish",
			Code:        "wish",
			BaseURL:     "https://www.wish.com",
			Logo:        "/logos/wish.png",
			Country:     "US",
			Language:    "en",
			Currency:    "USD",
			IsActive:    true,
			HotURL:      "https://www.wish.com/c/best-selling",
			CategoryURL: "https://www.wish.com/c/categories",
			ProductURL:  "https://www.wish.com/c/",
			Priority:    4,
		},
		{
			Name:        "AliExpress",
			Code:        "aliexpress",
			BaseURL:     "https://www.aliexpress.com",
			Logo:        "/logos/aliexpress.png",
			Country:     "CN",
			Language:    "en",
			Currency:    "USD",
			IsActive:    true,
			HotURL:      "https://www.aliexpress.com/p/best-selling/",
			CategoryURL: "https://www.aliexpress.com/category.html",
			ProductURL:  "https://www.aliexpress.com/item/",
			Priority:    7,
		},
	}

	for _, platform := range platforms {
		var existing models.Platform
		result := db.Where("code = ?", platform.Code).First(&existing)
		if result.RowsAffected == 0 {
			if err := db.Create(&platform).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) CreateProduct(product *models.Product) error {
	return d.DB.Create(product).Error
}

func (d *Database) GetProductByPlatformID(platformID, platform string) (*models.Product, error) {
	var product models.Product
	err := d.DB.Where("platform_id = ? AND platform = ?", platformID, platform).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (d *Database) UpdateProduct(product *models.Product) error {
	return d.DB.Save(product).Error
}

func (d *Database) GetHotProducts(platform string, limit int) ([]models.Product, error) {
	var products []models.Product
	query := d.DB.Where("is_hot = ?", true).Order("trending_score DESC")

	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	err := query.Limit(limit).Find(&products).Error
	return products, err
}

// GetProductsByCategory 按分类获取产品列表
// 参数：
//   category: 分类名称
//   limit: 返回数量限制
//   offset: 偏移量（用于分页）
// 返回：
//   []models.Product: 产品列表
//   error: 查询错误
func (d *Database) GetProductsByCategory(category string, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := d.DB.Where("category = ?", category).
		Order("sales_count DESC, trending_score DESC"). // 按销量和热度排序
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

// SearchProducts 搜索产品
// 支持标题、描述、标签的模糊搜索（不区分大小写）
// 参数：
//   keyword: 搜索关键词
//   limit: 返回数量限制
//   offset: 偏移量（用于分页）
// 返回：
//   []models.Product: 匹配的产品列表
//   error: 查询错误
func (d *Database) SearchProducts(keyword string, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	searchPattern := "%" + keyword + "%" // 构建模糊查询模式
	err := d.DB.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?", searchPattern, searchPattern, searchPattern).
		Order("trending_score DESC, sales_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

// CreatePriceHistory 创建价格历史记录
// 用于追踪产品价格变化趋势
func (d *Database) CreatePriceHistory(history *models.PriceHistory) error {
	return d.DB.Create(history).Error
}

// GetPriceHistory 获取产品的价格历史
// 参数：
//   productID: 产品ID
//   days: 查询天数（最近days天）
// 返回：
//   []models.PriceHistory: 价格历史列表
//   error: 查询错误
func (d *Database) GetPriceHistory(productID uint, days int) ([]models.PriceHistory, error) {
	var history []models.PriceHistory
	startDate := time.Now().AddDate(0, 0, -days) // 计算起始日期
	err := d.DB.Where("product_id = ? AND crawled_at >= ?", productID, startDate).
		Order("crawled_at ASC"). // 按时间升序排列
		Find(&history).Error
	return history, err
}

// CreateCrawlTask 创建爬取任务
// 任务用于异步处理爬取请求
func (d *Database) CreateCrawlTask(task *models.CrawlTask) error {
	return d.DB.Create(task).Error
}

// UpdateCrawlTask 更新爬取任务
// 通常用于更新任务状态（如pending→processing→completed/failed）
func (d *Database) UpdateCrawlTask(task *models.CrawlTask) error {
	return d.DB.Save(task).Error
}

// GetPendingCrawlTasks 获取待处理的爬取任务
// 按优先级降序、创建时间升序排列
// 参数：
//   limit: 返回数量限制
// 返回：
//   []models.CrawlTask: 待处理任务列表
//   error: 查询错误
func (d *Database) GetPendingCrawlTasks(limit int) ([]models.CrawlTask, error) {
	var tasks []models.CrawlTask
	err := d.DB.Where("status = ?", "pending").
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// GetPlatforms 获取所有启用的平台
// 按优先级降序排列
func (d *Database) GetPlatforms() ([]models.Platform, error) {
	var platforms []models.Platform
	err := d.DB.Where("is_active = ?", true).Order("priority DESC").Find(&platforms).Error
	return platforms, err
}

// GetPlatformByCode 根据平台代码获取平台信息
// 参数：
//   code: 平台代码（如 amazon, ebay, taobao）
// 返回：
//   *models.Platform: 平台信息
//   error: 查询错误（如平台不存在）
func (d *Database) GetPlatformByCode(code string) (*models.Platform, error) {
	var platform models.Platform
	err := d.DB.Where("code = ?", code).First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

// GetCategories 获取分类列表
// 参数：
//   platform: 平台代码（为空则返回所有平台的分类）
// 返回：
//   []models.Category: 分类列表
//   error: 查询错误
func (d *Database) GetCategories(platform string) ([]models.Category, error) {
	var categories []models.Category
	query := d.DB.Where("parent_id IS NULL") // 只获取一级分类

	if platform != "" {
		query = query.Where("platform = ? OR platform = ?", platform, "all")
	}

	err := query.Order("product_count DESC").Find(&categories).Error
	return categories, err
}

func (d *Database) CreateCategory(category *models.Category) error {
	return d.DB.Create(category).Error
}

func (d *Database) IncrementCategoryProductCount(categoryName string) error {
	return d.DB.Model(&models.Category{}).
		Where("name = ?", categoryName).
		UpdateColumn("product_count", gorm.Expr("product_count + ?", 1)).
		Error
}

func (d *Database) GetAnalytics(productID uint, days int) ([]models.Analytics, error) {
	var analytics []models.Analytics
	startDate := time.Now().AddDate(0, 0, -days)
	err := d.DB.Where("product_id = ? AND date >= ?", productID, startDate).
		Order("date DESC").
		Find(&analytics).Error
	return analytics, err
}

func (d *Database) CreateAnalytics(analytics *models.Analytics) error {
	return d.DB.Create(analytics).Error
}

func (d *Database) GetProductsByPlatform(platform string, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := d.DB.Where("platform = ?", platform).
		Order("trending_score DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (d *Database) GetAllProducts(limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	d.DB.Model(&models.Product{}).Count(&total)

	err := d.DB.Order("trending_score DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}
