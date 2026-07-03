package parser

import (
	"fmt"
	"hotdeal-tracker/internal/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProductParser interface {
	ParseProducts(html []byte, platform string) ([]models.Product, error)
	ParseCategories(html []byte, platform string) ([]models.Category, error)
	ParseProductDetail(html []byte, product *models.Product) error
}

type BaseParser struct{}

func (p *BaseParser) CleanText(text string) string {
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`[\n\r\t]`).ReplaceAllString(text, "")
	return text
}

func (p *BaseParser) ExtractNumber(text string) float64 {
	re := regexp.MustCompile(`[\d,.]+`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}
	match = strings.ReplaceAll(match, ",", "")
	match = strings.ReplaceAll(match, ",", "")
	if strings.Count(match, ".") > 1 {
		parts := strings.Split(match, ".")
		match = parts[0] + "." + strings.Join(parts[1:], "")
	}
	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}
	return price
}

func (p *BaseParser) ExtractDiscount(original, current float64) int {
	if original <= 0 || current <= 0 {
		return 0
	}
	discount := ((original - current) / original) * 100
	return int(discount)
}

type AmazonParser struct {
	BaseParser
}

func (ap *AmazonParser) ParseProducts(html []byte, platform string) ([]models.Product, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var products []models.Product

	doc.Find(".zg-item-immersion").Each(func(i int, s *goquery.Selection) {
		product := models.Product{
			Platform: platform,
		}

		titleSel := s.Find(".p13n-sc-truncated")
		if titleSel.Length() > 0 {
			product.Title = ap.CleanText(titleSel.Text())
		}

		linkSel := s.Find(".a-link-normal")
		if linkSel.Length() > 0 {
			href, _ := linkSel.Attr("href")
			product.ProductURL = "https://www.amazon.com" + href
			product.PlatformID = ap.ExtractAmazonProductID(href)
		}

		priceSel := s.Find(".p13n-sc-price")
		if priceSel.Length() > 0 {
			product.Price = ap.ExtractNumber(priceSel.Text())
		}

		ratingSel := s.Find(".a-icon-alt")
		if ratingSel.Length() > 0 {
			ratingStr := ratingSel.Text()
			rating, _ := strconv.ParseFloat(strings.Split(ratingStr, " ")[0], 64)
			product.Rating = rating
		}

		reviewsSel := s.Find(".a-size-small")
		if reviewsSel.Length() > 0 {
			product.ReviewCount = int(ap.ExtractNumber(reviewsSel.Text()))
		}

		product.Badge = ap.DetermineBadge(s)
		product.IsHot = true
		product.TrendingScore = float64(100 - i)

		products = append(products, product)
	})

	if len(products) == 0 {
		doc.Find("[data-component-type='s-search-result']").Each(func(i int, s *goquery.Selection) {
			product := models.Product{
				Platform: platform,
			}

			titleSel := s.Find("h2 a span")
			if titleSel.Length() > 0 {
				product.Title = ap.CleanText(titleSel.First().Text())
			}

			linkSel := s.Find("h2 a")
			if linkSel.Length() > 0 {
				href, _ := linkSel.Attr("href")
				product.ProductURL = "https://www.amazon.com" + href
				product.PlatformID = ap.ExtractAmazonProductID(href)
			}

			priceWholeSel := s.Find(".a-price-whole")
			if priceWholeSel.Length() > 0 {
				priceFracSel := s.Find(".a-price-fraction")
				whole := ap.CleanText(priceWholeSel.First().Text())
				var frac string
				if priceFracSel.Length() > 0 {
					frac = "." + ap.CleanText(priceFracSel.First().Text())
				} else {
					frac = ".00"
				}
				product.Price, _ = strconv.ParseFloat(whole+frac, 64)
			}

			ratingSel := s.Find("[aria-label*='out of']")
			if ratingSel.Length() > 0 {
				ratingStr := ratingSel.AttrOr("aria-label", "0 out of 5 stars")
				rating, _ := strconv.ParseFloat(strings.Split(ratingStr, " ")[0], 64)
				product.Rating = rating
			}

			reviewsSel := s.Find(".a-size-base")
			if reviewsSel.Length() > 0 {
				product.ReviewCount = int(ap.ExtractNumber(reviewsSel.First().Text()))
			}

			product.IsHot = true
			product.TrendingScore = float64(100 - i)

			products = append(products, product)
		})
	}

	return products, nil
}

