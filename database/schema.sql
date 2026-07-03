-- 脱单AI数据库DDL语句
-- MySQL 8.0+

-- 创建数据库
CREATE DATABASE IF NOT EXISTS dating_ai_db 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE dating_ai_db;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密）',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    role VARCHAR(20) NOT NULL DEFAULT 'USER' COMMENT '角色：USER, SUBSCRIBER, ADMIN',
    is_subscribed BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否订阅',
    subscription_end_date DATETIME COMMENT '订阅结束时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_subscription (is_subscribed, subscription_end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 知识表
CREATE TABLE IF NOT EXISTS knowledge (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '知识ID',
    title VARCHAR(200) NOT NULL COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    category VARCHAR(50) NOT NULL COMMENT '分类：约会技巧, 沟通技巧, 自我提升, 关系维护',
    is_published BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否发布',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category),
    INDEX idx_published (is_published),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识表';

-- 订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订阅ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    subscription_type VARCHAR(20) NOT NULL COMMENT '订阅类型：monthly, quarterly, annual',
    start_date DATETIME NOT NULL COMMENT '开始时间',
    end_date DATETIME NOT NULL COMMENT '结束时间',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' COMMENT '状态：ACTIVE, EXPIRED, CANCELLED',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅表';

-- 聊天记录表（可选）
CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role VARCHAR(20) NOT NULL COMMENT '角色：user, assistant',
    content TEXT NOT NULL COMMENT '消息内容',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天记录表';

-- 对话摘要表（长期记忆）
CREATE TABLE IF NOT EXISTS conversation_summaries (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '摘要ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    summary_text TEXT NOT NULL COMMENT '摘要文本',
    message_count_at_summary BIGINT NOT NULL DEFAULT 0 COMMENT '生成摘要时的消息总数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话摘要表（长期记忆，每个用户只保留最新一条）';

-- 知识浏览记录表（可选）
CREATE TABLE IF NOT EXISTS knowledge_views (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '浏览ID',
    user_id BIGINT COMMENT '用户ID（可为空，支持匿名浏览）',
    knowledge_id BIGINT NOT NULL COMMENT '知识ID',
    view_count INT NOT NULL DEFAULT 1 COMMENT '浏览次数',
    last_viewed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后浏览时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_knowledge_id (knowledge_id),
    INDEX idx_last_viewed_at (last_viewed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识浏览记录表';

-- 用户反馈表（可选）
CREATE TABLE IF NOT EXISTS user_feedback (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '反馈ID',
    user_id BIGINT COMMENT '用户ID（可为空，支持匿名反馈）',
    feedback_type VARCHAR(20) NOT NULL COMMENT '反馈类型：suggestion, bug, complaint, other',
    content TEXT NOT NULL COMMENT '反馈内容',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' COMMENT '状态：PENDING, PROCESSING, RESOLVED, CLOSED',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户反馈表';

-- 系统配置表（可选）
CREATE TABLE IF NOT EXISTS system_config (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    description VARCHAR(255) COMMENT '配置描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 插入默认管理员用户
INSERT INTO users (username, password, email, role, is_subscribed) 
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'admin@datingai.com', 'ADMIN', TRUE)
ON DUPLICATE KEY UPDATE username=username;

-- 插入默认测试用户
INSERT INTO users (username, password, email, role, is_subscribed) 
VALUES ('user', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'user@datingai.com', 'USER', FALSE)
ON DUPLICATE KEY UPDATE username=username;

-- 插入示例知识数据
INSERT INTO knowledge (title, content, category, is_published) VALUES
('如何提高自信心', '自信心是脱单的基础。首先，要接受自己的不完美，每个人都有自己的优点和缺点。其次，培养自己的兴趣爱好，让自己变得更加有趣。最后，多参与社交活动，积累成功的社交经验。', '自我提升', TRUE),
('第一次约会技巧', '第一次约会时，选择一个轻松的环境，如咖啡厅或公园。穿着得体，保持整洁。准备一些话题，避免尴尬的沉默。最重要的是，做自己，不要刻意伪装。', '约会技巧', TRUE),
('如何与异性有效沟通', '有效沟通的关键是倾听。不要只关注自己的表达，也要认真听取对方的想法。使用开放式问题，鼓励对方分享更多。保持眼神交流，展现你的真诚和兴趣。', '沟通技巧', TRUE),
('如何维护长期关系', '长期关系的维护需要双方的努力。保持良好的沟通，及时解决矛盾。给彼此足够的空间和信任。共同创造美好的回忆，定期安排约会时间。', '关系维护', TRUE)
ON DUPLICATE KEY UPDATE title=title;

-- 插入系统配置
INSERT INTO system_config (config_key, config_value, description) VALUES
('site_title', '脱单AI', '网站标题'),
('site_description', '智能脱单助手，让恋爱更简单', '网站描述'),
('max_free_messages', '10', '免费用户每日最大消息数'),
('subscription_price_monthly', '29', '月度订阅价格'),
('subscription_price_quarterly', '79', '季度订阅价格'),
('subscription_price_annual', '299', '年度订阅价格')
ON DUPLICATE KEY UPDATE config_key=config_key;

-- 创建视图：用户订阅状态视图
CREATE OR REPLACE VIEW user_subscription_status AS
SELECT 
    u.id AS user_id,
    u.username,
    u.email,
    u.is_subscribed,
    u.subscription_end_date,
    s.subscription_type,
    s.status AS subscription_status,
    CASE 
        WHEN u.subscription_end_date IS NULL THEN '未订阅'
        WHEN u.subscription_end_date < NOW() THEN '已过期'
        ELSE '订阅中'
    END AS subscription_status_text
FROM users u
LEFT JOIN subscriptions s ON u.id = s.user_id AND s.status = 'ACTIVE';

-- 创建存储过程：检查并更新过期订阅
DELIMITER //
CREATE PROCEDURE check_expired_subscriptions()
BEGIN
    UPDATE subscriptions 
    SET status = 'EXPIRED'
    WHERE end_date < NOW() AND status = 'ACTIVE';
    
    UPDATE users 
    SET is_subscribed = FALSE
    WHERE subscription_end_date < NOW() AND is_subscribed = TRUE;
END //
DELIMITER ;

-- 创建事件：每天自动检查过期订阅
CREATE EVENT IF NOT EXISTS check_subscriptions_event
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_TIMESTAMP
DO CALL check_expired_subscriptions();

-- 启用事件调度器
SET GLOBAL event_scheduler = ON;

-- 创建触发器：用户订阅时更新用户状态
DELIMITER //
CREATE TRIGGER update_user_subscription
AFTER INSERT ON subscriptions
FOR EACH ROW
BEGIN
    IF NEW.status = 'ACTIVE' THEN
        UPDATE users 
        SET is_subscribed = TRUE, 
            subscription_end_date = NEW.end_date
        WHERE id = NEW.user_id;
    END IF;
END //
DELIMITER ;

-- 授权（根据实际情况调整）
-- GRANT ALL PRIVILEGES ON dating_ai_db.* TO 'dating_user'@'localhost' IDENTIFIED BY 'your_password';
-- FLUSH PRIVILEGES;

-- 显示所有表
SHOW TABLES;

-- 显示表结构
-- DESC users;
-- DESC knowledge;
-- DESC subscriptions;
-- DESC chat_messages;
-- DESC knowledge_views;
-- DESC user_feedback;
-- DESC system_config;

-- 查询示例
-- 查询所有用户
-- SELECT * FROM users;

-- 查询所有已发布的知识
-- SELECT * FROM knowledge WHERE is_published = TRUE ORDER BY created_at DESC;

-- 查询用户订阅状态
-- SELECT * FROM user_subscription_status;

-- 查询活跃订阅
-- SELECT * FROM subscriptions WHERE status = 'ACTIVE';

-- 查询最近的知识浏览记录
-- SELECT * FROM knowledge_views ORDER BY last_viewed_at DESC LIMIT 10;
