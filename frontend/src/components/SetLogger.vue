<script setup lang="ts">
import { ref, computed } from 'vue'
import type { WorkoutExercise } from '../types'
import PlateCalculator from './PlateCalculator.vue'
import WarmupDisplay from './WarmupDisplay.vue'

const props = defineProps<{
  exercise: WorkoutExercise
}>()

const emit = defineEmits<{
  completeSet: [setId: number, reps: number]
}>()

const editingSetId = ref<number | null>(null)
const editReps = ref(5)

const allDone = computed(() =>
  props.exercise.sets?.every(s => s.completed) ?? false
)

const successCount = computed(() =>
  props.exercise.sets?.filter(s => s.completed && s.reps >= props.exercise.target_reps).length ?? 0
)

const isSuccess = computed(() =>
  allDone.value && successCount.value >= props.exercise.target_sets
)

function tapSet(setId: number, currentReps: number, completed: boolean) {
  if (completed) return
  editReps.value = props.exercise.target_reps
  editingSetId.value = setId
}

function confirmSet() {
  if (editingSetId.value !== null) {
    emit('completeSet', editingSetId.value, editReps.value)
    editingSetId.value = null
  }
}

function adjustReps(delta: number) {
  editReps.value = Math.max(0, Math.min(editReps.value + delta, props.exercise.target_reps))
}
</script>

<template>
  <div class="set-logger card" :class="{ success: isSuccess, failed: allDone && !isSuccess }">
    <div class="exercise-header">
      <div>
        <div class="exercise-name">{{ exercise.name }}</div>
        <div class="exercise-detail">
          {{ exercise.target_sets }}×{{ exercise.target_reps }} @ {{ exercise.weight }} kg
        </div>
      </div>
      <div v-if="allDone" class="status-badge" :class="{ success: isSuccess, fail: !isSuccess }">
        {{ isSuccess ? 'DONE' : 'FAIL' }}
      </div>
    </div>

    <WarmupDisplay :workingWeight="exercise.weight" />

    <PlateCalculator :weight="exercise.weight" />

    <div class="sets-grid">
      <button
        v-for="set in exercise.sets"
        :key="set.id"
        class="set-btn"
        :class="{
          done: set.completed && set.reps >= exercise.target_reps,
          partial: set.completed && set.reps < exercise.target_reps,
          active: editingSetId === set.id,
          pending: !set.completed && editingSetId !== set.id,
        }"
        @click="tapSet(set.id, set.reps, set.completed)"
      >
        <template v-if="set.completed">
          {{ set.reps }}
        </template>
        <template v-else>
          {{ set.set_number }}
        </template>
      </button>
    </div>

    <div v-if="editingSetId !== null" class="rep-editor">
      <div class="rep-label">Reps completed:</div>
      <div class="rep-controls">
        <button class="rep-btn" @click="adjustReps(-1)">−</button>
        <span class="rep-value">{{ editReps }}</span>
        <button class="rep-btn" @click="adjustReps(1)">+</button>
      </div>
      <button class="btn btn-primary confirm-btn" @click="confirmSet">
        Log Set
      </button>
    </div>
  </div>
</template>

<style scoped>
.set-logger {
  transition: border-color 0.3s;
  border: 1px solid transparent;
}
.set-logger.success { border-color: var(--success); }
.set-logger.failed { border-color: var(--fail); }

.exercise-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.exercise-name { font-size: 18px; font-weight: 700; }
.exercise-detail { font-size: 13px; color: var(--text-dim); margin-top: 2px; }

.status-badge {
  font-size: 11px; font-weight: 700;
  padding: 4px 8px; border-radius: 6px; letter-spacing: 1px;
}
.status-badge.success { background: var(--success); color: #000; }
.status-badge.fail { background: var(--fail); color: #000; }

.sets-grid { display: flex; gap: 8px; margin-top: 12px; }

.set-btn {
  flex: 1; height: 48px; border: none; border-radius: 8px;
  font-size: 16px; font-weight: 700; cursor: pointer; transition: all 0.2s;
}
.set-btn.pending { background: var(--surface-2); color: var(--text-dim); }
.set-btn.pending:active { transform: scale(0.95); }
.set-btn.active { background: var(--primary); color: white; transform: scale(1.05); }
.set-btn.done { background: var(--success); color: #000; }
.set-btn.partial { background: var(--warning); color: #000; }

.rep-editor {
  margin-top: 16px; padding-top: 16px; border-top: 1px solid var(--surface-2);
  display: flex; flex-direction: column; align-items: center; gap: 12px;
}
.rep-label { font-size: 13px; color: var(--text-dim); }
.rep-controls { display: flex; align-items: center; gap: 20px; }
.rep-btn {
  width: 44px; height: 44px; border: 1px solid var(--surface-2); border-radius: 50%;
  background: var(--surface-2); color: var(--text); font-size: 20px; font-weight: 700;
  cursor: pointer; display: flex; align-items: center; justify-content: center;
}
.rep-btn:active { background: var(--primary); }
.rep-value {
  font-size: 32px; font-weight: 800; min-width: 48px;
  text-align: center; font-variant-numeric: tabular-nums;
}
.confirm-btn { max-width: 200px; }
</style>
