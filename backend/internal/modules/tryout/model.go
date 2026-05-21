package tryout

import "time"

type Tryout struct {
	ID                 string     `json:"id"`
	Slug               string     `json:"slug"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Instructions       string     `json:"instructions"`
	Status             string     `json:"status"`
	DurationMinutes    int        `json:"duration_minutes"`
	StartsAt           *time.Time `json:"starts_at,omitempty"`
	EndsAt             *time.Time `json:"ends_at,omitempty"`
	ResultsPublishedAt *time.Time `json:"results_published_at,omitempty"`
	ShowLeaderboard    bool       `json:"show_leaderboard"`
	QuestionCount      int        `json:"question_count"`
}

type Attempt struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id,omitempty"`
	TryoutID          string     `json:"tryout_id,omitempty"`
	Status            string     `json:"status"`
	StartedAt         time.Time  `json:"started_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	LastSyncedAt      *time.Time `json:"last_synced_at,omitempty"`
	Version           int        `json:"version"`
	TotalQuestions    int        `json:"total_questions"`
	AnsweredQuestions int        `json:"answered_questions"`
	FlaggedQuestions  int        `json:"flagged_questions"`
	CorrectCount      int        `json:"correct_count"`
	WrongCount        int        `json:"wrong_count"`
	UnansweredCount   int        `json:"unanswered_count"`
	FinalScore        *float64   `json:"final_score,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type Question struct {
	ID           string           `json:"id"`
	Code         string           `json:"code"`
	QuestionType string           `json:"question_type"`
	PromptHTML   string           `json:"prompt_html"`
	ImageURL     string           `json:"image_url,omitempty"`
	DisplayOrder int              `json:"display_order"`
	Points       float64          `json:"points"`
	Options      []QuestionOption `json:"options,omitempty"`
}

type QuestionOption struct {
	Key          string `json:"key"`
	Text         string `json:"text"`
	DisplayOrder int    `json:"display_order"`
}

type AttemptAnswer struct {
	ID                string     `json:"id"`
	AttemptID         string     `json:"attempt_id"`
	QuestionID        string     `json:"question_id"`
	SelectedOptionKey *string    `json:"selected_option_key,omitempty"`
	AnswerText        *string    `json:"answer_text,omitempty"`
	IsFlagged         bool       `json:"is_flagged"`
	AnsweredAt        *time.Time `json:"answered_at,omitempty"`
	LastSavedAt       time.Time  `json:"last_saved_at"`
	ClientUpdatedAt   *time.Time `json:"client_updated_at,omitempty"`
	ServerVersion     int        `json:"server_version"`
	IsCorrect         *bool      `json:"is_correct,omitempty"`
	AwardedPoints     *float64   `json:"awarded_points,omitempty"`
}

type QuestionMeta struct {
	ID              string
	QuestionType    string
	ValidOptionKeys map[string]bool
}

type ScoringQuestion struct {
	ID               string
	QuestionType     string
	Points           float64
	CorrectOptionKey string
	Variants         map[string]bool
}

type ScoringAnswer struct {
	QuestionID        string
	SelectedOptionKey *string
	NormalizedAnswer  *string
	IsFlagged         bool
}

type ResultSummary struct {
	AttemptID         string     `json:"attempt_id"`
	Status            string     `json:"status"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	TotalQuestions    int        `json:"total_questions"`
	CorrectCount      int        `json:"correct_count"`
	WrongCount        int        `json:"wrong_count"`
	UnansweredCount   int        `json:"unanswered_count"`
	FinalScore        float64    `json:"final_score"`
	Rank              int        `json:"rank"`
	TotalParticipants int        `json:"total_participants"`
}

type LeaderboardEntry struct {
	Rank        int        `json:"rank"`
	UserID      string     `json:"user_id"`
	FullName    string     `json:"full_name"`
	FinalScore  float64    `json:"final_score"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}
