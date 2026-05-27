<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { Chart, registerables } from 'chart.js'
import type { BodyWeight } from '../types'
import { api } from '../api'

Chart.register(...registerables)

const entries = ref<BodyWeight[]>([])
const newWeight = ref<number | null>(null)
const loading = ref(true)
const chartCanvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

onMounted(async () => {
  await load()
})

async function load() {
  loading.value = true
  try {
    entries.value = await api.getBodyWeights(100)
  } finally {
    loading.value = false
  }
  await nextTick()
  renderChart()
}

async function addEntry() {
  if (!newWeight.value || newWeight.value <= 0) return
  await api.addBodyWeight(newWeight.value)
  newWeight.value = null
  await load()
}

async function removeEntry(id: number) {
  await api.deleteBodyWeight(id)
  await load()
}

function renderChart() {
  if (!chartCanvas.value || entries.value.length === 0) return
  if (chart) chart.destroy()

  const sorted = [...entries.value].reverse()

  chart = new Chart(chartCanvas.value, {
    type: 'line',
    data: {
      labels: sorted.map(e => e.date),
      datasets: [{
        label: 'Body Weight',
        data: sorted.map(e => e.weight),
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59,130,246,0.1)',
        borderWidth: 2,
        pointRadius: 4,
        pointBackgroundColor: '#3b82f6',
        fill: true,
        tension: 0.3,
      }],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: {
        x: {
          ticks: { color: '#888', maxRotation: 45 },
          grid: { color: 'rgba(255,255,255,0.05)' },
        },
        y: {
          ticks: { color: '#888', callback: (v) => v + ' kg' },
          grid: { color: 'rgba(255,255,255,0.05)' },
        },
      },
    },
  })
}
</script>

<template>
  <div class="bw-view">
    <h2 class="section-title">Body Weight</h2>

    <!-- Input -->
    <div class="card input-card">
      <div class="input-row">
        <input
          type="number"
          v-model.number="newWeight"
          placeholder="kg"
          min="20"
          max="300"
          step="0.1"
          inputmode="decimal"
          @keyup.enter="addEntry"
        >
        <button class="btn btn-primary add-btn" @click="addEntry">Log</button>
      </div>
    </div>

    <!-- Chart -->
    <div class="chart-container" v-if="entries.length > 0">
      <canvas ref="chartCanvas"></canvas>
    </div>

    <!-- History -->
    <div class="history">
      <div v-if="entries.length === 0" class="empty">Chưa có dữ liệu.</div>
      <div v-for="e in entries" :key="e.id" class="card entry">
        <div class="entry-info">
          <span class="entry-weight">{{ e.weight }} kg</span>
          <span class="entry-date">{{ e.date }}</span>
        </div>
        <button class="delete-btn" @click="removeEntry(e.id)">×</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bw-view {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
}

.input-card { margin-bottom: 0; }
.input-row {
  display: flex;
  gap: 8px;
}
.input-row input {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--surface-2);
  border-radius: 8px;
  background: var(--surface-2);
  color: var(--text);
  font-size: 18px;
  font-weight: 600;
  text-align: center;
  outline: none;
}
.input-row input:focus {
  border-color: var(--primary);
}
.add-btn {
  width: auto;
  padding: 12px 24px;
}

.chart-container {
  height: 200px;
  background: var(--surface);
  border-radius: var(--radius);
  padding: 12px;
}

.history { display: flex; flex-direction: column; gap: 6px; }
.empty {
  text-align: center;
  color: var(--text-dim);
  padding: 24px 0;
}

.entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0;
  padding: 12px 16px;
}
.entry-info { display: flex; flex-direction: column; }
.entry-weight { font-size: 16px; font-weight: 700; }
.entry-date { font-size: 12px; color: var(--text-dim); }
.delete-btn {
  border: none;
  background: transparent;
  color: var(--text-dim);
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
}
.delete-btn:active { color: var(--fail); }
</style>
