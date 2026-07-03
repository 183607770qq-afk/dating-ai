const app = getApp()

Page({
  data: {
    messages: [],
    inputMessage: '',
    loading: false,
    scrollTop: 0,
    isLoggedIn: false,
    welcomeSet: false
  },

  onLoad() {
    console.log('Chat page loaded')
    this.checkLoginStatus()
    this.loadHistory()
  },

  onShow() {
    this.checkLoginStatus()
  },

  /** 从后端加载最近的历史消息 */
  loadHistory() {
    const token = wx.getStorageSync('token')
    if (!token) {
      // 未登录，使用默认欢迎消息
      this.setWelcomeMessage()
      return
    }

    const that = this
    wx.request({
      url: `${app.globalData.baseUrl}/chat/recent?limit=20`,
      method: 'GET',
      header: {
        'Authorization': 'Bearer ' + token
      },
      success(res) {
        if (res.statusCode === 200 && res.data.messages && res.data.messages.length > 0) {
          console.log('Loaded history messages:', res.data.messages.length)
          that.setData({
            messages: res.data.messages,
            welcomeSet: true
          }, () => {
            that.scrollToBottom()
          })
        } else {
          // 无历史记录，使用默认欢迎消息
          that.setWelcomeMessage()
        }
      },
      fail() {
        console.log('Failed to load history, using welcome message')
        that.setWelcomeMessage()
      }
    })
  },

  setWelcomeMessage() {
    if (this.data.welcomeSet) {
      console.log('Welcome message already set, skipping')
      return
    }
    console.log('Setting welcome message')
    const welcomeMsg = {
      role: 'ai',
      content: '你好，我是你的关系顾问。你可以把聊天记录、约会经过或现在的纠结告诉我，我会帮你拆成更清晰的下一步。'
    }
    this.setData({
      messages: [welcomeMsg],
      welcomeSet: true
    }, () => {
      console.log('Welcome message set successfully:', this.data.messages)
    })
  },

  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    console.log('Token exists:', !!token)
    this.setData({
      isLoggedIn: !!token
    })
  },

  goBack() {
    wx.switchTab({
      url: '/pages/index/index'
    })
  },

  onInput(e) {
    this.setData({
      inputMessage: e.detail.value
    })
  },

  sendMessage() {
    const message = this.data.inputMessage.trim()
    if (!message || this.data.loading) {
      return
    }

    console.log('Sending message:', message)

    const userMsg = { role: 'user', content: message }
    const newMessages = [...this.data.messages, userMsg]
    
    this.setData({
      messages: newMessages,
      inputMessage: '',
      loading: true
    }, () => {
      this.scrollToBottom()
    })

    this.callApi(message)
  },

  callApi(message) {
    const that = this
    let receivedStreamChunk = false

    this.streamBuffer = ''
    this.streamDecoder = typeof TextDecoder !== 'undefined' ? new TextDecoder('utf-8') : null
    this.startAiMessage()

    const requestTask = wx.request({
      url: `${app.globalData.baseUrl}/llm/stream/advice`,
      method: 'POST',
      data: { question: message },
      header: {
        'Content-Type': 'application/json',
        'Accept': 'text/event-stream'
      },
      enableChunked: true,
      success(res) {
        console.log('API response received:', res)
        if (receivedStreamChunk) {
          that.completeStreaming()
          return
        }

        if (res.data && res.data.advice) {
          console.log('AI response received, length:', res.data.advice.length)
          that.appendToLastAiMessage(res.data.advice)
          that.completeStreaming()
        } else if (typeof res.data === 'string' && res.data) {
          that.handleStreamChunk(res.data)
          that.completeStreaming()
        } else {
          console.log('No advice found, using mock data')
          that.appendToLastAiMessage(that.getMockResponse(message))
          that.completeStreaming()
        }
      },
      fail(err) {
        console.error('API request failed:', err)
        that.appendToLastAiMessage(that.getMockResponse(message))
        that.completeStreaming()
      }
    })

    if (requestTask && requestTask.onChunkReceived) {
      requestTask.onChunkReceived((res) => {
        receivedStreamChunk = true
        const text = that.decodeChunk(res.data)
        that.handleStreamChunk(text)
      })
    }
  },

  startAiMessage() {
    const aiMsg = { role: 'ai', content: '' }
    this.setData({
      messages: [...this.data.messages, aiMsg]
    }, () => {
      this.scrollToBottom()
    })
  },

  appendToLastAiMessage(content) {
    if (!content) {
      return
    }

    const lastIndex = this.data.messages.length - 1
    const messages = this.data.messages.map((msg, idx) => {
      if (idx === lastIndex && msg.role === 'ai') {
        return { ...msg, content: msg.content + content }
      }
      return msg
    })

    this.setData({ messages }, () => {
      this.scrollToBottom()
    })
  },

  completeStreaming() {
    if (this.streamDecoder) {
      const rest = this.streamDecoder.decode()
      if (rest) {
        this.handleStreamChunk(rest)
      }
    }

    if (this.streamBuffer) {
      const rest = this.streamBuffer
      this.streamBuffer = ''
      this.handleStreamBlock(rest)
    }

    this.setData({ loading: false })
    this.streamBuffer = ''
    this.streamDecoder = null
  },

  handleStreamChunk(chunkText) {
    if (!chunkText) {
      return
    }

    this.streamBuffer = `${this.streamBuffer || ''}${chunkText}`
    const blocks = this.streamBuffer.split(/\r?\n\r?\n/)
    this.streamBuffer = blocks.pop() || ''

    blocks.forEach((block) => this.handleStreamBlock(block))
  },

  handleStreamBlock(block) {
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
      this.appendToLastAiMessage(parsed.content || parsed.advice || '')
    } catch (e) {
      this.appendToLastAiMessage(data)
    }
  },

  decodeChunk(arrayBuffer) {
    if (this.streamDecoder) {
      return this.streamDecoder.decode(arrayBuffer, { stream: true })
    }

    const bytes = new Uint8Array(arrayBuffer)
    let encoded = ''
    for (let i = 0; i < bytes.length; i += 1) {
      encoded += `%${bytes[i].toString(16).padStart(2, '0')}`
    }

    try {
      return decodeURIComponent(encoded)
    } catch (e) {
      return String.fromCharCode.apply(null, bytes)
    }
  },

  getMockResponse(message) {
    const responses = [
      '谢谢你的提问！作为你的AI情感顾问，我很乐意为你提供帮助。脱单需要耐心和方法，我们可以一起探讨适合你的策略。',
      '理解你的困惑，很多人在面对感情问题时都会感到迷茫。关键是要保持积极的心态，不断提升自己。',
      '好的，我来帮你分析一下。首先，建立自信是很重要的一步。你觉得自己在哪些方面可以做得更好呢？',
      '情感问题往往需要从自我认知开始。了解自己的需求和优势，才能更好地找到合适的伴侣。',
      '与人交往是一门艺术，需要学习和实践。不要害怕失败，每一次尝试都是成长的机会。'
    ]
    
    const keywords = {
      '你好': '你好！很高兴为你服务。有什么关于情感或脱单的问题想要问我吗？',
      '谢谢': '不客气！能帮到你我很开心。如果还有其他问题，随时可以问我。',
      '爱': '爱是一种需要学习的能力，包括理解、包容和付出。你想了解哪方面呢？',
      '单身': '单身状态其实是一个很好的自我提升时期。利用这段时间好好了解自己，为未来的关系做好准备。',
      '约会': '约会的关键是放松心态，做真实的自己。准备一些轻松的话题，保持良好的沟通氛围。',
      '表白': '表白需要勇气，但也要讲究时机和方式。确保对方也对你有好感，选择合适的场合很重要。'
    }

    for (const [keyword, response] of Object.entries(keywords)) {
      if (message.includes(keyword)) {
        return response
      }
    }

    return responses[Math.floor(Math.random() * responses.length)]
  },

  showTypingEffect(fullText) {
    console.log('Starting typing effect for text length:', fullText.length)
    
    const aiMsg = { role: 'ai', content: '' }
    const newMessages = [...this.data.messages, aiMsg]
    
    this.setData({
      messages: newMessages
    }, () => {
      let currentIndex = 0
      const typingInterval = setInterval(() => {
        if (currentIndex < fullText.length) {
          const charsToAdd = Math.min(2, fullText.length - currentIndex)
          const newContent = fullText.substring(0, currentIndex + charsToAdd)
          currentIndex += charsToAdd

          this.setData({
            messages: this.data.messages.map((msg, idx) => {
              if (idx === this.data.messages.length - 1) {
                return { ...msg, content: newContent }
              }
              return msg
            })
          })

          this.scrollToBottom()
        } else {
          clearInterval(typingInterval)
          this.setData({ loading: false })
        }
      }, 30)
    })
  },

  showErrorMessage() {
    const errorMsg = { role: 'ai', content: '抱歉，我暂时无法回答你的问题，请稍后再试。' }
    const newMessages = [...this.data.messages, errorMsg]
    this.setData({
      messages: newMessages,
      loading: false
    })
    this.scrollToBottom()
  },

  scrollToBottom() {
    setTimeout(() => {
      this.setData({ scrollTop: 99999 })
    }, 50)
  },

  onScrollToLower() {}
})
