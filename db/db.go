package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"stronglifts/models"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func New(path string) (*Database, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA foreign_keys=ON")
	db := &Database{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) Close() error { return db.conn.Close() }

func (db *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		google_id TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		nickname TEXT NOT NULL DEFAULT '',
		avatar_url TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		default_weight REAL NOT NULL,
		increment REAL NOT NULL DEFAULT 2.5
	);

	CREATE TABLE IF NOT EXISTS user_exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		exercise_id INTEGER NOT NULL REFERENCES exercises(id),
		current_weight REAL NOT NULL,
		fail_count INTEGER NOT NULL DEFAULT 0,
		UNIQUE(user_id, exercise_id)
	);

	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		type TEXT NOT NULL CHECK(type IN ('A','B')),
		date TEXT NOT NULL,
		completed INTEGER NOT NULL DEFAULT 0,
		notes TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS workout_exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_id INTEGER NOT NULL REFERENCES workouts(id),
		exercise_id INTEGER NOT NULL REFERENCES exercises(id),
		target_sets INTEGER NOT NULL,
		target_reps INTEGER NOT NULL,
		weight REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS workout_sets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_exercise_id INTEGER NOT NULL REFERENCES workout_exercises(id),
		set_number INTEGER NOT NULL,
		reps INTEGER NOT NULL DEFAULT 0,
		completed INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS body_weights (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		date TEXT NOT NULL,
		weight REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS app_settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`
	_, err := db.conn.Exec(schema)
	if err != nil {
		return err
	}
	return db.seedExercises()
}

func (db *Database) seedExercises() error {
	var count int
	db.conn.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if count > 0 {
		return nil
	}
	exercises := []struct {
		name    string
		weight  float64
		incr    float64
	}{
		{"Squat", 20, 2.5},
		{"Bench Press", 20, 2.5},
		{"Barbell Row", 30, 2.5},
		{"Overhead Press", 20, 2.5},
		{"Deadlift", 40, 5.0},
	}
	for _, e := range exercises {
		_, err := db.conn.Exec(
			"INSERT INTO exercises (name, default_weight, increment) VALUES (?, ?, ?)",
			e.name, e.weight, e.incr)
		if err != nil {
			return err
		}
	}
	return nil
}

// ======== Auth ========

func (db *Database) FindOrCreateUser(googleID, email, avatarURL string) (*models.User, bool, error) {
	var u models.User
	err := db.conn.QueryRow(
		"SELECT id, google_id, email, nickname, avatar_url, created_at FROM users WHERE google_id = ?",
		googleID).Scan(&u.ID, &u.GoogleID, &u.Email, &u.Nickname, &u.AvatarURL, &u.CreatedAt)
	if err == nil {
		// Update email/avatar in case they changed
		db.conn.Exec("UPDATE users SET email = ?, avatar_url = ? WHERE id = ?", email, avatarURL, u.ID)
		u.Email = email
		u.AvatarURL = avatarURL
		return &u, false, nil
	}

	now := time.Now().Format(time.RFC3339)
	res, err := db.conn.Exec(
		"INSERT INTO users (google_id, email, avatar_url, created_at) VALUES (?, ?, ?, ?)",
		googleID, email, avatarURL, now)
	if err != nil {
		return nil, false, err
	}
	id, _ := res.LastInsertId()
	u = models.User{ID: int(id), GoogleID: googleID, Email: email, AvatarURL: avatarURL, CreatedAt: now}

	// Initialize user exercises from template
	db.initUserExercises(u.ID)

	return &u, true, nil
}

func (db *Database) initUserExercises(userID int) {
	rows, _ := db.conn.Query("SELECT id, default_weight FROM exercises")
	if rows == nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var exID int
		var w float64
		rows.Scan(&exID, &w)
		db.conn.Exec(
			"INSERT OR IGNORE INTO user_exercises (user_id, exercise_id, current_weight) VALUES (?, ?, ?)",
			userID, exID, w)
	}
}

func (db *Database) SetNickname(userID int, nickname string) error {
	_, err := db.conn.Exec("UPDATE users SET nickname = ? WHERE id = ?", nickname, userID)
	return err
}

func (db *Database) GetUser(userID int) (*models.User, error) {
	var u models.User
	err := db.conn.QueryRow(
		"SELECT id, google_id, email, nickname, avatar_url, created_at FROM users WHERE id = ?",
		userID).Scan(&u.ID, &u.GoogleID, &u.Email, &u.Nickname, &u.AvatarURL, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (db *Database) CreateSession(userID int) (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	now := time.Now()
	expires := now.AddDate(0, 0, 30) // 30 days
	_, err := db.conn.Exec(
		"INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)",
		token, userID, now.Format(time.RFC3339), expires.Format(time.RFC3339))
	return token, err
}

func (db *Database) GetSession(token string) (int, error) {
	var userID int
	var expiresStr string
	err := db.conn.QueryRow(
		"SELECT user_id, expires_at FROM sessions WHERE token = ?", token,
	).Scan(&userID, &expiresStr)
	if err != nil {
		return 0, err
	}
	expires, _ := time.Parse(time.RFC3339, expiresStr)
	if time.Now().After(expires) {
		db.conn.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0, fmt.Errorf("session expired")
	}
	return userID, nil
}

func (db *Database) DeleteSession(token string) {
	db.conn.Exec("DELETE FROM sessions WHERE token = ?", token)
}

// ======== User Exercises ========

func (db *Database) GetUserExercises(userID int) ([]models.UserExercise, error) {
	rows, err := db.conn.Query(`
		SELECT ue.id, ue.user_id, ue.exercise_id, e.name, ue.current_weight, ue.fail_count, e.increment
		FROM user_exercises ue
		JOIN exercises e ON e.id = ue.exercise_id
		WHERE ue.user_id = ?
		ORDER BY e.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.UserExercise
	for rows.Next() {
		var ue models.UserExercise
		rows.Scan(&ue.ID, &ue.UserID, &ue.ExerciseID, &ue.Name, &ue.CurrentWeight, &ue.FailCount, &ue.Increment)
		result = append(result, ue)
	}
	return result, nil
}

