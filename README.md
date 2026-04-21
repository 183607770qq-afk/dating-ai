# 脱单AI - 智能脱单助手

一个基于AI技术的智能脱单助手小程序，提供专业的情感咨询和脱单建议。

## 项目简介

脱单AI是一个现代化的情感咨询平台，结合了大语言模型（DeepSeek）和向量数据库（Milvus）技术，为用户提供：

- 🤖 **AI情感咨询**：基于DeepSeek大模型，提供专业的情感建议
- 📚 **情感知识库**：丰富的情感知识内容，帮助用户提升社交技能
- 💳 **订阅服务**：解锁更多高级功能，获得专属脱单指导
- 🔐 **用户管理**：完整的用户认证和权限控制系统

## 技术栈

### 后端
- **框架**：Spring Boot 3.2
- **安全**：Spring Security + JWT
- **数据库**：MySQL 8.0
- **向量数据库**：Milvus 2.3
- **ORM**：Spring Data JPA
- **LLM**：DeepSeek API
- **构建工具**：Maven

### 前端
- **框架**：Vue 3 + Vite
- **状态管理**：Pinia
- **路由**：Vue Router
- **UI组件库**：Element Plus
- **HTTP客户端**：Axios
- **日期处理**：Day.js

## 系统要求

- **Java**：JDK 17+
- **Node.js**：18+
- **MySQL**：8.0+
- **Milvus**：2.3+
- **Maven**：3.8+

## 部署步骤

### 1. 环境准备

#### 1.1 安装Java环境
```bash
# 检查Java版本
java -version

# 如果没有安装，请安装JDK 17
# macOS
brew install openjdk@17

# Ubuntu/Debian
sudo apt update
sudo apt install openjdk-17-jdk
```

#### 1.2 安装Node.js环境
```bash
# 检查Node.js版本
node -v

# 如果没有安装，请安装Node.js 18+
# 使用nvm安装
nvm install 18
nvm use 18
```

#### 1.3 安装MySQL
```bash
# macOS
brew install mysql
brew services start mysql

# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql

# 创建数据库
mysql -u root -p
CREATE DATABASE dating_ai_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 1.4 安装Milvus（可选，用于向量检索）
```bash
# 使用Docker安装Milvus
docker-compose -f milvus-docker-compose.yml up -d

# 或者参考官方文档：https://milvus.io/docs/install_standalone-docker.md
```

### 2. 后端部署

#### 2.1 配置数据库连接
编辑 `backend/src/main/resources/application.properties`：

```properties
# 数据库配置
spring.datasource.url=jdbc:mysql://localhost:3306/dating_ai_db?useSSL=false&serverTimezone=UTC
spring.datasource.username=your_mysql_username
spring.datasource.password=your_mysql_password

# Milvus配置（可选）
milvus.host=localhost
milvus.port=19530

# DeepSeek API配置
llm.api.url=https://api.deepseek.com/v1/chat/completions
llm.api.key=your_deepseek_api_key

# JWT密钥配置
jwt.secret=your_jwt_secret_key
```

#### 2.2 构建后端项目
```bash
cd dating-ai/backend

# 使用Maven构建
./mvnw clean package -DskipTests

# 或者使用Maven Wrapper（Windows）
mvnw.cmd clean package -DskipTests
```

#### 2.3 运行后端服务
```bash
# 方式1：直接运行
./mvnw spring-boot:run

# 方式2：运行打包后的jar文件
java -jar target/dating-ai-backend-0.0.1-SNAPSHOT.jar

# 后台运行（Linux/macOS）
nohup java -jar target/dating-ai-backend-0.0.1-SNAPSHOT.jar > app.log 2>&1 &
```

后端服务默认运行在 `http://localhost:8080`

### 3. 前端部署

#### 3.1 安装依赖
```bash
cd dating-ai/frontend

# 安装npm依赖
npm install
```

#### 3.2 配置API地址（可选）
如果后端服务不在本地或端口不同，编辑 `frontend/src/main.js`：

```javascript
// 配置axios默认baseURL
axios.defaults.baseURL = 'http://your_backend_host:your_backend_port/api'
```

#### 3.3 开发环境运行
```bash
# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

#### 3.4 生产环境构建
```bash
# 构建生产版本
npm run build

