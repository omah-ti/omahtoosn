package tryout

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
)

type Service struct {
	pool *pgxpool.Pool
	repo *Repository
}

func NewService(pool *pgxpool.Pool, repo *Repository) *Service {
	return &Service{pool: pool, repo: repo}
}

func (s *Service) GetCurrentTryout(ctx context.Context, userID string) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}
	if err := s.syncExpiredAttempt(ctx, userID, tryout.ID); err != nil {
		return nil, err
	}
	var attemptData any
	attempt, err := s.repo.GetAttemptByUserTryout(ctx, s.pool, userID, tryout.ID)
	if err == nil {
		attemptData = compactAttempt(attempt)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, httpx.Internal("failed to fetch attempt state")
	}
	return map[string]any{
		"tryout":      tryout,
		"attempt":     attemptData,
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) StartCurrentTryout(ctx context.Context, userID string) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}
	if err := validateTryoutWindow(tryout); err != nil {
		return nil, err
	}

	tx, err := s.repo.Begin(ctx, s.pool)
	if err != nil {
		return nil, httpx.Internal("failed to start attempt transaction")
	}
	defer tx.Rollback(ctx)

	attempt, err := s.repo.LockAttemptByUserTryout(ctx, tx, userID, tryout.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			expiresAt := time.Now().UTC().Add(time.Duration(tryout.DurationMinutes) * time.Minute)
			if tryout.EndsAt != nil && tryout.EndsAt.Before(expiresAt) {
				expiresAt = *tryout.EndsAt
			}
			attempt, err = s.repo.CreateAttempt(ctx, tx, userID, tryout.ID, expiresAt, tryout.QuestionCount)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					attempt, err = s.repo.LockAttemptByUserTryout(ctx, tx, userID, tryout.ID)
					if err != nil {
						return nil, httpx.Internal("failed to recover concurrent attempt")
					}
				} else {
					return nil, httpx.Internal("failed to create attempt")
				}
			}
		} else {
			return nil, httpx.Internal("failed to fetch attempt")
		}
	}

	if attempt.Status == "ongoing" && time.Now().UTC().After(attempt.ExpiresAt) {
		attempt, err = s.finalizeAttemptTx(ctx, tx, attempt, true)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, httpx.Internal("failed to commit attempt transaction")
	}

	return map[string]any{
		"tryout":      tryout,
		"attempt":     compactAttempt(attempt),
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) GetCurrentAttempt(ctx context.Context, userID string) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}
	if err := s.syncExpiredAttempt(ctx, userID, tryout.ID); err != nil {
		return nil, err
	}
	attempt, err := s.repo.GetAttemptByUserTryout(ctx, s.pool, userID, tryout.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no attempt found")
		}
		return nil, httpx.Internal("failed to fetch attempt")
	}
	if attempt.Status != "ongoing" {
		return nil, httpx.Conflict("attempt is not active")
	}
	questions, err := s.repo.ListQuestions(ctx, s.pool, tryout.ID)
	if err != nil {
		return nil, httpx.Internal("failed to fetch questions")
	}
	answers, err := s.repo.ListAnswers(ctx, s.pool, attempt.ID)
	if err != nil {
		return nil, httpx.Internal("failed to fetch answers")
	}
	return map[string]any{
		"tryout":      tryout,
		"attempt":     compactAttempt(attempt),
		"questions":   questions,
		"answers":     answers,
		"navigator":   buildNavigator(questions, answers),
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) SaveCurrentAnswers(ctx context.Context, userID string, req SaveAnswersRequest) (map[string]any, error) {
	if len(req.Answers) == 0 {
		return nil, httpx.BadRequest("answers payload is required")
	}
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}

	tx, err := s.repo.Begin(ctx, s.pool)
	if err != nil {
		return nil, httpx.Internal("failed to start save transaction")
	}
	defer tx.Rollback(ctx)

	attempt, err := s.repo.LockAttemptByUserTryout(ctx, tx, userID, tryout.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no attempt found")
		}
		return nil, httpx.Internal("failed to fetch attempt")
	}
	if attempt.Status != "ongoing" {
		return nil, httpx.Conflict("attempt is not active")
	}
	if time.Now().UTC().After(attempt.ExpiresAt) {
		attempt, err = s.finalizeAttemptTx(ctx, tx, attempt, true)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, httpx.Internal("failed to commit expired attempt")
		}
		return nil, httpx.Conflict("attempt is expired")
	}

	questionIDs, err := distinctQuestionIDs(req.Answers)
	if err != nil {
		return nil, err
	}
	metas, err := s.repo.GetQuestionMetaByIDs(ctx, tx, tryout.ID, questionIDs)
	if err != nil {
		return nil, httpx.Internal("failed to validate questions")
	}
	if len(metas) != len(questionIDs) {
		return nil, httpx.BadRequest("one or more question_ids are invalid")
	}

	now := time.Now().UTC()
	for _, item := range req.Answers {
		selectedOptionKey, answerText, normalizedAnswer, answeredAt, err := normalizeSaveItem(item, metas[item.QuestionID], now)
		if err != nil {
			return nil, err
		}
		clientUpdatedAt := now
		if item.ClientUpdatedAt != nil {
			clientUpdatedAt = item.ClientUpdatedAt.UTC()
		}
		if err := s.repo.UpsertAnswer(ctx, tx, attempt.ID, item.QuestionID, selectedOptionKey, answerText, normalizedAnswer, item.IsFlagged, answeredAt, &clientUpdatedAt); err != nil {
			return nil, httpx.Internal("failed to save answer")
		}
	}

	attempt, err = s.repo.RefreshAttemptStats(ctx, tx, attempt.ID)
	if err != nil {
		return nil, httpx.Internal("failed to refresh attempt stats")
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, httpx.Internal("failed to commit save transaction")
	}
	return map[string]any{
		"attempt":     compactAttempt(attempt),
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) SubmitCurrentAttempt(ctx context.Context, userID string, req SubmitAttemptRequest) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}

	tx, err := s.repo.Begin(ctx, s.pool)
	if err != nil {
		return nil, httpx.Internal("failed to start submit transaction")
	}
	defer tx.Rollback(ctx)

	attempt, err := s.repo.LockAttemptByUserTryout(ctx, tx, userID, tryout.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no attempt found")
		}
		return nil, httpx.Internal("failed to fetch attempt")
	}

	if attempt.Status == "submitted" || attempt.Status == "auto_submitted" {
		if err := tx.Commit(ctx); err != nil {
			return nil, httpx.Internal("failed to finalize response")
		}
		result, err := s.repo.GetResultByUserTryout(ctx, s.pool, userID, tryout.ID)
		if err != nil {
			return nil, httpx.Internal("failed to fetch result")
		}
		return map[string]any{
			"attempt":     compactAttempt(attempt),
			"result":      result,
			"server_time": time.Now().UTC(),
		}, nil
	}

	if len(req.FinalAnswers) > 0 {
		questionIDs, err := distinctSubmitQuestionIDs(req.FinalAnswers)
		if err != nil {
			return nil, err
		}
		metas, err := s.repo.GetQuestionMetaByIDs(ctx, tx, tryout.ID, questionIDs)
		if err != nil {
			return nil, httpx.Internal("failed to validate final answers")
		}
		if len(metas) != len(questionIDs) {
			return nil, httpx.BadRequest("one or more final question_ids are invalid")
		}
		now := time.Now().UTC()
		for _, item := range req.FinalAnswers {
			normalizedItem := SaveAnswerItem{
				QuestionID:        item.QuestionID,
				SelectedOptionKey: item.SelectedOptionKey,
				AnswerText:        item.AnswerText,
				IsFlagged:         false,
				ClientUpdatedAt:   &now,
			}
			if item.IsFlagged != nil {
				normalizedItem.IsFlagged = *item.IsFlagged
			}
			selectedOptionKey, answerText, normalizedAnswer, answeredAt, err := normalizeSaveItem(normalizedItem, metas[item.QuestionID], now)
			if err != nil {
				return nil, err
			}
			if err := s.repo.UpsertAnswer(ctx, tx, attempt.ID, item.QuestionID, selectedOptionKey, answerText, normalizedAnswer, normalizedItem.IsFlagged, answeredAt, &now); err != nil {
				return nil, httpx.Internal("failed to save final answer")
			}
		}
	}

	attempt, err = s.finalizeAttemptTx(ctx, tx, attempt, time.Now().UTC().After(attempt.ExpiresAt))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, httpx.Internal("failed to commit submit transaction")
	}
	result, err := s.repo.GetResultByUserTryout(ctx, s.pool, userID, tryout.ID)
	if err != nil {
		return nil, httpx.Internal("failed to fetch result")
	}
	return map[string]any{
		"attempt":     compactAttempt(attempt),
		"result":      result,
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) GetCurrentResult(ctx context.Context, userID string) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}
	if err := s.syncExpiredAttempt(ctx, userID, tryout.ID); err != nil {
		return nil, err
	}
	result, err := s.repo.GetResultByUserTryout(ctx, s.pool, userID, tryout.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("result not found")
		}
		return nil, httpx.Internal("failed to fetch result")
	}
	return map[string]any{
		"tryout":      tryout,
		"result":      result,
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) GetCurrentLeaderboard(ctx context.Context, userID string, limit, offset int) (map[string]any, error) {
	tryout, err := s.repo.GetCurrentTryout(ctx, s.pool)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httpx.NotFound("no ongoing tryout")
		}
		return nil, httpx.Internal("failed to fetch current tryout")
	}
	if !tryout.ShowLeaderboard {
		return nil, httpx.Forbidden("leaderboard is disabled")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	entries, total, err := s.repo.ListLeaderboard(ctx, s.pool, tryout.ID, limit, offset)
	if err != nil {
		return nil, httpx.Internal("failed to fetch leaderboard")
	}
	var me any
	result, err := s.repo.GetResultByUserTryout(ctx, s.pool, userID, tryout.ID)
	if err == nil {
		me = result
	}
	return map[string]any{
		"tryout":      tryout,
		"entries":     entries,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
		"me":          me,
		"server_time": time.Now().UTC(),
	}, nil
}

