import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { api } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loggedIn = ref(false)
  const needNickname = ref(false)

  async function checkAuth() {
    try {
      user.value = await api.getMe()
      loggedIn.value = true
      needNickname.value = !user.value.nickname
    } catch {
      user.value = null
      loggedIn.value = false
    }
  }

  async function loginWithGoogle(credential: string) {
    const res = await api.googleAuth(credential)
    user.value = res.user
    loggedIn.value = true
    needNickname.value = res.need_nickname
  }

  async function setNickname(nickname: string) {
    await api.setNickname(nickname)
    if (user.value) user.value.nickname = nickname
    needNickname.value = false
  }

  async function logout() {
    await api.logout()
    user.value = null
    loggedIn.value = false
  }

  return { user, loggedIn, needNickname, checkAuth, loginWithGoogle, setNickname, logout }
})