func (db *Database) UpdateUserExercise(userID, exerciseID int, weight float64, failCount int) error {
	_, err := db.conn.Exec(
		"UPDATE user_exercises SET current_weight = ?, fail_count = ? WHERE user_id = ? AND exercise_id = ?",
		weight, failCount, userID, exerciseID)
	return err
}

// ======== Setup ========

func (db *Database) IsSetup(userID int) bool {
	var val string
	err := db.conn.QueryRow(
		"SELECT value FROM app_settings WHERE key = ?",
		fmt.Sprintf("setup_done_%d", userID)).Scan(&val)
	return err == nil && val == "true"
}

func (db *Database) MarkSetup(userID int) {
	db.conn.Exec(
		"INSERT OR REPLACE INTO app_settings (key, value) VALUES (?, 'true')",
		fmt.Sprintf("setup_done_%d", userID))
}

func CalculateStartingWeight(weight float64, reps int, increment float64) float64 {
	if reps <= 0 || weight <= 0 {
		return 20
	}
	oneRM := weight * (1.0 + float64(reps)/30.0)
	start := oneRM * 0.5
	start = math.Floor(start/increment) * increment
	if start < 20 {
		start = 20
	}
	return start
}

func (db *Database) SetStartingWeights(userID int, inputs []models.StartingWeightInput) ([]models.StartingWeightResult, error) {
	var results []models.StartingWeightResult
	for _, inp := range inputs {
		var name string
		var increment float64
		err := db.conn.QueryRow("SELECT name, increment FROM exercises WHERE id = ?", inp.ExerciseID).Scan(&name, &increment)
		if err != nil {
			return nil, fmt.Errorf("exercise %d not found", inp.ExerciseID)
		}
		oneRM := inp.Weight * (1.0 + float64(inp.Reps)/30.0)
		startWeight := CalculateStartingWeight(inp.Weight, inp.Reps, increment)
		db.conn.Exec(
			"UPDATE user_exercises SET current_weight = ?, fail_count = 0 WHERE user_id = ? AND exercise_id = ?",
			startWeight, userID, inp.ExerciseID)
		results = append(results, models.StartingWeightResult{
			ExerciseID: inp.ExerciseID, Name: name,
			Estimated1RM: math.Round(oneRM*10) / 10, StartingWeight: startWeight,
		})
	}
	return results, nil
}

// ======== Warmup ========

