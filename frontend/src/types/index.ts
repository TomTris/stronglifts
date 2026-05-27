export interface User {
  id: number
  email: string
  nickname: string
  avatar_url: string
  created_at: string
}

export interface AuthResponse {
  user: User
  need_nickname: boolean
}

export interface Exercise {
  id: number
  name: string
  increment: number
}

export interface UserExercise {
  id: number
  user_id: number
  exercise_id: number
  name: string
  current_weight: number
  fail_count: number
  increment: number
}

export interface Workout {
  id: number
  user_id: number
  type: 'A' | 'B'
  date: string
  completed: boolean
  notes: string
  exercises?: WorkoutExercise[]
}

export interface WorkoutExercise {
  id: number
  workout_id: number
  exercise_id: number
  name: string
  target_sets: number
  target_reps: number
  weight: number
  sets?: WorkoutSet[]
}

export interface WorkoutSet {
  id: number
  workout_exercise_id: number
  set_number: number
  reps: number
  completed: boolean
}

export interface ProgressEntry {
  date: string
  weight: number
}

export interface PlateResult {
  total_weight: number
  per_side: Plate[] | null
}

export interface Plate {
  weight: number
  count: number
}

export interface BodyWeight {
  id: number
  user_id: number
  date: string
  weight: number
}

export interface WarmupSet {
  set_number: number
  weight: number
  reps: number
}

export interface StartingWeightInput {
  exercise_id: number
  weight: number
  reps: number
}

export interface StartingWeightResult {
  exercise_id: number
  name: string
  estimated_1rm: number
  starting_weight: number
}

export interface LeaderboardEntry {
  rank: number
  user_id: number
  nickname: string
  avatar_url: string
  squat_weight: number
  bench_weight: number
  deadlift_weight: number
  total: number
  workout_count: number
}
