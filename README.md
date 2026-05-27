# StrongLifts 5×5 Tracker

Self-hosted workout tracker for the StrongLifts 5×5 program. Multi-user with Google login, leaderboard, PWA. Single binary, runs anywhere.

**Status:** Not deployed. Running locally.

---
## How to Run

### Requirements

- Go 1.22+
- Google OAuth Client ID (free, 5-minute setup)

### Create a Google OAuth Client ID

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Create or select a project
3. APIs & Services → Credentials → Create Credentials → **OAuth client ID**
4. Application type: **Web application**
5. Authorized JavaScript origins: add `http://localhost:8080`
6. Create → copy the Client ID

If the OAuth consent screen isn't configured yet:
- APIs & Services → OAuth consent screen → External
- Fill in App name and email → Save
- No additional scopes needed

### Configuration

Edit `.env`:

```env
GOOGLE_CLIENT_ID=xxxx.apps.googleusercontent.com
PORT=8080
DB_PATH=stronglifts.db
```

### Build and Run

```bash
make run
```

Open `http://localhost:8080`. Database is created automatically. Sign in with Google, pick a nickname, set up weights, start training.

### Use on Phone (Same WiFi)

Find your local IP:
```bash
ip addr | grep "inet " | grep -v 127.0.0.1
```

Open Chrome on your phone: `http://192.168.x.x:8080`

Only works when on the same WiFi network and the computer is running the server.

## What is StrongLifts 5×5

A barbell strength training program for beginners or anyone returning after a break. Five compound exercises, three sessions per week. Core principle: **progressive overload** — add weight every session, forcing the body to adapt continuously.

### The 5 Exercises

| Exercise | Primary Muscles |
|----------|----------------|
| Squat | Quads, glutes, core |
| Bench Press | Chest, front delts, triceps |
| Barbell Row | Back, biceps |
| Overhead Press | Shoulders, triceps |
| Deadlift | Lower back, glutes, hamstrings, full body |

### 2 Alternating Workouts

| Workout A | Workout B |
|-----------|-----------|
| Squat 5×5 | Squat 5×5 |
| Bench Press 5×5 | Overhead Press 5×5 |
| Barbell Row 5×5 | Deadlift 1×5 |

Three sessions per week, rest days in between. Week 1: A-B-A. Week 2: B-A-B. Repeat.

Squat appears every session. Deadlift is only 1 set × 5 reps because it's the heaviest lift and requires more recovery time.

### Progression Rules

- Complete all 5 reps for all sets → next session **+2.5 kg** (Deadlift +5 kg).
- Fail to hit 5 reps on any set → keep the same weight next session. Still do all 5 sets.
- Fail at the same weight **3 consecutive sessions** → **deload 10%** (reduce weight, work back up).

### Where to Start

Never lifted before: empty bar (20 kg). Add 2.5 kg every session.

Have lifting experience: enter your current weight × reps → app calculates 1RM → starts you at 50%.

### Warm-up Before Each Exercise

Before working weight, perform 2-4 progressively heavier warm-up sets:

| Set | Weight | Reps |
|-----|--------|------|
| 1 | Empty bar (20 kg) | 5 |
| 2 | 40% of working weight | 5 |
| 3 | 60% of working weight | 3 |
| 4 | 80% of working weight | 2 |

Only shows sets where weight > 20 kg and < working weight.

---

## What the App Does

Replaces the original StrongLifts app (ads, paid premium). Self-hosted, your data stays with you.

### Login & Multi-user

Sign in with Google. Each user has their own data. Choose a nickname after first login — displayed on the leaderboard.

### Setup Wizard

First time opening the app: choose "I've lifted before" or "I'm new." If experienced, the form comes pre-filled with sample values — adjust to your actual ability, app calculates starting weights automatically. If new, start from empty bar.

### Workout Logging

- Home screen shows next workout (A or B), exercise list with weights.
- Tap Start Workout.
- Each exercise shows warm-up sets (expanded by default) with plates needed.
- Tap a set circle → adjust reps → Log Set.
- Rest timer starts automatically: 90 seconds (all reps completed), 3 minutes (struggled), 5 minutes (heavy fail). Audio alert when done.
- Add notes if needed (back pain, poor sleep, etc.). Auto-saves.
- Finish Workout → weights update automatically for next session.

