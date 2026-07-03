const app = getApp()

Page({
  data: {
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    loading: false
  },

  onUsernameInput(e) {
    this.setData({
      username: e.detail.value
    })
  },

  onEmailInput(e) {
    this.setData({
      email: e.detail.value
    })
  },

  onPasswordInput(e) {
    this.setData({
      password: e.detail.value
    })
  },

  onConfirmPasswordInput(e) {
    this.setData({
      confirmPassword: e.detail.value
    })
  },

  get isFormValid() {
    const { username, email, password, confirmPassword } = this.data
    return username.trim() && 
           email.trim() && 
           password.trim() && 
           confirmPassword.trim() && 
           password === confirmPassword &&
           this.isValidEmail(email)
  },

  isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  },

  handleRegister() {
    const { username, email, password, confirmPassword } = this.data

    if (!username.trim()) {
      wx.showToast({ title: '请输入用户名', icon: 'none' })
      return
    }

    if (!email.trim()) {
      wx.showToast({ title: '请输入邮箱', icon: 'none' })
      return
    }

    if (!this.isValidEmail(email)) {
      wx.showToast({ title: '请输入有效的邮箱', icon: 'none' })
      return
    }

    if (!password.trim()) {
      wx.showToast({ title: '请输入密码', icon: 'none' })
      return
    }

    if (password !== confirmPassword) {
      wx.showToast({ title: '两次输入的密码不一致', icon: 'none' })
      return
    }

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.baseUrl}/auth/register`,
      method: 'POST',
      data: {
        username,
        email,
        password
      },
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.data) {
          wx.showToast({
            title: '注册成功',
            icon: 'success'
          })

          setTimeout(() => {
            wx.navigateTo({
              url: '/pages/login/login'
            })
          }, 1500)
        } else {
          wx.showToast({
            title: '注册失败',
            icon: 'none'
          })
        }
      },
      fail: (err) => {
        console.error('注册失败:', err)
        wx.showToast({
          title: '注册失败，请重试',
          icon: 'none'
        })
      },
      complete: () => {
        this.setData({ loading: false })
      }
    })
  },

  goBack() {
    wx.navigateBack()
  },

  goToLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  }
})