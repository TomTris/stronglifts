<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Workout } from '../types'
import { api } from '../api'

const history = ref<Workout[]>([])
const loading = ref(true)

onMounted(async () => {
  try { history.value = await api.getHistory(50) }
  finally { loading.value = false }
})

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short', day: 'numeric', month: 'short', year: 'numeric',
  })
}
</script>

<template>
  <div class="history-view">
    <h2 class="section-title">Workout History</h2>
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="history.length === 0" class="empty">No workouts yet.</div>
    <div v-else class="history-list">
      <div v-for="w in history" :key="w.id" class="card history-item">
        <div class="history-row">
          <div>
            <div class="history-type">Workout {{ w.type }}</div>
            <div class="history-date">{{ formatDate(w.date) }}</div>
            <div v-if="w.notes" class="history-notes">{{ w.notes }}</div>
          </div>
          <div class="history-status" :class="{ completed: w.completed, incomplete: !w.completed }">
            {{ w.completed ? 'Done' : 'Incomplete' }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-view { display: flex; flex-direction: column; gap: 12px; }
.section-title { font-size: 20px; font-weight: 700; }
.loading, .empty { text-align: center; color: var(--text-dim); padding: 32px 0; }
.history-list { display: flex; flex-direction: column; gap: 8px; }
.history-item { margin-bottom: 0; }
.history-row { display: flex; justify-content: space-between; align-items: flex-start; }
.history-type { font-size: 16px; font-weight: 600; }
.history-date { font-size: 13px; color: var(--text-dim); margin-top: 2px; }
.history-notes { font-size: 12px; color: var(--text-dim); margin-top: 4px; font-style: italic; max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.history-status { font-size: 12px; font-weight: 600; padding: 4px 10px; border-radius: 6px; white-space: nowrap; }
.history-status.completed { background: var(--success); color: #000; }
.history-status.incomplete { background: var(--surface-2); color: var(--text-dim); }
</style>
