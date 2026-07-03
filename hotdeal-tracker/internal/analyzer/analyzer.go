package analyzer

import (
	"fmt"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"sort"
	"time"
)

type Analyzer struct {
	db *database.Database
}

func NewAnalyzer(db *database.Database) *Analyzer {
	return &Analyzer{db: db}
}

type HotProduct struct {
	Product     models.Product `json:"product"`
	HotScore    float64        `json:"hot_score"`
	Trend       string         `json:"trend"`
	PriceChange float64        `json:"price_change"`
}

type CategoryStat struct {
	Category     string  `json:"category"`
	ProductCount int     `json:"product_count"`
	AvgPrice     float64 `json:"avg_price"`
	AvgRating    float64 `json:"avg_rating"`
	TopPlatform  string  `json:"top_platform"`
}

type PlatformStat struct {
	Platform     string  `json:"platform"`
	ProductCount int     `json:"product_count"`
	AvgPrice     float64 `json:"avg_price"`
	AvgRating    float64 `json:"avg_rating"`
	TotalSales   int     `json:"total_sales"`
}

type PriceTrend struct {
	Date         string  `json:"date"`
	AvgPrice     float64 `json:"avg_price"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	ProductCount int     `json:"product_count"`
}

func (a *Analyzer) GetHotProducts(limit int, platform string) ([]HotProduct, error) {
	products, err := a.db.GetHotProducts(platform, limit)
	if err != nil {
		return nil, err
	}

	var hotProducts []HotProduct
	for _, p := range products {
		history, _ := a.db.GetPriceHistory(p.ID, 30)

		hotScore := a.calculateHotScore(&p, history)
		trend := a.determineTrend(history)
		priceChange := a.calculatePriceChange(history)

		hotProducts = append(hotProducts, HotProduct{
			Product:     p,
			HotScore:    hotScore,
			Trend:       trend,
			PriceChange: priceChange,
		})
	}

	sort.Slice(hotProducts, func(i, j int) bool {
		return hotProducts[i].HotScore > hotProducts[j].HotScore
	})

	return hotProducts, nil
}

func (a *Analyzer) calculateHotScore(product *models.Product, history []models.PriceHistory) float64 {
	score := product.TrendingScore

	if product.IsHot {
		score += 50
	}

	if product.Rating >= 4.5 {
		score += 30
	} else if product.Rating >= 4.0 {
		score += 20
	} else if product.Rating >= 3.0 {
		score += 10
	}

	if product.SalesCount > 10000 {
		score += 40
	} else if product.SalesCount > 5000 {
		score += 30
	} else if product.SalesCount > 1000 {
		score += 20
	} else if product.SalesCount > 100 {
		score += 10
	}

	if len(history) >= 2 {
		latestPrice := history[len(history)-1].Price
		if latestPrice > 0 && product.Price < latestPrice {
			discount := ((latestPrice - product.Price) / latestPrice) * 100
			score += discount * 0.5
		}
	}

	if product.ReviewCount > 1000 {
		score += 20
	} else if product.ReviewCount > 100 {
		score += 10
	}

	return score
}

func (a *Analyzer) determineTrend(history []models.PriceHistory) string {
	if len(history) < 3 {
		return "stable"
	}

	recentPrices := make([]float64, 0)
	for i := len(history) - 3; i < len(history); i++ {
		recentPrices = append(recentPrices, history[i].Price)
	}

	first := recentPrices[0]
	last := recentPrices[len(recentPrices)-1]

	change := ((last - first) / first) * 100

	if change > 5 {
		return "rising"
	} else if change < -5 {
		return "falling"
	}

	return "stable"
}

func (a *Analyzer) calculatePriceChange(history []models.PriceHistory) float64 {
	if len(history) < 2 {
		return 0
	}

	oldest := history[0]
	newest := history[len(history)-1]

	if oldest.Price == 0 {
		return 0
	}

	return ((newest.Price - oldest.Price) / oldest.Price) * 100
}

func (a *Analyzer) GetCategoryStats() ([]CategoryStat, error) {
	var stats []CategoryStat

	categories, err := a.db.GetCategories("")
	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		products, err := a.db.GetProductsByCategory(cat.Name, 100, 0)
		if err != nil {
			continue
		}

		if len(products) == 0 {
			continue
		}

		var totalPrice, totalRating float64
		platformCounts := make(map[string]int)

		for _, p := range products {
			totalPrice += p.Price
			totalRating += p.Rating
			platformCounts[p.Platform]++
		}

		topPlatform := ""
		maxCount := 0
		for platform, count := range platformCounts {
			if count > maxCount {
				maxCount = count
				topPlatform = platform
			}
		}

		stats = append(stats, CategoryStat{
			Category:     cat.Name,
			ProductCount: len(products),
			AvgPrice:     totalPrice / float64(len(products)),
			AvgRating:    totalRating / float64(len(products)),
			TopPlatform:  topPlatform,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ProductCount > stats[j].ProductCount
	})

	return stats, nil
}

func (a *Analyzer) GetPlatformStats() ([]PlatformStat, error) {
	platforms, err := a.db.GetPlatforms()
	if err != nil {
		return nil, err
	}

	var stats []PlatformStat

	for _, platform := range platforms {
		products, err := a.db.GetProductsByPlatform(platform.Code, 1000, 0)
		if err != nil {
			continue
		}

		if len(products) == 0 {
			continue
		}

		var totalPrice, totalRating, totalSales float64

		for _, p := range products {
			totalPrice += p.Price
			totalRating += p.Rating
			totalSales += float64(p.SalesCount)
		}

		stats = append(stats, PlatformStat{
			Platform:     platform.Name,
			ProductCount: len(products),
			AvgPrice:     totalPrice / float64(len(products)),
			AvgRating:    totalRating / float64(len(products)),
			TotalSales:   int(totalSales),
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ProductCount > stats[j].ProductCount
	})

	return stats, nil
}

func (a *Analyzer) GetPriceTrends(productID uint, days int) ([]PriceTrend, error) {
	history, err := a.db.GetPriceHistory(productID, days)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no price history found")
	}

	dailyPrices := make(map[string][]float64)

	for _, h := range history {
		date := h.CrawledAt.Format("2006-01-02")
		dailyPrices[date] = append(dailyPrices[date], h.Price)
	}

	var trends []PriceTrend
	dates := make([]string, 0, len(dailyPrices))
	for date := range dailyPrices {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		prices := dailyPrices[date]
		if len(prices) == 0 {
			continue
		}

		var min, max, sum float64
		min = prices[0]
		max = prices[0]

		for _, p := range prices {
			if p < min {
				min = p
			}
			if p > max {
				max = p
			}
			sum += p
		}

		trends = append(trends, PriceTrend{
			Date:         date,
			AvgPrice:     sum / float64(len(prices)),
			MinPrice:     min,
			MaxPrice:     max,
			ProductCount: len(prices),
		})
	}

	return trends, nil
}

func (a *Analyzer) GetMarketInsights(platform string) (map[string]interface{}, error) {
	insights := make(map[string]interface{})

	products, err := a.db.GetHotProducts(platform, 100)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return insights, nil
	}

	var totalPrice, totalRating float64
	priceRanges := make(map[string]int)
	categories := make(map[string]int)
	platformCounts := make(map[string]int)

	for _, p := range products {
		totalPrice += p.Price
		totalRating += p.Rating
		categories[p.Category]++
		platformCounts[p.Platform]++

		priceRange := a.getPriceRange(p.Price)
		priceRanges[priceRange]++
	}

	insights["total_products"] = len(products)
	insights["avg_price"] = totalPrice / float64(len(products))
	insights["avg_rating"] = totalRating / float64(len(products))
	insights["price_distribution"] = priceRanges
	insights["top_categories"] = a.getTopN(categories, 5)
	insights["platform_distribution"] = platformCounts

	return insights, nil
}

func (a *Analyzer) getPriceRange(price float64) string {
	switch {
	case price < 10:
		return "$0-10"
	case price < 25:
		return "$10-25"
	case price < 50:
		return "$25-50"
	case price < 100:
		return "$50-100"
	case price < 250:
		return "$100-250"
	case price < 500:
		return "$250-500"
	default:
		return "$500+"
	}
}

func (a *Analyzer) getTopN(m map[string]int, n int) []map[string]interface{} {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var result []map[string]interface{}
	for i := 0; i < n && i < len(ss); i++ {
		result = append(result, map[string]interface{}{
			"name":  ss[i].Key,
			"count": ss[i].Value,
		})
	}

	return result
}

func (a *Analyzer) AnalyzeProduct(productID uint) (map[string]interface{}, error) {
	product, err := a.getProductByID(productID)
	if err != nil {
		return nil, err
	}

	history, err := a.db.GetPriceHistory(productID, 90)
	if err != nil {
		return nil, err
	}

	analytics, err := a.db.GetAnalytics(productID, 30)
	if err != nil {
		return nil, err
	}

	analysis := map[string]interface{}{
		"product":     product,
		"price_stats": a.calculatePriceStats(history),
		"engagement":  a.calculateEngagement(analytics),
		"recommendation": a.generateRecommendation(product, history),
	}

	return analysis, nil
}

func (a *Analyzer) getProductByID(id uint) (*models.Product, error) {
	return nil, nil
}

func (a *Analyzer) calculatePriceStats(history []models.PriceHistory) map[string]interface{} {
	if len(history) == 0 {
		return nil
	}

	var prices []float64
	for _, h := range history {
		prices = append(prices, h.Price)
	}

	min := prices[0]
	max := prices[0]
	sum := 0.0

	for _, p := range prices {
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
		sum += p
	}

	currentPrice := prices[len(prices)-1]
	avgPrice := sum / float64(len(prices))

	return map[string]interface{}{
		"current_price":    currentPrice,
		"average_price":    avgPrice,
		"min_price":        min,
		"max_price":        max,
		"price_volatility": max - min,
		"discount_from_avg": ((avgPrice - currentPrice) / avgPrice) * 100,
	}
}

func (a *Analyzer) calculateEngagement(analytics []models.Analytics) map[string]interface{} {
	if len(analytics) == 0 {
		return map[string]interface{}{
			"total_views":      0,
			"total_purchases":  0,
			"conversion_rate":  0,
		}
	}

	var totalViews, totalPurchases int
	var totalConversion float64

	for _, a := range analytics {
		totalViews += a.Views
		totalPurchases += a.Purchases
		totalConversion += a.ConversionRate
	}

	return map[string]interface{}{
		"total_views":       totalViews,
		"total_purchases":   totalPurchases,
		"avg_conversion_rate": totalConversion / float64(len(analytics)),
	}
}

func (a *Analyzer) generateRecommendation(product *models.Product, history []models.PriceHistory) string {
	if product == nil {
		return "No data available"
	}

	if product.Rating >= 4.5 && product.SalesCount > 1000 {
		return "Strong buy recommendation - High rating and strong sales"
	}

	if len(history) >= 5 {
		currentPrice := history[len(history)-1].Price
		var sum float64
		for _, h := range history {
			sum += h.Price
		}
		avgPrice := sum / float64(len(history))

		if currentPrice < avgPrice*0.9 {
			return "Good buy opportunity - Price is below average"
		} else if currentPrice > avgPrice*1.1 {
			return "Wait for better price - Currently above average"
		}
	}

	return "Hold - No significant price movement"
}

func (a *Analyzer) GetTrendingKeywords(days int, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 10
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	keywordCounts := make(map[string]int)

	products, _, err := a.db.GetAllProducts(1000, 0)
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		if p.CrawledAt.After(startDate) && p.CrawledAt.Before(endDate) {
			keywordCounts[p.Category]++
			if p.Tags != "" {
				tags := splitTags(p.Tags)
				for _, tag := range tags {
					keywordCounts[tag]++
				}
			}
		}
	}

	var keywords []map[string]interface{}
	for keyword, count := range keywordCounts {
		keywords = append(keywords, map[string]interface{}{
			"keyword": keyword,
			"count":   count,
		})
	}

	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i]["count"].(int) > keywords[j]["count"].(int)
	})

	if len(keywords) > limit {
		keywords = keywords[:limit]
	}

	return keywords, nil
}

func splitTags(tags string) []string {
	var result []string
	current := ""

	for _, c := range tags {
		if c == ',' || c == ';' || c == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}
