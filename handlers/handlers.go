package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"stronglifts/db"
	"stronglifts/models"
	"time"
)

type contextKey string
const userIDKey contextKey = "userID"

type Handler struct {
	db       *db.Database
	clientID string
}

func New(database *db.Database) *Handler {
	return &Handler{db: database, clientID: os.Getenv("GOOGLE_CLIENT_ID")}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Public config
	mux.HandleFunc("GET /api/config", h.getConfig)

	// Auth (no middleware)
	mux.HandleFunc("POST /api/auth/google", h.googleAuth)
	mux.HandleFunc("POST /api/auth/logout", h.logout)
	mux.HandleFunc("GET /api/auth/me", h.withAuth(h.getMe))
	mux.HandleFunc("POST /api/auth/nickname", h.withAuth(h.setNickname))

	// Setup
	mux.HandleFunc("GET /api/setup/status", h.withAuth(h.getSetupStatus))
	mux.HandleFunc("POST /api/setup/starting-weights", h.withAuth(h.setStartingWeights))

	// Exercises
	mux.HandleFunc("GET /api/exercises", h.withAuth(h.getExercises))
	mux.HandleFunc("PUT /api/exercises/{id}/weight", h.withAuth(h.updateExerciseWeight))

	// Workouts
	mux.HandleFunc("GET /api/workout/active", h.withAuth(h.getActiveWorkout))
	mux.HandleFunc("GET /api/workout/next-type", h.withAuth(h.getNextWorkoutType))
	mux.HandleFunc("POST /api/workout/start", h.withAuth(h.startWorkout))
	mux.HandleFunc("GET /api/workout/{id}", h.withAuth(h.getWorkout))
	mux.HandleFunc("POST /api/workout/{id}/complete", h.withAuth(h.completeWorkout))
	mux.HandleFunc("PUT /api/workout/{id}/notes", h.withAuth(h.updateNotes))
	mux.HandleFunc("DELETE /api/workout/{id}", h.withAuth(h.deleteWorkout))

	// Sets
	mux.HandleFunc("POST /api/set/{id}/complete", h.withAuth(h.completeSet))

	// Warmup / Plates
	mux.HandleFunc("GET /api/warmup/{weight}", h.getWarmupSets)
	mux.HandleFunc("GET /api/plates/{weight}", h.calculatePlates)

	// Progress / History
	mux.HandleFunc("GET /api/progress/{exerciseId}", h.withAuth(h.getProgress))
	mux.HandleFunc("GET /api/history", h.withAuth(h.getHistory))

	// Body weight
	mux.HandleFunc("GET /api/bodyweight", h.withAuth(h.getBodyWeights))
	mux.HandleFunc("POST /api/bodyweight", h.withAuth(h.addBodyWeight))
	mux.HandleFunc("DELETE /api/bodyweight/{id}", h.withAuth(h.deleteBodyWeight))

	// Leaderboard (public)
	mux.HandleFunc("GET /api/leaderboard", h.getLeaderboard)

	// Export
	mux.HandleFunc("GET /api/export", h.withAuth(h.exportData))
}

// ======== Auth Middleware ========

func (h *Handler) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := ""
		if c, err := r.Cookie("session"); err == nil {
			token = c.Value
		}
		if token == "" {
			token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		}
		if token == "" {
			jsonError(w, "unauthorized", 401)
			return
		}
		userID, err := h.db.GetSession(token)
		if err != nil {
			jsonError(w, "unauthorized", 401)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func getUserID(r *http.Request) int {
	return r.Context().Value(userIDKey).(int)
}

// ======== Helpers ========

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ======== Config ========

func (h *Handler) getConfig(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]string{
		"google_client_id": h.clientID,
	})
}

// ======== Auth Handlers ========

