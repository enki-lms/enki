-- Insert a test course
INSERT INTO
    courses (name, institution)
VALUES (
        'CS101: Introduction to Computer Science',
        'Demo Institution'
    );

-- Insert a test problem group
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
            LIMIT 1
        ),
        'practice',
        'Week 1: Basics',
        'Input/Output and simple arithmetic.'
    );

-- Insert a test problem (A+B)
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
            LIMIT 1
        ),
        'A + B',
        'Read two integers from input and print their sum.',
        'Read two integers from input and print their sum. The input contains one line with two space-separated integers.',
        1000,
        128
    );

-- Insert tests for A+B
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
            LIMIT 1
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
            LIMIT 1
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
            LIMIT 1
        ),
        '100 200',
        '300',
        10
    );

-- Insert a challenging problem (Factorial)
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
            LIMIT 1
        ),
        'Factorial',
        'Read an integer N from input and print N!.',
        'Read an integer N from input and print N!. Input constraints: 0 <= N <= 10.',
        1000,
        128
    );

-- Insert tests for Factorial
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
            LIMIT 1
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
            LIMIT 1
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
            LIMIT 1
        ),
        '3',
        '6',
        10
    );