CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email text NOT NULL,
    institution text NOT NULL,
    full_name text NOT NULL,
    given_name text NOT NULL,
    role user_role NOT NULL
);

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name text NOT NULL,
    subject text NOT NULL,
    institution text NOT NULL,
    owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE user_courses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    UNIQUE (user_id, course_id)
);

CREATE TYPE comp_sci_problem_type AS ENUM ('exam', 'practice', 'turtle');

CREATE TABLE comp_sci_problem_group (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type comp_sci_problem_type NOT NULL,
    course_id BIGINT NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    name text NOT NULL,
    description text NOT NULL
);

CREATE TABLE comp_sci_problems (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    group_id BIGINT NOT NULL REFERENCES comp_sci_problem_group (id) ON DELETE CASCADE,
    type comp_sci_problem_type NOT NULL DEFAULT 'practice',
    name text NOT NULL,
    description text NOT NULL,
    problem_text text NOT NULL,
    time_limit_ms INT,
    memory_limit_mb INT
);

CREATE TABLE comp_sci_test_cases (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    problem_id BIGINT NOT NULL REFERENCES comp_sci_problems (id) ON DELETE CASCADE,
    input text NOT NULL,
    output text NOT NULL,
    correct_points INT NOT NULL,
    -- For turtle problems: URL to the ideal image
    image_url text
);

CREATE TABLE comp_sci_submissions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL REFERENCES comp_sci_problems (id) ON DELETE CASCADE,
    code text NOT NULL,
    score INT NOT NULL,
    max_score INT NOT NULL,
    passed_tests INT NOT NULL,
    total_tests INT NOT NULL,
    results_json text NOT NULL
);

-- =====================
-- Quiz Problem Types
-- =====================

CREATE TYPE quiz_problem_type AS ENUM (
    'open_ended',      -- Free-text response
    'true_false',      -- True/False question
    'mcq_single',      -- Multiple choice, one correct answer
    'mcq_multi',       -- Multiple choice, multiple correct answers
    'fill_blank'       -- Fill in the blank
);

CREATE TABLE quiz_problem_groups (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type comp_sci_problem_type NOT NULL DEFAULT 'practice',
    course_id BIGINT NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    name text NOT NULL,
    description text NOT NULL
);

CREATE TABLE quiz_problems (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    group_id BIGINT NOT NULL REFERENCES quiz_problem_groups (id) ON DELETE CASCADE,
    problem_type quiz_problem_type NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    problem_text text NOT NULL,
    points INT NOT NULL DEFAULT 1,
    -- For fill_blank: comma-separated acceptable answers
    correct_answer text
);

CREATE TABLE quiz_problem_options (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    problem_id BIGINT NOT NULL REFERENCES quiz_problems (id) ON DELETE CASCADE,
    option_text text NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT false,
    display_order INT NOT NULL DEFAULT 0
);

CREATE TABLE quiz_submissions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL REFERENCES quiz_problems (id) ON DELETE CASCADE,
    answer_text text,                    -- For open_ended, fill_blank
    selected_options BIGINT[],           -- For true_false, mcq_single, mcq_multi
    is_correct BOOLEAN,                  -- Auto-graded result (NULL for open_ended)
    score INT NOT NULL DEFAULT 0,        -- Earned points
    max_score INT NOT NULL,              -- Maximum possible points
    feedback text                        -- Teacher/AI feedback for open_ended
);

-- =====================
-- Exam Sessions
-- =====================

CREATE TYPE exam_session_status AS ENUM ('pending', 'active', 'ended');

CREATE TYPE exam_student_status AS ENUM ('assigned', 'active', 'submitted', 'discontinued');

-- Exam session created by teacher
CREATE TABLE exam_sessions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    problem_group_type text NOT NULL, -- 'comp_sci' or 'quiz'
    problem_group_id BIGINT NOT NULL, -- References either group table
    opened_by BIGINT NOT NULL REFERENCES users (id),
    duration_minutes INT NOT NULL,
    status exam_session_status NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP -- When teacher clicked Start
);

-- Students assigned to exam session (each has their own timer)
CREATE TABLE exam_session_students (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    session_id BIGINT NOT NULL REFERENCES exam_sessions (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status exam_student_status NOT NULL DEFAULT 'assigned',
    joined_at TIMESTAMP, -- When they connected via WebSocket
    ends_at TIMESTAMP, -- joined_at + duration (per-student timer)
    submitted_at TIMESTAMP, -- When they submitted
    auto_submitted BOOLEAN DEFAULT false,
    UNIQUE (session_id, user_id)
);

-- Work-in-progress saved periodically for recovery and auto-submit
CREATE TABLE exam_work_in_progress (
    id BIGSERIAL PRIMARY KEY,
    session_student_id BIGINT NOT NULL REFERENCES exam_session_students (id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL,
    problem_type text NOT NULL, -- 'comp_sci' or 'quiz'
    current_answer text NOT NULL, -- Code or answer JSON
    saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (
        session_student_id,
        problem_id,
        problem_type
    )
);
-- =====================
-- Course Materials
-- =====================

CREATE TABLE course_materials (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    course_id BIGINT NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    title text NOT NULL,
    file_url text NOT NULL,
    file_type text NOT NULL,
    uploaded_by BIGINT REFERENCES users (id) ON DELETE SET NULL
);