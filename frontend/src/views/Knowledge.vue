<template>
  <div class="knowledge">
    <!-- 头部导航 -->
    <header class="header">
      <div class="container">
        <div class="logo">
          <h1>脱单AI</h1>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/knowledge" class="nav-link active">情感知识</router-link>
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

    <!-- 知识页面内容 -->
    <section class="knowledge-content">
      <div class="container">
        <h2>情感知识库</h2>
        <p>探索丰富的情感知识，帮助你更好地理解和处理人际关系</p>

        <!-- 分类筛选 -->
        <div class="category-filter">
          <el-button 
            v-for="category in categories" 
            :key="category"
            :class="{ active: currentCategory === category }"
            @click="filterByCategory(category)"
          >
            {{ category }}
          </el-button>
        </div>

        <!-- 知识列表 -->
        <div class="knowledge-list">
          <div v-if="knowledgeStore.loading" class="loading">
            <el-spinner type="primary" size="large" />
          </div>
          <div v-else-if="knowledgeStore.knowledgeList.length === 0" class="empty">
            <el-empty description="暂无知识内容" />
          </div>
          <div v-else class="knowledge-grid">
            <div 
              v-for="knowledge in knowledgeStore.knowledgeList" 
              :key="knowledge.id"
              class="knowledge-card"
              @click="viewKnowledge(knowledge.id)"
            >
              <div class="knowledge-card-header">
                <span class="category-tag">{{ knowledge.category }}</span>
              </div>
              <h3>{{ knowledge.title }}</h3>
              <p class="knowledge-content">{{ truncateContent(knowledge.content, 100) }}</p>
              <div class="knowledge-card-footer">
                <span class="publish-date">{{ formatDate(knowledge.createdAt) }}</span>
                <el-button type="text" size="small">查看详情</el-button>
              </div>
            </div>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { useKnowledgeStore } from '../store/knowledge'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()
const knowledgeStore = useKnowledgeStore()

const categories = ['全部', '约会技巧', '沟通技巧', '自我提升', '关系维护']
const currentCategory = ref('全部')

const handleLogout = () => {
  userStore.logout()
}

const filterByCategory = async (category) => {
  currentCategory.value = category
  if (category === '全部') {
    await knowledgeStore.fetchAllKnowledge()
  } else {
    await knowledgeStore.fetchKnowledgeByCategory(category)
  }
}

const viewKnowledge = (id) => {
  router.push(`/knowledge/${id}`)
}

const truncateContent = (content, length) => {
  if (content.length <= length) {
    return content
  }
  return content.substring(0, length) + '...'
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(async () => {
  await knowledgeStore.fetchAllKnowledge()
})
</script>

<style scoped>
.knowledge {
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

.nav-link.active {
  color: #409eff;
  font-weight: 600;
  border-bottom: 2px solid #409eff;
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

.knowledge-content {
  flex: 1;
  padding: 60px 0;
  background-color: #f5f7fa;
}

.knowledge-content h2 {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #333;
}

.knowledge-content p {
  font-size: 16px;
  color: #666;
  margin-bottom: 40px;
}

.category-filter {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.category-filter .el-button {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
}

.category-filter .el-button.active {
  background-color: #409eff;
  color: white;
  border-color: #409eff;
}

.knowledge-list {
  margin-top: 30px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.empty {
  padding: 60px 0;
  text-align: center;
}

.knowledge-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.knowledge-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.knowledge-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
}

.knowledge-card-header {
  margin-bottom: 15px;
}

.category-tag {
  background-color: #ecf5ff;
  color: #409eff;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.knowledge-card h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #333;
  line-height: 1.4;
}

.knowledge-content {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin-bottom: 20px;
}

.knowledge-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #999;
}

.publish-date {
  flex: 1;
}

.knowledge-card-footer .el-button {
  padding: 0;
  color: #409eff;
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
  
  .knowledge-grid {
    grid-template-columns: 1fr;
  }
}
</style>
