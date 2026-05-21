package tryout

import "time"

type SaveAnswersRequest struct {
	AttemptVersion int              `json:"attempt_version"`
	Answers        []SaveAnswerItem `json:"answers"`
}

type SaveAnswerItem struct {
	QuestionID        string     `json:"question_id"`
	SelectedOptionKey *string    `json:"selected_option_key"`
	AnswerText        *string    `json:"answer_text"`
	IsFlagged         bool       `json:"is_flagged"`
	ClientUpdatedAt   *time.Time `json:"client_updated_at"`
}

type SubmitAttemptRequest struct {
	AttemptVersion int                `json:"attempt_version"`
	FinalAnswers   []SubmitAnswerItem `json:"final_answers"`
}

type SubmitAnswerItem struct {
	QuestionID        string  `json:"question_id"`
	SelectedOptionKey *string `json:"selected_option_key"`
	AnswerText        *string `json:"answer_text"`
	IsFlagged         *bool   `json:"is_flagged"`
}
