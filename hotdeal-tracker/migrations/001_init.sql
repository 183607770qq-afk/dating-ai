-- Create hotdeal_tracker database
CREATE DATABASE hotdeal_tracker;

-- Connect to the database
\c hotdeal_tracker;

-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Platforms table
CREATE TABLE IF NOT EXISTS platforms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    base_url TEXT,
    logo VARCHAR(255),
    country VARCHAR(10),
    language VARCHAR(20),
    currency VARCHAR(10),
    is_active BOOLEAN DEFAULT TRUE,
    hot_url TEXT,
    category_url TEXT,
    product_url TEXT,
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    parent_id INT REFERENCES categories(id),
    platform VARCHAR(50),
    icon VARCHAR(255),
    description TEXT,
    product_count INT DEFAULT 0,
    hot_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    platform_id VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    image_url TEXT,
    product_url TEXT NOT NULL,
    price DECIMAL(10,2),
    original_price DECIMAL(10,2),
    currency VARCHAR(10) DEFAULT 'USD',
    sales_count INT DEFAULT 0,
    review_count INT DEFAULT 0,
    rating DECIMAL(3,2),
    category VARCHAR(100),
    tags TEXT,
    badge VARCHAR(50),
    shop_name VARCHAR(200),
    shop_id VARCHAR(100),
    is_hot BOOLEAN DEFAULT FALSE,
    trending_score DECIMAL(10,2) DEFAULT 0,
    crawled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_platform_product UNIQUE (platform_id, platform)
);

-- Price history table
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    original_price DECIMAL(10,2),
    sales_count INT,
    crawled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crawl tasks table
CREATE TABLE IF NOT EXISTS crawl_tasks (
    id SERIAL PRIMARY KEY,
    platform_id INT REFERENCES platforms(id),
    category_id INT REFERENCES categories(id),
    keyword VARCHAR(200),
    url TEXT NOT NULL,
    type VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    priority INT DEFAULT 0,
    error_msg TEXT,
    retry_count INT DEFAULT 0,
    product_count INT DEFAULT 0,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analytics table
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    views INT DEFAULT 0,
    unique_views INT DEFAULT 0,
    cart_adds INT DEFAULT 0,
    purchases INT DEFAULT 0,
    revenue DECIMAL(12,2) DEFAULT 0,
    conversion_rate DECIMAL(5,3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_products_platform ON products(platform);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_is_hot ON products(is_hot);
CREATE INDEX IF NOT EXISTS idx_products_trending ON products(trending_score DESC);
CREATE INDEX IF NOT EXISTS idx_price_history_product ON price_history(product_id);
CREATE INDEX IF NOT EXISTS idx_price_history_date ON price_history(crawled_at);
CREATE INDEX IF NOT EXISTS idx_analytics_product_date ON analytics(product_id, date);
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_status ON crawl_tasks(status);
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_priority ON crawl_tasks(priority DESC);