### Plate Calculator

Each exercise shows a visual of plates to load on each side of the bar, color-coded: 20 kg red, 10 kg blue, 5 kg yellow, 2.5 kg green, 1.25 kg gray.

### Progress Charts

Line chart of weight over time for each exercise. Switch between exercises via tabs.

### Body Weight

Log daily body weight, chart over time.

### Leaderboard

Rank tab — all users ranked by total of Squat + Bench Press + Deadlift working weights. Shows avatar, nickname, per-exercise breakdown, total workout count.

### Export / Backup

Settings → Export → JSON file containing all exercises, workouts, and body weights.

### PWA

On Android: Chrome → "Add to Home screen" → runs like a native app. Static assets cached offline (API calls still need server connectivity).



# Architecture
# StrongLifts 5×5 — Technical Documentation

---

## Architecture

Single binary. Go backend embeds the Vue frontend. One process serves both the API and static files. No reverse proxy, no containers, no microservices.

```
Browser (phone/desktop)
    │
    │  HTTP
    │
Go server (:8080)
    ├── /api/*        → REST handlers
    ├── /             → Vue SPA (embedded)
    └── SQLite file   → stronglifts.db
```

---

## Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Backend | Go 1.22+ | Single binary output, `go:embed` bundles frontend, `net/http` sufficient for REST, no framework needed |
| Database | SQLite via `modernc.org/sqlite` | Pure Go driver — no CGO required, cross-compiles cleanly, zero config, single file |
| Auth | Google Identity Services | Server-side ID token verification, no complex OAuth flow, no stored passwords |
| Session | Random token + cookie | 64-char hex, HttpOnly, SameSite=Lax, 30-day expiry. Server-side session table |
| Frontend | Vue 3 + TypeScript | Composition API, type-safe, reactive. Pinia for state, Vue Router for SPA |
| Charts | Chart.js 4 | Lightweight, responsive, sufficient for line charts |
| Build tool | Vite 5 | Fast builds, good HMR, dev proxy |
| PWA | Service Worker + Manifest | Static asset caching, installable on Android/iOS |

---

## Project Structure

