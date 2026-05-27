import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import LoginView from './views/LoginView.vue'
import NicknameView from './views/NicknameView.vue'
import HomeView from './views/HomeView.vue'
import WorkoutView from './views/WorkoutView.vue'
import HistoryView from './views/HistoryView.vue'
import ProgressView from './views/ProgressView.vue'
import SetupView from './views/SetupView.vue'
import BodyWeightView from './views/BodyWeightView.vue'
import SettingsView from './views/SettingsView.vue'
import LeaderboardView from './views/LeaderboardView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: LoginView, meta: { public: true } },
    { path: '/nickname', component: NicknameView },
    { path: '/', component: HomeView },
    { path: '/workout', component: WorkoutView },
    { path: '/history', component: HistoryView },
    { path: '/progress', component: ProgressView },
    { path: '/setup', component: SetupView },
    { path: '/bodyweight', component: BodyWeightView },
    { path: '/settings', component: SettingsView },
    { path: '/leaderboard', component: LeaderboardView },
  ],
})

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)

// Auth guard
import { useAuthStore } from './stores/auth'
let authChecked = false

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // Only check auth once on app load
  if (!authChecked) {
    authChecked = true
    await auth.checkAuth()
  }

  if (!auth.loggedIn && !to.meta.public) {
    return '/login'
  }
  if (auth.loggedIn && to.path === '/login') {
    return '/'
  }
})

// Listen for auth expiry (session expired mid-use)
window.addEventListener('auth-expired', () => {
  const auth = useAuthStore()
  auth.loggedIn = false
  auth.user = null
  authChecked = true // don't re-check, we know it's expired
  router.push('/login')
})

app.mount('#app')

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').catch(() => {})
  })
}
