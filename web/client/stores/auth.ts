import { defineStore } from 'pinia'
import type { AuthResponse, User } from '~/types/api'

const tokenKey = 'vidfixa.auth.token'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '' as string,
    user: null as User | null,
    ready: false,
    initializing: null as Promise<void> | null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token && state.user),
    isAdmin: (state) => state.user?.role === 'admin',
  },
  actions: {
    async initialize() {
      if (!import.meta.client || this.ready) return
      if (this.initializing) return this.initializing
      this.initializing = (async () => {
        this.token = localStorage.getItem(tokenKey) || ''
        try {
          if (this.token) await this.refreshUser()
        } catch {
          this.clearSession()
        } finally {
          this.ready = true
          this.initializing = null
        }
      })()
      return this.initializing
    },
    async refreshUser() {
      if (!this.token) {
        this.user = null
        return
      }
      const config = useRuntimeConfig()
      const result = await $fetch<AuthResponse>('/auth/me', {
        baseURL: config.public.apiBaseUrl,
        headers: { Authorization: `Bearer ${this.token}` },
        credentials: 'include',
      })
      this.user = result.user
    },
    async login(email: string, password: string) {
      const config = useRuntimeConfig()
      const result = await $fetch<AuthResponse>('/auth/login', {
        baseURL: config.public.apiBaseUrl,
        method: 'POST',
        body: { email, password },
        credentials: 'include',
      })
      if (!result.token) throw new Error('The API did not return an access token.')
      this.token = result.token
      localStorage.setItem(tokenKey, result.token)
      await this.refreshUser()
      this.ready = true
      return this.user!
    },
    async register(fullName: string, email: string, password: string) {
      const config = useRuntimeConfig()
      const result = await $fetch<AuthResponse>('/auth/register', {
        baseURL: config.public.apiBaseUrl,
        method: 'POST',
        body: { full_name: fullName, email, password },
        credentials: 'include',
      })
      if (!result.token) throw new Error('The API did not return an access token.')
      this.token = result.token
      localStorage.setItem(tokenKey, result.token)
      await this.refreshUser()
      this.ready = true
      return this.user!
    },
    clearSession() {
      this.token = ''
      this.user = null
      this.ready = true
      if (import.meta.client) localStorage.removeItem(tokenKey)
    },
    logout() {
      this.clearSession()
      return navigateTo('/login')
    },
  },
})
