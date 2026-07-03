package crawler

import (
	"fmt"
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"hotdeal-tracker/pkg/httpclient"
	"hotdeal-tracker/pkg/parser"
	"log"
	"strings"
	"sync"
	"time"
)

// Crawler 爬虫核心模块
// 负责从各电商平台爬取热销产品数据
// 支持多平台并发爬取，带有反爬机制
type Crawler struct {
	db        *database.Database      // 数据库连接
	client    *httpclient.HTTPClient  // HTTP客户端（支持反爬）
	cfg       *config.CrawlerConfig   // 爬虫配置
	parser    parser.ProductParser    // 产品解析器
	wg        sync.WaitGroup          // 等待组，用于并发控制
	taskQueue chan models.CrawlTask   // 任务队列
	mu        sync.Mutex              // 互斥锁，保护共享资源
}

// NewCrawler 创建爬虫实例
// 参数：
//   db: 数据库实例
//   cfg: 爬虫配置
// 返回：
//   *Crawler: 爬虫实例
func NewCrawler(db *database.Database, cfg *config.CrawlerConfig) *Crawler {
	return &Crawler{
		db:        db,
		client:    httpclient.NewHTTPClient(cfg), // 创建支持反爬的HTTP客户端
		cfg:       cfg,
		parser:    nil,
		taskQueue: make(chan models.CrawlTask, 100), // 缓冲队列，最多100个任务
	}
}

// Start 启动爬虫服务
// 为每个启用的平台启动一个并发爬虫协程
func (c *Crawler) Start() error {
	// 获取所有配置的平台
	platforms, err := c.db.GetPlatforms()
	if err != nil {
		return fmt.Errorf("failed to get platforms: %w", err)
	}

	// 为每个启用的平台启动并发爬取
	for _, platform := range platforms {
		if !c.isPlatformEnabled(platform.Code) {
			log.Printf("Platform %s is disabled, skipping", platform.Code)
			continue
		}

		c.wg.Add(1)
		go c.crawlPlatform(platform) // 并发爬取每个平台
	}

	go c.processTaskQueue()

	return nil
}

// isPlatformEnabled 检查平台是否启用
// 如果配置中没有指定平台列表，则默认全部启用
func (c *Crawler) isPlatformEnabled(platformCode string) bool {
	if c.cfg.Platforms == nil {
		return true
	}

	enabled, ok := c.cfg.Platforms[platformCode]
	return ok && enabled
}

// crawlPlatform 爬取单个平台的热销产品
// 流程：
// 1. 创建爬取任务并加入队列
// 2. 爬取分类信息
// 3. 添加请求延迟，避免被封IP
func (c *Crawler) crawlPlatform(platform models.Platform) {
	defer c.wg.Done() // 任务完成时通知等待组

	log.Printf("Starting to crawl platform: %s", platform.Name)

	// 获取热销页面URL
	hotURL := platform.HotURL
	if hotURL == "" {
		log.Printf("No hot URL for platform %s, using base URL", platform.Name)
		hotURL = platform.BaseURL
	}

	// 创建爬取任务
	task := models.CrawlTask{
		PlatformID: platform.ID,
		URL:        hotURL,
		Type:       "hot",        // 任务类型：热销产品
		Status:     "pending",    // 任务状态：待处理
		Priority:   platform.Priority, // 优先级
	}

	if err := c.db.CreateCrawlTask(&task); err != nil {
		log.Printf("Failed to create crawl task for %s: %v", platform.Name, err)
		return
	}

	// 将任务加入队列
	c.taskQueue <- task

	// 爬取分类信息
	c.crawlCategories(platform)

	// 添加请求延迟，模拟人类行为
	time.Sleep(time.Duration(c.cfg.Delay) * time.Millisecond)
}

