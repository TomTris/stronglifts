import type {
  User, AuthResponse, UserExercise, Workout, ProgressEntry, PlateResult,
  BodyWeight, WarmupSet, StartingWeightInput, StartingWeightResult, LeaderboardEntry
} from '@/types'

const BASE = '/api'

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE}${path}`, { credentials: 'include' })
  if (res.status === 401) {
    // Don't trigger auth-expired for the auth check itself
    if (!path.startsWith('/auth/')) {
      window.dispatchEvent(new Event('auth-expired'))
    }
    throw new Error('unauthorized')
  }
  if (!res.ok) throw new Error(`GET ${path}: ${res.statusText}`)
  return res.json()
}

async function post<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST', credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (res.status === 401) {
    if (!path.startsWith('/auth/')) {
      window.dispatchEvent(new Event('auth-expired'))
    }
    throw new Error('unauthorized')
  }
  if (!res.ok) throw new Error(`POST ${path}: ${res.statusText}`)
  return res.json()
}

async function put<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    method: 'PUT', credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`PUT ${path}: ${res.statusText}`)
  return res.json()
}

async function del<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE}${path}`, { method: 'DELETE', credentials: 'include' })
  if (!res.ok) throw new Error(`DELETE ${path}: ${res.statusText}`)
  return res.json()
}

export const api = {
  // Auth
  googleAuth: (credential: string) => post<AuthResponse>('/auth/google', { credential }),
  logout: () => post<{ status: string }>('/auth/logout'),
  getMe: () => get<User>('/auth/me'),
  setNickname: (nickname: string) => post<{ status: string }>('/auth/nickname', { nickname }),

  // Setup
  getSetupStatus: () => get<{ setup_done: boolean }>('/setup/status'),
  setStartingWeights: (inputs: StartingWeightInput[]) => post<StartingWeightResult[]>('/setup/starting-weights', inputs),

  // Exercises
  getExercises: () => get<UserExercise[]>('/exercises'),
  updateWeight: (id: number, weight: number) => put<{ status: string }>(`/exercises/${id}/weight`, { weight }),

  // Workouts
  getActiveWorkout: () => get<Workout | null>('/workout/active'),
  getNextType: () => get<{ type: string }>('/workout/next-type'),
  startWorkout: () => post<Workout>('/workout/start'),
  getWorkout: (id: number) => get<Workout>(`/workout/${id}`),
  completeWorkout: (id: number) => post<{ status: string }>(`/workout/${id}/complete`),
  updateNotes: (id: number, notes: string) => put<{ status: string }>(`/workout/${id}/notes`, { notes }),
  deleteWorkout: (id: number) => del<{ status: string }>(`/workout/${id}`),

  // Sets
  completeSet: (id: number, reps: number) => post<{ status: string }>(`/set/${id}/complete`, { reps }),

  // Warmup / Plates
  getWarmupSets: (weight: number) => get<WarmupSet[]>(`/warmup/${weight}`),
  getPlates: (weight: number) => get<PlateResult>(`/plates/${weight}`),

  // Progress
  getProgress: (exerciseId: number) => get<ProgressEntry[]>(`/progress/${exerciseId}`),
  getHistory: (limit?: number) => get<Workout[]>(`/history${limit ? `?limit=${limit}` : ''}`),

  // Body weight
  getBodyWeights: (limit?: number) => get<BodyWeight[]>(`/bodyweight${limit ? `?limit=${limit}` : ''}`),
  addBodyWeight: (weight: number) => post<BodyWeight>('/bodyweight', { weight }),
  deleteBodyWeight: (id: number) => del<{ status: string }>(`/bodyweight/${id}`),

  // Leaderboard
  getLeaderboard: () => get<LeaderboardEntry[]>('/leaderboard'),

  // Export
  getExportUrl: () => `${BASE}/export`,
}
