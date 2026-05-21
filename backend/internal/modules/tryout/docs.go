package tryout

import "time"

type AttemptSummary struct {
	ID                string     `json:"id"`
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
}

type NavigatorItem struct {
	QuestionID   string `json:"question_id"`
	DisplayOrder int    `json:"display_order"`
	IsAnswered   bool   `json:"is_answered"`
	IsFlagged    bool   `json:"is_flagged"`
}

type CurrentTryoutData struct {
	Tryout     Tryout          `json:"tryout"`
	Attempt    *AttemptSummary `json:"attempt,omitempty"`
	ServerTime time.Time       `json:"server_time"`
}

type CurrentTryoutSuccessResponse struct {
	Success   bool              `json:"success" example:"true"`
	Message   string            `json:"message" example:"current tryout fetched"`
	RequestID string            `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      CurrentTryoutData `json:"data"`
}

type StartCurrentTryoutData struct {
	Tryout     Tryout         `json:"tryout"`
	Attempt    AttemptSummary `json:"attempt"`
	ServerTime time.Time      `json:"server_time"`
}

type StartCurrentTryoutSuccessResponse struct {
	Success   bool                   `json:"success" example:"true"`
	Message   string                 `json:"message" example:"attempt prepared"`
	RequestID string                 `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      StartCurrentTryoutData `json:"data"`
}

type CurrentAttemptData struct {
	Tryout     Tryout          `json:"tryout"`
	Attempt    AttemptSummary  `json:"attempt"`
	Questions  []Question      `json:"questions"`
	Answers    []AttemptAnswer `json:"answers"`
	Navigator  []NavigatorItem `json:"navigator"`
	ServerTime time.Time       `json:"server_time"`
}

type CurrentAttemptSuccessResponse struct {
	Success   bool               `json:"success" example:"true"`
	Message   string             `json:"message" example:"current attempt fetched"`
	RequestID string             `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      CurrentAttemptData `json:"data"`
}

type SaveCurrentAnswersData struct {
	Attempt    AttemptSummary `json:"attempt"`
	ServerTime time.Time      `json:"server_time"`
}

type SaveCurrentAnswersSuccessResponse struct {
	Success   bool                   `json:"success" example:"true"`
	Message   string                 `json:"message" example:"answers saved"`
	RequestID string                 `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      SaveCurrentAnswersData `json:"data"`
}

type SubmitCurrentAttemptData struct {
	Attempt    AttemptSummary `json:"attempt"`
	Result     ResultSummary  `json:"result"`
	ServerTime time.Time      `json:"server_time"`
}

type SubmitCurrentAttemptSuccessResponse struct {
	Success   bool                     `json:"success" example:"true"`
	Message   string                   `json:"message" example:"attempt submitted"`
	RequestID string                   `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      SubmitCurrentAttemptData `json:"data"`
}

type CurrentResultData struct {
	Tryout     Tryout        `json:"tryout"`
	Result     ResultSummary `json:"result"`
	ServerTime time.Time     `json:"server_time"`
}

type CurrentResultSuccessResponse struct {
	Success   bool              `json:"success" example:"true"`
	Message   string            `json:"message" example:"result fetched"`
	RequestID string            `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      CurrentResultData `json:"data"`
}

type CurrentLeaderboardData struct {
	Tryout     Tryout             `json:"tryout"`
	Entries    []LeaderboardEntry `json:"entries"`
	Total      int                `json:"total"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
	Me         *ResultSummary     `json:"me,omitempty"`
	ServerTime time.Time          `json:"server_time"`
}

type CurrentLeaderboardSuccessResponse struct {
	Success   bool                   `json:"success" example:"true"`
	Message   string                 `json:"message" example:"leaderboard fetched"`
	RequestID string                 `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      CurrentLeaderboardData `json:"data"`
}
