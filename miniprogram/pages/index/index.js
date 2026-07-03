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

  goToChat() {
    wx.switchTab({
      url: '/pages/chat/chat'
    })
  },

  goToKnowledge() {
    wx.navigateTo({
      url: '/pages/knowledge/knowledge'
    })
  },

  goToSubscription() {
    wx.navigateTo({
      url: '/pages/subscription/subscription'
    })
  },

  goToProfile() {
    wx.switchTab({
      url: '/pages/profile/profile'
    })
  },

  goToLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  }
})
