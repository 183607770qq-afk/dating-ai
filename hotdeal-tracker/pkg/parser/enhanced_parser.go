package parser

import (
	"fmt"
	"hotdeal-tracker/internal/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// EnhancedParser 增强版产品解析器
// 支持从HTML和JSON两种格式中提取产品数据
// 采用CSS选择器+正则表达式混合策略，提高多平台兼容性
type EnhancedParser struct {
	BaseParser // 嵌入基础解析器，继承其方法（如CleanText, ParsePrice等）
}

// ParseProducts 解析产品列表
// 解析策略：
// 1. 首先尝试使用预设的CSS选择器提取产品（适用于静态HTML页面）
// 2. 如果第一步失败，尝试从包含"product"或"item"类名的元素中提取
// 3. 如果提取到的产品价格大部分为0，尝试从页面内嵌的JSON数据中提取（适用于动态渲染页面）
//
// 参数：
//   html: 网页HTML内容
//   platform: 平台标识（如 ebay, amazon, taobao）
// 返回：
//   []models.Product: 产品列表
//   error: 解析错误
func (ep *EnhancedParser) ParseProducts(html []byte, platform string) ([]models.Product, error) {
	// 将HTML字符串转换为goquery文档对象
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var products []models.Product

	// 步骤1: 尝试预设的CSS选择器
	// 这些选择器覆盖了主流电商平台的产品列表结构
	selectors := []string{
		".zg-item-immersion",          // Amazon热销榜单
		"[data-component-type='s-search-result']", // Amazon搜索结果
		".product-item",               // 通用产品项
		".product-card",               // 卡片式产品
		".item",                       // 通用项
		".goods-item",                 // 淘宝/京东商品项
		".list-item",                  // 列表项
		".search-item",                // 搜索结果项
		".item-box",                   // 商品盒子
		".product-list li",            // 产品列表
		".grid-item",                  // 网格项
		".result-item",                // 结果项
		".item-container",             // 商品容器
	}

	for _, selector := range selectors {
		found := false
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			found = true
			product := ep.parseProductItem(s, platform, len(products)+i)
			if product.Title != "" || product.ProductURL != "" {
				products = append(products, product)
			}
		})
		// 找到产品后立即停止尝试其他选择器
		if found && len(products) > 0 {
			break
		}
	}

	// 步骤2: 如果没有找到产品，尝试更宽泛的选择器
	if len(products) == 0 {
		doc.Find("div, article, li").Each(func(i int, s *goquery.Selection) {
			class := s.AttrOr("class", "")
			if strings.Contains(class, "product") || strings.Contains(class, "item") {
				product := ep.parseProductItem(s, platform, i)
				if product.Title != "" && product.ProductURL != "" {
					products = append(products, product)
				}
			}
		})
	}

	// 步骤3: 如果产品价格大部分为0（说明CSS解析失败），尝试从JSON中提取
	// 动态渲染页面（如eBay）的价格数据通常内嵌在JavaScript的JSON对象中
	if len(products) == 0 || ep.hasNoPrices(products) {
		jsonProducts := ep.parseFromJSON(html, platform)
		if len(jsonProducts) > 0 {
			products = append(products, jsonProducts...)
		}
	}

	return products, nil
}

// hasNoPrices 判断产品列表是否大部分价格为0
// 当80%以上的产品价格为0时，说明CSS解析可能失败，需要切换到JSON解析策略
func (ep *EnhancedParser) hasNoPrices(products []models.Product) bool {
	if len(products) == 0 {
		return true
	}
	zeroPriceCount := 0
	for _, p := range products {
		if p.Price == 0 {
			zeroPriceCount++
		}
	}
	// 如果超过80%的产品价格为0，返回true
	return float64(zeroPriceCount)/float64(len(products)) > 0.8
}

