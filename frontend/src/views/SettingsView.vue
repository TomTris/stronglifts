<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useWorkoutStore } from '../stores/workout'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const store = useWorkoutStore()
const auth = useAuthStore()
const router = useRouter()

function exportData() { window.open(api.getExportUrl(), '_blank') }

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="settings-view">
    <h2 class="section-title">Settings</h2>

    <!-- User info -->
    <div class="card user-card" v-if="auth.user">
      <div class="user-row">
        <img v-if="auth.user.avatar_url" :src="auth.user.avatar_url" class="user-avatar" referrerpolicy="no-referrer">
        <div class="user-info">
          <div class="user-name">{{ auth.user.nickname || auth.user.email }}</div>
          <div class="user-email">{{ auth.user.email }}</div>
        </div>
      </div>
    </div>

    <div class="card setting-item" @click="router.push('/setup')">
      <div class="setting-label">Chỉnh Starting Weight</div>
      <div class="setting-desc">Nhập lại số liệu hiện tại</div>
    </div>

    <div class="card setting-item" @click="router.push('/history')">
      <div class="setting-label">Workout History</div>
      <div class="setting-desc">Xem tất cả buổi tập trước</div>
    </div>

    <div class="card setting-item" @click="router.push('/bodyweight')">
      <div class="setting-label">Body Weight</div>
      <div class="setting-desc">Theo dõi cân nặng</div>
    </div>

    <div class="card setting-item" @click="exportData">
      <div class="setting-label">Export / Backup</div>
      <div class="setting-desc">Tải dữ liệu dạng JSON</div>
    </div>

    <h3 class="sub-title">Current Weights</h3>
    <div v-for="ex in store.exercises" :key="ex.id" class="card exercise-card">
      <div class="ex-row">
        <div class="ex-name">{{ ex.name }}</div>
        <div class="ex-weight">{{ ex.current_weight }} kg</div>
      </div>
      <div class="ex-meta">
        <span>+{{ ex.increment }} kg/session</span>
        <span v-if="ex.fail_count > 0" class="fail-text">{{ ex.fail_count }}/3 fails</span>
      </div>
    </div>

    <button class="btn btn-danger" style="margin-top: 16px;" @click="handleLogout">Logout</button>
  </div>
</template>

<style scoped>
.settings-view { display: flex; flex-direction: column; gap: 12px; }
.section-title { font-size: 20px; font-weight: 700; }
.sub-title { font-size: 16px; font-weight: 600; margin-top: 8px; }

.user-card { margin-bottom: 0; }
.user-row { display: flex; align-items: center; gap: 12px; }
.user-avatar { width: 44px; height: 44px; border-radius: 50%; }
.user-info { flex: 1; }
.user-name { font-weight: 600; font-size: 16px; }
.user-email { font-size: 12px; color: var(--text-dim); }

.setting-item { cursor: pointer; margin-bottom: 0; transition: background 0.2s; }
.setting-item:active { background: var(--surface-2); }
.setting-label { font-weight: 600; font-size: 15px; }
.setting-desc { font-size: 13px; color: var(--text-dim); margin-top: 2px; }
.exercise-card { margin-bottom: 0; }
.ex-row { display: flex; justify-content: space-between; align-items: center; }
.ex-name { font-weight: 600; }
.ex-weight { font-size: 18px; font-weight: 700; font-variant-numeric: tabular-nums; }
.ex-meta { display: flex; gap: 12px; margin-top: 4px; font-size: 12px; color: var(--text-dim); }
.fail-text { color: var(--warning); }
</style>
