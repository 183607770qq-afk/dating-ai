package api

import (
	"hotdeal-tracker/internal/analyzer"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler API 请求处理器
// 负责处理所有 HTTP 请求，调用相应的业务逻辑
type Handler struct {
	db       *database.Database      // 数据库连接
	analyzer *analyzer.Analyzer      // 数据分析器
}

// NewHandler 创建 API 处理器实例
// 参数：
//   db: 数据库实例
// 返回：
//   *Handler: API 处理器实例
func NewHandler(db *database.Database) *Handler {
	return &Handler{
		db:       db,
		analyzer: analyzer.NewAnalyzer(db), // 创建数据分析器
	}
}

// Response API 统一响应格式
type Response struct {
	Success bool        `json:"success"`  // 请求是否成功
	Message string      `json:"message,omitempty"` // 错误信息（失败时）
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Meta    *Meta       `json:"meta,omitempty"`    // 分页元数据
}

// Meta 分页元数据
type Meta struct {
	Page       int   `json:"page"`        // 当前页码
	Limit      int   `json:"limit"`       // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
}

// Success 返回成功响应
// 参数：
//   c: gin上下文
//   data: 响应数据
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta 返回带分页信息的成功响应
// 参数：
//   c: gin上下文
//   data: 响应数据
//   meta: 分页元数据
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error 返回错误响应
// 参数：
//   c: gin上下文
//   status: HTTP状态码
//   message: 错误信息
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

// GetHotProducts 获取热销产品列表
// Query参数：
//   limit: 返回数量（默认20）
//   platform: 平台代码（可选，如 amazon, ebay）
func (h *Handler) GetHotProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	platform := c.Query("platform")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// 调用分析器获取热销产品
	hotProducts, err := h.analyzer.GetHotProducts(limit, platform)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get hot products: "+err.Error())
		return
	}

	Success(c, hotProducts)
}

// GetProducts 获取产品列表（支持多条件筛选）
// Query参数：
//   page: 页码（默认1）
//   limit: 每页数量（默认20，最大100）
//   platform: 平台代码（可选）
//   category: 分类名称（可选）
//   keyword: 搜索关键词（可选）
func (h *Handler) GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	platform := c.Query("platform")
	category := c.Query("category")
	keyword := c.Query("keyword")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// 参数校验
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var products []models.Product
	var total int64
	var err error

	// 根据不同条件查询产品
	if keyword != "" {
		// 优先按关键词搜索
		products, err = h.db.SearchProducts(keyword, limit, offset)
		if err == nil {
			total = int64(len(products))
		}
	} else if category != "" {
		// 按分类查询
		products, err = h.db.GetProductsByCategory(category, limit, offset)
		if err == nil {
			total = int64(len(products))
		}
	} else if platform != "" {
		// 按平台查询
		products, err = h.db.GetProductsByPlatform(platform, limit, offset)
		if err == nil {
			total = int64(len(products))
		}
	} else {
		// 获取所有产品
		products, total, err = h.db.GetAllProducts(limit, offset)
	}

	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get products: "+err.Error())
		return
	}

	SuccessWithMeta(c, products, &Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)), // 计算总页数
	})
}

// GetProduct 获取单个产品详情
// Path参数：
//   id: 产品ID
func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	analysis, err := h.analyzer.AnalyzeProduct(uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, "Product not found")
		return
	}

	Success(c, analysis)
}

// GetCategories 获取分类列表
// Query参数：
//   platform: 平台代码（可选，如 amazon, ebay）
func (h *Handler) GetCategories(c *gin.Context) {
	platform := c.Query("platform")

	categories, err := h.db.GetCategories(platform)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get categories: "+err.Error())
		return
	}

	Success(c, categories)
}

// GetPlatforms 获取所有启用的平台列表
func (h *Handler) GetPlatforms(c *gin.Context) {
	platforms, err := h.db.GetPlatforms()
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get platforms: "+err.Error())
		return
	}

	Success(c, platforms)
}

// GetCategoryStats 获取分类统计数据
// 返回各分类的产品数量、平均价格等统计信息
func (h *Handler) GetCategoryStats(c *gin.Context) {
	stats, err := h.analyzer.GetCategoryStats()
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get category stats: "+err.Error())
		return
	}

	Success(c, stats)
}

// GetPlatformStats 获取平台统计数据
// 返回各平台的产品数量、平均价格等统计信息
func (h *Handler) GetPlatformStats(c *gin.Context) {
	stats, err := h.analyzer.GetPlatformStats()
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get platform stats: "+err.Error())
		return
	}

	Success(c, stats)
}

// GetPriceTrends 获取产品价格趋势
// Path参数：
//   id: 产品ID
// Query参数：
//   days: 查询天数（默认30天）
func (h *Handler) GetPriceTrends(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)

	trends, err := h.analyzer.GetPriceTrends(uint(id), days)
	if err != nil {
		Error(c, http.StatusNotFound, "No price history found")
		return
	}

	Success(c, trends)
}

// GetMarketInsights 获取市场洞察数据
// Query参数：
//   platform: 平台代码（可选，如 amazon, ebay）
// 返回市场整体趋势、价格分布、热销品类等分析数据
func (h *Handler) GetMarketInsights(c *gin.Context) {
	platform := c.Query("platform")

	insights, err := h.analyzer.GetMarketInsights(platform)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get market insights: "+err.Error())
		return
	}

	Success(c, insights)
}

// GetTrendingKeywords 获取热门搜索关键词
// Query参数：
//   days: 查询天数（默认7天）
//   limit: 返回数量（默认10）
func (h *Handler) GetTrendingKeywords(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	limitStr := c.DefaultQuery("limit", "10")

	days, _ := strconv.Atoi(daysStr)
	limit, _ := strconv.Atoi(limitStr)

	keywords, err := h.analyzer.GetTrendingKeywords(days, limit)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to get trending keywords: "+err.Error())
		return
	}

	Success(c, keywords)
}

func (h *Handler) SearchProducts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		Error(c, http.StatusBadRequest, "Search keyword is required")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	products, err := h.db.SearchProducts(keyword, limit, offset)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Failed to search products: "+err.Error())
		return
	}

	Success(c, products)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	Success(c, map[string]string{
		"status":  "healthy",
		"service": "hotdeal-tracker",
	})
}
