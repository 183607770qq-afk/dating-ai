import { defineStore } from 'pinia'
import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

axios.defaults.baseURL = API_BASE_URL

// 检查本地存储中的token
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: token,
    isLoggedIn: !!token,
    isSubscribed: false,
    subscriptionEndDate: null
  }),
  getters: {
    getUser: (state) => state.user,
    getToken: (state) => state.token,
    getIsLoggedIn: (state) => state.isLoggedIn,
    getIsSubscribed: (state) => state.isSubscribed
  },
  actions: {
    async login(username, password) {
      try {
        const response = await axios.post('/auth/login', {
          username,
          password
        })
        const token = response.data.token
        localStorage.setItem('token', token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        this.token = token
        this.isLoggedIn = true
        await this.getSubscriptionStatus()
        return response.data
      } catch (error) {
        throw error
      }
    },
    async register(userData) {
      try {
        const response = await axios.post('/auth/register', userData)
        return response.data
      } catch (error) {
        throw error
      }
    },
    async getSubscriptionStatus() {
      try {
        const response = await axios.get('/subscription/status')
        this.isSubscribed = response.data.isSubscribed
        this.subscriptionEndDate = response.data.subscriptionEndDate
      } catch (error) {
        console.error('Error getting subscription status:', error)
      }
    },
    async subscribe(subscriptionType) {
      try {
        const response = await axios.post('/subscription/subscribe', {
          subscriptionType
        })
        this.isSubscribed = true
        this.subscriptionEndDate = response.data.subscriptionEndDate
        return response.data
      } catch (error) {
        throw error
      }
    },
    logout() {
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
      this.user = null
      this.token = null
      this.isLoggedIn = false
      this.isSubscribed = false
      this.subscriptionEndDate = null
    }
  }
})
