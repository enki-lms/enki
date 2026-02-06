-- =====================
-- Users Seed Data
-- =====================

INSERT INTO
    users (
        email,
        institution,
        full_name,
        given_name,
        role
    )
VALUES (
        'teacher@enki.com',
        'Enki University',
        'Professor Enki',
        'Enki',
        'teacher'
    ),
    (
        'student@enki.com',
        'Enki University',
        'Student Learner',
        'Student',
        'student'
    ),
    (
        'admin@enki.com',
        'Enki University',
        'Admin User',
        'Admin',
        'admin'
    );

-- =====================
-- Courses Seed Data
-- =====================

INSERT INTO
    courses (name, subject, institution)
VALUES (
        'CS101: Introduction to Computer Science',
        'Computer Science',
        'Enki University'
    );

-- =====================
-- Enrollments Seed Data
-- =====================

INSERT INTO
    user_courses (user_id, course_id)
VALUES (
        (
            SELECT id
            FROM users
            WHERE
                email = 'student@enki.com'
        ),
        (
            SELECT id
            FROM courses
            WHERE
                name = 'CS101: Introduction to Computer Science'
        )
    );

INSERT INTO
    user_courses (user_id, course_id)
VALUES (
        (
            SELECT id
            FROM users
            WHERE
                email = 'teacher@enki.com'
        ),
        (
            SELECT id
            FROM courses
            WHERE
                name = 'CS101: Introduction to Computer Science'
        )
    );

-- =====================
-- CS Problem Groups
-- =====================

-- Practice Group
INSERT INTO
    comp_sci_problem_group (
        course_id,
        type,
        name,
        description
    )
VALUES (
        (
            SELECT id
            FROM courses
            WHERE
                name = 'CS101: Introduction to Computer Science'
        ),
        'practice',
        'Week 1: Basics',
        'Introduction to input/output and simple arithmetic.'
    );

-- Exam Group
INSERT INTO
    comp_sci_problem_group (
        course_id,
        type,
        name,
        description
    )
VALUES (
        (
            SELECT id
            FROM courses
            WHERE
                name = 'CS101: Introduction to Computer Science'
        ),
        'exam',
        'Midterm Exam',
        'Midterm examination covering basic concepts.'
    );

-- =====================
-- CS Problems
-- =====================

-- Problem 1: A + B (Practice)
INSERT INTO
    comp_sci_problems (
        group_id,
        name,
        description,
        problem_text,
        time_limit_ms,
        memory_limit_mb
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problem_group
            WHERE
                name = 'Week 1: Basics'
        ),
        'A + B',
        'Read two integers and print their sum.',
        'Read two integers from input and print their sum. The input contains one line with two space-separated integers.',
        1000,
        128
    );

-- Problem 2: Factorial (Practice)
INSERT INTO
    comp_sci_problems (
        group_id,
        name,
        description,
        problem_text,
        time_limit_ms,
        memory_limit_mb
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problem_group
            WHERE
                name = 'Week 1: Basics'
        ),
        'Factorial',
        'Calculate the factorial of N.',
        'Read an integer N (0 <= N <= 10) from input and print N!. Example: Input 5, Output 120.',
        1000,
        128
    );

-- Problem 3: Fibonacci (Exam)
INSERT INTO
    comp_sci_problems (
        group_id,
        name,
        description,
        problem_text,
        time_limit_ms,
        memory_limit_mb
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problem_group
            WHERE
                name = 'Midterm Exam'
        ),
        'Fibonacci Sequence',
        'Print the Nth Fibonacci number.',
        'Read an integer N (0 <= N <= 30) from input and print the Nth Fibonacci number. F(0) = 0, F(1) = 1.',
        1000,
        128
    );

-- =====================
-- CS Test Cases
-- =====================

-- Tests for A + B
INSERT INTO
    comp_sci_test_cases (
        problem_id,
        input,
        output,
        correct_points
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'A + B'
        ),
        '3 5',
        '8',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'A + B'
        ),
        '10 -5',
        '5',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'A + B'
        ),
        '100 200',
        '300',
        10
    );

-- Tests for Factorial
INSERT INTO
    comp_sci_test_cases (
        problem_id,
        input,
        output,
        correct_points
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Factorial'
        ),
        '5',
        '120',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Factorial'
        ),
        '0',
        '1',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Factorial'
        ),
        '3',
        '6',
        10
    );

-- Tests for Fibonacci
INSERT INTO
    comp_sci_test_cases (
        problem_id,
        input,
        output,
        correct_points
    )
VALUES (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Fibonacci Sequence'
        ),
        '0',
        '0',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Fibonacci Sequence'
        ),
        '1',
        '1',
        10
    ),
    (
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'Fibonacci Sequence'
        ),
        '10',
        '55',
        10
    );

-- =====================
-- Quiz Problem Groups
-- =====================

INSERT INTO
    quiz_problem_groups (
        course_id,
        type,
        name,
        description
    )
