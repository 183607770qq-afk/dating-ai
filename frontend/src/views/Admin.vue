<template>
  <div class="admin">
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
            <router-link to="/admin" class="nav-link active">管理后台</router-link>
            <button @click="handleLogout" class="nav-link">退出</button>
          </div>
        </nav>
      </div>
    </header>

    <!-- 管理员页面内容 -->
    <section class="admin-content">
      <div class="container">
        <div v-if="!userStore.getIsLoggedIn" class="login-prompt">
          <el-empty description="请先登录" />
          <router-link to="/login" class="btn primary">立即登录</router-link>
        </div>
        <div v-else class="admin-container">
          <div class="admin-header">
            <h2>管理后台</h2>
          </div>
          <div class="admin-tabs">
            <el-tabs v-model="activeTab">
              <el-tab-pane label="用户管理" name="users">
                <div class="users-management">
                  <el-button type="primary" class="add-user-btn">添加用户</el-button>
                  <el-table :data="users" style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column prop="username" label="用户名" />
                    <el-table-column prop="email" label="邮箱" />
                    <el-table-column prop="role" label="角色" />
                    <el-table-column prop="isSubscribed" label="订阅状态" />
                    <el-table-column prop="createdAt" label="创建时间" />
                    <el-table-column label="操作" width="150">
                      <template #default="scope">
                        <el-button size="small" type="primary">编辑</el-button>
                        <el-button size="small" type="danger">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>
              <el-tab-pane label="知识管理" name="knowledge">
                <div class="knowledge-management">
                  <el-button type="primary" class="add-knowledge-btn">添加知识</el-button>
                  <el-table :data="knowledgeList" style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column prop="title" label="标题" />
                    <el-table-column prop="category" label="分类" />
                    <el-table-column prop="isPublished" label="发布状态" />
                    <el-table-column prop="createdAt" label="创建时间" />
                    <el-table-column label="操作" width="150">
                      <template #default="scope">
                        <el-button size="small" type="primary">编辑</el-button>
                        <el-button size="small" type="danger">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>
              <el-tab-pane label="订阅管理" name="subscriptions">
                <div class="subscriptions-management">
                  <el-table :data="subscriptions" style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column prop="username" label="用户名" />
                    <el-table-column prop="subscriptionType" label="订阅类型" />
                    <el-table-column prop="startDate" label="开始时间" />
                    <el-table-column prop="endDate" label="结束时间" />
                    <el-table-column prop="status" label="状态" />
                    <el-table-column label="操作" width="150">
                      <template #default="scope">
                        <el-button size="small" type="primary">编辑</el-button>
                        <el-button size="small" type="danger">取消</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>
              <el-tab-pane label="系统设置" name="settings">
                <div class="system-settings">
                  <el-form>
                    <el-form-item label="网站标题">
                      <el-input v-model="systemSettings.siteTitle" />
                    </el-form-item>
                    <el-form-item label="网站描述">
                      <el-input v-model="systemSettings.siteDescription" type="textarea" />
                    </el-form-item>
                    <el-form-item label="DeepSeek API Key">
                      <el-input v-model="systemSettings.deepSeekApiKey" type="password" />
                    </el-form-item>
                    <el-form-item label="Milvus配置">
                      <el-input v-model="systemSettings.milvusHost" placeholder="Milvus主机" />
                      <el-input v-model="systemSettings.milvusPort" placeholder="Milvus端口" style="margin-top: 10px" />
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
import { ref } from 'vue'
import { useUserStore } from '../store/user'

const userStore = useUserStore()
const activeTab = ref('users')

const users = ref([
  { id: 1, username: 'user', email: 'user@example.com', role: 'USER', isSubscribed: false, createdAt: '2026-01-01' },
  { id: 2, username: 'admin', email: 'admin@example.com', role: 'ADMIN', isSubscribed: true, createdAt: '2026-01-01' }
])

const knowledgeList = ref([
  { id: 1, title: '如何提高自信心', category: '自我提升', isPublished: true, createdAt: '2026-01-01' },
  { id: 2, title: '第一次约会技巧', category: '约会技巧', isPublished: true, createdAt: '2026-01-02' }
])

const subscriptions = ref([
  { id: 1, username: 'admin', subscriptionType: '月度订阅', startDate: '2026-01-01', endDate: '2026-02-01', status: 'ACTIVE' }
])

const systemSettings = ref({
  siteTitle: '脱单AI',
  siteDescription: '智能脱单助手，让恋爱更简单',
  deepSeekApiKey: 'your-api-key',
  milvusHost: 'localhost',
  milvusPort: '19530'
})

const handleLogout = () => {
  userStore.logout()
}
</script>

<style scoped>
.admin {
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

.admin-content {
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

.admin-container {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.admin-header {
  margin-bottom: 30px;
}

.admin-header h2 {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.admin-tabs {
  width: 100%;
}

.users-management, .knowledge-management, .subscriptions-management, .system-settings {
  padding: 20px 0;
}

.add-user-btn, .add-knowledge-btn {
  margin-bottom: 20px;
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
  
  .admin-container {
    padding: 20px;
  }
}
</style>