func (s *Service) syncExpiredAttempt(ctx context.Context, userID, tryoutID string) error {
	tx, err := s.repo.Begin(ctx, s.pool)
	if err != nil {
		return httpx.Internal("failed to start expiry transaction")
	}
	defer tx.Rollback(ctx)
	attempt, err := s.repo.LockAttemptByUserTryout(ctx, tx, userID, tryoutID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return httpx.Internal("failed to fetch attempt state")
	}
	if attempt.Status != "ongoing" || !time.Now().UTC().After(attempt.ExpiresAt) {
		return tx.Commit(ctx)
	}
	if _, err := s.finalizeAttemptTx(ctx, tx, attempt, true); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return httpx.Internal("failed to commit expiry transaction")
	}
	return nil
}

func (s *Service) finalizeAttemptTx(ctx context.Context, tx pgx.Tx, attempt Attempt, auto bool) (Attempt, error) {
	questions, err := s.repo.GetScoringQuestions(ctx, tx, attempt.TryoutID)
	if err != nil {
		return Attempt{}, httpx.Internal("failed to fetch scoring questions")
	}
	answers, err := s.repo.GetScoringAnswers(ctx, tx, attempt.ID)
	if err != nil {
		return Attempt{}, httpx.Internal("failed to fetch scoring answers")
	}
	answeredQuestions := 0
	correctCount := 0
	wrongCount := 0
	unansweredCount := 0
	finalScore := 0.0

	for _, question := range questions {
		answer, exists := answers[question.ID]
		if !exists || isBlankAnswer(answer) {
			unansweredCount++
			continue
		}
		answeredQuestions++
		isCorrect := false
		switch question.QuestionType {
		case "multiple_choice":
			isCorrect = answer.SelectedOptionKey != nil && strings.EqualFold(strings.TrimSpace(*answer.SelectedOptionKey), strings.TrimSpace(question.CorrectOptionKey))
		case "short_text":
			isCorrect = answer.NormalizedAnswer != nil && question.Variants[*answer.NormalizedAnswer]
		}
		if isCorrect {
			correctCount++
			finalScore += question.Points
			flag := true
			if err := s.repo.UpdateAnswerScore(ctx, tx, attempt.ID, question.ID, &flag, question.Points); err != nil {
				return Attempt{}, httpx.Internal("failed to update answer score")
			}
		} else {
			wrongCount++
			flag := false
			if err := s.repo.UpdateAnswerScore(ctx, tx, attempt.ID, question.ID, &flag, 0); err != nil {
				return Attempt{}, httpx.Internal("failed to update answer score")
			}
		}
	}
	status := "submitted"
	if auto {
		status = "auto_submitted"
	}
	return s.repo.FinalizeAttempt(ctx, tx, attempt.ID, status, len(questions), answeredQuestions, correctCount, wrongCount, unansweredCount, math.Round(finalScore*100)/100)
}

