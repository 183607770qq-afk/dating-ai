# HotDeal Tracker - 全网热销品追踪与分析系统

一个使用 Go 语言开发的高性能全网热销品爬取和数据分析系统，支持国内外多个电商平台，提供 RESTful API 和可视化 Dashboard。

## 功能特性

### 🔍 数据爬取
- 支持 10+ 主流电商平台（Amazon、eBay、淘宝、京东、拼多多、天猫、Shopee、Lazada、AliExpress、Wish 等）
- 支持关键词搜索和分类爬取
- 自动处理反爬机制（代理池、User-Agent 轮换、请求延迟）
- Chrome DevTools Protocol 支持（动态页面爬取）
- 定时任务自动更新数据

### 📊 数据分析
- 热销商品排行和趋势分析
- 价格走势追踪和历史对比
- 平台分布和类目统计
- 市场洞察和关键词趋势
- 个性化购买建议

### 🚀 API 接口
- RESTful API 设计
- 完整的 CRUD 操作
- 分页、筛选、搜索支持
- CORS 跨域支持

### 📈 可视化 Dashboard
- 实时数据展示
- 图表统计分析
- 灵活的数据筛选
- 响应式设计

## 技术栈

### 后端
- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM 库
- **PostgreSQL** - 主数据库
- **Redis** - 缓存（可选）
- **goquery** - HTML 解析
- **chromedp** - Chrome DevTools Protocol

### 前端
- HTML5 + CSS3 + JavaScript
- Chart.js - 图表库
- Axios - HTTP 客户端

## 项目结构

```
hotdeal-tracker/
├── cmd/
│   ├── api/           # API 服务入口
│   ├── scraper/       # 爬虫服务入口
│   └── dashboard/     # Dashboard 服务入口
├── internal/
│   ├── config/        # 配置管理
│   ├── crawler/       # 爬虫核心逻辑
│   ├── database/      # 数据库层
│   ├── models/        # 数据模型
│   ├── analyzer/      # 数据分析
│   └── api/           # API 处理器
├── pkg/
│   ├── httpclient/    # HTTP 客户端
│   └── parser/        # HTML 解析器
├── web/
│   └── public/        # 前端资源
├── migrations/        # 数据库迁移
├── scripts/          # 工具脚本
└── config.yaml       # 配置文件
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 14+
- Redis 7+（可选）

### 1. 克隆项目

```bash
cd hotdeal-tracker
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE hotdeal_tracker;
```

修改 `config.yaml` 配置文件：

```yaml
database:
  host: "localhost"
  port: 5432
  user: "root"
  password: "root"
  name: "hotdeal_tracker"
  sslmode: "disable"
```

### 4. 运行数据库迁移

```bash
psql -U postgres -d hotdeal_tracker -f migrations/001_init.sql
```

### 5. 启动服务

**启动 API 服务：**
```bash
go run cmd/api/main.go -config ./config.yaml
```

**启动爬虫服务：**
```bash
go run cmd/scraper/main.go -config ./config.yaml
```

**一次性爬取：**
```bash
go run cmd/scraper/main.go -once -keyword "bestseller" -platform "amazon"
```

### 6. 访问 Dashboard

打开浏览器访问：`http://localhost:8080/web/public/index.html`

## API 文档

### 产品接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/products` | 获取产品列表 |
| GET | `/api/v1/products/hot` | 获取热销产品 |
| GET | `/api/v1/products/search` | 搜索产品 |
| GET | `/api/v1/products/:id` | 获取产品详情 |
| GET | `/api/v1/products/:id/trends` | 获取价格趋势 |

### 类目接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/categories` | 获取类目列表 |
| GET | `/api/v1/categories/stats` | 获取类目统计 |

### 平台接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/platforms` | 获取平台列表 |
| GET | `/api/v1/platforms/stats` | 获取平台统计 |

### 数据分析接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/insights/market` | 获取市场洞察 |
| GET | `/api/v1/insights/keywords` | 获取热门关键词 |

## 配置说明

详细配置项说明：

```yaml
server:
  port: "8080"           # 服务端口
  mode: "debug"          # 运行模式 (debug/release)
  read_timeout: 60       # 读取超时(秒)
  write_timeout: 60      # 写入超时(秒)

database:
  host: "localhost"      # 数据库主机
  port: 5432             # 数据库端口
  user: "postgres"       # 数据库用户
  password: "postgres"   # 数据库密码
  name: "hotdeal_tracker" # 数据库名
  sslmode: "disable"     # SSL模式

redis:
  host: "localhost"      # Redis主机
  port: 6379             # Redis端口
  password: ""           # Redis密码
  db: 0                  # 数据库编号

crawler:
  user_agent: "..."      # User-Agent字符串
  timeout: 30           # 请求超时(秒)
  max_depth: 3           # 最大爬取深度
  delay: 1000            # 请求间隔(毫秒)
  concurrent: 5          # 并发数
  retry_times: 3         # 重试次数
  enable_headless: true  # 启用无头浏览器
  proxy_pool: []         # 代理池
  platforms:             # 启用的平台
    amazon: true
    ebay: true
    # ...
```

## 数据模型

### Product（产品）

```go
type Product struct {
    ID            uint      // 主键
    PlatformID    string    // 平台产品ID
    Platform      string    // 平台名称
    Title         string    // 产品标题
    Description   string    // 产品描述
    ImageURL      string    // 图片URL
    ProductURL    string    // 产品URL
    Price         float64   // 当前价格
    OriginalPrice float64   // 原价
    Currency      string    // 货币单位
    SalesCount    int       // 销量
    ReviewCount   int       // 评论数
    Rating        float64   // 评分
    Category      string    // 类别
    Tags          string    // 标签
    Badge         string    // 徽章
    IsHot         bool      // 是否热销
    TrendingScore float64   // 趋势分数
    CrawledAt     time.Time // 爬取时间
}
```

## 开发指南

### 添加新平台

1. 在 `pkg/parser/parser.go` 中实现 `ProductParser` 接口
2. 在 `GetParser()` 函数中添加平台解析器
3. 在数据库中添加入口 URL
4. 在 `config.yaml` 中启用该平台

### 运行测试

```bash
go test ./...
```

### 构建发布

```bash
# 构建 API 服务
go build -o bin/api cmd/api/main.go

# 构建爬虫服务
go build -o bin/scraper cmd/scraper/main.go
```

## License

MIT License
