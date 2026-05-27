<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkoutStore } from '../stores/workout'

const store = useWorkoutStore()
const router = useRouter()

const workoutAExercises = ['Squat', 'Bench Press', 'Barbell Row']
const workoutBExercises = ['Squat', 'Overhead Press', 'Deadlift']

const nextExercises = computed(() =>
  store.nextType === 'A' ? workoutAExercises : workoutBExercises
)

function getExerciseWeight(name: string): number {
  const ex = store.exercises.find(e => e.name === name)
  return ex ? ex.current_weight : 0
}

function getExerciseFails(name: string): number {
  const ex = store.exercises.find(e => e.name === name)
  return ex ? ex.fail_count : 0
}

function getSetsReps(name: string): string {
  return name === 'Deadlift' ? '1×5' : '5×5'
}

async function handleStart() {
  if (store.activeWorkout) { router.push('/workout'); return }
  await store.startWorkout()
  router.push('/workout')
}
</script>

<template>
  <div class="home">
    <div class="next-label">
      <span v-if="store.activeWorkout" class="active-badge">WORKOUT IN PROGRESS</span>
      <span v-else>NEXT WORKOUT</span>
    </div>
    <div class="workout-type">Workout {{ store.activeWorkout?.type || store.nextType }}</div>

    <div class="exercise-list">
      <div v-for="name in nextExercises" :key="name" class="card exercise-card">
        <div class="exercise-row">
          <div class="exercise-info">
            <div class="exercise-name">{{ name }}</div>
            <div class="exercise-scheme">{{ getSetsReps(name) }}</div>
          </div>
          <div class="exercise-weight-block">
            <div class="exercise-weight">{{ getExerciseWeight(name) }} kg</div>
            <div v-if="getExerciseFails(name) > 0" class="fail-indicator">
              {{ getExerciseFails(name) }}/3 fails
            </div>
          </div>
        </div>
      </div>
    </div>

    <button class="btn btn-primary" @click="handleStart" :disabled="store.loading">
      {{ store.activeWorkout ? 'Continue Workout' : 'Start Workout' }}
    </button>
  </div>
</template>

<style scoped>
.home { display: flex; flex-direction: column; gap: 16px; }
.next-label { text-align: center; font-size: 12px; font-weight: 600; letter-spacing: 2px; color: var(--text-dim); text-transform: uppercase; }
.active-badge { color: var(--warning); }
.workout-type { text-align: center; font-size: 32px; font-weight: 800; }
.exercise-list { display: flex; flex-direction: column; gap: 8px; margin: 8px 0; }
.exercise-card { margin-bottom: 0; }
.exercise-row { display: flex; justify-content: space-between; align-items: center; }
.exercise-name { font-size: 16px; font-weight: 600; }
.exercise-scheme { font-size: 13px; color: var(--text-dim); margin-top: 2px; }
.exercise-weight-block { text-align: right; }
.exercise-weight { font-size: 20px; font-weight: 700; font-variant-numeric: tabular-nums; }
.fail-indicator { font-size: 12px; color: var(--warning); margin-top: 2px; }
</style>