func compactAttempt(attempt Attempt) map[string]any {
	return map[string]any{
		"id":                 attempt.ID,
		"status":             attempt.Status,
		"started_at":         attempt.StartedAt,
		"expires_at":         attempt.ExpiresAt,
		"submitted_at":       attempt.SubmittedAt,
		"last_synced_at":     attempt.LastSyncedAt,
		"version":            attempt.Version,
		"total_questions":    attempt.TotalQuestions,
		"answered_questions": attempt.AnsweredQuestions,
		"flagged_questions":  attempt.FlaggedQuestions,
		"correct_count":      attempt.CorrectCount,
		"wrong_count":        attempt.WrongCount,
		"unanswered_count":   attempt.UnansweredCount,
		"final_score":        attempt.FinalScore,
	}
}

func validateTryoutWindow(tryout Tryout) error {
	now := time.Now().UTC()
	if tryout.StartsAt != nil && now.Before(*tryout.StartsAt) {
		return httpx.Forbidden("tryout has not started")
	}
	if tryout.EndsAt != nil && now.After(*tryout.EndsAt) {
		return httpx.Conflict("tryout is closed")
	}
	return nil
}

func distinctQuestionIDs(items []SaveAnswerItem) ([]string, error) {
	seen := map[string]bool{}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		id := strings.TrimSpace(item.QuestionID)
		if id == "" {
			return nil, httpx.BadRequest("question_id is required")
		}
		if seen[id] {
			return nil, httpx.BadRequest("duplicate question_id in payload")
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

func distinctSubmitQuestionIDs(items []SubmitAnswerItem) ([]string, error) {
	seen := map[string]bool{}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		id := strings.TrimSpace(item.QuestionID)
		if id == "" {
			return nil, httpx.BadRequest("question_id is required")
		}
		if seen[id] {
			return nil, httpx.BadRequest("duplicate question_id in payload")
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

func normalizeSaveItem(item SaveAnswerItem, meta QuestionMeta, now time.Time) (*string, *string, *string, *time.Time, error) {
	var selectedOptionKey *string
	var answerText *string
	var normalizedAnswer *string
	var answeredAt *time.Time

	selectedValue := ""
	if item.SelectedOptionKey != nil {
		selectedValue = strings.ToUpper(strings.TrimSpace(*item.SelectedOptionKey))
	}
	answerValue := ""
	if item.AnswerText != nil {
		answerValue = strings.TrimSpace(*item.AnswerText)
	}
	if selectedValue != "" && answerValue != "" {
		return nil, nil, nil, nil, httpx.BadRequest("selected_option_key and answer_text cannot be filled together")
	}

	switch meta.QuestionType {
	case "multiple_choice":
		if answerValue != "" {
			return nil, nil, nil, nil, httpx.BadRequest("answer_text is invalid for multiple_choice")
		}
		if selectedValue != "" {
			if !meta.ValidOptionKeys[selectedValue] {
				return nil, nil, nil, nil, httpx.BadRequest("selected_option_key is invalid")
			}
			selectedOptionKey = &selectedValue
			answeredAt = &now
		}
	case "short_text":
		if selectedValue != "" {
			return nil, nil, nil, nil, httpx.BadRequest("selected_option_key is invalid for short_text")
		}
		if answerValue != "" {
			answerText = &answerValue
			normalized := normalizeAnswer(answerValue)
			normalizedAnswer = &normalized
			answeredAt = &now
		}
	default:
		return nil, nil, nil, nil, httpx.BadRequest("unsupported question type")
	}
	return selectedOptionKey, answerText, normalizedAnswer, answeredAt, nil
}

func normalizeAnswer(input string) string {
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(input)))
	return strings.Join(parts, " ")
}

func buildNavigator(questions []Question, answers []AttemptAnswer) []map[string]any {
	answerMap := map[string]AttemptAnswer{}
	for _, answer := range answers {
		answerMap[answer.QuestionID] = answer
	}
	navigator := make([]map[string]any, 0, len(questions))
	for _, question := range questions {
		answer, exists := answerMap[question.ID]
		isAnswered := false
		isFlagged := false
		if exists {
			isFlagged = answer.IsFlagged
			if answer.SelectedOptionKey != nil && strings.TrimSpace(*answer.SelectedOptionKey) != "" {
				isAnswered = true
			}
			if answer.AnswerText != nil && strings.TrimSpace(*answer.AnswerText) != "" {
				isAnswered = true
			}
		}
		navigator = append(navigator, map[string]any{
			"question_id":   question.ID,
			"display_order": question.DisplayOrder,
			"is_answered":   isAnswered,
			"is_flagged":    isFlagged,
		})
	}
	return navigator
}

func isBlankAnswer(answer ScoringAnswer) bool {
	return (answer.SelectedOptionKey == nil || strings.TrimSpace(*answer.SelectedOptionKey) == "") && (answer.NormalizedAnswer == nil || strings.TrimSpace(*answer.NormalizedAnswer) == "")
}
