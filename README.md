# 脱单AI - 智能脱单助手

一个基于AI技术的智能脱单助手，支持Web端和微信小程序，提供专业的情感咨询和脱单建议。

## 项目简介

脱单AI是一个现代化的情感咨询平台，结合了大语言模型（DeepSeek/Ollama）和向量数据库（Milvus）技术，为用户提供：

- 🤖 **AI情感咨询**：基于大语言模型，提供专业的情感建议，支持**流式输出**（打字机效果）
- 📚 **情感知识库**：丰富的情感知识内容，帮助用户提升社交技能
- 💳 **订阅服务**：解锁更多高级功能，获得专属脱单指导
- 🔐 **用户管理**：完整的用户认证和权限控制系统
- 📱 **微信小程序**：支持微信小程序端访问，随时随地获取情感咨询

## 技术栈

### 后端
- **框架**：Spring Boot 3.2
- **安全**：Spring Security + JWT
- **数据库**：MySQL 8.0
- **向量数据库**：Milvus 2.3
- **ORM**：Spring Data JPA
- **LLM**：DeepSeek API / Ollama
- **AI编排**：LangChain4j
- **流式响应**：SSE (Server-Sent Events)
- **构建工具**：Maven

### 前端（Web）
- **框架**：Vue 3 + Vite
- **状态管理**：Pinia
- **路由**：Vue Router
- **UI组件库**：Element Plus
- **HTTP客户端**：Axios
- **日期处理**：Day.js

### 前端（微信小程序）
- **框架**：微信小程序原生框架
- **UI设计**：自定义组件
- **流式响应**：wx.request stream API

## 项目结构

```
dating-ai/
├── backend/                    # 后端服务
│   ├── src/main/java/
│   │   └── com/datingai/
│   │       ├── controller/     # REST API控制器
│   │       │   ├── LLMController.java          # 普通LLM接口
│   │       │   └── StreamLLMController.java    # 流式LLM接口（新增）
│   │       ├── llm/           # LLM服务层
│   │       │   ├── DatingAdviceAssistant.java  # LangChain4j AI Service接口
│   │       │   ├── LLMService.java             # 普通LLM服务
│   │       │   └── StreamLLMService.java       # 流式LLM服务（新增）
│   │       ├── config/
│   │       │   └── LangChain4jConfig.java      # LangChain4j模型和AI Service配置
│   │       └── ...            # 其他模块
│   └── pom.xml
├── frontend/                   # Web前端
│   ├── src/
│   │   ├── views/
│   │   │   └── Chat.vue       # Web端聊天页面
│   │   └── ...
│   └── package.json
├── miniprogram/               # 微信小程序（新增）
│   ├── pages/
│   │   ├── index/             # 首页
│   │   └── chat/              # 聊天页面（流式输出）
│   ├── app.js
│   ├── app.json
│   └── app.wxss
└── README.md
```

## 系统要求

- **Java**：JDK 17+
- **Node.js**：18+
- **MySQL**：8.0+
- **Milvus**：2.3+（可选）
- **Maven**：3.8+
- **微信开发者工具**（用于小程序开发）

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

#### 1.4 安装Ollama（推荐本地开发）
```bash
# 安装Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 拉取LLM模型
ollama pull qwen3-vl:8b

# 启动Ollama服务
ollama serve
```

### 2. 后端部署

#### 2.1 配置数据库连接
编辑 `backend/src/main/resources/application.properties`：

```properties
# 数据库配置
spring.datasource.url=jdbc:mysql://localhost:3306/dating_ai_db?useSSL=false&serverTimezone=UTC
spring.datasource.username=your_mysql_username
spring.datasource.password=your_mysql_password

# LLM配置 - 使用Ollama（推荐本地开发）
llm.provider=langchain4j
llm.api.url=http://localhost:11434/api/chat
llm.base-url=http://localhost:11434
llm.model=qwen3-vl:4b
llm.temperature=0.7
llm.timeout=PT60S

# LLM配置 - 使用DeepSeek（服务器使用）
# llm.api.url=https://api.deepseek.com/v1/chat/completions
# llm.api.key=your_deepseek_api_key
# llm.model=deepseek-chat

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

### 2.4 LangChain4j 学习入口

本项目普通问答接口 `/api/llm/advice` 已经接入 LangChain4j，核心代码建议按下面顺序阅读：

1. `backend/pom.xml`
   - 引入 `langchain4j` 和 `langchain4j-ollama`。
   - `langchain4j` 提供 AI Service、Prompt 注解等核心能力。
   - `langchain4j-ollama` 提供 Ollama 模型适配器。

2. `backend/src/main/java/com/datingai/config/LangChain4jConfig.java`
   - 创建 `ChatModel`，这是 LangChain4j 对聊天模型的统一抽象。
   - 创建 `DatingAdviceAssistant`，LangChain4j 会为这个接口动态生成实现类。

3. `backend/src/main/java/com/datingai/llm/DatingAdviceAssistant.java`
   - 使用 `@SystemMessage` 定义 AI 角色、边界和长期规则。
   - 使用 `@UserMessage` 定义用户问题模板。

4. `backend/src/main/java/com/datingai/llm/LLMService.java`
   - Controller 只调用 `LLMService`，不用关心底层模型实现。
   - `llm.provider=langchain4j` 时走 LangChain4j。
   - `llm.provider=legacy` 时走原来的手写 HTTP 请求。

5. `backend/src/main/java/com/datingai/llm/LegacyDatingAdviceClient.java`
   - 保留原始 HTTP 调用方式，方便和 LangChain4j 对照学习。
   - 也给流式接口继续保留 `llm.api.url` 这类旧配置。

如果想临时切回原来的手写 HTTP 实现，把配置改成：

```properties
llm.provider=legacy
```

可以用下面命令测试：

```bash
curl -X POST http://localhost:8080/api/llm/advice \
  -H "Content-Type: application/json" \
  -d '{"question":"第一次约会不知道聊什么，怎么自然开场？"}'