```
stronglifts/
├── main.go                     # Entry point
│                                 - .env loader (built-in, no library)
│                                 - go:embed frontend/dist/*
│                                 - SPA fallback routing
│                                 - HTTP server start
│
├── go.mod                      # modernc.org/sqlite dependency
├── .env.example                # Config template
├── .gitignore
├── Makefile
│
├── models/
│   └── models.go               # All shared structs
│                                 - User, Session
│                                 - Exercise, UserExercise
│                                 - Workout, WorkoutExercise, WorkoutSet
│                                 - BodyWeight, ProgressEntry
│                                 - PlateResult, WarmupSet
│                                 - LeaderboardEntry
│                                 - StartingWeightInput/Result
│                                 - AuthResponse, ExportData
│
├── db/
│   └── db.go                   # Database layer (single file, ~500 lines)
│                                 - Schema migration (CREATE IF NOT EXISTS)
│                                 - Seed data (5 exercises)
│                                 - Auth: FindOrCreateUser, sessions, nickname
│                                 - User exercises: CRUD, starting weight calc
│                                 - Workouts: create, complete, progression logic
│                                 - Warmup: percentage-based calculation
│                                 - Body weight: CRUD
│                                 - Leaderboard: aggregate query
│                                 - Export: full data dump
│
├── handlers/
│   └── handlers.go             # HTTP layer (single file, ~400 lines)
│                                 - Route registration
│                                 - Auth middleware (cookie/Bearer token)
│                                 - 24 endpoints
│                                 - JSON request/response helpers
│
└── frontend/
    ├── index.html              # HTML entry, PWA meta tags
    ├── package.json            # Vue, Pinia, Vue Router, Chart.js
    ├── vite.config.ts          # Path alias @/, dev proxy /api → :8080
    ├── tsconfig.json
    │
    ├── public/                 # Copied to dist/ as-is
    │   ├── manifest.json       # PWA manifest
    │   ├── sw.js               # Service worker (cache-first static, network-first API)
    │   ├── icon-192.png
    │   └── icon-512.png
    │
    ├── dist/                   # Pre-built, embedded into Go binary
    │
    └── src/
        ├── main.ts             # App bootstrap
        │                         - Router setup (10 routes)
        │                         - Auth guard (check once, redirect /login)
        │                         - auth-expired event listener
        │                         - Service worker registration
        │
        ├── App.vue             # Root component
        │                         - Top bar, bottom nav (5 tabs)
        │                         - Nav hidden when not logged in
        │                         - Watch auth → init workout store
        │
        ├── env.d.ts            # Vue/Vite type declarations
        │
        ├── api/
        │   └── index.ts        # HTTP client
        │                         - fetch wrapper with credentials: include
        │                         - 401 handling (dispatch event, skip for /auth/*)
        │                         - get/post/put/del helpers
        │                         - All 24 endpoint functions
        │
        ├── stores/
        │   ├── auth.ts         # Auth state
        │   │                     - user, loggedIn, needNickname
        │   │                     - checkAuth, loginWithGoogle, setNickname, logout
        │   │
        │   └── workout.ts      # Workout state
        │                         - exercises, activeWorkout, nextType, setupDone
        │                         - CRUD actions, init()
        │
        ├── types/
        │   └── index.ts        # TypeScript interfaces (mirror Go models)
        │
        ├── views/
        │   ├── LoginView.vue       # Google Sign-In button
        │   │                         - Fetch client ID from /api/config
        │   │                         - Load Google Identity Services SDK
        │   │                         - Handle credential → loginWithGoogle
        │   │
        │   ├── NicknameView.vue    # Post-login nickname entry
        │   │
        │   ├── HomeView.vue        # Dashboard
        │   │                         - Next workout type (A/B)
        │   │                         - Exercise list with current weights
        │   │                         - Start/Continue button
        │   │
        │   ├── SetupView.vue       # Starting weight wizard
        │   │                         - Pre-filled form values
        │   │                         - 1RM calculation results
        │   │
        │   ├── WorkoutView.vue     # Active workout
        │   │                         - Exercise list with SetLogger components
        │   │                         - Rest timer integration
        │   │                         - Notes textarea (auto-save)
        │   │                         - Finish/Cancel with confirmation
        │   │
        │   ├── ProgressView.vue    # Charts per exercise
        │   │                         - Tab selector for exercises
        │   │                         - Chart.js line chart
        │   │                         - Current stats (weight, fails, sessions)
        │   │
        │   ├── BodyWeightView.vue  # Body weight log
        │   │
        │   ├── HistoryView.vue     # Past workouts list
        │   │
        │   ├── LeaderboardView.vue # Rankings
        │   │                         - Medals for top 3
        │   │                         - S/B/D breakdown
        │   │                         - Highlight current user
        │   │
        │   └── SettingsView.vue    # User info, links, export, logout
        │
        └── components/
            ├── SetLogger.vue       # Per-exercise workout UI
            │                         - Sets grid (tap to log)
            │                         - Rep editor (+/- buttons)
            │                         - Success/fail status badge
            │                         - Integrates WarmupDisplay + PlateCalculator
            │
            ├── WarmupDisplay.vue   # Warm-up sets
            │                         - Expanded by default
            │                         - Percentage-based calculation
            │                         - Bar-weight-only message
            │                         - Plate display per warmup set
            │
            ├── RestTimer.vue       # Countdown timer
            │                         - SVG ring animation
            │                         - Audio beep on complete
            │                         - Skip button
            │
            └── PlateCalculator.vue # Plate visual
                                      - Color-coded per weight
                                      - Height scaled to plate weight
```

---

## Database Schema

SQLite, single file, 8 tables. WAL mode enabled.

### Entity Relationship

```
users ──┬── sessions
        ├── user_exercises ──── exercises (template)
        ├── workouts ──── workout_exercises ──── workout_sets
        │                       └── exercises
        └── body_weights

app_settings (key-value, standalone)
```

### Tables

**users** — Google-authenticated accounts.

| Column | Type | Constraint |
|--------|------|------------|
| id | INTEGER | PK, autoincrement |
| google_id | TEXT | UNIQUE, NOT NULL |
| email | TEXT | NOT NULL |
| nickname | TEXT | Default '' |
| avatar_url | TEXT | Default '' |
| created_at | TEXT | ISO 8601 |

**sessions** — Server-side session tokens.