# 构建后的文件在 dist/ 目录中
```

#### 3.5 部署到Web服务器
```bash
# 使用Nginx部署示例
# 1. 将 dist/ 目录中的文件复制到Nginx的web目录
cp -r dist/* /var/www/html/dating-ai/

# 2. 配置Nginx（/etc/nginx/sites-available/dating-ai）
server {
    listen 80;
    server_name your_domain.com;
    root /var/www/html/dating-ai;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# 3. 启用配置并重启Nginx
sudo ln -s /etc/nginx/sites-available/dating-ai /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 4. Docker部署（推荐）

#### 4.1 创建Dockerfile
创建 `backend/Dockerfile`：

```dockerfile
FROM openjdk:17-jdk-slim

WORKDIR /app

COPY target/dating-ai-backend-0.0.1-SNAPSHOT.jar app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
```

创建 `frontend/Dockerfile`：

```dockerfile
FROM node:18-alpine as build

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

#### 4.2 创建docker-compose.yml
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: dating_ai_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  milvus:
    image: milvusdb/milvus:v2.3.4
    ports:
      - "19530:19530"
    volumes:
      - milvus_data:/var/lib/milvus

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - SPRING_DATASOURCE_URL=jdbc:mysql://mysql:3306/dating_ai_db
      - SPRING_DATASOURCE_USERNAME=root
      - SPRING_DATASOURCE_PASSWORD=root_password
      - MILVUS_HOST=milvus
      - LLM_API_KEY=your_deepseek_api_key
    depends_on:
      - mysql
      - milvus

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  mysql_data:
  milvus_data:
```

#### 4.3 使用Docker Compose部署
```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down
```

### 5. 云服务部署

#### 5.1 阿里云ECS部署
```bash
# 1. 购买ECS实例（推荐2核4G以上配置）
# 2. 配置安全组，开放80、8080端口
# 3. 连接服务器并安装环境
ssh root@your_server_ip

# 安装Docker和Docker Compose
curl -fsSL https://get.docker.com | sh
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 上传项目文件并部署
docker-compose up -d
```

#### 5.2 腾讯云部署
类似阿里云，使用腾讯云服务器部署。

### 6. 配置SSL证书（HTTPS）

#### 6.1 使用Let's Encrypt免费证书
```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your_domain.com

# 自动续期
sudo certbot renew --dry-run
```

#### 6.2 手动配置SSL
```bash
# 在Nginx配置中添加SSL配置
server {
    listen 443 ssl;
    server_name your_domain.com;
    
    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    
    # 其他配置...
}

# HTTP重定向到HTTPS
server {
    listen 80;
    server_name your_domain.com;
    return 301 https://$server_name$request_uri;
}
```

## 环境变量配置

### 后端环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SPRING_DATASOURCE_URL` | MySQL连接URL | `jdbc:mysql://localhost:3306/dating_ai_db` |
| `SPRING_DATASOURCE_USERNAME` | MySQL用户名 | `root` |
| `SPRING_DATASOURCE_PASSWORD` | MySQL密码 | - |
| `MILVUS_HOST` | Milvus主机地址 | `localhost` |
| `MILVUS_PORT` | Milvus端口 | `19530` |
| `LLM_API_URL` | DeepSeek API地址 | `https://api.deepseek.com/v1/chat/completions` |
| `LLM_API_KEY` | DeepSeek API密钥 | - |
| `JWT_SECRET` | JWT密钥 | - |

### 前端环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_API_BASE_URL` | 后端API地址 | `http://localhost:8080/api` |

## 测试账号

系统默认初始化两个测试账号：

- **普通用户**：用户名 `user`，密码 `password`
- **管理员用户**：用户名 `admin`，密码 `password`

## 常见问题

### 1. 数据库连接失败
- 检查MySQL服务是否启动
- 检查数据库配置是否正确
- 检查防火墙设置

### 2. Milvus连接失败
- 检查Milvus服务是否启动
- 检查Milvus配置是否正确
- 如不需要向量检索，可暂时禁用Milvus相关功能

### 3. DeepSeek API调用失败
- 检查API密钥是否正确
- 检查网络连接
- 检查API调用额度

### 4. 前端无法访问后端API
- 检查CORS配置
- 检查API地址配置
- 检查防火墙设置

## 维护与更新

### 后端更新
```bash
# 拉取最新代码
git pull origin main

# 重新构建
cd backend
./mvnw clean package -DskipTests

# 重启服务
# 如果使用Docker
docker-compose up -d --build backend

# 如果直接运行
# 先停止旧服务，然后重新启动
```

### 前端更新
```bash
# 拉取最新代码
git pull origin main

# 重新构建
cd frontend
npm install
npm run build

# 部署到Web服务器
# 如果使用Docker
docker-compose up -d --build frontend

# 如果手动部署，复制dist目录到Web服务器
```

## 性能优化

### 后端优化
- 启用数据库连接池
- 配置Redis缓存
- 使用CDN加速静态资源
- 启用Gzip压缩

### 前端优化
- 启用代码分割
- 使用懒加载
- 优化图片资源
- 启用Service Worker

## 安全建议

1. **修改默认密码**：部署后立即修改默认测试账号密码
2. **使用HTTPS**：生产环境必须使用HTTPS
3. **定期更新依赖**：及时更新第三方依赖库
4. **配置防火墙**：只开放必要的端口
5. **定期备份**：定期备份数据库和重要数据
6. **监控日志**：配置日志监控和告警

## 技术支持

如有问题，请通过以下方式联系：

- 邮箱：contact@datingai.com
- 电话：123-456-7890
- GitHub Issues：[项目地址]/issues

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

---

**注意**：部署前请确保已阅读并理解所有配置项，根据实际环境进行相应调整。