func GetWarmupSets(workingWeight float64) []models.WarmupSet {
	barWeight := 20.0
	if workingWeight <= barWeight {
		return nil
	}
	var sets []models.WarmupSet
	setNum := 1
	sets = append(sets, models.WarmupSet{SetNumber: setNum, Weight: barWeight, Reps: 5})
	setNum++
	steps := []struct {
		pct  float64
		reps int
	}{{0.4, 5}, {0.6, 3}, {0.8, 2}}
	for _, s := range steps {
		w := roundWeight(workingWeight * s.pct)
		if w > barWeight && w < workingWeight {
			sets = append(sets, models.WarmupSet{SetNumber: setNum, Weight: w, Reps: s.reps})
			setNum++
		}
	}
	return sets
}

func roundWeight(w float64) float64 { return math.Round(w/2.5) * 2.5 }

// ======== Workouts ========

func (db *Database) GetNextWorkoutType(userID int) string {
	var lastType string
	err := db.conn.QueryRow(
		"SELECT type FROM workouts WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID,
	).Scan(&lastType)
	if err != nil || lastType == "B" {
		return "A"
	}
	return "B"
}

func (db *Database) CreateWorkout(userID int, wType string) (*models.Workout, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now().Format(time.RFC3339)
	res, err := tx.Exec("INSERT INTO workouts (user_id, type, date) VALUES (?, ?, ?)", userID, wType, now)
	if err != nil {
		return nil, err
	}
	workoutID, _ := res.LastInsertId()

	type exDef struct {
		name       string
		sets, reps int
	}
	var exDefs []exDef
	if wType == "A" {
		exDefs = []exDef{{"Squat", 5, 5}, {"Bench Press", 5, 5}, {"Barbell Row", 5, 5}}
	} else {
		exDefs = []exDef{{"Squat", 5, 5}, {"Overhead Press", 5, 5}, {"Deadlift", 1, 5}}
	}

	workout := &models.Workout{ID: int(workoutID), UserID: userID, Type: wType, Date: time.Now()}

	for _, ed := range exDefs {
		var exID int
		var curWeight float64
		err := tx.QueryRow(`
			SELECT e.id, ue.current_weight FROM exercises e
			JOIN user_exercises ue ON ue.exercise_id = e.id
			WHERE e.name = ? AND ue.user_id = ?`, ed.name, userID,
		).Scan(&exID, &curWeight)
		if err != nil {
			return nil, fmt.Errorf("exercise %s not found for user: %w", ed.name, err)
		}

		res, err := tx.Exec(
			"INSERT INTO workout_exercises (workout_id, exercise_id, target_sets, target_reps, weight) VALUES (?, ?, ?, ?, ?)",
			workoutID, exID, ed.sets, ed.reps, curWeight)
		if err != nil {
			return nil, err
		}
		weID, _ := res.LastInsertId()
		we := models.WorkoutExercise{
			ID: int(weID), WorkoutID: int(workoutID), ExerciseID: exID,
			Name: ed.name, TargetSets: ed.sets, TargetReps: ed.reps, Weight: curWeight,
		}
		for s := 1; s <= ed.sets; s++ {
			tx.Exec("INSERT INTO workout_sets (workout_exercise_id, set_number) VALUES (?, ?)", weID, s)
			we.Sets = append(we.Sets, models.WorkoutSet{WorkoutExerciseID: int(weID), SetNumber: s})
		}
		workout.Exercises = append(workout.Exercises, we)
	}
	return workout, tx.Commit()
}

