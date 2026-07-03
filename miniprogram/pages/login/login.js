const app = getApp()

Page({
  data: {
    username: '',
    password: '',
    loading: false
  },

  onLoad() {
    // 检查是否已登录
    const token = wx.getStorageSync('token')
    if (token) {
      wx.switchTab({
        url: '/pages/index/index'
      })
    }
  },

  onUsernameInput(e) {
    this.setData({
      username: e.detail.value
    })
  },

  onPasswordInput(e) {
    this.setData({
      password: e.detail.value
    })
  },

  handleLogin() {
    const { username, password } = this.data
    
    if (!username.trim() || !password.trim()) {
      wx.showToast({
        title: '请输入用户名和密码',
        icon: 'none'
      })
      return
    }

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.baseUrl}/auth/login`,
      method: 'POST',
      data: {
        username,
        password
      },
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.data && res.data.token) {
          wx.setStorageSync('token', res.data.token)
          wx.setStorageSync('username', username)
          
          wx.showToast({
            title: '登录成功',
            icon: 'success'
          })

          setTimeout(() => {
            wx.switchTab({
              url: '/pages/index/index'
            })
          }, 1500)
        } else {
          wx.showToast({
            title: '登录失败',
            icon: 'none'
          })
        }
      },
      fail: (err) => {
        console.error('登录失败:', err)
        wx.showModal({
          title: '网络错误',
          content: '请确保已在开发者工具中开启"不校验域名"模式，或在微信公众平台配置服务器域名',
          showCancel: false
        })
      },
      complete: () => {
        this.setData({ loading: false })
      }
    })
  },

  goToRegister() {
    wx.navigateTo({
      url: '/pages/register/register'
    })
  }
})