package api

import "time"

// Auth
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Apps
type App struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	BundleID           string  `json:"bundle_id"`
	Platform           string  `json:"platform"`
	TaskCount          int     `json:"task_count"`
	LastAggregateScore *int    `json:"last_aggregate_score"`
	LastRunAt          *string `json:"last_run_at"`
}

// Tasks
type Task struct {
	ID             string  `json:"id"`
	AppID          string  `json:"app_id"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	MaxSteps       int     `json:"max_steps"`
	ScoreThreshold *int    `json:"score_threshold"`
}

// Runs
type TriggerRunRequest struct {
	DeviceType string   `json:"device_type,omitempty"`
	TaskIDs    []string `json:"task_ids,omitempty"`
	CommitHash string   `json:"commit_hash,omitempty"`
	PRNumber   *int     `json:"pr_number,omitempty"`
}

type Run struct {
	ID             string     `json:"id"`
	AppID          string     `json:"app_id"`
	Status         string     `json:"status"`
	DeviceType     string     `json:"device_type"`
	CommitHash     *string    `json:"commit_hash"`
	TaskCount      int        `json:"task_count"`
	CompletedCount int        `json:"completed_count"`
	AggregateScore *int       `json:"aggregate_score"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

// Task Results
type Scores struct {
	Discoverability    int `json:"discoverability"`
	Efficiency         int `json:"efficiency"`
	NavigationClarity  int `json:"navigation_clarity"`
	FeedbackClarity    int `json:"feedback_clarity"`
	ConfirmationClarity int `json:"confirmation_clarity"`
	InterruptionImpact int `json:"interruption_impact"`
}

type Finding struct {
	Category       string `json:"category"`
	Severity       string `json:"severity"`
	StepNumber     int    `json:"step_number"`
	Observation    string `json:"observation"`
	Recommendation string `json:"recommendation"`
}

type TaskResult struct {
	ID              string    `json:"id"`
	RunID           string    `json:"run_id"`
	TaskID          string    `json:"task_id"`
	TaskDescription *string   `json:"task_description"`
	Succeeded       bool      `json:"succeeded"`
	Steps           int       `json:"steps"`
	TotalMs         *int      `json:"total_ms"`
	AggregateScore  int       `json:"aggregate_score"`
	Scores          Scores    `json:"scores"`
	Findings        []Finding `json:"findings"`
}

// Specs / Playbooks
type Spec struct {
	ID         string    `json:"id"`
	RunID      string    `json:"run_id"`
	TaskID     string    `json:"task_id"`
	Name       *string   `json:"name"`
	StepCount  *int      `json:"step_count"`
	Difficulty *string   `json:"difficulty"`
	Platform   string    `json:"platform"`
	CreatedAt  time.Time `json:"created_at"`
}

// Playback
type PlaybackResponse struct {
	PlaybackRunID string `json:"playback_run_id"`
	Status        string `json:"status"`
	StreamURL     string `json:"stream_url"`
}

type PlaybackResult struct {
	ID             string `json:"id"`
	SpecID         string `json:"spec_id"`
	Status         string `json:"status"`
	StepsCompleted *int   `json:"steps_completed"`
	TotalSteps     *int   `json:"total_steps"`
	FailureReason  *string `json:"failure_reason"`
	ElapsedMs      *int   `json:"elapsed_ms"`
}

// SSE Event
type StepEvent struct {
	Type      string `json:"type"`
	RunID     string `json:"run_id"`
	TaskID    string `json:"task_id"`
	Step      int    `json:"step"`
	Action    string `json:"action"`
	Label     string `json:"label"`
	Outcome   string `json:"outcome"`
	ElapsedMs int    `json:"elapsed_ms"`
	Status    string `json:"status"`

	// For completed events
	Succeeded      *bool `json:"succeeded,omitempty"`
	AggregateScore *int  `json:"aggregate_score,omitempty"`
}
