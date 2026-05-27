<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const error = ref('')

declare const google: any

onMounted(async () => {
  // Get Google Client ID from server
  try {
    const res = await fetch('/api/config')
    const config = await res.json()
    if (!config.google_client_id) {
      error.value = 'GOOGLE_CLIENT_ID not configured on server'
      return
    }
    loadGoogleScript(config.google_client_id)
  } catch {
    error.value = 'Cannot reach server'
  }
})

function loadGoogleScript(clientId: string) {
  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.onload = () => initGoogle(clientId)
  script.onerror = () => { error.value = 'Cannot load Google Sign-In' }
  document.head.appendChild(script)
}

function initGoogle(clientId: string) {
  google.accounts.id.initialize({
    client_id: clientId,
    callback: handleCredentialResponse,
  })
  google.accounts.id.renderButton(
    document.getElementById('google-btn'),
    { theme: 'filled_black', size: 'large', width: 300, text: 'signin_with' }
  )
}

async function handleCredentialResponse(response: { credential: string }) {
  try {
    await auth.loginWithGoogle(response.credential)
    if (auth.needNickname) {
      router.push('/nickname')
    } else {
      router.push('/')
    }
  } catch (err: any) {
    error.value = 'Login failed: ' + err.message
  }
}
</script>

<template>
  <div class="login">
    <div class="login-box">
      <h1 class="title">StrongLifts 5×5</h1>
      <p class="subtitle">Barbell Strength Tracker</p>
      <div class="google-wrapper">
        <div id="google-btn"></div>
      </div>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </div>
</template>

<style scoped>
.login { display: flex; align-items: center; justify-content: center; min-height: 80vh; }
.login-box { text-align: center; display: flex; flex-direction: column; align-items: center; gap: 16px; }
.title { font-size: 32px; font-weight: 800; }
.subtitle { font-size: 14px; color: var(--text-dim); letter-spacing: 1px; }
.google-wrapper { margin-top: 24px; }
.error { color: var(--fail); font-size: 13px; margin-top: 8px; max-width: 280px; }
</style>