// crawlCategories 爬取平台分类信息
// 获取平台的分类导航，用于后续按分类爬取产品
func (c *Crawler) crawlCategories(platform models.Platform) {
	// 如果没有配置分类URL，跳过
	if platform.CategoryURL == "" {
		return
	}

	// 获取分类页面HTML
	html, err := c.client.Get(platform.CategoryURL, nil)
	if err != nil {
		log.Printf("Failed to fetch categories for %s: %v", platform.Name, err)
		return
	}

	// 使用解析器提取分类信息
	c.parser = parser.GetParser(platform.Code)
	categories, err := c.parser.ParseCategories(html, platform.Code)
	if err != nil {
		log.Printf("Failed to parse categories for %s: %v", platform.Name, err)
		return
	}

	// 保存分类并创建爬取任务
	for _, category := range categories {
		category.Platform = platform.Code

		if err := c.db.CreateCategory(&category); err != nil {
			log.Printf("Failed to save category %s: %v", category.Name, err)
		}

		// 构建分类页面URL并创建爬取任务
		categoryURL := c.buildCategoryURL(platform, category)
		if categoryURL != "" {
			task := models.CrawlTask{
				PlatformID: platform.ID,
				CategoryID: &category.ID,
				Keyword:    category.Name,
				URL:        categoryURL,
				Type:       "category", // 任务类型：分类产品
				Status:     "pending",
				Priority:   5,
			}

			if err := c.db.CreateCrawlTask(&task); err != nil {
				log.Printf("Failed to create task for category %s: %v", category.Name, err)
			}

			c.taskQueue <- task
		}
	}
}

func (c *Crawler) buildCategoryURL(platform models.Platform, category models.Category) string {
	switch platform.Code {
	case "amazon":
		return fmt.Sprintf("https://www.amazon.com/s?k=%s&s=review-rank", category.Name)
	case "ebay":
		return fmt.Sprintf("https://www.ebay.com/sch/i.html?_nkw=%s&scl=1&mkcid=1&mkrid=711-53200-19255-0", category.Name)
	case "taobao":
		return fmt.Sprintf("https://s.taobao.com/search?q=%s", category.Name)
	case "jd":
		return fmt.Sprintf("https://search.jd.com/Search?keyword=%s&enc=utf-8", category.Name)
	case "aliexpress":
		return fmt.Sprintf("https://www.aliexpress.com/wholesale?SearchText=%s", category.Name)
	default:
		return ""
	}
}

func (c *Crawler) processTaskQueue() {
	concurrent := c.cfg.Concurrent
	if concurrent <= 0 {
		concurrent = 3
	}

	semaphore := make(chan struct{}, concurrent)

	for task := range c.taskQueue {
		semaphore <- struct{}{}

		c.wg.Add(1)
		go func(t models.CrawlTask) {
			defer c.wg.Done()
			defer func() { <-semaphore }()

			c.executeTask(t)
		}(task)
	}
}

func (c *Crawler) executeTask(task models.CrawlTask) {
	startTime := time.Now()
	task.Status = "running"
	task.StartedAt = &startTime

	if err := c.db.UpdateCrawlTask(&task); err != nil {
		log.Printf("Failed to update task status: %v", err)
	}

	platform, err := c.db.GetPlatformByCode(c.getPlatformCodeByID(task.PlatformID))
	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("Failed to get platform: %v", err)
		c.db.UpdateCrawlTask(&task)
		return
	}

	c.parser = parser.GetParser(platform.Code)

	var html []byte

	html, err = c.client.Get(task.URL, nil)

	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("Failed to fetch URL: %v", err)
		task.RetryCount++

		if task.RetryCount < c.cfg.RetryTimes {
			task.Status = "pending"
		}

		c.db.UpdateCrawlTask(&task)
		return
	}

	products, err := c.parser.ParseProducts(html, platform.Code)
	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = fmt.Sprintf("Failed to parse products: %v", err)
		c.db.UpdateCrawlTask(&task)
		return
	}

	productCount := 0
	for _, p := range products {
		p.PlatformID = c.ensurePlatformID(p.PlatformID, p.ProductURL, platform.Code)
		p.CrawledAt = time.Now()

		existing, _ := c.db.GetProductByPlatformID(p.PlatformID, platform.Code)
		if existing != nil {
			p.ID = existing.ID
			p.TrendingScore = c.calculateTrendingScore(existing, &p)
			if err := c.db.UpdateProduct(&p); err != nil {
				log.Printf("Failed to update product: %v", err)
			}

			history := models.PriceHistory{
				ProductID:     existing.ID,
				Price:         p.Price,
				OriginalPrice: p.OriginalPrice,
				SalesCount:    p.SalesCount,
				CrawledAt:    time.Now(),
			}
			c.db.CreatePriceHistory(&history)
		} else {
			if err := c.db.CreateProduct(&p); err != nil {
				log.Printf("Failed to create product: %v", err)
			}
		}

		productCount++
	}

	completedTime := time.Now()
	task.Status = "completed"
	task.CompletedAt = &completedTime
	task.ProductCount = productCount

	if err := c.db.UpdateCrawlTask(&task); err != nil {
		log.Printf("Failed to update task status: %v", err)
	}

	log.Printf("Task completed: %s, crawled %d products from %s", task.URL, productCount, platform.Name)

	time.Sleep(time.Duration(c.cfg.Delay) * time.Millisecond)
}

