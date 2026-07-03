Page({
  data: {
    topics: [
      {
        icon: '💬',
        title: '聊天开场',
        desc: '从共同场景切入，避免一上来就查户口'
      },
      {
        icon: '☕',
        title: '约会准备',
        desc: '地点、节奏和话题都提前留一点余量'
      },
      {
        icon: '💝',
        title: '表达好感',
        desc: '用具体感受代替夸张承诺，降低对方压力'
      },
      {
        icon: '🧭',
        title: '关系推进',
        desc: '观察回应、尊重边界，让互动更稳定'
      }
    ]
  },

  openTopic(e) {
    const topic = this.data.topics[e.currentTarget.dataset.index]
    wx.showModal({
      title: topic.title,
      content: `${topic.desc}\n\n更完整的知识内容正在整理中，可以先去咨询页让 AI 按你的场景给建议。`,
      confirmText: '去咨询',
      cancelText: '稍后',
      success: (res) => {
        if (res.confirm) {
          wx.switchTab({
            url: '/pages/chat/chat'
          })
        }
      }
    })
  }
})
