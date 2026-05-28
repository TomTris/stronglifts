<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkoutStore } from '../stores/workout'
import { api } from '../api'
import SetLogger from '../components/SetLogger.vue'
import RestTimer from '../components/RestTimer.vue'

const store = useWorkoutStore()
const router = useRouter()

const showConfirmFinish = ref(false)
const showConfirmCancel = ref(false)
const notes = ref('')
const notesSaved = ref(false)

const workout = computed(() => store.activeWorkout)

if (workout.value?.notes) {
  notes.value = workout.value.notes
}

const allExercisesDone = computed(() =>
  workout.value?.exercises?.every(ex =>
    ex.sets?.every(s => s.completed)
  ) ?? false
)

async function handleCompleteSet(setId: number, reps: number) {
  const exercise = workout.value?.exercises?.find(ex =>
    ex.sets?.some(s => s.id === setId)
  )
  if (!exercise) return

  await store.completeSet(setId, reps)

  let duration = 90
  if (reps >= exercise.target_reps) {
    duration = 90
  } else if (reps >= 3) {
    duration = 180
  } else {
    duration = 300
  }

  const updatedExercise = store.activeWorkout?.exercises?.find(ex => ex.id === exercise.id)
  const allSetsDone = updatedExercise?.sets?.every(s => s.completed) ?? false

  if (!allSetsDone) {
    store.startTimer(duration)
  }
}

function timerDone() {
  // Timer already cleared itself via store.clearTimer()
}

let notesTimer: ReturnType<typeof setTimeout> | null = null
function onNotesInput() {
  notesSaved.value = false
  if (notesTimer) clearTimeout(notesTimer)
  notesTimer = setTimeout(async () => {
    if (workout.value) {
      await api.updateNotes(workout.value.id, notes.value)
      notesSaved.value = true
    }
  }, 1000)
}

async function finishWorkout() {
  if (workout.value && notes.value) {
    await api.updateNotes(workout.value.id, notes.value)
  }
  await store.finishWorkout()
  router.push('/')
}

async function cancelWorkout() {
  await store.cancelWorkout()
  router.push('/')
}
</script>

<template>
  <div class="workout-view" v-if="workout">
    <div class="workout-header">
      <div class="workout-title">Workout {{ workout.type }}</div>
      <div class="workout-date">
        {{ new Date(workout.date).toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'short' }) }}
      </div>
    </div>

    <RestTimer v-if="store.timerActive" @done="timerDone" />

    <div class="exercises" v-show="!store.timerActive">
      <SetLogger
        v-for="ex in workout.exercises"
        :key="ex.id"
        :exercise="ex"
        @completeSet="handleCompleteSet"
      />
    </div>

    <div class="card notes-card" v-show="!store.timerActive">
      <label class="notes-label">
        Notes
        <span v-if="notesSaved" class="saved-indicator">saved</span>
      </label>
      <textarea
        v-model="notes"
        @input="onNotesInput"
        placeholder="Session notes (back pain, poor sleep, ...)"
        rows="3"
      ></textarea>
    </div>

    <div class="actions" v-show="!store.timerActive">
      <button
        v-if="allExercisesDone && !showConfirmFinish"
        class="btn btn-primary"
        @click="showConfirmFinish = true"
      >
        Finish Workout
      </button>

      <div v-if="showConfirmFinish" class="confirm-box card">
        <p>Complete workout? Weights will update automatically.</p>
        <div class="confirm-actions">
          <button class="btn btn-primary" @click="finishWorkout">Confirm</button>
          <button class="btn btn-secondary" @click="showConfirmFinish = false">Cancel</button>
        </div>
      </div>

      <button
        v-if="!showConfirmCancel"
        class="btn btn-danger"
        style="margin-top: 8px;"
        @click="showConfirmCancel = true"
      >
        Cancel Workout
      </button>

      <div v-if="showConfirmCancel" class="confirm-box card">
        <p>Cancel workout? All data from this session will be deleted.</p>
        <div class="confirm-actions">
          <button class="btn btn-danger" @click="cancelWorkout">Delete</button>
          <button class="btn btn-secondary" @click="showConfirmCancel = false">Keep</button>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="no-workout">
    <p>No active workout.</p>
    <button class="btn btn-primary" @click="router.push('/')">Go to Home</button>
  </div>
</template>

<style scoped>
.workout-view { display: flex; flex-direction: column; gap: 16px; }
.workout-header { text-align: center; }
.workout-title { font-size: 24px; font-weight: 800; }
.workout-date { font-size: 13px; color: var(--text-dim); margin-top: 4px; }
.exercises { display: flex; flex-direction: column; gap: 12px; }
.notes-card { margin-bottom: 0; }
.notes-label {
  display: flex; justify-content: space-between; align-items: center;
  font-size: 12px; font-weight: 600; color: var(--text-dim);
  text-transform: uppercase; letter-spacing: 1px; margin-bottom: 8px;
}
.saved-indicator { color: var(--success); font-size: 11px; text-transform: none; }
textarea {
  width: 100%; padding: 10px; border: 1px solid var(--surface-2);
  border-radius: 8px; background: var(--surface-2); color: var(--text);
  font-size: 14px; font-family: inherit; resize: vertical; outline: none;
}
textarea:focus { border-color: var(--primary); }
.actions { margin-top: 8px; }
.confirm-box { margin-top: 8px; }
.confirm-box p { font-size: 14px; margin-bottom: 12px; }
.confirm-actions { display: flex; gap: 8px; }
.confirm-actions .btn { flex: 1; padding: 10px; font-size: 14px; }
.no-workout { text-align: center; padding: 48px 0; }
.no-workout p { color: var(--text-dim); margin-bottom: 16px; }
</style>
