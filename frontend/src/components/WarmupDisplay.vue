<script setup lang="ts">
import { ref, watch } from 'vue'
import type { WarmupSet } from '../types'
import { api } from '../api'
import PlateCalculator from './PlateCalculator.vue'

const props = defineProps<{ workingWeight: number }>()
const warmupSets = ref<WarmupSet[]>([])
const expanded = ref(true)
const isBarWeight = ref(false)
const loadError = ref(false)

watch(() => props.workingWeight, async (w) => {
  loadError.value = false
  if (!w || w <= 20) { warmupSets.value = []; isBarWeight.value = true; return }
  isBarWeight.value = false
  try {
    const result = await api.getWarmupSets(w)
    warmupSets.value = result || []
  } catch { warmupSets.value = []; loadError.value = true }
}, { immediate: true })
</script>

<template>
  <div class="warmup" v-if="warmupSets.length > 0">
    <button class="warmup-toggle" @click="expanded = !expanded">
      <span class="warmup-header">
        <span class="warmup-badge">WARM-UP</span>
        <span class="warmup-summary">{{ warmupSets.length }} sets before working weight</span>
      </span>
      <span class="warmup-arrow" :class="{ open: expanded }">▾</span>
    </button>
    <div v-if="expanded" class="warmup-list">
      <div v-for="s in warmupSets" :key="s.set_number" class="warmup-set">
        <span class="ws-num">{{ s.set_number }}</span>
        <div class="ws-info"><span class="ws-weight">{{ s.weight }} kg</span><span class="ws-reps">× {{ s.reps }} reps</span></div>
        <PlateCalculator :weight="s.weight" />
      </div>
    </div>
  </div>
  <div class="warmup warmup-none" v-else-if="isBarWeight">
    <span class="warmup-badge">WARM-UP</span>
    <span class="warmup-note">Empty bar (20 kg) — no warm-up needed</span>
  </div>
  <div class="warmup warmup-error" v-else-if="loadError">
    <span class="warmup-note">Failed to load warm-up (weight: {{ workingWeight }} kg)</span>
  </div>
</template>

<style scoped>
.warmup { margin: 8px 0; background: rgba(233,69,96,0.06); border: 1px solid rgba(233,69,96,0.15); border-radius: 8px; padding: 10px 12px; }
.warmup-none { display: flex; align-items: center; gap: 10px; }
.warmup-error { border-color: rgba(248,113,113,0.3); background: rgba(248,113,113,0.06); }
.warmup-toggle { display: flex; align-items: center; justify-content: space-between; width: 100%; padding: 0; border: none; background: transparent; color: var(--text); cursor: pointer; }
.warmup-header { display: flex; flex-direction: column; align-items: flex-start; gap: 2px; }
.warmup-badge { font-size: 11px; font-weight: 700; letter-spacing: 1.5px; color: var(--primary); flex-shrink: 0; }
.warmup-summary { font-size: 12px; color: var(--text-dim); }
.warmup-note { font-size: 12px; color: var(--text-dim); }
.warmup-arrow { font-size: 14px; color: var(--text-dim); transition: transform 0.2s; }
.warmup-arrow.open { transform: rotate(180deg); }
.warmup-list { display: flex; flex-direction: column; gap: 6px; margin-top: 10px; padding-top: 10px; border-top: 1px solid rgba(233,69,96,0.1); }
.warmup-set { display: flex; align-items: center; gap: 10px; padding: 8px 10px; background: var(--surface-2); border-radius: 6px; font-size: 14px; flex-wrap: wrap; }
.ws-num { width: 24px; height: 24px; display: flex; align-items: center; justify-content: center; background: var(--surface); border-radius: 50%; font-size: 12px; font-weight: 700; color: var(--text-dim); flex-shrink: 0; }
.ws-info { display: flex; align-items: center; gap: 6px; min-width: 100px; }
.ws-weight { font-weight: 700; }
.ws-reps { color: var(--text-dim); font-size: 13px; }
</style>
