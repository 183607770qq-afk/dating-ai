<template>
  <div class="chat">
    <!-- 头部导航 -->
    <header class="header">
      <div class="container">
        <div class="logo">
          <h1>脱单AI</h1>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/knowledge" class="nav-link">情感知识</router-link>
          <router-link to="/chat" class="nav-link active">AI咨询</router-link>
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

    <!-- 聊天界面 -->
    <section class="chat-content">
      <div class="container">
        <div class="chat-container">
          <div class="chat-header">
            <h2>AI情感顾问</h2>
            <p>基于DeepSeek大模型，为你提供专业的情感建议</p>
          </div>
          <div class="chat-messages" ref="chatMessages">
            <!-- 系统消息 -->
            <div class="message system-message">
              <div class="message-content">
                <p>你好！我是你的AI情感顾问，有什么可以帮助你的吗？</p>
              </div>
            </div>
            <!-- 用户消息 -->
            <div 
              v-for="(message, index) in messages" 
              :key="index"
              :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']"
            >
              <div class="message-avatar">
                <img 
                  :src="message.role === 'user' ? userAvatar : aiAvatar" 
                  :alt="message.role === 'user' ? '用户' : 'AI'"
                />
              </div>
              <div class="message-content">
                <p>{{ message.content }}</p>
              </div>
            </div>
            <!-- 加载中消息 -->
            <div v-if="loading && waitingForFirstChunk" class="message ai-message">
              <div class="message-avatar">
                <img :src="aiAvatar" alt="AI" />
              </div>
              <div class="message-content">
                <el-spinner type="primary" size="small" />
                <span>AI正在思考...</span>
              </div>
            </div>
          </div>
          <div class="chat-input">
            <el-input
              v-model="inputMessage"
              placeholder="输入你的问题..."
              @keyup.enter="sendMessage"
              type="textarea"
              :rows="3"
              resize="none"
            ></el-input>
            <el-button type="primary" @click="sendMessage" :loading="loading">发送</el-button>
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
import { ref, nextTick, onMounted, watch } from 'vue'
import { useUserStore } from '../store/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const API_BASE_URL = 'http://localhost:8080/api'
const chatMessages = ref(null)
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const waitingForFirstChunk = ref(false)
const historyLoaded = ref(false)

const userAvatar = 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20icon&image_size=square'
const aiAvatar = 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=AI%20assistant%20avatar%20icon&image_size=square'

const handleLogout = () => {
  userStore.logout()
}

// 加载历史消息
const loadHistory = async () => {
  const token = localStorage.getItem('token')
  if (!token) return

  try {
    const response = await fetch(`${API_BASE_URL}/chat/recent?limit=20`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) return

    const data = await response.json()
    if (data.messages && data.messages.length > 0) {
      messages.value = data.messages
      historyLoaded.value = true
      await nextTick()
      scrollToBottom()
    }
  } catch (e) {
    console.error('Failed to load chat history:', e)
  }
}

onMounted(() => {
  loadHistory()
})

// 登录状态变化时重新加载历史
watch(() => userStore.getIsLoggedIn, (loggedIn) => {
  if (loggedIn && !historyLoaded.value) {
    loadHistory()
  }
})

const sendMessage = async () => {
  if (!inputMessage.value.trim() || loading.value) {
    return
  }

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: inputMessage.value
  })

  // 清空输入框
  const message = inputMessage.value
  inputMessage.value = ''

  // 滚动到底部
  await nextTick()
  scrollToBottom()

  // 显示加载状态
  loading.value = true
  waitingForFirstChunk.value = true

  messages.value.push({
    role: 'ai',
    content: ''
  })
  const aiMessageIndex = messages.value.length - 1

  await nextTick()
  scrollToBottom()

  try {
    await streamAdvice(message, aiMessageIndex)
  } catch (error) {
    messages.value[aiMessageIndex].content = '抱歉，我暂时无法回答你的问题，请稍后再试。'
    ElMessage.error('获取AI回复失败')
  } finally {
    // 隐藏加载状态
    loading.value = false
    waitingForFirstChunk.value = false
    // 滚动到底部
    await nextTick()
    scrollToBottom()
  }
}

const streamAdvice = async (question, aiMessageIndex) => {
  const headers = {
    'Content-Type': 'application/json',
    'Accept': 'text/event-stream'
  }

  const token = localStorage.getItem('token')
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}/llm/stream/advice`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ question })
  })

  if (!response.ok) {
    throw new Error(`Stream request failed: ${response.status}`)
  }

  if (!response.body) {
    const fallback = await response.text()
    await appendStreamContent(aiMessageIndex, fallback)
    return
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buffer = ''

  while (true) {
    const { value, done } = await reader.read()
    if (done) {
      break
    }

    buffer += decoder.decode(value, { stream: true })
    buffer = await consumeStreamBuffer(buffer, aiMessageIndex)
  }

  buffer += decoder.decode()
  if (buffer.trim()) {
    await consumeStreamBlock(buffer, aiMessageIndex)
  }
}

const consumeStreamBuffer = async (buffer, aiMessageIndex) => {
  const blocks = buffer.split(/\r?\n\r?\n/)
  const rest = blocks.pop() || ''
  for (const block of blocks) {
    await consumeStreamBlock(block, aiMessageIndex)
  }
  return rest
}

const consumeStreamBlock = async (block, aiMessageIndex) => {
  const dataLines = block
    .split(/\r?\n/)
    .filter((line) => line.startsWith('data:'))
    .map((line) => line.substring(5).trim())

  const data = dataLines.length ? dataLines.join('\n') : block.trim()
  if (!data || data === '[DONE]') {
    return
  }

  try {
    const parsed = JSON.parse(data)
    await appendStreamContent(aiMessageIndex, parsed.content || parsed.advice || '')
  } catch (error) {
    await appendStreamContent(aiMessageIndex, data)
  }
}

const appendStreamContent = async (aiMessageIndex, content) => {
  if (!content) {
    return
  }

  waitingForFirstChunk.value = false
  messages.value[aiMessageIndex].content += content
  await nextTick()
  scrollToBottom()
}

const scrollToBottom = () => {
  if (chatMessages.value) {
    chatMessages.value.scrollTop = chatMessages.value.scrollHeight
  }
}
</script>

<style scoped>
.chat {
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

.chat-content {
  flex: 1;
  padding: 60px 0;
  background-color: #f5f7fa;
}

.chat-container {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  height: 600px;
  display: flex;
  flex-direction: column;
}

.chat-header {
  background-color: #409eff;
  color: white;
  padding: 20px;
  text-align: center;
}

.chat-header h2 {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 5px;
}

.chat-header p {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message {
  display: flex;
  gap: 10px;
  max-width: 80%;
}

.system-message {
  align-self: center;
  max-width: 90%;
}

.system-message .message-content {
  background-color: #f0f0f0;
  color: #666;
  border-radius: 12px;
  padding: 10px 15px;
  font-size: 14px;
}

.user-message {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.ai-message {
  align-self: flex-start;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-content {
  flex: 1;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.5;
}

.user-message .message-content {
  background-color: #409eff;
  color: white;
  border-bottom-right-radius: 4px;
}

.ai-message .message-content {
  background-color: #f0f0f0;
  color: #333;
  border-bottom-left-radius: 4px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat-input {
  padding: 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.chat-input .el-input {
  flex: 1;
}

.chat-input .el-button {
  align-self: flex-end;
  height: 40px;
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
  
  .chat-container {
    height: 500px;
  }
  
  .message {
    max-width: 90%;
  }
}
</style>
