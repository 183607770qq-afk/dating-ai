<template>
  <div class="home">
    <!-- 头部导航 -->
    <header class="header">
      <div class="container">
        <div class="logo">
          <h1>脱单AI</h1>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/knowledge" class="nav-link">情感知识</router-link>
          <router-link to="/chat" class="nav-link">AI咨询</router-link>
          <router-link to="/subscription" class="nav-link">订阅服务</router-link>
          <div v-if="!userStore.getIsLoggedIn" class="nav-auth">
            <router-link to="/login" class="nav-link">登录</router-link>
            <router-link to="/register" class="nav-link primary">注册</router-link>
          </div>
          <div v-else class="nav-user">
            <router-link to="/profile" class="nav-link">个人中心</router-link>
            <button @click="handleLogout" class="nav-link">退出</button>
          </div>
        </nav>
      </div>
    </header>

    <!-- 英雄区域 -->
    <section class="hero">
      <div class="container">
        <div class="hero-content">
          <h2>智能脱单助手</h2>
          <p>基于AI技术，为你提供专业的情感咨询和脱单建议</p>
          <div class="hero-buttons">
            <router-link to="/chat" class="btn primary">立即咨询</router-link>
            <router-link to="/subscription" class="btn secondary">订阅服务</router-link>
          </div>
        </div>
        <div class="hero-image">
          <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=modern%20dating%20app%20interface%20with%20AI%20assistant&image_size=landscape_16_9" alt="脱单AI" />
        </div>
      </div>
    </section>

    <!-- 功能介绍 -->
    <section class="features">
      <div class="container">
        <h3>核心功能</h3>
        <div class="features-grid">
          <router-link to="/chat" class="feature-item">
            <div class="feature-icon">
              <i class="el-icon-message"></i>
            </div>
            <h4>AI情感咨询</h4>
            <p>基于DeepSeek大模型，为你提供专业的情感建议</p>
          </router-link>
          <router-link to="/knowledge" class="feature-item">
            <div class="feature-icon">
              <i class="el-icon-collection"></i>
            </div>
            <h4>情感知识库</h4>
            <p>丰富的情感知识，帮助你更好地理解和处理人际关系</p>
          </router-link>
          <router-link to="/subscription" class="feature-item">
            <div class="feature-icon">
              <i class="el-icon-document"></i>
            </div>
            <h4>个性化建议</h4>
            <p>根据你的具体情况，提供定制化的脱单策略</p>
          </router-link>
          <router-link to="/subscription" class="feature-item">
            <div class="feature-icon">
              <i class="el-icon-lock"></i>
            </div>
            <h4>隐私保护</h4>
            <p>严格保护你的个人信息和聊天记录</p>
          </router-link>
        </div>
      </div>
    </section>

    <!-- 订阅服务 -->
    <section class="subscription-section">
      <div class="container">
        <h3>订阅服务</h3>
        <p>解锁更多高级功能，获得专属的脱单指导</p>
        <router-link to="/subscription" class="btn primary">了解详情</router-link>
      </div>
    </section>

    <!-- 底部 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-info">
            <h4>脱单AI</h4>
            <p>智能脱单助手，让恋爱更简单</p>
          </div>
          <div class="footer-links">
            <h5>快速链接</h5>
            <ul>
              <li><router-link to="/">首页</router-link></li>
              <li><router-link to="/knowledge">情感知识</router-link></li>
              <li><router-link to="/chat">AI咨询</router-link></li>
              <li><router-link to="/subscription">订阅服务</router-link></li>
            </ul>
          </div>
          <div class="footer-contact">
            <h5>联系我们</h5>
            <p>邮箱: contact@datingai.com</p>
            <p>电话: 123-456-7890</p>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; 2026 脱单AI. 保留所有权利.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { useUserStore } from '../store/user'

const userStore = useUserStore()

const handleLogout = () => {
  userStore.logout()
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.header .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 70px;
}

.logo h1 {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin: 0;
}

.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-size: 16px;
  padding: 8px 16px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.nav-link:hover {
  color: #409eff;
}

.nav-link.primary {
  background-color: #409eff;
  color: white;
}

.nav-link.primary:hover {
  background-color: #66b1ff;
  color: white;
}

.nav-auth, .nav-user {
  display: flex;
  gap: 10px;
  align-items: center;
}

.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 100px 0;
}

.hero .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 40px;
}

.hero-content {
  flex: 1;
}

.hero-content h2 {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 20px;
  line-height: 1.2;
}

.hero-content p {
  font-size: 18px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.hero-buttons {
  display: flex;
  gap: 20px;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  cursor: pointer;
  border: none;
}

.btn.primary {
  background-color: white;
  color: #409eff;
}

.btn.primary:hover {
  background-color: #f5f7fa;
}

.btn.secondary {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid white;
}

.btn.secondary:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.hero-image {
  flex: 1;
  max-width: 500px;
}

.hero-image img {
  width: 100%;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.features {
  padding: 100px 0;
  background-color: #f5f7fa;
}

.features h3 {
  text-align: center;
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 60px;
  color: #333;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.feature-item {
  background-color: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  text-align: center;
  text-decoration: none;
  cursor: pointer;
  display: block;
}

.feature-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
}

.feature-icon {
  width: 60px;
  height: 60px;
  background-color: #ecf5ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 24px;
  color: #409eff;
}

.feature-item h4 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #333;
}

.feature-item p {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.subscription-section {
  background-color: #409eff;
  color: white;
  padding: 80px 0;
  text-align: center;
}

.subscription-section h3 {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 10px;
}

.subscription-section p {
  font-size: 16px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.footer {
  background-color: #333;
  color: white;
  padding: 60px 0 30px;
  margin-top: auto;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 40px;
  margin-bottom: 40px;
}

.footer-info h4 {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 15px;
}

.footer-info p {
  font-size: 14px;
  opacity: 0.8;
  line-height: 1.5;
}

.footer-links h5, .footer-contact h5 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 15px;
}

.footer-links ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links li {
  margin-bottom: 10px;
}

.footer-links a {
  color: white;
  text-decoration: none;
  font-size: 14px;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.footer-links a:hover {
  opacity: 1;
}

.footer-contact p {
  font-size: 14px;
  opacity: 0.8;
  margin-bottom: 10px;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 30px;
  text-align: center;
  font-size: 14px;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .hero .container {
    flex-direction: column;
    text-align: center;
  }
  
  .hero-content h2 {
    font-size: 36px;
  }
  
  .hero-buttons {
    justify-content: center;
  }
  
  .nav {
    gap: 10px;
  }
  
  .nav-link {
    font-size: 14px;
    padding: 6px 12px;
  }
}
</style>
