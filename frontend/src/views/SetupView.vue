<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkoutStore } from '../stores/workout'
import type { StartingWeightInput, StartingWeightResult } from '../types'
import { api } from '../api'

const store = useWorkoutStore()
const router = useRouter()

const mode = ref<'choose' | 'input' | 'results'>('choose')
const results = ref<StartingWeightResult[]>([])
const submitting = ref(false)

const inputs = ref([
  { exercise_id: 1, name: 'Squat', weight: 80, reps: 5 },
  { exercise_id: 2, name: 'Bench Press', weight: 60, reps: 2 },
  { exercise_id: 3, name: 'Barbell Row', weight: 30, reps: 10 },
  { exercise_id: 4, name: 'Overhead Press', weight: 30, reps: 10 },
  { exercise_id: 5, name: 'Deadlift', weight: 92, reps: 3 },
])

const syncedInputs = computed(() => {
  if (store.exercises.length === 0) return inputs.value
  return inputs.value.map(inp => {
    const ex = store.exercises.find(e => e.name === inp.name)
    return ex ? { ...inp, exercise_id: ex.exercise_id } : inp
  })
})

async function useDefaults() {
  await api.setStartingWeights([])
  store.setupDone = true
  await store.fetchExercises()
  router.push('/')
}

async function submit() {
  submitting.value = true
  try {
    const payload: StartingWeightInput[] = syncedInputs.value
      .filter(i => i.weight > 0 && i.reps > 0)
      .map(i => ({ exercise_id: i.exercise_id, weight: i.weight, reps: i.reps }))

    if (payload.length === 0) { await useDefaults(); return }

    results.value = await api.setStartingWeights(payload)
    mode.value = 'results'
    store.setupDone = true
    await store.fetchExercises()
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="setup">
    <h1 class="title">Setup</h1>

    <div v-if="mode === 'choose'" class="choose">
      <p class="desc">Bạn đã từng tập barbell chưa?</p>
      <button class="btn btn-primary" @click="mode = 'input'">
        Đã tập — nhập số liệu hiện tại
      </button>
      <button class="btn btn-secondary" style="margin-top:12px" @click="useDefaults">
        Chưa — bắt đầu từ bar trống
      </button>
    </div>

    <div v-if="mode === 'input'" class="input-form">
      <p class="desc">
        Nhập weight và số rep bạn làm được gần nhất. Chỉnh số nếu khác. Bỏ trống = mặc định.
      </p>

      <div v-for="inp in inputs" :key="inp.name" class="card exercise-input">
        <div class="exercise-label">{{ inp.name }}</div>
        <div class="input-row">
          <div class="input-group">
            <label>KG</label>
            <input type="number" v-model.number="inp.weight" min="0" step="2.5" inputmode="decimal">
          </div>
          <div class="input-group">
            <label>REPS</label>
            <input type="number" v-model.number="inp.reps" min="0" max="30" inputmode="numeric">
          </div>
        </div>
      </div>

      <button class="btn btn-primary" @click="submit" :disabled="submitting">Tính Starting Weight</button>
    </div>

    <div v-if="mode === 'results'" class="results">
      <p class="desc">Starting weights đã được tính:</p>
      <div v-for="r in results" :key="r.exercise_id" class="card result-card">
        <div class="result-name">{{ r.name }}</div>
        <div class="result-detail">
          <span>1RM ≈ {{ r.estimated_1rm }} kg</span>
          <span class="result-arrow">→</span>
          <span class="result-start">Start: {{ r.starting_weight }} kg</span>
        </div>
      </div>
      <button class="btn btn-primary" @click="router.push('/')">Bắt đầu tập</button>
    </div>
  </div>
</template>

<style scoped>
.setup { display: flex; flex-direction: column; gap: 16px; padding-top: 24px; }
.title { text-align: center; font-size: 28px; font-weight: 800; }
.desc { font-size: 14px; color: var(--text-dim); line-height: 1.6; margin-bottom: 8px; }
.exercise-input { margin-bottom: 0; }
.exercise-label { font-weight: 600; font-size: 15px; margin-bottom: 8px; }
.input-row { display: flex; gap: 12px; }
.input-group { flex: 1; display: flex; flex-direction: column; gap: 4px; }
.input-group label { font-size: 11px; color: var(--text-dim); text-transform: uppercase; letter-spacing: 1px; }
.input-group input {
  padding: 10px 12px; border: 1px solid var(--surface-2); border-radius: 8px;
  background: var(--surface-2); color: var(--text); font-size: 16px; font-weight: 600;
  text-align: center; outline: none; width: 100%;
}
.input-group input:focus { border-color: var(--primary); }
.result-card { margin-bottom: 0; }
.result-name { font-weight: 600; font-size: 15px; }
.result-detail { display: flex; align-items: center; gap: 8px; margin-top: 4px; font-size: 13px; color: var(--text-dim); }
.result-arrow { color: var(--primary); font-weight: 700; }
.result-start { color: var(--success); font-weight: 700; }
</style>
