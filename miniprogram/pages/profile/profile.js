const app = getApp()

Page({
  data: {
    isLoggedIn: false,
    username: '',
    avatarText: ''
  },

  onLoad() {
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
  },

  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    const username = wx.getStorageSync('username')
    
    this.setData({
      isLoggedIn: !!token,
      username: username || '',
      avatarText: username ? username.charAt(0).toUpperCase() : ''
    })
  },

  goToLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  },

  goToChat() {
    wx.switchTab({
      url: '/pages/chat/chat'
    })
  },

  goToHistory() {
    wx.switchTab({
      url: '/pages/chat/chat'
    })
  },

  handleLogout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          wx.removeStorageSync('token')
          wx.removeStorageSync('username')
          
          this.setData({
            isLoggedIn: false,
            username: '',
            avatarText: ''
          })
          
          wx.showToast({
            title: '已退出登录',
            icon: 'success'
          })
        }
      }
    })
  }
})
