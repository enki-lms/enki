-- =====================
-- Users
-- =====================

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: ListUsersByInstitution :many
SELECT * FROM users WHERE institution = $1 ORDER BY id;

-- name: ListStudentsByInstitution :many
SELECT *
FROM users
WHERE
    institution = $1
    AND role = 'student'
ORDER BY id;

-- name: CreateUser :one
INSERT INTO
    users (
        email,
        institution,
        full_name,
        given_name,
        role
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    institution = $3,
    full_name = $4,
    given_name = $5,
    role = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- =====================
-- Courses
-- =====================

-- name: GetCourse :one
SELECT * FROM courses WHERE id = $1 LIMIT 1;

-- name: ListCourses :many
SELECT * FROM courses ORDER BY id;

-- name: ListCoursesByInstitution :many
SELECT * FROM courses WHERE institution = $1 ORDER BY id;

-- name: CreateCourse :one
INSERT INTO
    courses (
        name,
        subject,
        institution,
        owner_id
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCourse :one
UPDATE courses
SET
    name = $2,
    subject = $3,
    institution = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: ListCoursesByUser :many
SELECT c.*
FROM courses c
    JOIN user_courses uc ON c.id = uc.course_id
WHERE
    uc.user_id = $1
ORDER BY c.id;

-- name: ListCoursesByOwner :many
SELECT * FROM courses WHERE owner_id = $1 ORDER BY id;

-- =====================
-- User Courses (Enrollments)
-- =====================

-- name: GetUserCourse :one
SELECT * FROM user_courses WHERE id = $1 LIMIT 1;

-- name: ListUserCoursesByCourse :many
SELECT * FROM user_courses WHERE course_id = $1 ORDER BY id;

-- name: ListUserCoursesByUser :many
SELECT * FROM user_courses WHERE user_id = $1 ORDER BY id;

-- name: CreateUserCourse :one
INSERT INTO
    user_courses (user_id, course_id)
VALUES ($1, $2) RETURNING *;

-- name: DeleteUserCourse :exec
DELETE FROM user_courses WHERE id = $1;

-- name: DeleteUserCourseByUserAndCourse :exec
DELETE FROM user_courses WHERE user_id = $1 AND course_id = $2;

-- =====================
-- Computer Science Problem Groups
-- =====================

-- name: GetCompSciProblemGroup :one
SELECT * FROM comp_sci_problem_group WHERE id = $1 LIMIT 1;

-- name: ListCompSciProblemGroups :many
SELECT * FROM comp_sci_problem_group ORDER BY id;

-- name: ListCompSciProblemGroupsByCourse :many
SELECT *
FROM comp_sci_problem_group
WHERE
    course_id = $1
ORDER BY id;

-- name: ListCompSciProblemGroupsByType :many
SELECT * FROM comp_sci_problem_group WHERE type = $1 ORDER BY id;

-- name: CreateCompSciProblemGroup :one
INSERT INTO
    comp_sci_problem_group (
        type,
        course_id,
        name,
        description
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCompSciProblemGroup :one
UPDATE comp_sci_problem_group
SET
    type = $2,
    course_id = $3,
    name = $4,
    description = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteCompSciProblemGroup :exec
DELETE FROM comp_sci_problem_group WHERE id = $1;

-- =====================
-- Computer Science Problems
-- =====================

-- name: GetCompSciProblem :one
SELECT * FROM comp_sci_problems WHERE id = $1 LIMIT 1;

-- name: ListCompSciProblems :many
SELECT * FROM comp_sci_problems ORDER BY id;

-- name: ListCompSciProblemsByGroup :many
SELECT * FROM comp_sci_problems WHERE group_id = $1 ORDER BY id;

-- name: CreateCompSciProblem :one
INSERT INTO
    comp_sci_problems (
        group_id,
        name,
        description,
        problem_text,
        time_limit_ms,
        memory_limit_mb
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateCompSciProblem :one
UPDATE comp_sci_problems
SET
    group_id = $2,
    name = $3,
    description = $4,
    problem_text = $5,
    time_limit_ms = $6,
    memory_limit_mb = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteCompSciProblem :exec
DELETE FROM comp_sci_problems WHERE id = $1;

-- =====================
-- Computer Science Test Cases
-- =====================

-- name: GetCompSciTestCase :one
SELECT * FROM comp_sci_test_cases WHERE id = $1 LIMIT 1;

-- name: ListCompSciTestCases :many
SELECT * FROM comp_sci_test_cases ORDER BY id;

-- name: ListCompSciTestCasesByProblem :many
SELECT * FROM comp_sci_test_cases WHERE problem_id = $1 ORDER BY id;

-- name: CreateCompSciTestCase :one
INSERT INTO
    comp_sci_test_cases (
        problem_id,
        input,
        output,
        correct_points
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCompSciTestCase :one
UPDATE comp_sci_test_cases
SET
    problem_id = $2,
    input = $3,
    output = $4,
    correct_points = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteCompSciTestCase :exec
DELETE FROM comp_sci_test_cases WHERE id = $1;

-- =====================
-- Computer Science Submissions
-- =====================

-- name: GetCompSciSubmission :one
SELECT * FROM comp_sci_submissions WHERE id = $1 LIMIT 1;

-- name: ListCompSciSubmissionsByProblem :many
SELECT *
FROM comp_sci_submissions
WHERE
    problem_id = $1
ORDER BY created_at DESC;

-- name: ListCompSciSubmissionsByUser :many
SELECT *
FROM comp_sci_submissions
WHERE
    user_id = $1
ORDER BY created_at DESC;

-- name: ListCompSciSubmissionsByUserAndProblem :many
SELECT *
FROM comp_sci_submissions
WHERE
    user_id = $1
    AND problem_id = $2
ORDER BY created_at DESC;

-- name: ListCompSciSubmissionsByUserWithLimit :many
SELECT *
FROM comp_sci_submissions
WHERE
    user_id = $1
ORDER BY created_at DESC
LIMIT CASE WHEN sqlc.narg('limit_count')::int IS NULL THEN 2147483647 ELSE sqlc.narg('limit_count')::int END;

-- name: CreateCompSciSubmission :one
INSERT INTO
    comp_sci_submissions (
        user_id,
        problem_id,
        code,
        score,
        max_score,
        passed_tests,
        total_tests,
        results_json
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING *;

-- =====================
-- Quiz Problem Groups
-- =====================

-- name: GetQuizProblemGroup :one
SELECT * FROM quiz_problem_groups WHERE id = $1 LIMIT 1;

-- name: ListQuizProblemGroups :many
SELECT * FROM quiz_problem_groups ORDER BY id;

-- name: ListQuizProblemGroupsByCourse :many
SELECT * FROM quiz_problem_groups WHERE course_id = $1 ORDER BY id;

-- name: CreateQuizProblemGroup :one
INSERT INTO
    quiz_problem_groups (
        type,
        course_id,
        name,
        description
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateQuizProblemGroup :one
UPDATE quiz_problem_groups
SET
    type = $2,
    course_id = $3,
    name = $4,
    description = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteQuizProblemGroup :exec
DELETE FROM quiz_problem_groups WHERE id = $1;

-- =====================
-- Quiz Problems
-- =====================

-- name: GetQuizProblem :one
SELECT * FROM quiz_problems WHERE id = $1 LIMIT 1;

-- name: ListQuizProblems :many
SELECT * FROM quiz_problems ORDER BY id;

-- name: ListQuizProblemsByGroup :many
SELECT * FROM quiz_problems WHERE group_id = $1 ORDER BY id;

-- name: CreateQuizProblem :one
INSERT INTO
    quiz_problems (
        group_id,
        problem_type,
        name,
        description,
        problem_text,
        points,
        correct_answer
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateQuizProblem :one
UPDATE quiz_problems
SET
    group_id = $2,
    problem_type = $3,
    name = $4,
    description = $5,
    problem_text = $6,
    points = $7,
    correct_answer = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteQuizProblem :exec
DELETE FROM quiz_problems WHERE id = $1;

-- =====================
-- Quiz Problem Options
-- =====================

-- name: GetQuizProblemOption :one
SELECT * FROM quiz_problem_options WHERE id = $1 LIMIT 1;

-- name: ListQuizProblemOptions :many
SELECT * FROM quiz_problem_options ORDER BY id;

-- name: ListQuizProblemOptionsByProblem :many
SELECT *
FROM quiz_problem_options
WHERE
    problem_id = $1
ORDER BY display_order;

-- name: CreateQuizProblemOption :one
INSERT INTO
    quiz_problem_options (
        problem_id,
        option_text,
        is_correct,
        display_order
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateQuizProblemOption :one
UPDATE quiz_problem_options
SET
    problem_id = $2,
    option_text = $3,
    is_correct = $4,
    display_order = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteQuizProblemOption :exec
DELETE FROM quiz_problem_options WHERE id = $1;

-- name: DeleteQuizProblemOptionsByProblem :exec
DELETE FROM quiz_problem_options WHERE problem_id = $1;

-- =====================
-- Quiz Submissions
-- =====================

-- name: GetQuizSubmission :one
SELECT * FROM quiz_submissions WHERE id = $1 LIMIT 1;

-- name: ListQuizSubmissionsByProblem :many
SELECT *
FROM quiz_submissions
WHERE
    problem_id = $1
ORDER BY created_at DESC;

-- name: ListQuizSubmissionsByUser :many
SELECT *
FROM quiz_submissions
WHERE
    user_id = $1
ORDER BY created_at DESC;

-- name: ListQuizSubmissionsByUserAndProblem :many
SELECT *
FROM quiz_submissions
WHERE
    user_id = $1
    AND problem_id = $2
ORDER BY created_at DESC;

-- name: ListQuizSubmissionsByUserWithLimit :many
SELECT *
FROM quiz_submissions
WHERE
    user_id = $1
ORDER BY created_at DESC
LIMIT CASE WHEN sqlc.narg('limit_count')::int IS NULL THEN 2147483647 ELSE sqlc.narg('limit_count')::int END;

-- name: CreateQuizSubmission :one
INSERT INTO
    quiz_submissions (
        user_id,
        problem_id,
        answer_text,
        selected_options,
        is_correct,
        score,
        max_score,
        feedback
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING *;

-- name: UpdateQuizSubmissionFeedback :one
UPDATE quiz_submissions
SET
    is_correct = $2,
    score = $3,
    feedback = $4
WHERE
    id = $1 RETURNING *;

-- =====================
-- Exam Sessions
-- =====================

-- name: GetExamSession :one
SELECT * FROM exam_sessions WHERE id = $1 LIMIT 1;

-- name: ListExamSessionsByGroup :many
SELECT *
FROM exam_sessions
WHERE
    problem_group_type = $1
    AND problem_group_id = $2
ORDER BY created_at DESC;

-- name: ListExamSessionsByTeacher :many
SELECT *
FROM exam_sessions
WHERE
    opened_by = $1
ORDER BY created_at DESC;

-- name: ListActiveExamSessions :many
SELECT *
FROM exam_sessions
WHERE
    status = 'active'
ORDER BY created_at DESC;

-- name: CreateExamSession :one
INSERT INTO
    exam_sessions (
        problem_group_type,
        problem_group_id,
        opened_by,
        duration_minutes
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: StartExamSession :one
UPDATE exam_sessions
SET
    status = 'active',
    started_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: EndExamSession :one
UPDATE exam_sessions
SET
    status = 'ended',
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: DeleteExamSession :exec
DELETE FROM exam_sessions WHERE id = $1;

-- =====================
-- Exam Session Students
-- =====================

-- name: GetExamSessionStudent :one
SELECT * FROM exam_session_students WHERE id = $1 LIMIT 1;

-- name: GetExamSessionStudentBySessionAndUser :one
SELECT *
FROM exam_session_students
WHERE
    session_id = $1
    AND user_id = $2
LIMIT 1;

-- name: ListExamSessionStudents :many
SELECT *
FROM exam_session_students
WHERE
    session_id = $1
ORDER BY id;

-- name: ListActiveExamSessionsForUser :many
SELECT es.*
FROM
    exam_sessions es
    JOIN exam_session_students ess ON es.id = ess.session_id
WHERE
    ess.user_id = $1
    AND es.status = 'active'
    AND ess.status IN ('assigned', 'active')
ORDER BY es.created_at DESC;

-- name: CreateExamSessionStudent :one
INSERT INTO
    exam_session_students (session_id, user_id)
VALUES ($1, $2) RETURNING *;

-- name: JoinExamSession :one
UPDATE exam_session_students
SET
    status = 'active',
    joined_at = CURRENT_TIMESTAMP,
    ends_at = $2
WHERE
    id = $1 RETURNING *;

-- name: SubmitExamSession :one
UPDATE exam_session_students
SET
    status = 'submitted',
    submitted_at = CURRENT_TIMESTAMP,
    auto_submitted = $2
WHERE
    id = $1 RETURNING *;

-- name: DiscontinueExamSessionStudent :one
UPDATE exam_session_students
SET
    status = 'discontinued'
WHERE
    id = $1 RETURNING *;

-- name: ListExamStudentsNeedingAutoSubmit :many
SELECT *
FROM exam_session_students
WHERE
    status = 'active'
    AND ends_at < CURRENT_TIMESTAMP;

-- =====================
-- Exam Work In Progress
-- =====================

-- name: GetExamWorkInProgress :one
SELECT *
FROM exam_work_in_progress
WHERE
    session_student_id = $1
    AND problem_id = $2
    AND problem_type = $3
LIMIT 1;

-- name: ListExamWorkInProgressByStudent :many
SELECT * FROM exam_work_in_progress WHERE session_student_id = $1;

-- name: UpsertExamWorkInProgress :one
INSERT INTO
    exam_work_in_progress (
        session_student_id,
        problem_id,
        problem_type,
        current_answer
    )
VALUES ($1, $2, $3, $4) ON CONFLICT (
        session_student_id,
        problem_id,
        problem_type
    ) DO
UPDATE
SET
    current_answer = $4,
    saved_at = CURRENT_TIMESTAMP RETURNING *;

-- name: DeleteExamWorkInProgressByStudent :exec
DELETE FROM exam_work_in_progress WHERE session_student_id = $1;