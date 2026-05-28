import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserExercise, Workout } from '@/types'
import { api } from '@/api'

export const useWorkoutStore = defineStore('workout', () => {
  const exercises = ref<UserExercise[]>([])
  const activeWorkout = ref<Workout | null>(null)
  const nextType = ref<string>('A')
  const loading = ref(false)
  const setupDone = ref(false)

  // Timer state — lives in store so it survives tab switches
  const timerEndTime = ref<number | null>(null) // unix ms when timer expires
  const timerDuration = ref(0) // original duration in seconds

  const timerActive = computed(() => {
    if (!timerEndTime.value) return false
    return Date.now() < timerEndTime.value
  })

  function startTimer(seconds: number) {
    timerDuration.value = seconds
    timerEndTime.value = Date.now() + seconds * 1000
  }

  function clearTimer() {
    timerEndTime.value = null
    timerDuration.value = 0
  }

  async function fetchExercises() {
    exercises.value = await api.getExercises()
  }

  async function fetchActiveWorkout() {
    try { activeWorkout.value = await api.getActiveWorkout() }
    catch { activeWorkout.value = null }
  }

  async function fetchNextType() {
    const res = await api.getNextType()
    nextType.value = res.type
  }

  async function checkSetup() {
    const res = await api.getSetupStatus()
    setupDone.value = res.setup_done
  }

  async function startWorkout() {
    loading.value = true
    try { activeWorkout.value = await api.startWorkout() }
    finally { loading.value = false }
  }

  async function completeSet(setId: number, reps: number) {
    await api.completeSet(setId, reps)
    if (activeWorkout.value)
      activeWorkout.value = await api.getWorkout(activeWorkout.value.id)
  }

  async function finishWorkout() {
    if (!activeWorkout.value) return
    await api.completeWorkout(activeWorkout.value.id)
    activeWorkout.value = null
    clearTimer()
    await fetchExercises()
    await fetchNextType()
  }

  async function cancelWorkout() {
    if (!activeWorkout.value) return
    await api.deleteWorkout(activeWorkout.value.id)
    activeWorkout.value = null
    clearTimer()
  }

  async function init() {
    await Promise.all([checkSetup(), fetchExercises(), fetchActiveWorkout(), fetchNextType()])
  }

  return {
    exercises, activeWorkout, nextType, loading, setupDone,
    timerEndTime, timerDuration, timerActive, startTimer, clearTimer,
    fetchExercises, fetchActiveWorkout, fetchNextType, checkSetup,
    startWorkout, completeSet, finishWorkout, cancelWorkout, init,
  }
})
