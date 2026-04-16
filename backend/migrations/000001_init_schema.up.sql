CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE tryout_status AS ENUM (
    'draft',
    'scheduled',
    'ongoing',
    'finished',
    'archived'
);

CREATE TYPE question_type AS ENUM (
    'multiple_choice',
    'short_text'
);

CREATE TYPE attempt_status AS ENUM (
    'ongoing',
    'submitted',
    'auto_submitted'
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE,
    full_name VARCHAR(120) NOT NULL,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    device_id TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_sessions_expiry_check CHECK (expires_at > created_at)
);

CREATE INDEX idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX idx_auth_sessions_active_user_id ON auth_sessions(user_id) WHERE revoked_at IS NULL;

CREATE TABLE tryouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(80) NOT NULL UNIQUE,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    instructions TEXT,
    status tryout_status NOT NULL DEFAULT 'draft',
    duration_minutes INTEGER NOT NULL,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    results_published_at TIMESTAMPTZ,
    max_attempts INTEGER NOT NULL DEFAULT 1,
    show_leaderboard BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tryouts_duration_check CHECK (duration_minutes > 0),
    CONSTRAINT tryouts_max_attempts_check CHECK (max_attempts > 0),
    CONSTRAINT tryouts_window_check CHECK (starts_at IS NULL OR ends_at IS NULL OR starts_at < ends_at)
);

CREATE UNIQUE INDEX uq_tryouts_single_ongoing ON tryouts(status) WHERE status = 'ongoing';

CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tryout_id UUID NOT NULL REFERENCES tryouts(id) ON DELETE CASCADE,
    code VARCHAR(80) NOT NULL,
    question_type question_type NOT NULL,
    prompt_html TEXT NOT NULL,
    image_url TEXT,
    display_order INTEGER NOT NULL,
    points NUMERIC(10,2) NOT NULL DEFAULT 1.00,
    explanation_html TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT questions_display_order_check CHECK (display_order > 0),
    CONSTRAINT questions_points_check CHECK (points >= 0)
);

CREATE UNIQUE INDEX uq_questions_tryout_code ON questions(tryout_id, code);
CREATE UNIQUE INDEX uq_questions_tryout_display_order ON questions(tryout_id, display_order);
CREATE INDEX idx_questions_tryout_id ON questions(tryout_id);

CREATE TABLE question_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_key VARCHAR(8) NOT NULL,
    option_text TEXT NOT NULL,
    display_order INTEGER NOT NULL,
    CONSTRAINT question_options_display_order_check CHECK (display_order > 0),
    CONSTRAINT question_options_option_key_check CHECK (char_length(option_key) > 0)
);

CREATE UNIQUE INDEX uq_question_options_key ON question_options(question_id, option_key);
CREATE UNIQUE INDEX uq_question_options_order ON question_options(question_id, display_order);
CREATE INDEX idx_question_options_question_id ON question_options(question_id);

CREATE TABLE question_answer_keys (
    question_id UUID PRIMARY KEY REFERENCES questions(id) ON DELETE CASCADE,
    correct_option_key VARCHAR(8) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT question_answer_keys_option_key_check CHECK (char_length(correct_option_key) > 0)
);

CREATE TABLE question_short_answer_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer_text TEXT NOT NULL,
    normalized_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT question_short_answer_variants_text_check CHECK (char_length(normalized_text) > 0)
);

CREATE UNIQUE INDEX uq_question_short_answer_variants_norm
    ON question_short_answer_variants(question_id, normalized_text);
CREATE INDEX idx_question_short_answer_variants_question_id
    ON question_short_answer_variants(question_id);

CREATE TABLE attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tryout_id UUID NOT NULL REFERENCES tryouts(id) ON DELETE CASCADE,
    status attempt_status NOT NULL DEFAULT 'ongoing',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    submitted_at TIMESTAMPTZ,
    last_synced_at TIMESTAMPTZ,
    version INTEGER NOT NULL DEFAULT 1,
    total_questions INTEGER NOT NULL DEFAULT 0,
    answered_questions INTEGER NOT NULL DEFAULT 0,
    flagged_questions INTEGER NOT NULL DEFAULT 0,
    correct_count INTEGER NOT NULL DEFAULT 0,
    wrong_count INTEGER NOT NULL DEFAULT 0,
    unanswered_count INTEGER NOT NULL DEFAULT 0,
    final_score NUMERIC(10,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT attempts_version_check CHECK (version > 0),
    CONSTRAINT attempts_total_questions_check CHECK (total_questions >= 0),
    CONSTRAINT attempts_answered_questions_check CHECK (answered_questions >= 0),
    CONSTRAINT attempts_flagged_questions_check CHECK (flagged_questions >= 0),
    CONSTRAINT attempts_correct_count_check CHECK (correct_count >= 0),
    CONSTRAINT attempts_wrong_count_check CHECK (wrong_count >= 0),
    CONSTRAINT attempts_unanswered_count_check CHECK (unanswered_count >= 0),
    CONSTRAINT attempts_expiry_check CHECK (expires_at > started_at),
    CONSTRAINT attempts_submission_state_check CHECK (
        (status = 'ongoing' AND submitted_at IS NULL)
        OR (status IN ('submitted', 'auto_submitted') AND submitted_at IS NOT NULL)
    )
);

CREATE UNIQUE INDEX uq_attempts_user_tryout ON attempts(user_id, tryout_id);
CREATE INDEX idx_attempts_tryout_status ON attempts(tryout_id, status);
CREATE INDEX idx_attempts_leaderboard
    ON attempts(tryout_id, final_score DESC, submitted_at ASC)
    WHERE status IN ('submitted', 'auto_submitted');

CREATE TABLE attempt_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_key VARCHAR(8),
    answer_text TEXT,
    normalized_answer TEXT,
    is_flagged BOOLEAN NOT NULL DEFAULT FALSE,
    answered_at TIMESTAMPTZ,
    last_saved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    client_updated_at TIMESTAMPTZ,
    server_version INTEGER NOT NULL DEFAULT 1,
    is_correct BOOLEAN,
    awarded_points NUMERIC(10,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT attempt_answers_server_version_check CHECK (server_version > 0),
    CONSTRAINT attempt_answers_single_value_check CHECK (
        NOT (selected_option_key IS NOT NULL AND answer_text IS NOT NULL)
    )
);

CREATE UNIQUE INDEX uq_attempt_answers_attempt_question ON attempt_answers(attempt_id, question_id);
CREATE INDEX idx_attempt_answers_attempt_id ON attempt_answers(attempt_id);
CREATE INDEX idx_attempt_answers_question_id ON attempt_answers(question_id);

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_tryouts_set_updated_at
BEFORE UPDATE ON tryouts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_questions_set_updated_at
BEFORE UPDATE ON questions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_question_answer_keys_set_updated_at
BEFORE UPDATE ON question_answer_keys
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_attempts_set_updated_at
BEFORE UPDATE ON attempts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_attempt_answers_set_updated_at
BEFORE UPDATE ON attempt_answers
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
