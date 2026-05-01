DROP TRIGGER IF EXISTS trg_attempt_answers_set_updated_at ON attempt_answers;
DROP TRIGGER IF EXISTS trg_attempts_set_updated_at ON attempts;
DROP TRIGGER IF EXISTS trg_question_answer_keys_set_updated_at ON question_answer_keys;
DROP TRIGGER IF EXISTS trg_questions_set_updated_at ON questions;
DROP TRIGGER IF EXISTS trg_tryouts_set_updated_at ON tryouts;
DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;

DROP TABLE IF EXISTS attempt_answers;
DROP TABLE IF EXISTS attempts;
DROP TABLE IF EXISTS question_short_answer_variants;
DROP TABLE IF EXISTS question_answer_keys;
DROP TABLE IF EXISTS question_options;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS tryouts;
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TYPE IF EXISTS attempt_status;
DROP TYPE IF EXISTS question_type;
DROP TYPE IF EXISTS tryout_status;
