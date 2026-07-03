Page({
  data: {
    isLoggedIn: false,
    benefits: [
      '更多 AI 咨询额度',
      '约会复盘与回复建议',
      '长期关系成长计划',
      '优先体验新功能'
    ]
  },

  onLoad() {
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
  },

  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    this.setData({
      isLoggedIn: !!token
    })
  },

  handleSubscribe() {
    if (!this.data.isLoggedIn) {
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }

    wx.showToast({
      title: '支付功能开发中',
      icon: 'none'
    })
  },

  goToChat() {
    wx.switchTab({
      url: '/pages/chat/chat'
    })
  }
})