```

### 3. Web前端部署

#### 3.1 安装依赖
```bash
cd dating-ai/frontend

# 安装npm依赖
npm install
```

#### 3.2 配置API地址（可选）
如果后端服务不在本地或端口不同，编辑相关配置文件。

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

### 4. 微信小程序开发

#### 4.1 打开小程序项目
1. 打开微信开发者工具
2. 选择"导入项目"
3. 选择 `dating-ai/miniprogram` 目录
4. 配置AppID（测试可使用测试号）

#### 4.2 配置API地址
编辑 `miniprogram/app.js`：

```javascript
App({
  globalData: {
    userInfo: null,
    baseUrl: 'http://localhost:8080/api'  // 后端API地址
  }
})
```

#### 4.3 预览和调试
- 在微信开发者工具中点击"预览"
- 使用微信扫码在手机上预览
- 调试模式下可查看控制台日志

#### 4.4 发布小程序
1. 在微信开发者工具中点击"上传"
2. 登录微信公众平台提交审核
3. 审核通过后发布

## API接口说明

### 普通问答接口
```
POST /api/llm/advice
Content-Type: application/json

{
  "question": "你的情感问题"
}

响应：
{
  "advice": "AI的回答内容"
}
```

### 流式问答接口（新增）
```
POST /api/llm/stream/advice
Content-Type: application/json
Accept: text/event-stream

{
  "question": "你的情感问题"
}

响应（SSE格式）：
data: {"content": "AI回复的一部分"}
data: {"content": "AI回复的下一部分"}
data: [DONE]
```

## 流式输出实现原理

### 后端实现
- 使用 `SseEmitter` 创建Server-Sent Events连接
- 流式读取LLM服务的响应
- 逐块将内容推送给客户端

### 前端实现（Web）
- 使用 `EventSource` 或 `axios` 的流式响应功能
- 实时更新UI，实现打字机效果

### 小程序实现
- 使用 `wx.request` 的 `responseType: 'stream'`
- 通过 `getReader()` 读取数据流
- 逐块解析SSE格式并更新界面

## 环境变量配置

### 后端环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SPRING_DATASOURCE_URL` | MySQL连接URL | `jdbc:mysql://localhost:3306/dating_ai_db` |
| `SPRING_DATASOURCE_USERNAME` | MySQL用户名 | `root` |
| `SPRING_DATASOURCE_PASSWORD` | MySQL密码 | - |
| `LLM_API_URL` | LLM API地址 | `http://localhost:11434/api/chat` |
| `LLM_MODEL` | LLM模型名称 | `qwen3-vl:8b` |
| `LLM_API_KEY` | DeepSeek API密钥（使用DeepSeek时） | - |
| `JWT_SECRET` | JWT密钥 | - |

### 前端环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_API_BASE_URL` | 后端API地址 | `http://localhost:8080/api` |

### 小程序配置
| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `baseUrl` | 后端API地址 | `http://localhost:8080/api` |

## 测试账号

系统默认初始化两个测试账号：

- **普通用户**：用户名 `user`，密码 `password`
- **管理员用户**：用户名 `admin`，密码 `password`

## 常见问题

### 1. 数据库连接失败
- 检查MySQL服务是否启动
- 检查数据库配置是否正确
- 检查防火墙设置

### 2. Ollama连接失败
- 检查Ollama服务是否启动（`ollama serve`）
- 检查模型是否已拉取（`ollama list`）

### 3. 流式输出不工作
- 检查后端是否正确配置了流式API
- 检查网络连接
- 检查浏览器控制台是否有错误

### 4. 小程序无法访问后端API
- 检查小程序的服务器域名配置
- 检查后端CORS配置
- 使用测试号时需开启"不校验域名"模式

### 5. DeepSeek API调用失败
- 检查API密钥是否正确
- 检查网络连接
- 检查API调用额度

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
```

### 小程序更新
```bash
# 拉取最新代码
git pull origin main

# 在微信开发者工具中上传新版本
```

## 性能优化

### 后端优化
- 启用数据库连接池
- 配置Redis缓存
- 使用CDN加速静态资源
- 启用Gzip压缩
- SSE连接复用

### 前端优化
- 启用代码分割
- 使用懒加载
- 优化图片资源
- 流式响应的UI优化

### 小程序优化
- 图片懒加载
- 分包加载
- 缓存策略优化

## 安全建议

1. **修改默认密码**：部署后立即修改默认测试账号密码
2. **使用HTTPS**：生产环境必须使用HTTPS
3. **配置域名白名单**：小程序需配置服务器域名
4. **定期更新依赖**：及时更新第三方依赖库
5. **配置防火墙**：只开放必要的端口
6. **定期备份**：定期备份数据库和重要数据
7. **监控日志**：配置日志监控和告警

## 技术支持

如有问题，请通过以下方式联系：

- 邮箱：contact@datingai.com
- 电话：123-456-7890
- GitHub Issues：[项目地址]/issues

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

---

**注意**：部署前请确保已阅读并理解所有配置项，根据实际环境进行相应调整。