// parseProductItem 从单个HTML元素中解析产品信息
// 依次提取：标题、链接、价格、评分、评论数、销量、图片
//
// 参数：
//   s: goquery选择器对象，代表一个产品元素
//   platform: 平台标识
//   index: 产品在列表中的位置（用于计算热度分数）
// 返回：
//   models.Product: 解析后的产品对象
func (ep *EnhancedParser) parseProductItem(s *goquery.Selection, platform string, index int) models.Product {
	// 初始化产品对象，设置平台和热度标记
	product := models.Product{
		Platform:      platform,
		IsHot:         true,
		TrendingScore: float64(100 - index), // 位置越靠前，热度分数越高
	}

	// 1. 提取标题
	titleSelectors := []string{
		"h2", "h3", ".title", ".product-title", ".goods-title", ".item-title",
		"[class*='title']", ".name", ".product-name", ".goods-name",
		".listing-title", ".item-name", ".title-text",
	}
	for _, sel := range titleSelectors {
		titleSel := s.Find(sel)
		if titleSel.Length() > 0 {
			title := ep.CleanText(titleSel.First().Text())
			if title != "" && len(title) > 5 { // 标题至少5个字符才有效
				product.Title = title
				break
			}
		}
	}

	// 2. 提取产品链接
	linkSelectors := []string{
		"a", ".product-link", ".item-link", "[href*='/dp/']", "[href*='/product/']",
	}
	for _, sel := range linkSelectors {
		linkSel := s.Find(sel)
		if linkSel.Length() > 0 {
			href, exists := linkSel.First().Attr("href")
			if exists && href != "" {
				product.ProductURL = ep.fixURL(href, platform) // 修复相对URL
				product.PlatformID = ExtractProductID(product.ProductURL, platform) // 提取产品ID
				break
			}
		}
	}

	// 3. 提取价格
	priceSelectors := []string{
		".price", ".product-price", ".item-price", ".goods-price",
		"[class*='price']", ".p-price", ".price-current",
		".deal-price", ".sales-price", ".cost",
	}
	for _, sel := range priceSelectors {
		priceSel := s.Find(sel)
		if priceSel.Length() > 0 {
			priceText := ep.CleanText(priceSel.First().Text())
			if priceText != "" {
				price, currency := ParsePrice(priceText) // 解析价格和货币
				product.Price = price
				product.Currency = currency
				break
			}
		}
	}

	// 4. 提取评分
	ratingSelectors := []string{
		".rating", ".star-rating", "[aria-label*='out of']",
		".score", ".product-rating", ".item-rating",
	}
	for _, sel := range ratingSelectors {
		ratingSel := s.Find(sel)
		if ratingSel.Length() > 0 {
			ratingText := ep.CleanText(ratingSel.First().Text())
			if ratingText != "" {
				product.Rating = ep.ExtractRating(ratingText)
				break
			}
		}
	}

	// 5. 提取评论数
	reviewSelectors := []string{
		".review-count", ".reviews", ".review-num", ".item-review",
		"[class*='review']", ".rating-count",
	}
	for _, sel := range reviewSelectors {
		reviewSel := s.Find(sel)
		if reviewSel.Length() > 0 {
			reviewText := ep.CleanText(reviewSel.First().Text())
			product.ReviewCount = int(ep.ExtractNumber(reviewText))
			if product.ReviewCount > 0 {
				break
			}
		}
	}

	// 6. 提取销量
	salesSelectors := []string{
		".sales", ".sales-count", ".sold", ".volume",
		"[class*='sales']", "[class*='sold']",
	}
	for _, sel := range salesSelectors {
		salesSel := s.Find(sel)
		if salesSel.Length() > 0 {
			salesText := ep.CleanText(salesSel.First().Text())
			product.SalesCount = int(ep.ExtractNumber(salesText))
			if product.SalesCount > 0 {
				break
			}
		}
	}

	// 7. 提取图片URL
	imageSelectors := []string{
		"img", ".product-image", ".item-image", ".goods-img",
	}
	for _, sel := range imageSelectors {
		imgSel := s.Find(sel)
		if imgSel.Length() > 0 {
			src, exists := imgSel.First().Attr("src")
			if exists && src != "" {
				product.ImageURL = src
				break
			}
		}
	}

	return product
}

