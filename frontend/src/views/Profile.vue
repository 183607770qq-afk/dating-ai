<template>
  <div class="profile">
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
            <router-link to="/profile" class="nav-link active">个人中心</router-link>
            <button @click="handleLogout" class="nav-link">退出</button>
          </div>
        </nav>
      </div>
    </header>

    <!-- 个人中心内容 -->
    <section class="profile-content">
      <div class="container">
        <div v-if="!userStore.getIsLoggedIn" class="login-prompt">
          <el-empty description="请先登录" />
          <router-link to="/login" class="btn primary">立即登录</router-link>
        </div>
        <div v-else class="profile-container">
          <div class="profile-header">
            <h2>个人中心</h2>
          </div>
          <div class="profile-tabs">
            <el-tabs v-model="activeTab">
              <el-tab-pane label="个人信息" name="info">
                <div class="profile-info">
                  <el-form :model="userInfo" label-width="120px">
                    <el-form-item label="用户名">
                      <el-input v-model="userInfo.username" disabled />
                    </el-form-item>
                    <el-form-item label="邮箱">
                      <el-input v-model="userInfo.email" disabled />
                    </el-form-item>
                    <el-form-item label="注册时间">
                      <el-input v-model="userInfo.createdAt" disabled />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary">编辑信息</el-button>
                    </el-form-item>
                  </el-form>
                </div>
              </el-tab-pane>
              <el-tab-pane label="订阅状态" name="subscription">
                <div class="subscription-info">
                  <div v-if="userStore.getIsSubscribed" class="status-active">
                    <el-icon class="status-icon"><CircleCheckFilled /></el-icon>
                    <h3>您当前是订阅用户</h3>
                    <p>订阅类型: {{ subscriptionType }}</p>
                    <p>订阅到期时间: {{ formatDate(userStore.subscriptionEndDate) }}</p>
                    <el-button type="primary">管理订阅</el-button>
                  </div>
                  <div v-else class="status-inactive">
                    <el-icon class="status-icon"><InfoFilled /></el-icon>
                    <h3>您当前不是订阅用户</h3>
                    <p>订阅后即可解锁所有高级功能</p>
                    <router-link to="/subscription" class="btn primary">立即订阅</router-link>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="历史记录" name="history">
                <div class="history-info">
                  <div v-if="historyMessages.length === 0 && !historyLoading">
                    <el-empty description="暂无历史记录" />
                  </div>
                  <div v-else class="history-list">
                    <div 
                      v-for="msg in historyMessages" 
                      :key="msg.id" 
                      :class="['history-item', msg.role === 'user' ? 'user-msg' : 'ai-msg']"
                    >
                      <div class="history-role">{{ msg.role === 'user' ? '你' : 'AI顾问' }}</div>
                      <div class="history-content">{{ msg.content }}</div>
                      <div class="history-time">{{ formatTime(msg.createdAt) }}</div>
                    </div>
                  </div>
                  <div v-if="historyTotal > historyMessages.length" class="history-load-more">
                    <el-button @click="loadMoreHistory" :loading="historyLoading">加载更多</el-button>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="设置" name="settings">
                <div class="settings-info">
                  <el-form>
                    <el-form-item label="修改密码">
                      <el-button type="primary">修改密码</el-button>
                    </el-form-item>
                    <el-form-item label="通知设置">
                      <el-switch v-model="notificationSettings.email" label="邮箱通知" />
                      <el-switch v-model="notificationSettings.sms" label="短信通知" />
                    </el-form-item>
                    <el-form-item label="隐私设置">
                      <el-switch v-model="privacySettings.publicProfile" label="公开个人资料" />
                      <el-switch v-model="privacySettings.showOnlineStatus" label="显示在线状态" />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary">保存设置</el-button>
                    </el-form-item>
                  </el-form>
                </div>
              </el-tab-pane>
            </el-tabs>
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
import { ref, watch } from 'vue'
import { useUserStore } from '../store/user'
import { CircleCheckFilled, InfoFilled } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const API_BASE_URL = 'http://localhost:8080/api'
const userStore = useUserStore()
const activeTab = ref('info')

const userInfo = ref({
  username: 'user',
  email: 'user@example.com',
  createdAt: '2026-01-01'
})

const subscriptionType = ref('月度订阅')

const notificationSettings = ref({
  email: true,
  sms: false
})

const privacySettings = ref({
  publicProfile: false,
  showOnlineStatus: true
})

const handleLogout = () => {
  userStore.logout()
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// ---- 历史记录相关 ----
const historyMessages = ref([])
const historyTotal = ref(0)
const historyPage = ref(0)
const historyLoading = ref(false)

const loadHistory = async (page = 0) => {
  const token = localStorage.getItem('token')
  if (!token) return

  historyLoading.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/chat/history?page=${page}&size=20`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) {
      historyLoading.value = false
      return
    }
    const data = await response.json()
    if (page === 0) {
      historyMessages.value = data.messages || []
    } else {
      historyMessages.value.push(...(data.messages || []))
    }
    historyTotal.value = data.total || 0
    historyPage.value = page
  } catch (e) {
    console.error('Failed to load history:', e)
  } finally {
    historyLoading.value = false
  }
}

const loadMoreHistory = () => {
  loadHistory(historyPage.value + 1)
}

const formatTime = (timeStr) => {
  return dayjs(timeStr).format('MM-DD HH:mm')
}

// 切换到历史记录 Tab 时自动加载
watch(activeTab, (tab) => {
  if (tab === 'history' && historyMessages.value.length === 0) {
    loadHistory(0)
  }
})
</script>

<style scoped>
.profile {
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

.profile-content {
  flex: 1;
  padding: 60px 0;
  background-color: #f5f7fa;
}

.login-prompt {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 60px 40px;
  text-align: center;
  max-width: 500px;
  margin: 0 auto;
}

.login-prompt .btn {
  margin-top: 30px;
  padding: 10px 30px;
}

.profile-container {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.profile-header {
  margin-bottom: 30px;
}

.profile-header h2 {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.profile-tabs {
  width: 100%;
}

.profile-info, .subscription-info, .history-info, .settings-info {
  padding: 20px 0;
}

.status-active {
  color: #67c23a;
  text-align: center;
  padding: 40px 0;
}

.status-inactive {
  color: #e6a23c;
  text-align: center;
  padding: 40px 0;
}

.status-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.status-active h3, .status-inactive h3 {
  margin-bottom: 15px;
}

.status-active p, .status-inactive p {
  margin-bottom: 10px;
  color: #666;
}

.btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  cursor: pointer;
  border: none;
  display: inline-block;
}

.btn.primary {
  background-color: #409eff;
  color: white;
}

.btn.primary:hover {
  background-color: #66b1ff;
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
  
  .profile-container {
    padding: 20px;
  }
}

/* 历史记录样式 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  padding: 12px 16px;
  border-radius: 8px;
  background-color: #f5f7fa;
}

.history-item.user-msg {
  background-color: #ecf5ff;
  border-left: 3px solid #409eff;
}

.history-item.ai-msg {
  background-color: #f0f9eb;
  border-left: 3px solid #67c23a;
}

.history-role {
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.history-content {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  word-break: break-word;
  white-space: pre-wrap;
}

.history-time {
  font-size: 12px;
  color: #c0c4cc;
  margin-top: 8px;
  text-align: right;
}

.history-load-more {
  text-align: center;
  padding: 16px 0;
}
</style>
