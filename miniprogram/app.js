App({
  onLaunch() {
    console.log('App Launch')
    this.setupRequestInterceptor()
  },
  onShow() {
    console.log('App Show')
    this.checkTokenExpiry()
  },
  onHide() {
    console.log('App Hide')
  },
  globalData: {
    userInfo: null,
    baseUrl: 'http://localhost:8080/api'
  },

  setupRequestInterceptor() {
    const that = this
    const originalRequest = wx.request

    wx.request = function (options) {
      const token = wx.getStorageSync('token')
      if (token && options.header) {
        options.header['Authorization'] = `Bearer ${token}`
      }

      const originalSuccess = options.success
      options.success = function (res) {
        if (res.data && res.data.code === 401) {
          const message = res.data.message
          if (message === 'token_expired' || message === 'invalid_token') {
            that.handleTokenExpired()
            return
          }
        }
        if (originalSuccess) {
          originalSuccess(res)
        }
      }

      const originalFail = options.fail
      options.fail = function (err) {
        if (err.statusCode === 401) {
          that.handleTokenExpired()
          return
        }
        if (originalFail) {
          originalFail(err)
        }
      }

      return originalRequest(options)
    }
  },

  handleTokenExpired() {
    console.log('Token expired, logging out...')
    wx.removeStorageSync('token')
    wx.removeStorageSync('username')
    wx.removeStorageSync('userInfo')

    wx.showModal({
      title: '登录过期',
      content: '您的登录已过期，请重新登录',
      showCancel: false,
      success: function () {
        wx.reLaunch({
          url: '/pages/login/login'
        })
      }
    })
  },

  checkTokenExpiry() {
    const token = wx.getStorageSync('token')
    if (token) {
      const payload = this.decodeToken(token)
      if (payload && payload.exp) {
        const expiryTime = payload.exp * 1000
        const currentTime = Date.now()
        if (currentTime > expiryTime) {
          this.handleTokenExpired()
        }
      }
    }
  },

  decodeToken(token) {
    try {
      const base64Url = token.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
      }).join(''))
      return JSON.parse(jsonPayload)
    } catch (e) {
      console.error('Failed to decode token:', e)
      return null
    }
  }
})