VALUES (
        (
            SELECT id
            FROM courses
            WHERE
                name = 'CS101: Introduction to Computer Science'
        ),
        'practice',
        'Week 1 Quiz',
        'Test your understanding of basic concepts.'
    );

-- =====================
-- Quiz Problems
-- =====================

-- Open-ended
INSERT INTO
    quiz_problems (
        group_id,
        problem_type,
        name,
        description,
        problem_text,
        points
    )
VALUES (
        (
            SELECT id
            FROM quiz_problem_groups
            WHERE
                name = 'Week 1 Quiz'
        ),
        'open_ended',
        'Explain Variables',
        'Open-ended question about variables.',
        'In your own words, explain what a variable is and why it is useful in programming.',
        5
    );

-- True/False
INSERT INTO
    quiz_problems (
        group_id,
        problem_type,
        name,
        description,
        problem_text,
        points
    )
VALUES (
        (
            SELECT id
            FROM quiz_problem_groups
            WHERE
                name = 'Week 1 Quiz'
        ),
        'true_false',
        'Python is Interpreted',
        'True/False about Python.',
        'Python is an interpreted language.',
        2
    );

INSERT INTO
    quiz_problem_options (
        problem_id,
        option_text,
        is_correct,
        display_order
    )
VALUES (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Python is Interpreted'
        ),
        'True',
        true,
        1
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Python is Interpreted'
        ),
        'False',
        false,
        2
    );

-- MCQ Single
INSERT INTO
    quiz_problems (
        group_id,
        problem_type,
        name,
        description,
        problem_text,
        points
    )
VALUES (
        (
            SELECT id
            FROM quiz_problem_groups
            WHERE
                name = 'Week 1 Quiz'
        ),
        'mcq_single',
        'Data Type Question',
        'MCQ about data types.',
        'Which of the following is NOT a primitive data type in Python?',
        3
    );

INSERT INTO
    quiz_problem_options (
        problem_id,
        option_text,
        is_correct,
        display_order
    )
VALUES (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Data Type Question'
        ),
        'int',
        false,
        1
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Data Type Question'
        ),
        'float',
        false,
        2
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Data Type Question'
        ),
        'array',
        true,
        3
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Data Type Question'
        ),
        'bool',
        false,
        4
    );

-- MCQ Multi
INSERT INTO
    quiz_problems (
        group_id,
        problem_type,
        name,
        description,
        problem_text,
        points
    )
VALUES (
        (
            SELECT id
            FROM quiz_problem_groups
            WHERE
                name = 'Week 1 Quiz'
        ),
        'mcq_multi',
        'Valid Identifiers',
        'MCQ with multiple correct answers.',
        'Which of the following are valid Python variable names? (Select all that apply)',
        4
    );

INSERT INTO
    quiz_problem_options (
        problem_id,
        option_text,
        is_correct,
        display_order
    )
VALUES (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Valid Identifiers'
        ),
        'my_variable',
        true,
        1
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Valid Identifiers'
        ),
        '2fast',
        false,
        2
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Valid Identifiers'
        ),
        '_private',
        true,
        3
    ),
    (
        (
            SELECT id
            FROM quiz_problems
            WHERE
                name = 'Valid Identifiers'
        ),
        'class',
        false,
        4
    );

-- Fill in the blank
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
VALUES (
        (
            SELECT id
            FROM quiz_problem_groups
            WHERE
                name = 'Week 1 Quiz'
        ),
        'fill_blank',
        'Loop Keyword',
        'Fill in the blank about loops.',
        'The keyword used to create a loop that iterates over a sequence in Python is _____.',
        2,
        'for,For,FOR'
    );

-- =====================
-- Submissions Seed Data
-- =====================

-- Student submission for A + B (Correct)
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
        (
            SELECT id
            FROM users
            WHERE
                email = 'student@enki.com'
        ),
        (
            SELECT id
            FROM comp_sci_problems
            WHERE
                name = 'A + B'
        ),
        'a, b = map(int, input().split())\nprint(a + b)',
        30,
        30,
        3,
        3,
        '[{"passed": true}, {"passed": true}, {"passed": true}]'
    );

-- =====================
-- Exam Sessions Seed Data
-- =====================

-- Create an exam session for the Midterm Exam
INSERT INTO
    exam_sessions (
        problem_group_type,
        problem_group_id,
        opened_by,
        duration_minutes,
        status
    )
VALUES (
        'comp_sci',
        (
            SELECT id
            FROM comp_sci_problem_group
            WHERE
                name = 'Midterm Exam'
        ),
        (
            SELECT id
            FROM users
            WHERE
                email = 'teacher@enki.com'
        ),
        60,
        'active'
    );

-- Assign the student to the exam session
INSERT INTO
    exam_session_students (session_id, user_id, status)
VALUES (
        (
            SELECT id
            FROM exam_sessions
            LIMIT 1
        ),
        (
            SELECT id
            FROM users
            WHERE
                email = 'student@enki.com'
        ),
        'assigned'
    );