func (h *Handler) googleAuth(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Credential string `json:"credential"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", 400)
		return
	}

	// Verify ID token with Google
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", body.Credential))
	if err != nil || resp.StatusCode != 200 {
		jsonError(w, "invalid google token", 401)
		return
	}
	defer resp.Body.Close()
	tokenBody, _ := io.ReadAll(resp.Body)

	var claims struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
		Aud     string `json:"aud"`
	}
	json.Unmarshal(tokenBody, &claims)

	if h.clientID != "" && claims.Aud != h.clientID {
		jsonError(w, "token audience mismatch", 401)
		return
	}

	user, isNew, err := h.db.FindOrCreateUser(claims.Sub, claims.Email, claims.Picture)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	token, err := h.db.CreateSession(user.ID)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 3600,
	})

	jsonResponse(w, models.AuthResponse{
		User:         *user,
		NeedNickname: isNew || user.Nickname == "",
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session"); err == nil {
		h.db.DeleteSession(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: "", Path: "/", MaxAge: -1,
	})
	jsonResponse(w, map[string]string{"status": "ok"})
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.db.GetUser(getUserID(r))
	if err != nil {
		jsonError(w, "user not found", 404)
		return
	}
	jsonResponse(w, user)
}

func (h *Handler) setNickname(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", 400)
		return
	}
	nick := strings.TrimSpace(body.Nickname)
	if nick == "" || len(nick) > 20 {
		jsonError(w, "nickname 1-20 characters", 400)
		return
	}
	h.db.SetNickname(getUserID(r), nick)
	jsonResponse(w, map[string]string{"status": "ok"})
}

// ======== Setup ========

func (h *Handler) getSetupStatus(w http.ResponseWriter, r *http.Request) {
	done := h.db.IsSetup(getUserID(r))
	jsonResponse(w, map[string]bool{"setup_done": done})
}

func (h *Handler) setStartingWeights(w http.ResponseWriter, r *http.Request) {
	var inputs []models.StartingWeightInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		jsonError(w, "invalid body", 400)
		return
	}
	uid := getUserID(r)
	results, err := h.db.SetStartingWeights(uid, inputs)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	h.db.MarkSetup(uid)
	jsonResponse(w, results)
}

// ======== Exercises ========

func (h *Handler) getExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.db.GetUserExercises(getUserID(r))
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, exercises)
}

func (h *Handler) updateExerciseWeight(w http.ResponseWriter, r *http.Request) {
	exID, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ Weight float64 `json:"weight"` }
	json.NewDecoder(r.Body).Decode(&body)
	h.db.UpdateUserExercise(getUserID(r), exID, body.Weight, 0)
	jsonResponse(w, map[string]string{"status": "ok"})
}

// ======== Workouts ========

func (h *Handler) getActiveWorkout(w http.ResponseWriter, r *http.Request) {
	workout, err := h.db.GetActiveWorkout(getUserID(r))
	if err != nil {
		jsonResponse(w, nil)
		return
	}
	jsonResponse(w, workout)
}

func (h *Handler) getNextWorkoutType(w http.ResponseWriter, r *http.Request) {
	t := h.db.GetNextWorkoutType(getUserID(r))
	jsonResponse(w, map[string]string{"type": t})
}

func (h *Handler) startWorkout(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if active, _ := h.db.GetActiveWorkout(uid); active != nil {
		jsonError(w, "active workout exists", 400)
		return
	}
	workout, err := h.db.CreateWorkout(uid, h.db.GetNextWorkoutType(uid))
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, workout)
}

func (h *Handler) getWorkout(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	workout, err := h.db.GetWorkout(id)
	if err != nil {
		jsonError(w, "not found", 404)
		return
	}
	if workout.UserID != getUserID(r) {
		jsonError(w, "forbidden", 403)
		return
	}
	jsonResponse(w, workout)
}

func (h *Handler) completeWorkout(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	h.db.CompleteWorkout(id)
	jsonResponse(w, map[string]string{"status": "ok"})
}

func (h *Handler) updateNotes(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ Notes string `json:"notes"` }
	json.NewDecoder(r.Body).Decode(&body)
	h.db.UpdateWorkoutNotes(id, body.Notes)
	jsonResponse(w, map[string]string{"status": "ok"})
}

func (h *Handler) deleteWorkout(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	h.db.DeleteWorkout(id)
	jsonResponse(w, map[string]string{"status": "ok"})
}

func (h *Handler) completeSet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ Reps int `json:"reps"` }
	json.NewDecoder(r.Body).Decode(&body)
	h.db.CompleteSet(id, body.Reps)
	jsonResponse(w, map[string]string{"status": "ok"})
}

// ======== Warmup / Plates (public) ========

func (h *Handler) getWarmupSets(w http.ResponseWriter, r *http.Request) {
	weight, _ := strconv.ParseFloat(strings.TrimSpace(r.PathValue("weight")), 64)
	sets := db.GetWarmupSets(weight)
	if sets == nil { sets = []models.WarmupSet{} }
	jsonResponse(w, sets)
}

func (h *Handler) calculatePlates(w http.ResponseWriter, r *http.Request) {
	weight, _ := strconv.ParseFloat(strings.TrimSpace(r.PathValue("weight")), 64)
	barWeight := 20.0
	perSide := (weight - barWeight) / 2.0
	if perSide < 0 {
		jsonResponse(w, models.PlateResult{TotalWeight: barWeight})
		return
	}
	available := []float64{20, 10, 5, 2.5, 1.25}
	var plates []models.Plate
	remaining := perSide
	for _, p := range available {
		if remaining >= p {
			count := int(math.Floor(remaining / p))
			plates = append(plates, models.Plate{Weight: p, Count: count})
			remaining -= float64(count) * p
		}
	}
	jsonResponse(w, models.PlateResult{TotalWeight: weight, PerSide: plates})
}

// ======== Progress / History ========

func (h *Handler) getProgress(w http.ResponseWriter, r *http.Request) {
	exID, _ := strconv.Atoi(r.PathValue("exerciseId"))
	entries, _ := h.db.GetProgress(getUserID(r), exID)
	if entries == nil { entries = []models.ProgressEntry{} }
	jsonResponse(w, entries)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 { limit = n }
	}
	history, _ := h.db.GetWorkoutHistory(getUserID(r), limit)
	if history == nil { history = []models.Workout{} }
	jsonResponse(w, history)
}

// ======== Body Weight ========

func (h *Handler) getBodyWeights(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 { limit = n }
	}
	bws, _ := h.db.GetBodyWeights(getUserID(r), limit)
	if bws == nil { bws = []models.BodyWeight{} }
	jsonResponse(w, bws)
}

func (h *Handler) addBodyWeight(w http.ResponseWriter, r *http.Request) {
	var body struct{ Weight float64 `json:"weight"` }
	json.NewDecoder(r.Body).Decode(&body)
	if body.Weight <= 0 {
		jsonError(w, "weight must be positive", 400)
		return
	}
	bw, err := h.db.AddBodyWeight(getUserID(r), body.Weight)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, bw)
}

func (h *Handler) deleteBodyWeight(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	h.db.DeleteBodyWeight(id)
	jsonResponse(w, map[string]string{"status": "ok"})
}

// ======== Leaderboard ========

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.db.GetLeaderboard()
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	if entries == nil { entries = []models.LeaderboardEntry{} }
	jsonResponse(w, entries)
}

// ======== Export ========

func (h *Handler) exportData(w http.ResponseWriter, r *http.Request) {
	data, _ := h.db.ExportAll(getUserID(r))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=stronglifts_backup_%s.json", time.Now().Format("2006-01-02")))
	json.NewEncoder(w).Encode(data)
}
