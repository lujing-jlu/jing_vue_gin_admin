import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login } from '@/api/user'
import type { LoginRequest } from '@/api/user'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const loginAction = async (loginForm: LoginRequest) => {
    try {
      const res = await login(loginForm)
      token.value = res.token
      userInfo.value = res.user
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.user))
      router.push('/')
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    router.push('/login')
  }

  return {
    token,
    userInfo,
    loginAction,
    logout
  }
}) 