func (db *Database) GetWorkout(id int) (*models.Workout, error) {
	var w models.Workout
	var dateStr string
	var completed int
	err := db.conn.QueryRow(
		"SELECT id, user_id, type, date, completed, notes FROM workouts WHERE id = ?", id,
	).Scan(&w.ID, &w.UserID, &w.Type, &dateStr, &completed, &w.Notes)
	if err != nil {
		return nil, err
	}
	w.Date, _ = time.Parse(time.RFC3339, dateStr)
	w.Completed = completed == 1

	rows, err := db.conn.Query(`
		SELECT we.id, we.workout_id, we.exercise_id, e.name, we.target_sets, we.target_reps, we.weight
		FROM workout_exercises we JOIN exercises e ON e.id = we.exercise_id
		WHERE we.workout_id = ? ORDER BY we.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var we models.WorkoutExercise
		rows.Scan(&we.ID, &we.WorkoutID, &we.ExerciseID, &we.Name, &we.TargetSets, &we.TargetReps, &we.Weight)
		setRows, _ := db.conn.Query(
			"SELECT id, workout_exercise_id, set_number, reps, completed FROM workout_sets WHERE workout_exercise_id = ? ORDER BY set_number", we.ID)
		if setRows != nil {
			for setRows.Next() {
				var s models.WorkoutSet
				var comp int
				setRows.Scan(&s.ID, &s.WorkoutExerciseID, &s.SetNumber, &s.Reps, &comp)
				s.Completed = comp == 1
				we.Sets = append(we.Sets, s)
			}
			setRows.Close()
		}
		w.Exercises = append(w.Exercises, we)
	}
	return &w, nil
}

func (db *Database) CompleteSet(setID int, reps int) error {
	_, err := db.conn.Exec("UPDATE workout_sets SET reps = ?, completed = 1 WHERE id = ?", reps, setID)
	return err
}

func (db *Database) UpdateWorkoutNotes(workoutID int, notes string) error {
	_, err := db.conn.Exec("UPDATE workouts SET notes = ? WHERE id = ?", notes, workoutID)
	return err
}

func (db *Database) CompleteWorkout(workoutID int) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	tx.QueryRow("SELECT user_id FROM workouts WHERE id = ?", workoutID).Scan(&userID)
	tx.Exec("UPDATE workouts SET completed = 1 WHERE id = ?", workoutID)

	rows, _ := tx.Query(
		"SELECT exercise_id, target_sets, target_reps, weight FROM workout_exercises WHERE workout_id = ?", workoutID)
	type exR struct{ exID, tSets, tReps int; weight float64 }
	var results []exR
	for rows.Next() {
		var r exR
		rows.Scan(&r.exID, &r.tSets, &r.tReps, &r.weight)
		results = append(results, r)
	}
	rows.Close()

	for _, r := range results {
		var successSets int
		tx.QueryRow(`
			SELECT COUNT(*) FROM workout_sets ws
			JOIN workout_exercises we ON we.id = ws.workout_exercise_id
			WHERE we.workout_id = ? AND we.exercise_id = ? AND ws.completed = 1 AND ws.reps >= ?`,
			workoutID, r.exID, r.tReps).Scan(&successSets)

		var curWeight float64
		var failCount int
		var increment float64
		tx.QueryRow(`
			SELECT ue.current_weight, ue.fail_count, e.increment
			FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id
			WHERE ue.user_id = ? AND ue.exercise_id = ?`, userID, r.exID,
		).Scan(&curWeight, &failCount, &increment)

		if successSets >= r.tSets {
			tx.Exec("UPDATE user_exercises SET current_weight = ?, fail_count = 0 WHERE user_id = ? AND exercise_id = ?",
				curWeight+increment, userID, r.exID)
		} else {
			failCount++
			if failCount >= 3 {
				nw := math.Floor(curWeight*0.9/2.5) * 2.5
				if nw < 20 { nw = 20 }
				tx.Exec("UPDATE user_exercises SET current_weight = ?, fail_count = 0 WHERE user_id = ? AND exercise_id = ?",
					nw, userID, r.exID)
			} else {
				tx.Exec("UPDATE user_exercises SET fail_count = ? WHERE user_id = ? AND exercise_id = ?",
					failCount, userID, r.exID)
			}
		}
	}
	return tx.Commit()
}

func (db *Database) GetActiveWorkout(userID int) (*models.Workout, error) {
	var id int
	err := db.conn.QueryRow(
		"SELECT id FROM workouts WHERE user_id = ? AND completed = 0 ORDER BY id DESC LIMIT 1", userID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return db.GetWorkout(id)
}

func (db *Database) DeleteWorkout(id int) error {
	tx, _ := db.conn.Begin()
	defer tx.Rollback()
	tx.Exec("DELETE FROM workout_sets WHERE workout_exercise_id IN (SELECT id FROM workout_exercises WHERE workout_id = ?)", id)
	tx.Exec("DELETE FROM workout_exercises WHERE workout_id = ?", id)
	tx.Exec("DELETE FROM workouts WHERE id = ?", id)
	return tx.Commit()
}

func (db *Database) GetWorkoutHistory(userID, limit int) ([]models.Workout, error) {
	rows, err := db.conn.Query(
		"SELECT id, user_id, type, date, completed, notes FROM workouts WHERE user_id = ? ORDER BY id DESC LIMIT ?",
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Workout
	for rows.Next() {
		var w models.Workout
		var dateStr string
		var completed int
		rows.Scan(&w.ID, &w.UserID, &w.Type, &dateStr, &completed, &w.Notes)
		w.Date, _ = time.Parse(time.RFC3339, dateStr)
		w.Completed = completed == 1
		result = append(result, w)
	}
	return result, nil
}

func (db *Database) GetProgress(userID, exerciseID int) ([]models.ProgressEntry, error) {
	rows, err := db.conn.Query(`
		SELECT w.date, we.weight FROM workout_exercises we
		JOIN workouts w ON w.id = we.workout_id
		WHERE w.user_id = ? AND we.exercise_id = ? AND w.completed = 1
		ORDER BY w.date`, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.ProgressEntry
	for rows.Next() {
		var p models.ProgressEntry
		var dateStr string
		rows.Scan(&dateStr, &p.Weight)
		t, _ := time.Parse(time.RFC3339, dateStr)
		p.Date = t.Format("2006-01-02")
		result = append(result, p)
	}
	return result, nil
}

// ======== Body Weight ========

func (db *Database) AddBodyWeight(userID int, weight float64) (*models.BodyWeight, error) {
	date := time.Now().Format("2006-01-02")
	db.conn.Exec("DELETE FROM body_weights WHERE user_id = ? AND date = ?", userID, date)
	res, err := db.conn.Exec("INSERT INTO body_weights (user_id, date, weight) VALUES (?, ?, ?)", userID, date, weight)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.BodyWeight{ID: int(id), UserID: userID, Date: date, Weight: weight}, nil
}

func (db *Database) GetBodyWeights(userID, limit int) ([]models.BodyWeight, error) {
	rows, err := db.conn.Query(
		"SELECT id, user_id, date, weight FROM body_weights WHERE user_id = ? ORDER BY date DESC LIMIT ?",
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.BodyWeight
	for rows.Next() {
		var bw models.BodyWeight
		rows.Scan(&bw.ID, &bw.UserID, &bw.Date, &bw.Weight)
		result = append(result, bw)
	}
	return result, nil
}

func (db *Database) DeleteBodyWeight(id int) error {
	_, err := db.conn.Exec("DELETE FROM body_weights WHERE id = ?", id)
	return err
}

// ======== Leaderboard ========

func (db *Database) GetLeaderboard() ([]models.LeaderboardEntry, error) {
	rows, err := db.conn.Query(`
		SELECT
			u.id, u.nickname, u.avatar_url,
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Squat'), 0),
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Bench Press'), 0),
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Deadlift'), 0),
			(SELECT COUNT(*) FROM workouts w WHERE w.user_id = u.id AND w.completed = 1)
		FROM users u
		WHERE u.nickname != ''
		ORDER BY (
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Squat'), 0) +
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Bench Press'), 0) +
			COALESCE((SELECT ue.current_weight FROM user_exercises ue JOIN exercises e ON e.id = ue.exercise_id WHERE ue.user_id = u.id AND e.name = 'Deadlift'), 0)
		) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e models.LeaderboardEntry
		rows.Scan(&e.UserID, &e.Nickname, &e.AvatarURL, &e.SquatWeight, &e.BenchWeight, &e.DeadliftWeight, &e.WorkoutCount)
		e.Total = e.SquatWeight + e.BenchWeight + e.DeadliftWeight
		e.Rank = rank
		rank++
		result = append(result, e)
	}
	return result, nil
}

// ======== Export ========

func (db *Database) ExportAll(userID int) (*models.ExportData, error) {
	exercises, _ := db.GetUserExercises(userID)
	ids := []int{}
	rows, _ := db.conn.Query("SELECT id FROM workouts WHERE user_id = ? ORDER BY id", userID)
	if rows != nil {
		for rows.Next() {
			var id int
			rows.Scan(&id)
			ids = append(ids, id)
		}
		rows.Close()
	}
	var workouts []models.Workout
	for _, id := range ids {
		w, err := db.GetWorkout(id)
		if err == nil {
			workouts = append(workouts, *w)
		}
	}
	bws, _ := db.GetBodyWeights(userID, 10000)
	return &models.ExportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		Exercises: exercises, Workouts: workouts, BodyWeights: bws,
	}, nil
}
