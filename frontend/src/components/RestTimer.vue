<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useWorkoutStore } from '../stores/workout'

const store = useWorkoutStore()
const emit = defineEmits<{ done: [] }>()

const remaining = ref(0)
let interval: ReturnType<typeof setInterval> | null = null

function updateRemaining() {
  if (!store.timerEndTime) {
    remaining.value = 0
    return
  }
  const left = Math.ceil((store.timerEndTime - Date.now()) / 1000)
  if (left <= 0) {
    remaining.value = 0
    stop()
    beep()
    store.clearTimer()
    emit('done')
  } else {
    remaining.value = left
  }
}

function start() {
  if (interval) clearInterval(interval)
  updateRemaining()
  interval = setInterval(updateRemaining, 250)
}

function stop() {
  if (interval) { clearInterval(interval); interval = null }
}

function skip() {
  stop()
  store.clearTimer()
  emit('done')
}

function beep() {
  try {
    const ctx = new AudioContext()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.frequency.value = 880
    gain.gain.value = 0.3
    osc.start()
    osc.stop(ctx.currentTime + 0.3)
  } catch {}
}

const display = computed(() => {
  const m = Math.floor(remaining.value / 60)
  const s = remaining.value % 60
  return `${m}:${s.toString().padStart(2, '0')}`
})

const progress = computed(() => {
  if (!store.timerDuration) return 0
  return ((store.timerDuration - remaining.value) / store.timerDuration) * 100
})

onMounted(start)
onUnmounted(stop)
</script>

<template>
  <div class="timer">
    <div class="timer-ring">
      <svg viewBox="0 0 120 120" class="timer-svg">
        <circle cx="60" cy="60" r="52" class="ring-bg" />
        <circle
          cx="60" cy="60" r="52"
          class="ring-progress"
          :style="{
            strokeDasharray: `${2 * Math.PI * 52}`,
            strokeDashoffset: `${2 * Math.PI * 52 * (1 - progress / 100)}`
          }"
        />
      </svg>
      <div class="timer-display">{{ display }}</div>
    </div>
    <div class="timer-label">REST</div>
    <button class="btn btn-secondary timer-skip" @click="skip">Skip</button>
  </div>
</template>

<style scoped>
.timer { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 16px 0; }
.timer-ring { position: relative; width: 120px; height: 120px; }
.timer-svg { transform: rotate(-90deg); }
.ring-bg { fill: none; stroke: var(--surface-2); stroke-width: 6; }
.ring-progress { fill: none; stroke: var(--primary); stroke-width: 6; stroke-linecap: round; transition: stroke-dashoffset 0.3s linear; }
.timer-display { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; font-size: 28px; font-weight: 700; font-variant-numeric: tabular-nums; }
.timer-label { font-size: 12px; font-weight: 600; letter-spacing: 2px; color: var(--text-dim); }
.timer-skip { width: auto; padding: 8px 32px; font-size: 14px; }
</style>
