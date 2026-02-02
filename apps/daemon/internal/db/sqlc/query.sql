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
INSERT INTO courses (name, institution) VALUES ($1, $2) RETURNING *;

-- name: UpdateCourse :one
UPDATE courses
SET
    name = $2,
    institution = $3,
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