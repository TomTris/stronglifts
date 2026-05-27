<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useWorkoutStore } from './stores/workout'

const auth = useAuthStore()
const store = useWorkoutStore()
const router = useRouter()

watch(() => auth.loggedIn, async (logged) => {
  if (logged) {
    await store.init()
    if (!store.setupDone) router.push('/setup')
  }
})
</script>

<template>
  <div class="app">
    <header class="top-bar" v-if="auth.loggedIn">
      <h1>StrongLifts 5×5</h1>
    </header>
    <main class="content">
      <router-view />
    </main>
    <nav class="bottom-nav" v-if="auth.loggedIn">
      <router-link to="/" class="nav-item">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"/></svg>
        <span>Home</span>
      </router-link>
      <router-link to="/workout" class="nav-item">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="10" rx="1"/><line x1="12" y1="3" x2="12" y2="7"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
        <span>Workout</span>
      </router-link>
      <router-link to="/leaderboard" class="nav-item">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
        <span>Rank</span>
      </router-link>
      <router-link to="/progress" class="nav-item">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        <span>Progress</span>
      </router-link>
      <router-link to="/settings" class="nav-item">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
        <span>Settings</span>
      </router-link>
    </nav>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
:root {
  --bg: #0f0f1a; --surface: #1a1a2e; --surface-2: #252540;
  --primary: #e94560; --primary-dim: #c73652;
  --text: #eee; --text-dim: #888;
  --success: #4ade80; --fail: #f87171; --warning: #fbbf24;
  --radius: 12px;
}
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--bg); color: var(--text);
  -webkit-font-smoothing: antialiased; overscroll-behavior: none;
}
.app { display: flex; flex-direction: column; min-height: 100vh; max-width: 480px; margin: 0 auto; }
.top-bar { padding: 16px 20px; text-align: center; border-bottom: 1px solid var(--surface-2); }
.top-bar h1 { font-size: 18px; font-weight: 700; letter-spacing: 1px; text-transform: uppercase; }
.content { flex: 1; padding: 16px; padding-bottom: 80px; overflow-y: auto; }
.bottom-nav {
  position: fixed; bottom: 0; left: 50%; transform: translateX(-50%);
  width: 100%; max-width: 480px; display: flex;
  background: var(--surface); border-top: 1px solid var(--surface-2);
  padding: 6px 0; padding-bottom: env(safe-area-inset-bottom, 6px); z-index: 100;
}
.nav-item {
  flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px;
  text-decoration: none; color: var(--text-dim); font-size: 10px; font-weight: 500;
  padding: 4px 0; transition: color 0.2s;
}
.nav-item.router-link-exact-active, .nav-item.router-link-active { color: var(--primary); }
.btn {
  display: inline-flex; align-items: center; justify-content: center;
  padding: 14px 24px; border: none; border-radius: var(--radius);
  font-size: 16px; font-weight: 600; cursor: pointer; transition: all 0.2s; width: 100%;
}
.btn-primary { background: var(--primary); color: white; }
.btn-primary:active { background: var(--primary-dim); transform: scale(0.98); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { background: var(--surface-2); color: var(--text); }
.btn-secondary:active { background: var(--surface); }
.btn-danger { background: transparent; color: var(--fail); border: 1px solid var(--fail); }
.card { background: var(--surface); border-radius: var(--radius); padding: 16px; margin-bottom: 12px; }
</style>