| Column | Type | Constraint |
|--------|------|------------|
| token | TEXT | PK, 64-char hex |
| user_id | INTEGER | FK → users.id |
| created_at | TEXT | |
| expires_at | TEXT | created_at + 30 days |

**exercises** — Template table. 5 rows, seeded on first run. Never modified per user.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| name | TEXT | UNIQUE (Squat, Bench Press, ...) |
| default_weight | REAL | Starting weight for new users |
| increment | REAL | 2.5 or 5.0 |

**user_exercises** — Per-user exercise state. Created when a user first logs in.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| user_id | INTEGER | FK → users.id |
| exercise_id | INTEGER | FK → exercises.id |
| current_weight | REAL | Current working weight |
| fail_count | INTEGER | 0-2, resets on success or deload |
| | | UNIQUE(user_id, exercise_id) |

**workouts** — One row per training session.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| user_id | INTEGER | FK → users.id |
| type | TEXT | CHECK('A','B') |
| date | TEXT | ISO 8601 |
| completed | INTEGER | 0 or 1 |
| notes | TEXT | Free-form |

**workout_exercises** — Snapshot of exercises in a workout.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| workout_id | INTEGER | FK → workouts.id |
| exercise_id | INTEGER | FK → exercises.id |
| target_sets | INTEGER | 5 (or 1 for Deadlift) |
| target_reps | INTEGER | 5 |
| weight | REAL | Frozen at workout creation time |

**workout_sets** — Individual set results.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| workout_exercise_id | INTEGER | FK → workout_exercises.id |
| set_number | INTEGER | 1-5 |
| reps | INTEGER | Actual reps performed |
| completed | INTEGER | 0 or 1 |

**body_weights** — Daily weight log.

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| user_id | INTEGER | FK → users.id |
| date | TEXT | YYYY-MM-DD |
| weight | REAL | kg |

**app_settings** — Key-value store. Currently only `setup_done_{userID}`.

---

## Auth Flow

```
 1. Browser loads /login
 2. Frontend fetches GET /api/config → gets GOOGLE_CLIENT_ID
 3. Google Identity Services SDK loads from accounts.google.com
 4. User clicks "Sign in with Google" → Google popup
 5. Google returns ID token (JWT) to frontend callback
 6. Frontend POSTs token to /api/auth/google
 7. Backend verifies token via Google tokeninfo endpoint
 8. Backend creates or finds user in DB (upsert by google_id)
 9. Backend creates session (crypto/rand 32 bytes → 64-char hex, 30-day expiry)
10. Backend sets HttpOnly cookie "session" (SameSite=Lax, path=/)
11. Frontend redirects to /nickname (new user) or / (returning user)
12. All subsequent API calls include cookie automatically
13. Auth middleware extracts cookie → looks up session → injects userID into context
```

Session validation: token exists in sessions table + current time < expires_at. Expired sessions are deleted on access.

---

## Progression Algorithm

Runs when `POST /api/workout/{id}/complete` is called. Executes inside a single SQLite transaction.

```
for each exercise in the workout:
    count = sets where (completed = true AND reps >= target_reps)

    if count >= target_sets:
        # Success
        current_weight += increment    (2.5 or 5.0)
        fail_count = 0

    else:
        # Fail
        fail_count += 1

        if fail_count >= 3:
            # Deload
            current_weight = floor(current_weight × 0.9 / 2.5) × 2.5
            if current_weight < 20: current_weight = 20
            fail_count = 0
```

---

## Warmup Calculation

```
Input:  working_weight (float64)
Output: []WarmupSet

if working_weight <= 20:
    return empty    # Already at bar weight, nothing lighter to warm up with

sets = [{weight: 20, reps: 5}]    # Always: empty bar

for (pct, reps) in [(0.4, 5), (0.6, 3), (0.8, 2)]:
    w = round(working_weight × pct / 2.5) × 2.5
    if w > 20 AND w < working_weight:
        sets.append({weight: w, reps: reps})

return sets
```

---

## Starting Weight Calculation

Epley formula:

```
1RM = weight × (1 + reps / 30)
starting_weight = floor(1RM × 0.5 / increment) × increment
minimum = 20 kg
```

