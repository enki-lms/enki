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
    institution text NOT NULL
);

CREATE TABLE user_courses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    UNIQUE (user_id, course_id)
);

CREATE TYPE comp_sci_problem_type AS ENUM ('exam', 'practice');

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
    correct_points INT NOT NULL
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