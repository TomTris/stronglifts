<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { Chart, registerables } from 'chart.js'
import type { UserExercise, ProgressEntry } from '../types'
import { api } from '../api'

Chart.register(...registerables)
const exercises = ref<UserExercise[]>([])
const selectedId = ref<number>(1)
const entries = ref<ProgressEntry[]>([])
const chartCanvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

onMounted(async () => {
  exercises.value = await api.getExercises()
  if (exercises.value.length > 0) selectedId.value = exercises.value[0].exercise_id
})

watch(selectedId, async () => {
  entries.value = await api.getProgress(selectedId.value)
  await nextTick(); renderChart()
}, { immediate: false })

watch(exercises, async (val) => {
  if (val.length > 0) {
    entries.value = await api.getProgress(selectedId.value)
    await nextTick(); renderChart()
  }
}, { once: true })

function renderChart() {
  if (!chartCanvas.value) return
  if (chart) chart.destroy()
  const ex = exercises.value.find(e => e.exercise_id === selectedId.value)
  chart = new Chart(chartCanvas.value, {
    type: 'line',
    data: {
      labels: entries.value.map(e => e.date),
      datasets: [{ label: ex?.name || 'Weight', data: entries.value.map(e => e.weight), borderColor: '#e94560', backgroundColor: 'rgba(233,69,96,0.1)', borderWidth: 2, pointRadius: 4, pointBackgroundColor: '#e94560', fill: true, tension: 0.3 }],
    },
    options: {
      responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: '#888', maxRotation: 45 }, grid: { color: 'rgba(255,255,255,0.05)' } },
        y: { ticks: { color: '#888', callback: (v) => v + ' kg' }, grid: { color: 'rgba(255,255,255,0.05)' } },
      },
    },
  })
}
</script>

<template>
  <div class="progress-view">
    <h2 class="section-title">Progress</h2>
    <div class="exercise-tabs">
      <button v-for="ex in exercises" :key="ex.exercise_id" class="tab-btn" :class="{ active: selectedId === ex.exercise_id }" @click="selectedId = ex.exercise_id">{{ ex.name }}</button>
    </div>
    <div class="card stats-card" v-if="exercises.length">
      <div class="stat"><div class="stat-label">CURRENT</div><div class="stat-value">{{ exercises.find(e => e.exercise_id === selectedId)?.current_weight ?? 0 }} kg</div></div>
      <div class="stat"><div class="stat-label">FAILS</div><div class="stat-value">{{ exercises.find(e => e.exercise_id === selectedId)?.fail_count ?? 0 }}/3</div></div>
      <div class="stat"><div class="stat-label">SESSIONS</div><div class="stat-value">{{ entries.length }}</div></div>
    </div>
    <div class="chart-container">
      <canvas ref="chartCanvas"></canvas>
      <div v-if="entries.length === 0" class="chart-empty">No data yet. Complete at least 1 workout.</div>
    </div>
  </div>
</template>

<style scoped>
.progress-view { display: flex; flex-direction: column; gap: 16px; }
.section-title { font-size: 20px; font-weight: 700; }
.exercise-tabs { display: flex; gap: 6px; overflow-x: auto; -webkit-overflow-scrolling: touch; scrollbar-width: none; padding-bottom: 4px; }
.exercise-tabs::-webkit-scrollbar { display: none; }
.tab-btn { padding: 8px 14px; border: none; border-radius: 8px; background: var(--surface); color: var(--text-dim); font-size: 13px; font-weight: 600; white-space: nowrap; cursor: pointer; transition: all 0.2s; }
.tab-btn.active { background: var(--primary); color: white; }
.stats-card { display: flex; justify-content: space-around; }
.stat { text-align: center; }
.stat-label { font-size: 11px; color: var(--text-dim); text-transform: uppercase; letter-spacing: 1px; }
.stat-value { font-size: 20px; font-weight: 700; margin-top: 4px; font-variant-numeric: tabular-nums; }
.chart-container { position: relative; height: 280px; background: var(--surface); border-radius: var(--radius); padding: 16px; }
.chart-empty { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; color: var(--text-dim); font-size: 14px; text-align: center; padding: 20px; }
</style>