Example: Squat 80 kg × 6 reps → 1RM = 80 × 1.2 = 96 kg → start = floor(48 / 2.5) × 2.5 = 47.5 kg.

---

## Leaderboard Query

```sql
SELECT
    u.id, u.nickname, u.avatar_url,
    squat.current_weight,
    bench.current_weight,
    deadlift.current_weight,
    (SELECT COUNT(*) FROM workouts WHERE user_id = u.id AND completed = 1)
FROM users u
WHERE u.nickname != ''
ORDER BY (squat + bench + deadlift) DESC
```

Subqueries join `user_exercises` with `exercises` by name. Only users with a nickname appear.

---

## API Endpoints

24 endpoints total. 5 public, 19 require session cookie.

### Public

| Method | Path | Purpose |
|--------|------|---------|
| GET | /api/config | Returns google_client_id |
| POST | /api/auth/google | Exchange Google ID token for session |
| GET | /api/warmup/{weight} | Calculate warmup sets |
| GET | /api/plates/{weight} | Calculate plate breakdown |
| GET | /api/leaderboard | All users ranked by total |

### Authenticated

| Method | Path | Purpose |
|--------|------|---------|
| GET | /api/auth/me | Current user info |
| POST | /api/auth/nickname | Set nickname (1-20 chars) |
| POST | /api/auth/logout | Clear session |
| GET | /api/setup/status | Check if setup wizard completed |
| POST | /api/setup/starting-weights | Calculate and set starting weights |
| GET | /api/exercises | User's 5 exercises with current weights |
| PUT | /api/exercises/{id}/weight | Manual weight override |
| GET | /api/workout/active | Current incomplete workout or null |
| GET | /api/workout/next-type | "A" or "B" |
| POST | /api/workout/start | Create new workout with sets |
| GET | /api/workout/{id} | Full workout detail |
| POST | /api/workout/{id}/complete | Finish workout, run progression |
| PUT | /api/workout/{id}/notes | Update notes |
| DELETE | /api/workout/{id} | Cancel/delete workout |
| POST | /api/set/{id}/complete | Log a set with reps |
| GET | /api/progress/{exerciseId} | Weight history for charting |
| GET | /api/history | Past workouts list |
| GET | /api/bodyweight | Body weight history |
| POST | /api/bodyweight | Log today's weight |
| DELETE | /api/bodyweight/{id} | Delete entry |
| GET | /api/export | Full JSON data dump |

---

## Frontend Routing

| Path | View | Auth | Purpose |
|------|------|------|---------|
| /login | LoginView | Public | Google Sign-In |
| /nickname | NicknameView | Required | First-time nickname |
| / | HomeView | Required | Dashboard, start workout |
| /workout | WorkoutView | Required | Active workout logging |
| /setup | SetupView | Required | Starting weight wizard |
| /progress | ProgressView | Required | Charts per exercise |
| /bodyweight | BodyWeightView | Required | Body weight log |
| /history | HistoryView | Required | Past workouts |
| /leaderboard | LeaderboardView | Required | Rankings |
| /settings | SettingsView | Required | User info, export, logout |

Auth guard runs once on app load (`authChecked` flag prevents repeated API calls). Failed check → redirect to `/login`. Success → initialize workout store.

---

## .env Loader

Built-in, no external library. ~25 lines in `main.go`.

- Reads `.env` from working directory (optional — no error if missing)
- Parses `KEY=VALUE` lines
- Supports `#` comments, blank lines, single/double quotes around values
- Does NOT override existing environment variables
- Environment variables always take priority over `.env`

---

## Build Pipeline

```
1. Frontend: npm run build-only
   → Vite bundles Vue app into frontend/dist/ (HTML + JS + CSS)

2. Backend: go build -o stronglifts .
   → go:embed packs frontend/dist/* into the binary
   → modernc.org/sqlite compiled in (pure Go, no CGO)

3. Output: single ~15MB binary
   → Contains entire app: HTTP server + frontend + SQLite driver
   → Zero runtime dependencies
   → Runs on any OS/arch Go supports
```

Cross-compile:
```bash
GOOS=linux GOARCH=arm64 go build -o stronglifts-arm .
GOOS=windows GOARCH=amd64 go build -o stronglifts.exe .
GOOS=darwin GOARCH=arm64 go build -o stronglifts-mac .
```
