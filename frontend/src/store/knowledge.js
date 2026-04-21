import { defineStore } from 'pinia'
import axios from 'axios'

export const useKnowledgeStore = defineStore('knowledge', {
  state: () => ({
    knowledgeList: [],
    currentKnowledge: null,
    loading: false,
    error: null
  }),
  getters: {
    getKnowledgeList: (state) => state.knowledgeList,
    getCurrentKnowledge: (state) => state.currentKnowledge,
    getLoading: (state) => state.loading,
    getError: (state) => state.error
  },
  actions: {
    async fetchAllKnowledge() {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get('/knowledge/all')
        this.knowledgeList = response.data
      } catch (error) {
        this.error = 'Failed to fetch knowledge'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    async fetchKnowledgeByCategory(category) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get(`/knowledge/category/${category}`)
        this.knowledgeList = response.data
      } catch (error) {
        this.error = 'Failed to fetch knowledge by category'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    async fetchKnowledgeById(id) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get(`/knowledge/${id}`)
        this.currentKnowledge = response.data
      } catch (error) {
        this.error = 'Failed to fetch knowledge detail'
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    async createKnowledge(knowledge) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post('/knowledge/create', knowledge)
        this.knowledgeList.push(response.data)
        return response.data
      } catch (error) {
        this.error = 'Failed to create knowledge'
        console.error(error)
        throw error
      } finally {
        this.loading = false
      }
    },
    async updateKnowledge(id, knowledge) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.put(`/knowledge/${id}`, knowledge)
        const index = this.knowledgeList.findIndex(item => item.id === id)
        if (index !== -1) {
          this.knowledgeList[index] = response.data
        }
        if (this.currentKnowledge && this.currentKnowledge.id === id) {
          this.currentKnowledge = response.data
        }
        return response.data
      } catch (error) {
        this.error = 'Failed to update knowledge'
        console.error(error)
        throw error
      } finally {
        this.loading = false
      }
    },
    async deleteKnowledge(id) {
      this.loading = true
      this.error = null
      try {
        await axios.delete(`/knowledge/${id}`)
        this.knowledgeList = this.knowledgeList.filter(item => item.id !== id)
        if (this.currentKnowledge && this.currentKnowledge.id === id) {
          this.currentKnowledge = null
        }
      } catch (error) {
        this.error = 'Failed to delete knowledge'
        console.error(error)
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
