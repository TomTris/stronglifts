<script setup lang="ts">
import { ref, watch } from 'vue'
import type { PlateResult } from '../types'
import { api } from '../api'

const props = defineProps<{
  weight: number
}>()

const plates = ref<PlateResult | null>(null)

async function load() {
  try {
    plates.value = await api.getPlates(props.weight)
  } catch {
    plates.value = null
  }
}

watch(() => props.weight, load, { immediate: true })

const plateColors: Record<number, string> = {
  20: '#c0392b',
  10: '#2980b9',
  5: '#f1c40f',
  2.5: '#27ae60',
  1.25: '#95a5a6',
}
</script>

<template>
  <div class="plates" v-if="plates && plates.per_side?.length">
    <div class="plates-label">Per side:</div>
    <div class="plate-list">
      <template v-for="p in plates.per_side" :key="p.weight">
        <div
          v-for="i in p.count"
          :key="`${p.weight}-${i}`"
          class="plate"
          :style="{ background: plateColors[p.weight] || '#666', height: `${Math.max(24, p.weight * 2)}px` }"
        >
          {{ p.weight }}
        </div>
      </template>
    </div>
  </div>
  <div class="plates" v-else-if="plates">
    <div class="plates-label">Empty bar ({{ plates.total_weight }} kg)</div>
  </div>
</template>

<style scoped>
.plates {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
}

.plates-label {
  font-size: 12px;
  color: var(--text-dim);
  white-space: nowrap;
}

.plate-list {
  display: flex;
  align-items: center;
  gap: 3px;
}

.plate {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  color: white;
}
</style>