func (ap *AmazonParser) ExtractAmazonProductID(url string) string {
	re := regexp.MustCompile(`/dp/([A-Z0-9]{10})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (ap *AmazonParser) DetermineBadge(s *goquery.Selection) string {
	badgeSel := s.Find(".BadgeName")
	if badgeSel.Length() > 0 {
		badge := ap.CleanText(badgeSel.Text())
		return badge
	}
	return "#1 Best Seller"
}

func (ap *AmazonParser) ParseCategories(html []byte, platform string) ([]models.Category, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var categories []models.Category

	doc.Find("a.category").Each(func(i int, s *goquery.Selection) {
		category := models.Category{
			Platform: platform,
		}

		nameSel := s.Find(".category-name")
		if nameSel.Length() > 0 {
			category.Name = ap.CleanText(nameSel.Text())
		} else {
			category.Name = ap.CleanText(s.Text())
		}

		href, _ := s.Attr("href")
		category.Slug = strings.TrimPrefix(href, "/")

		categories = append(categories, category)
	})

	return categories, nil
}

func (ap *AmazonParser) ParseProductDetail(html []byte, product *models.Product) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return err
	}

	descSel := doc.Find("#productDescription")
	if descSel.Length() > 0 {
		product.Description = ap.CleanText(descSel.Text())
	}

	imgSel := doc.Find("#landingImage")
	if imgSel.Length() > 0 {
		product.ImageURL, _ = imgSel.Attr("src")
	}

	priceSel := doc.Find("#priceblock_ourprice")
	if priceSel.Length() > 0 {
		product.Price = ap.ExtractNumber(priceSel.Text())
	} else {
		priceSel = doc.Find("#priceblock_dealprice")
		if priceSel.Length() > 0 {
			product.Price = ap.ExtractNumber(priceSel.Text())
		}
	}

	return nil
}

type GenericParser struct {
	BaseParser
}

func (gp *GenericParser) ParseProducts(html []byte, platform string) ([]models.Product, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var products []models.Product

	doc.Find(".product-item, .product-card, .item, [class*='product']").Each(func(i int, s *goquery.Selection) {
		product := models.Product{
			Platform: platform,
		}

		titleSel := s.Find("h3, .title, .product-title, [class*='title']")
		if titleSel.Length() > 0 {
			product.Title = gp.CleanText(titleSel.First().Text())
		}

		linkSel := s.Find("a")
		if linkSel.Length() > 0 {
			href, _ := linkSel.Attr("href")
			product.ProductURL = href
		}

		priceSel := s.Find(".price, [class*='price']")
		if priceSel.Length() > 0 {
			product.Price = gp.ExtractNumber(priceSel.First().Text())
		}

		product.IsHot = true
		product.TrendingScore = float64(100 - i)

		products = append(products, product)
	})

	return products, nil
}

func (gp *GenericParser) ParseCategories(html []byte, platform string) ([]models.Category, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var categories []models.Category

	doc.Find("a.category-link, .category a, nav a").Each(func(i int, s *goquery.Selection) {
		category := models.Category{
			Platform: platform,
			Name:     gp.CleanText(s.Text()),
		}

		href, _ := s.Attr("href")
		category.Slug = href

		categories = append(categories, category)
	})

	return categories, nil
}

func (gp *GenericParser) ParseProductDetail(html []byte, product *models.Product) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return err
	}

	descSel := doc.Find("meta[name='description']")
	if descSel.Length() > 0 {
		product.Description, _ = descSel.Attr("content")
	}

	imgSel := doc.Find("meta[property='og:image']")
	if imgSel.Length() > 0 {
		product.ImageURL, _ = imgSel.Attr("content")
	}

	return nil
}

func ParsePrice(priceStr string) (float64, string) {
	priceStr = strings.TrimSpace(priceStr)

	currencyMap := map[string]string{
		"$":  "USD",
		"€":  "EUR",
		"£":  "GBP",
		"¥":  "CNY",
		"₹":  "INR",
		"R$": "BRL",
		"A$": "AUD",
		"C$": "CAD",
	}

	var currency string
	var amount string

	for symbol, curr := range currencyMap {
		if strings.HasPrefix(priceStr, symbol) {
			currency = curr
			amount = strings.TrimPrefix(priceStr, symbol)
			break
		}
	}

	if currency == "" {
		re := regexp.MustCompile(`^([\d,.]+)`)
		match := re.FindString(priceStr)
		if match != "" {
			amount = match
			currency = "USD"
		}
	}

	amount = regexp.MustCompile(`[^\d.]`).ReplaceAllString(amount, "")
	price, _ := strconv.ParseFloat(amount, 64)

	return price, currency
}

func ExtractProductID(url string, platform string) string {
	switch strings.ToLower(platform) {
	case "amazon":
		re := regexp.MustCompile(`/dp/([A-Z0-9]{10})`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "ebay":
		re := regexp.MustCompile(`/itm/(\d+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "taobao":
		re := regexp.MustCompile(`id=(\d+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "jd":
		re := regexp.MustCompile(`/(\d+)\.html`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	case "tmall":
		re := regexp.MustCompile(`id=(\d+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return fmt.Sprintf("unknown-%d", len(url))
}
