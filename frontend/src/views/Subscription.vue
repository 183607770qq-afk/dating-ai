<template>
  <div class="subscription">
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
          <router-link to="/subscription" class="nav-link active">订阅服务</router-link>
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

    <!-- 订阅页面内容 -->
    <section class="subscription-content">
      <div class="container">
        <h2>订阅服务</h2>
        <p>解锁更多高级功能，获得专属的脱单指导</p>

        <!-- 订阅状态 -->
        <div v-if="userStore.getIsLoggedIn" class="subscription-status">
          <div v-if="userStore.getIsSubscribed" class="status-active">
            <el-icon class="status-icon"><CheckCircle /></el-icon>
            <h3>您当前是订阅用户</h3>
            <p>订阅到期时间: {{ formatDate(userStore.subscriptionEndDate) }}</p>
          </div>
          <div v-else class="status-inactive">
            <el-icon class="status-icon"><InfoFilled /></el-icon>
            <h3>您当前不是订阅用户</h3>
            <p>订阅后即可解锁所有高级功能</p>
          </div>
        </div>

        <!-- 订阅计划 -->
        <div class="subscription-plans">
          <div class="plan-card">
            <div class="plan-header">
              <h3>月度订阅</h3>
              <div class="plan-price">
                <span class="price">¥29</span>
                <span class="period">/月</span>
              </div>
            </div>
            <div class="plan-features">
              <ul>
                <li><el-icon class="feature-icon"><Check /></el-icon> 无限次AI咨询</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 访问所有情感知识</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 个性化脱单计划</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 优先客服支持</li>
              </ul>
            </div>
            <div class="plan-footer">
              <el-button 
                v-if="userStore.getIsLoggedIn" 
                type="primary" 
                @click="subscribe('monthly')"
                :disabled="userStore.getIsSubscribed"
              >
                {{ userStore.getIsSubscribed ? '已订阅' : '立即订阅' }}
              </el-button>
              <el-button 
                v-else 
                type="primary" 
                @click="router.push('/login')"
              >
                登录后订阅
              </el-button>
            </div>
          </div>

          <div class="plan-card popular">
            <div class="popular-badge">推荐</div>
            <div class="plan-header">
              <h3>季度订阅</h3>
              <div class="plan-price">
                <span class="price">¥79</span>
                <span class="period">/季度</span>
              </div>
              <div class="plan-saving">节省20%</div>
            </div>
            <div class="plan-features">
              <ul>
                <li><el-icon class="feature-icon"><Check /></el-icon> 无限次AI咨询</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 访问所有情感知识</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 个性化脱单计划</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 优先客服支持</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 专属情感导师指导</li>
              </ul>
            </div>
            <div class="plan-footer">
              <el-button 
                v-if="userStore.getIsLoggedIn" 
                type="primary" 
                @click="subscribe('quarterly')"
                :disabled="userStore.getIsSubscribed"
              >
                {{ userStore.getIsSubscribed ? '已订阅' : '立即订阅' }}
              </el-button>
              <el-button 
                v-else 
                type="primary" 
                @click="router.push('/login')"
              >
                登录后订阅
              </el-button>
            </div>
          </div>

          <div class="plan-card">
            <div class="plan-header">
              <h3>年度订阅</h3>
              <div class="plan-price">
                <span class="price">¥299</span>
                <span class="period">/年</span>
              </div>
              <div class="plan-saving">节省40%</div>
            </div>
            <div class="plan-features">
              <ul>
                <li><el-icon class="feature-icon"><Check /></el-icon> 无限次AI咨询</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 访问所有情感知识</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 个性化脱单计划</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 优先客服支持</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 专属情感导师指导</li>
                <li><el-icon class="feature-icon"><Check /></el-icon> 线下活动优先参与</li>
              </ul>
            </div>
            <div class="plan-footer">
              <el-button 
                v-if="userStore.getIsLoggedIn" 
                type="primary" 
                @click="subscribe('annual')"
                :disabled="userStore.getIsSubscribed"
              >
                {{ userStore.getIsSubscribed ? '已订阅' : '立即订阅' }}
              </el-button>
              <el-button 
                v-else 
                type="primary" 
                @click="router.push('/login')"
              >
                登录后订阅
              </el-button>
            </div>
          </div>
        </div>

        <!-- 常见问题 -->
        <div class="faq">
          <h3>常见问题</h3>
          <el-collapse>
            <el-collapse-item title="订阅后可以取消吗？">
              <p>是的，您可以随时取消订阅。取消后，您的订阅将在当前订阅周期结束后失效。</p>
            </el-collapse-item>
            <el-collapse-item title="订阅包含哪些内容？">
              <p>订阅后，您将获得无限次AI咨询、访问所有情感知识、个性化脱单计划、优先客服支持等高级功能。</p>
            </el-collapse-item>
            <el-collapse-item title="如何更新我的支付信息？">
              <p>您可以在个人中心页面更新您的支付信息。</p>
            </el-collapse-item>
            <el-collapse-item title="订阅费用会自动续费吗？">
              <p>是的，订阅会自动续费，您可以在个人中心页面关闭自动续费功能。</p>
            </el-collapse-item>
          </el-collapse>
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
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { ElMessage } from 'element-plus'
import { CheckCircle, InfoFilled, Check } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const handleLogout = () => {
  userStore.logout()
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const subscribe = async (subscriptionType) => {
  try {
    await userStore.subscribe(subscriptionType)
    ElMessage.success('订阅成功')
  } catch (error) {
    ElMessage.error('订阅失败，请稍后重试')
  }
}
</script>

<style scoped>
.subscription {
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

.subscription-content {
  flex: 1;
  padding: 60px 0;
  background-color: #f5f7fa;
}

.subscription-content h2 {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #333;
  text-align: center;
}

.subscription-content p {
  font-size: 16px;
  color: #666;
  margin-bottom: 40px;
  text-align: center;
}

.subscription-status {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 30px;
  margin-bottom: 40px;
  text-align: center;
}

.status-active {
  color: #67c23a;
}

.status-inactive {
  color: #e6a23c;
}

.status-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.subscription-plans {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  margin-bottom: 60px;
}

.plan-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 30px;
  position: relative;
  transition: all 0.3s ease;
}

.plan-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
}

.plan-card.popular {
  border: 2px solid #409eff;
}

.popular-badge {
  position: absolute;
  top: 0;
  right: 0;
  background-color: #409eff;
  color: white;
  padding: 5px 15px;
  border-bottom-left-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.plan-header {
  text-align: center;
  margin-bottom: 30px;
}

.plan-header h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 15px;
  color: #333;
}

.plan-price {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.plan-price .period {
  font-size: 16px;
  font-weight: normal;
  color: #666;
}

.plan-saving {
  font-size: 14px;
  color: #67c23a;
  font-weight: 600;
}

.plan-features {
  margin-bottom: 30px;
}

.plan-features ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.plan-features li {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  font-size: 14px;
  color: #666;
}

.feature-icon {
  color: #67c23a;
}

.plan-footer {
  text-align: center;
}

.plan-footer .el-button {
  width: 100%;
  height: 40px;
  font-size: 16px;
}

.faq {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.faq h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  color: #333;
  text-align: center;
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
  text-align: left;
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
  text-align: left;
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
  
  .subscription-plans {
    grid-template-columns: 1fr;
  }
}
</style>
