<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { LeaderboardEntry } from '../types'
import { api } from '../api'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const entries = ref<LeaderboardEntry[]>([])
const loading = ref(true)

onMounted(async () => {
  try { entries.value = await api.getLeaderboard() }
  finally { loading.value = false }
})

function getMedal(rank: number): string {
  if (rank === 1) return '🥇'
  if (rank === 2) return '🥈'
  if (rank === 3) return '🥉'
  return `${rank}`
}
</script>

<template>
  <div class="lb-view">
    <h2 class="section-title">Leaderboard</h2>
    <p class="section-desc">Tổng Squat + Bench + Deadlift (working weight)</p>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="entries.length === 0" class="empty">Chưa có ai trên bảng xếp hạng.</div>

    <div v-else class="lb-list">
      <div
        v-for="e in entries" :key="e.user_id"
        class="card lb-row"
        :class="{ me: e.user_id === auth.user?.id }"
      >
        <div class="lb-rank">{{ getMedal(e.rank) }}</div>

        <div class="lb-avatar" v-if="e.avatar_url">
          <img :src="e.avatar_url" :alt="e.nickname" referrerpolicy="no-referrer">
        </div>
        <div class="lb-avatar placeholder" v-else>{{ e.nickname[0]?.toUpperCase() }}</div>

        <div class="lb-info">
          <div class="lb-name">{{ e.nickname }}</div>
          <div class="lb-detail">
            S {{ e.squat_weight }} · B {{ e.bench_weight }} · D {{ e.deadlift_weight }}
          </div>
          <div class="lb-workouts">{{ e.workout_count }} workouts</div>
        </div>

        <div class="lb-total">{{ e.total }} <span class="lb-unit">kg</span></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.lb-view { display: flex; flex-direction: column; gap: 12px; }
.section-title { font-size: 20px; font-weight: 700; }
.section-desc { font-size: 12px; color: var(--text-dim); margin-top: -8px; }
.loading, .empty { text-align: center; color: var(--text-dim); padding: 32px 0; }

.lb-list { display: flex; flex-direction: column; gap: 8px; }

.lb-row {
  display: flex; align-items: center; gap: 12px;
  margin-bottom: 0; padding: 12px 14px;
  transition: border-color 0.2s;
  border: 1px solid transparent;
}
.lb-row.me { border-color: var(--primary); }

.lb-rank {
  font-size: 18px; font-weight: 800; min-width: 32px; text-align: center;
}

.lb-avatar {
  width: 36px; height: 36px; border-radius: 50%; overflow: hidden; flex-shrink: 0;
}
.lb-avatar img { width: 100%; height: 100%; object-fit: cover; }
.lb-avatar.placeholder {
  display: flex; align-items: center; justify-content: center;
  background: var(--surface-2); color: var(--text-dim);
  font-weight: 700; font-size: 16px;
}

.lb-info { flex: 1; min-width: 0; }
.lb-name { font-weight: 600; font-size: 15px; }
.lb-detail { font-size: 12px; color: var(--text-dim); margin-top: 2px; }
.lb-workouts { font-size: 11px; color: var(--text-dim); }

.lb-total { font-size: 22px; font-weight: 800; font-variant-numeric: tabular-nums; text-align: right; }
.lb-unit { font-size: 12px; font-weight: 400; color: var(--text-dim); }
</style>