// ExtractRating 从文本中提取评分（0-5分）
// 支持多种格式：
//   - "4.5 stars"
//   - "rating 4.8"
//   - "4.5 out of 5"
//   - "4.5分"
//   - "评分4.8"
func (ep *EnhancedParser) ExtractRating(text string) float64 {
	// 模式1: 匹配包含stars/rating/out of/分/评分等关键词的文本
	re := regexp.MustCompile(`(\d+\.?\d*)\s*(?:stars?|rating|out of|分|评分)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		rating, _ := strconv.ParseFloat(matches[1], 64)
		return rating
	}

	// 模式2: 直接提取数字，然后验证是否在0-5范围内
	re2 := regexp.MustCompile(`(\d+\.?\d*)`)
	matches = re2.FindStringSubmatch(text)
	if len(matches) > 1 {
		rating, _ := strconv.ParseFloat(matches[1], 64)
		if rating >= 0 && rating <= 5 {
			return rating
		}
	}

	return 0
}

// fixURL 修复相对URL，转换为完整的绝对URL
// 处理三种情况：
// 1. 协议相对URL: "//xxx.com/path" -> "https://xxx.com/path"
// 2. 相对路径: "/path" 或 "path" -> "https://www.platform.com/path"
// 3. 绝对URL: 保持不变
func (ep *EnhancedParser) fixURL(href, platform string) string {
	// 处理协议相对URL（如 //www.ebay.com/xxx）
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	// 如果不是完整的HTTP URL，需要补充基础URL
	if !strings.HasPrefix(href, "http") {
		baseURLs := map[string]string{
			"amazon":     "https://www.amazon.com",
			"ebay":       "https://www.ebay.com",
			"taobao":     "https://www.taobao.com",
			"jd":         "https://www.jd.com",
			"tmall":      "https://www.tmall.com",
			"pdd":        "https://www.pinduoduo.com",
			"aliexpress": "https://www.aliexpress.com",
			"shopee":     "https://www.shopee.com",
			"lazada":     "https://www.lazada.com",
			"wish":       "https://www.wish.com",
		}

		if base, ok := baseURLs[platform]; ok {
			if strings.HasPrefix(href, "/") {
				return base + href
			}
			return base + "/" + href
		}
	}

	return href
}

// ParseCategories 从页面中提取分类信息
// 用于构建网站导航结构，支持后续按分类爬取
func (ep *EnhancedParser) ParseCategories(html []byte, platform string) ([]models.Category, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var categories []models.Category

	// 分类选择器列表
	catSelectors := []string{
		"a.category", ".category-item a", ".nav-item a",
		".menu-item a", "[href*='/category/']", "[href*='/browse/']",
	}

	for _, selector := range catSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			name := ep.CleanText(s.Text())
			href, _ := s.Attr("href")

			// 过滤无效分类（名称至少2个字符，URL不为空）
			if name != "" && len(name) > 2 && href != "" {
				categories = append(categories, models.Category{
					Name:     name,
					Platform: platform,
					Slug:     href,
				})
			}
		})

		if len(categories) > 0 {
			break
		}
	}

	return categories, nil
}

// ParseProductDetail 解析产品详情页
// 补充产品的描述和图片信息
// 参数 product 是指针类型，直接修改原对象
func (ep *EnhancedParser) ParseProductDetail(html []byte, product *models.Product) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return err
	}

	// 提取产品描述
	descSelectors := []string{
		"#productDescription", ".product-desc", ".description",
		"meta[name='description']", ".detail-desc",
	}

	for _, sel := range descSelectors {
		descSel := doc.Find(sel)
		if descSel.Length() > 0 {
			if sel == "meta[name='description']" {
				product.Description, _ = descSel.Attr("content")
			} else {
				product.Description = ep.CleanText(descSel.Text())
			}
			if product.Description != "" {
				break
			}
		}
	}

	// 提取产品主图
	imgSelectors := []string{
		"#landingImage", ".main-image", "meta[property='og:image']",
		".product-img", "[data-src]",
	}

	for _, sel := range imgSelectors {
		imgSel := doc.Find(sel)
		if imgSel.Length() > 0 {
			if strings.Contains(sel, "meta") {
				product.ImageURL, _ = imgSel.Attr("content")
			} else {
				product.ImageURL, _ = imgSel.Attr("src")
				if product.ImageURL == "" {
					// 如果src为空，尝试data-src属性（懒加载图片）
					product.ImageURL, _ = imgSel.Attr("data-src")
				}
			}
			if product.ImageURL != "" {
				break
			}
		}
	}

	return nil
}

// GetParser 获取产品解析器实例
// 目前统一返回 EnhancedParser，后续可扩展为根据平台返回不同解析器
func GetParser(platform string) ProductParser {
	return &EnhancedParser{}
}

// parseFromJSON 从页面内嵌的JSON数据中提取产品信息
// 适用于动态渲染页面（如eBay），这类页面的产品数据通常以JSON格式内嵌在JavaScript中
// 
// 解析策略：
// 1. 检测页面是否包含"ListingItemCard"关键字（eBay产品卡片标识）
// 2. 使用正则表达式提取JSON对象
// 3. 从JSON中提取产品URL、标题、价格等信息
func (ep *EnhancedParser) parseFromJSON(html []byte, platform string) []models.Product {
	var products []models.Product
	htmlStr := string(html)

	// eBay平台使用ListingItemCard作为产品卡片的标识
	if strings.Contains(htmlStr, "ListingItemCard") {
		// 正则匹配包含ListingItemCard的JSON对象
		re := regexp.MustCompile(`"model":\s*(\{[^}]+"_type":"ListingItemCard"[^}]+\})`)
		matches := re.FindAllStringSubmatch(htmlStr, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				// 从匹配到的JSON片段中提取各字段
				urlMatch := regexp.MustCompile(`"URL":"([^"]+)"`).FindStringSubmatch(match[1])
				titleMatch := regexp.MustCompile(`"text":"([^"]+)"`).FindStringSubmatch(match[1])
				priceMatch := regexp.MustCompile(`"displayPrice":\s*\{"_type":"TextualDisplayValue","value":\{"value":([\d.]+)`)
				priceMatches := priceMatch.FindStringSubmatch(match[1])

				url := ""
				if len(urlMatch) >= 2 {
					url = urlMatch[1]
				}

				title := fmt.Sprintf("Product %d", i+1)
				if len(titleMatch) >= 2 {
					title = titleMatch[1]
				}

				price := 0.0
				if len(priceMatches) >= 2 {
					price, _ = strconv.ParseFloat(priceMatches[1], 64)
				}

				product := models.Product{
					Title:         title,
					ProductURL:    url,
					Platform:      platform,
					Price:         price,
					Currency:      "CNY",
					IsHot:         true,
					TrendingScore: float64(100 - i),
				}

				// 只添加有效产品（价格>0且URL不为空）
				if price > 0 && url != "" {
					products = append(products, product)
				}
			}
		}
	}

	return products
}