func (c *Crawler) getPlatformCodeByID(platformID uint) string {
	return ""
}

func (c *Crawler) ensurePlatformID(platformID, url, platform string) string {
	if platformID != "" {
		return platformID
	}
	return parser.ExtractProductID(url, platform)
}

func (c *Crawler) calculateTrendingScore(existing *models.Product, new *models.Product) float64 {
	score := existing.TrendingScore

	salesDiff := float64(new.SalesCount - existing.SalesCount)
	if salesDiff > 0 {
		score += salesDiff * 0.3
	}

	priceDiff := existing.Price - new.Price
	if priceDiff > 0 {
		discount := (priceDiff / existing.Price) * 100
		score += discount * 0.2
	}

	if new.Rating > existing.Rating {
		score += (new.Rating - existing.Rating) * 10
	}

	reviewDiff := float64(new.ReviewCount - existing.ReviewCount)
	if reviewDiff > 0 {
		score += reviewDiff * 0.1
	}

	return score
}

func (c *Crawler) Stop() {
	close(c.taskQueue)
	c.wg.Wait()
}

func (c *Crawler) CrawlKeyword(keyword string, platformCode string) ([]models.Product, error) {
	platform, err := c.db.GetPlatformByCode(platformCode)
	if err != nil {
		return nil, fmt.Errorf("platform not found: %w", err)
	}

	url := c.buildSearchURL(*platform, keyword)

	html, err := c.client.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	c.parser = parser.GetParser(platformCode)
	products, err := c.parser.ParseProducts(html, platformCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	for i := range products {
		products[i].PlatformID = c.ensurePlatformID(products[i].PlatformID, products[i].ProductURL, platformCode)
		products[i].CrawledAt = time.Now()
		products[i].Category = keyword

		if err := c.db.CreateProduct(&products[i]); err != nil {
			log.Printf("Failed to save product: %v", err)
		}
	}

	return products, nil
}

func (c *Crawler) buildSearchURL(platform models.Platform, keyword string) string {
	keyword = strings.ReplaceAll(keyword, " ", "+")

	switch platform.Code {
	case "amazon":
		return fmt.Sprintf("https://www.amazon.com/s?k=%s&s=review-rank", keyword)
	case "ebay":
		return fmt.Sprintf("https://www.ebay.com/sch/i.html?_nkw=%s", keyword)
	case "taobao":
		return fmt.Sprintf("https://s.taobao.com/search?q=%s", keyword)
	case "jd":
		return fmt.Sprintf("https://search.jd.com/Search?keyword=%s&enc=utf-8", keyword)
	case "tmall":
		return fmt.Sprintf("https://list.tmall.com/search_product.htm?q=%s", keyword)
	case "aliexpress":
		return fmt.Sprintf("https://www.aliexpress.com/wholesale?SearchText=%s", keyword)
	case "shopee":
		return fmt.Sprintf("https://www.shopee.com.my/search?keyword=%s", keyword)
	case "lazada":
		return fmt.Sprintf("https://www.lazada.com/shop/%s", keyword)
	default:
		return platform.BaseURL + "/search?q=" + keyword
	}
}

// CrawlHotProducts 爬取平台热销产品
// 流程：
// 1. 获取热销页面HTML
// 2. 使用解析器提取产品信息
// 3. 保存产品到数据库（更新已有产品或创建新产品）
// 4. 记录价格历史
//
// 参数：
//   platform: 平台信息
// 返回：
//   []models.Product: 爬取到的产品列表
//   error: 爬取错误
func (c *Crawler) CrawlHotProducts(platform models.Platform) ([]models.Product, error) {
	log.Printf("Crawling hot products from %s", platform.Name)
	log.Printf("URL: %s", platform.HotURL)

	// 获取对应平台的解析器
	c.parser = parser.GetParser(platform.Code)

	// 获取热销页面HTML
	html, err := c.client.Get(platform.HotURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	log.Printf("HTML length: %d bytes", len(html))

	// 使用解析器提取产品信息
	products, err := c.parser.ParseProducts(html, platform.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products: %w", err)
	}

	log.Printf("Parser returned %d products", len(products))

	// 处理每个产品
	for i, p := range products {
		// 确保PlatformID有效
		p.PlatformID = c.ensurePlatformID(p.PlatformID, p.ProductURL, platform.Code)
		p.CrawledAt = time.Now()
		p.IsHot = true
		p.TrendingScore = float64(100 - i) // 位置越靠前，热度越高

		// 检查产品是否已存在
		existing, _ := c.db.GetProductByPlatformID(p.PlatformID, platform.Code)
		if existing != nil {
			// 更新已有产品
			p.ID = existing.ID
			p.TrendingScore = c.calculateTrendingScore(existing, &p)
			if err := c.db.UpdateProduct(&p); err != nil {
				log.Printf("Failed to update product: %v", err)
			}

			// 记录价格历史
			history := models.PriceHistory{
				ProductID:     existing.ID,
				Price:         p.Price,
				OriginalPrice: p.OriginalPrice,
				SalesCount:    p.SalesCount,
				CrawledAt:    time.Now(),
			}
			c.db.CreatePriceHistory(&history)
		} else {
			// 创建新产品
			if err := c.db.CreateProduct(&p); err != nil {
				log.Printf("Failed to create product: %v", err)
			}
		}
	}

	log.Printf("Crawled %d hot products from %s", len(products), platform.Name)
	return products, nil
}

func (c *Crawler) CrawlHotProductsWithHTML(html []byte, platform models.Platform) ([]models.Product, error) {
	log.Printf("Parsing hot products from saved HTML for %s", platform.Name)

	c.parser = parser.GetParser(platform.Code)

	log.Printf("HTML length: %d bytes", len(html))

	products, err := c.parser.ParseProducts(html, platform.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products: %w", err)
	}

	log.Printf("Parser returned %d products", len(products))

	for i, p := range products {
		p.PlatformID = c.ensurePlatformID(p.PlatformID, p.ProductURL, platform.Code)
		p.CrawledAt = time.Now()
		p.IsHot = true
		p.TrendingScore = float64(100 - i)

		if p.Price > 0 && p.ProductURL != "" {
			existing, _ := c.db.GetProductByPlatformID(p.PlatformID, platform.Code)
			if existing != nil {
				p.ID = existing.ID
				p.TrendingScore = c.calculateTrendingScore(existing, &p)
				if err := c.db.UpdateProduct(&p); err != nil {
					log.Printf("Failed to update product: %v", err)
				}

				history := models.PriceHistory{
					ProductID:     existing.ID,
					Price:         p.Price,
					OriginalPrice: p.OriginalPrice,
					SalesCount:    p.SalesCount,
					CrawledAt:    time.Now(),
				}
				c.db.CreatePriceHistory(&history)
			} else {
				if err := c.db.CreateProduct(&p); err != nil {
					log.Printf("Failed to create product: %v", err)
				}
			}
		}
	}

	log.Printf("Crawled %d hot products from %s", len(products), platform.Name)
	return products, nil
}

func (c *Crawler) CrawlProductDetail(productURL string, platformCode string) (*models.Product, error) {
	html, err := c.client.Get(productURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	c.parser = parser.GetParser(platformCode)
	products, err := c.parser.ParseProducts(html, platformCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no products found")
	}

	product := &products[0]
	product.PlatformID = c.ensurePlatformID(product.PlatformID, productURL, platformCode)

	err = c.parser.ParseProductDetail(html, product)
	if err != nil {
		log.Printf("Failed to parse product detail: %v", err)
	}

	return product, nil
}
