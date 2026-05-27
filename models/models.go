package models

import "time"

type User struct {
	ID        int    `json:"id"`
	GoogleID  string `json:"-"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}

type Exercise struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Increment float64 `json:"increment"`
}

type UserExercise struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	ExerciseID    int     `json:"exercise_id"`
	Name          string  `json:"name"`
	CurrentWeight float64 `json:"current_weight"`
	FailCount     int     `json:"fail_count"`
	Increment     float64 `json:"increment"`
}

type Workout struct {
	ID        int               `json:"id"`
	UserID    int               `json:"user_id"`
	Type      string            `json:"type"`
	Date      time.Time         `json:"date"`
	Completed bool              `json:"completed"`
	Notes     string            `json:"notes"`
	Exercises []WorkoutExercise `json:"exercises,omitempty"`
}

type WorkoutExercise struct {
	ID         int          `json:"id"`
	WorkoutID  int          `json:"workout_id"`
	ExerciseID int          `json:"exercise_id"`
	Name       string       `json:"name"`
	TargetSets int          `json:"target_sets"`
	TargetReps int          `json:"target_reps"`
	Weight     float64      `json:"weight"`
	Sets       []WorkoutSet `json:"sets,omitempty"`
}

type WorkoutSet struct {
	ID                int  `json:"id"`
	WorkoutExerciseID int  `json:"workout_exercise_id"`
	SetNumber         int  `json:"set_number"`
	Reps              int  `json:"reps"`
	Completed         bool `json:"completed"`
}

type ProgressEntry struct {
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
}

type PlateResult struct {
	TotalWeight float64 `json:"total_weight"`
	PerSide     []Plate `json:"per_side"`
}

type Plate struct {
	Weight float64 `json:"weight"`
	Count  int     `json:"count"`
}

type BodyWeight struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
}

type WarmupSet struct {
	SetNumber int     `json:"set_number"`
	Weight    float64 `json:"weight"`
	Reps      int     `json:"reps"`
}

type StartingWeightInput struct {
	ExerciseID int     `json:"exercise_id"`
	Weight     float64 `json:"weight"`
	Reps       int     `json:"reps"`
}

type StartingWeightResult struct {
	ExerciseID     int     `json:"exercise_id"`
	Name           string  `json:"name"`
	Estimated1RM   float64 `json:"estimated_1rm"`
	StartingWeight float64 `json:"starting_weight"`
}

type ExportData struct {
	ExportedAt  string       `json:"exported_at"`
	Exercises   []UserExercise `json:"exercises"`
	Workouts    []Workout    `json:"workouts"`
	BodyWeights []BodyWeight `json:"body_weights"`
}

type LeaderboardEntry struct {
	Rank           int     `json:"rank"`
	UserID         int     `json:"user_id"`
	Nickname       string  `json:"nickname"`
	AvatarURL      string  `json:"avatar_url"`
	SquatWeight    float64 `json:"squat_weight"`
	BenchWeight    float64 `json:"bench_weight"`
	DeadliftWeight float64 `json:"deadlift_weight"`
	Total          float64 `json:"total"`
	WorkoutCount   int     `json:"workout_count"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	NeedNickname bool   `json:"need_nickname"`
}
