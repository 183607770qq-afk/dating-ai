<template>
  <div class="knowledge-detail">
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

    <!-- 知识详情内容 -->
    <section class="detail-content">
      <div class="container">
        <div v-if="knowledgeStore.loading" class="loading">
          <el-spinner type="primary" size="large" />
        </div>
        <div v-else-if="!knowledgeStore.currentKnowledge" class="not-found">
          <el-empty description="知识内容不存在" />
        </div>
        <div v-else class="knowledge-detail-card">
          <div class="knowledge-header">
            <span class="category-tag">{{ knowledgeStore.currentKnowledge.category }}</span>
            <h2>{{ knowledgeStore.currentKnowledge.title }}</h2>
            <div class="publish-info">
              <span class="publish-date">{{ formatDate(knowledgeStore.currentKnowledge.createdAt) }}</span>
            </div>
          </div>
          <div class="knowledge-body">
            <div v-html="formatContent(knowledgeStore.currentKnowledge.content)"></div>
          </div>
          <div class="knowledge-footer">
            <el-button type="primary" @click="goBack">返回列表</el-button>
          </div>
        </div>
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
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../store/user'
import { useKnowledgeStore } from '../store/knowledge'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const knowledgeStore = useKnowledgeStore()

const handleLogout = () => {
  userStore.logout()
}

const goBack = () => {
  router.push('/knowledge')
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatContent = (content) => {
  // 简单的内容格式化，实际项目中可能需要更复杂的处理
  return content
    .replace(/\n/g, '<br>')
    .replace(/### (.*?)/g, '<h3>$1</h3>')
    .replace(/## (.*?)/g, '<h2>$1</h2>')
    .replace(/# (.*?)/g, '<h1>$1</h1>')
}

onMounted(async () => {
  const id = route.params.id
  await knowledgeStore.fetchKnowledgeById(id)
})
</script>

<style scoped>
.knowledge-detail {
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
  max-width: 800px;
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

.detail-content {
  flex: 1;
  padding: 60px 0;
  background-color: #f5f7fa;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.not-found {
  padding: 60px 0;
  text-align: center;
}

.knowledge-detail-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.knowledge-header {
  margin-bottom: 30px;
}

.category-tag {
  background-color: #ecf5ff;
  color: #409eff;
  padding: 6px 16px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 15px;
  display: inline-block;
}

.knowledge-header h2 {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #333;
  line-height: 1.3;
}

.publish-info {
  font-size: 14px;
  color: #999;
}

.knowledge-body {
  margin-bottom: 40px;
  line-height: 1.8;
  color: #333;
}

.knowledge-body h1, .knowledge-body h2, .knowledge-body h3 {
  margin-top: 30px;
  margin-bottom: 20px;
  color: #333;
}

.knowledge-body h1 {
  font-size: 24px;
}

.knowledge-body h2 {
  font-size: 20px;
}

.knowledge-body h3 {
  font-size: 18px;
}

.knowledge-body p {
  margin-bottom: 20px;
}

.knowledge-footer {
  text-align: center;
  padding-top: 30px;
  border-top: 1px solid #f0f0f0;
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
  .nav {
    gap: 10px;
  }
  
  .nav-link {
    font-size: 14px;
    padding: 6px 12px;
  }
  
  .knowledge-detail-card {
    padding: 20px;
  }
  
  .knowledge-header h2 {
    font-size: 24px;
  }
}
</style>
