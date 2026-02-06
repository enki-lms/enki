export interface Problem {
	id: number;
	name: string;
	description: string;
	problem_text: string;
	group_id: number;
	time_limit_ms: number | null;
	memory_limit_mb: number | null;
}

export interface TestCaseResult {
	test_case_id: number;
	passed: boolean;
	points: number;
	expected?: string;
	actual?: string;
	error?: string;
	timed_out?: boolean;
}

export interface SubmissionResult {
	problem_id: number;
	total_test_cases: number;
	passed: number;
	failed: number;
	score: number;
	max_score: number;
	results: TestCaseResult[];
}

export interface SubmissionResponse {
	submission_id: number;
	result: SubmissionResult;
}

export interface CodeSubmission {
	id: number;
	created_at: string;
	user_id: number;
	problem_id: number;
	code: string;
	score: number;
	max_score: number;
	passed_tests: number;
	total_tests: number;
	results_json: string;
}

export interface QuizSubmission {
	id: number;
	created_at: string;
	user_id: number;
	problem_id: number;
	answer_text: string | null;
	selected_options: number[];
	is_correct: boolean | null;
	score: number;
	max_score: number;
	feedback: string | null;
	user_name?: string;
	user_email?: string;
}

export interface Course {
	id: number;
	name: string;
	subject: string;
	institution: string;
	owner_id: number;
	created_at: string;
	updated_at: string;
}

export interface ProblemGroup {
	id: number;
	course_id?: number;
	name: string;
	description: string;
	type: "exam" | "practice";
}

// Teacher Admin Types
export interface Student {
	id: number;
	full_name: string;
	email: string;
}

export interface Enrollment {
	id: number;
	user_id: number;
	course_id: number;
	user?: Student;
}

export interface TestCase {
	id: number;
	problem_id: number;
	input: string;
	output: string;
	correct_points: number;
}

export interface QuizGroup {
	id: number;
	course_id?: number;
	name: string;
	description: string;
	type: "exam" | "practice";
}

export interface QuizOption {
	id?: number;
	option_text: string;
	is_correct: boolean;
	display_order: number;
}

export type QuizProblemType = "open_ended" | "true_false" | "mcq_single" | "mcq_multi" | "fill_blank";

export interface QuizProblem {
	id: number;
	group_id: number;
	problem_type: QuizProblemType;
	name: string;
	description: string;
	problem_text: string;
	points: number;
	correct_answer?: string;
	options?: QuizOption[];
}

export interface ExamSession {
	id: number;
	problem_group_type: "comp_sci" | "quiz";
	problem_group_id: number;
	duration_minutes: number;
	status: "pending" | "active" | "ended";
	started_at?: string;
	ended_at?: string;
	created_at: string;
}

export interface ExamStudent {
	id: number;
	session_id: number;
	student_id: number;
	status: "active" | "submitted" | "discontinued";
	started_at?: string;
	submitted_at?: string;
	student?: Student;
}

// Input types for create/update operations
export interface CourseInput {
	name: string;
	subject: string;
}

export interface ProblemGroupInput {
	type: "exam" | "practice";
	name: string;
	description: string;
}

export interface ProblemInput {
	name: string;
	description: string;
	problem_text: string;
	time_limit_ms?: number;
	memory_limit_mb?: number;
}

export interface TestCaseInput {
	input: string;
	output: string;
	correct_points: number;
}

export interface QuizGroupInput {
	type: "exam" | "practice";
	name: string;
	description: string;
}

export interface QuizProblemInput {
	problem_type: QuizProblemType;
	name: string;
	description: string;
	problem_text: string;
	points: number;
	correct_answer?: string;
	options?: Omit<QuizOption, "id">[];
}

export interface ExamSessionInput {
	problem_group_type: "comp_sci" | "quiz";
	problem_group_id: number;
	duration_minutes: number;
	student_ids: number[];
}

export interface GradeInput {
	is_correct?: boolean;
	score: number;
	feedback?: string;
}

// Helper to handle API errors
async function handleResponse<T>(response: Response): Promise<T> {
	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: response.statusText }));
		throw new Error(error.error || response.statusText);
	}
	return response.json();
}

const getCookie = (name: string): string | null => {
	const value = `; ${document.cookie}`;
	const parts = value.split(`; ${name}=`);
	if (parts.length === 2) return parts.pop()?.split(';').shift() || null;
	return null;
};

// Helper to handle API calls with auth
const fetchWithAuth = async (url: string, options: RequestInit = {}): Promise<Response> => {
	// Cookies are automatically sent by the browser for same-origin (proxied) requests
	// or if credentials: 'include' is set. 
	// Since we are proxying, we might rely on the cookie being there.
	// However, if we need to manually read a token from a cookie and send it as Bearer, we can.
	// The backend says "Bearer <token> (or cookie)", so cookie should suffice if standard.
	// Let's ensure we send credentials: 'include' just in case.
	
	const token = getCookie('token');
	const headers = {
		...options.headers,
		...(token ? { 'Authorization': `Bearer ${token}` } : {}),
	};

	return fetch(url, { ...options, headers });
};

export const api = {
	// ============ Courses ============
	getCourses: async (): Promise<Course[]> => {
		const response = await fetchWithAuth("/api/courses/teaching"); // Changed to specific endpoint for teachers
		const data = await handleResponse<Course[] | null>(response);
		return data || [];
	},

	// Student endpoint - returns courses the student is enrolled in
	getEnrolledCourses: async (): Promise<Course[]> => {
		const response = await fetchWithAuth("/api/courses");
		const data = await handleResponse<Course[] | null>(response);
		return data || [];
	},

	getCourse: async (id: number): Promise<Course> => {
		const response = await fetchWithAuth(`/api/courses/${id}`);
		return handleResponse<Course>(response);
	},

	createCourse: async (input: CourseInput): Promise<Course> => {
		const response = await fetchWithAuth("/api/courses", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<Course>(response);
	},

	updateCourse: async (id: number, input: CourseInput): Promise<Course> => {
		const response = await fetchWithAuth(`/api/courses/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<Course>(response);
	},

	deleteCourse: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/courses/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Enrollments ============
	getCourseEnrollments: async (courseId: number): Promise<Enrollment[]> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/enrollments`);
		const data = await handleResponse<Enrollment[] | null>(response);
		return data || [];
	},

	enrollStudent: async (courseId: number, userId: number): Promise<Enrollment> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/enrollments`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ user_id: userId }),
		});
		return handleResponse<Enrollment>(response);
	},

	unenrollStudent: async (courseId: number, userId: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/enrollments/${userId}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Students ============
	getSchoolStudents: async (): Promise<Student[]> => {
		const response = await fetchWithAuth("/api/users/students");
		const data = await handleResponse<Student[] | null>(response);
		return data || [];
	},

	// ============ Problem Groups (CS) ============
	getCourseProblemGroups: async (courseId: number): Promise<ProblemGroup[]> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/problem-groups`);
		const data = await handleResponse<ProblemGroup[] | null>(response);
		return data || [];
	},

	getProblemGroup: async (id: number): Promise<ProblemGroup> => {
		const response = await fetchWithAuth(`/api/problem-groups/${id}`);
		return handleResponse<ProblemGroup>(response);
	},

	createProblemGroup: async (courseId: number, input: ProblemGroupInput): Promise<ProblemGroup> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/problem-groups`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<ProblemGroup>(response);
	},

	updateProblemGroup: async (id: number, input: ProblemGroupInput): Promise<ProblemGroup> => {
		const response = await fetchWithAuth(`/api/problem-groups/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<ProblemGroup>(response);
	},

	deleteProblemGroup: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/problem-groups/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Problems (CS) ============
	getGroupProblems: async (groupId: number): Promise<Problem[]> => {
		const response = await fetchWithAuth(`/api/problem-groups/${groupId}/problems`);
		const data = await handleResponse<Problem[] | null>(response);
		return data || [];
	},

	getProblem: async (id: number): Promise<Problem> => {
		const response = await fetchWithAuth(`/api/problems/${id}`);
		return handleResponse<Problem>(response);
	},

	createProblem: async (groupId: number, input: ProblemInput): Promise<Problem> => {
		const response = await fetchWithAuth(`/api/problem-groups/${groupId}/problems`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<Problem>(response);
	},

	updateProblem: async (id: number, input: ProblemInput): Promise<Problem> => {
		const response = await fetchWithAuth(`/api/problems/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<Problem>(response);
	},

	deleteProblem: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/problems/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Test Cases (CS) ============
	getProblemTestCases: async (problemId: number): Promise<TestCase[]> => {
		const response = await fetchWithAuth(`/api/problems/${problemId}/test-cases`);
		const data = await handleResponse<TestCase[] | null>(response);
		return data || [];
	},

	getTestCase: async (id: number): Promise<TestCase> => {
		const response = await fetchWithAuth(`/api/test-cases/${id}`);
		return handleResponse<TestCase>(response);
	},

	createTestCase: async (problemId: number, input: TestCaseInput): Promise<TestCase> => {
		const response = await fetchWithAuth(`/api/problems/${problemId}/test-cases`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<TestCase>(response);
	},

	updateTestCase: async (id: number, input: TestCaseInput): Promise<TestCase> => {
		const response = await fetchWithAuth(`/api/test-cases/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<TestCase>(response);
	},

	deleteTestCase: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/test-cases/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Quiz Groups ============
	getCourseQuizGroups: async (courseId: number): Promise<QuizGroup[]> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/quiz-groups`);
		const data = await handleResponse<QuizGroup[] | null>(response);
		return data || [];
	},

	getQuizGroup: async (id: number): Promise<QuizGroup> => {
		const response = await fetchWithAuth(`/api/quiz-groups/${id}`);
		return handleResponse<QuizGroup>(response);
	},

	createQuizGroup: async (courseId: number, input: QuizGroupInput): Promise<QuizGroup> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/quiz-groups`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<QuizGroup>(response);
	},

	updateQuizGroup: async (id: number, input: QuizGroupInput): Promise<QuizGroup> => {
		const response = await fetchWithAuth(`/api/quiz-groups/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<QuizGroup>(response);
	},

	deleteQuizGroup: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/quiz-groups/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Quiz Problems ============
	getQuizGroupProblems: async (groupId: number): Promise<QuizProblem[]> => {
		const response = await fetchWithAuth(`/api/quiz-groups/${groupId}/problems`);
		const data = await handleResponse<QuizProblem[] | null>(response);
		return data || [];
	},

	getQuizProblem: async (id: number): Promise<QuizProblem> => {
		const response = await fetchWithAuth(`/api/quiz-problems/${id}`);
		return handleResponse<QuizProblem>(response);
	},

	createQuizProblem: async (groupId: number, input: QuizProblemInput): Promise<QuizProblem> => {
		const response = await fetchWithAuth(`/api/quiz-groups/${groupId}/problems`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<QuizProblem>(response);
	},

	updateQuizProblem: async (id: number, input: QuizProblemInput): Promise<QuizProblem> => {
		const response = await fetchWithAuth(`/api/quiz-problems/${id}`, {
			method: "PUT",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<QuizProblem>(response);
	},

	deleteQuizProblem: async (id: number): Promise<void> => {
		const response = await fetchWithAuth(`/api/quiz-problems/${id}`, {
			method: "DELETE",
		});
		await handleResponse<void>(response);
	},

	// ============ Exam Sessions ============
	getExamSessions: async (): Promise<ExamSession[]> => {
		const response = await fetchWithAuth("/api/exams/sessions");
		const data = await handleResponse<ExamSession[] | null>(response);
		return data || [];
	},

	getExamSession: async (id: number): Promise<ExamSession> => {
		const response = await fetchWithAuth(`/api/exams/sessions/${id}`);
		return handleResponse<ExamSession>(response);
	},

	createExamSession: async (input: ExamSessionInput): Promise<ExamSession> => {
		const response = await fetchWithAuth("/api/exams/sessions", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<ExamSession>(response);
	},

	startExamSession: async (id: number): Promise<ExamSession> => {
		const response = await fetchWithAuth(`/api/exams/sessions/${id}/start`, {
			method: "POST",
		});
		return handleResponse<ExamSession>(response);
	},

	endExamSession: async (id: number): Promise<ExamSession> => {
		const response = await fetchWithAuth(`/api/exams/sessions/${id}/end`, {
			method: "POST",
		});
		return handleResponse<ExamSession>(response);
	},

	getExamSessionStudents: async (sessionId: number): Promise<ExamStudent[]> => {
		const response = await fetchWithAuth(`/api/exams/sessions/${sessionId}/students`);
		const data = await handleResponse<ExamStudent[] | null>(response);
		return data || [];
	},

	discontinueStudent: async (sessionId: number, studentId: number): Promise<void> => {
		const response = await fetchWithAuth(
			`/api/exams/sessions/${sessionId}/students/${studentId}/discontinue`,
			{ method: "POST" }
		);
		await handleResponse<void>(response);
	},

	// ============ Grading ============
	gradeQuizSubmission: async (submissionId: number, input: GradeInput): Promise<QuizSubmission> => {
		const response = await fetchWithAuth(`/api/quiz-submissions/${submissionId}/grade`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(input),
		});
		return handleResponse<QuizSubmission>(response);
	},

	getQuizProblemSubmissions: async (problemId: number): Promise<QuizSubmission[]> => {
		const response = await fetchWithAuth(`/api/quiz-problems/${problemId}/all-submissions`);
		const data = await handleResponse<QuizSubmission[] | null>(response);
		return data || [];
	},

	// ============ Student-facing (existing) ============
	submitSolution: async (problemId: string, code: string): Promise<SubmissionResponse> => {
		const response = await fetchWithAuth(`/api/problems/${problemId}/submit`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ code }),
		});
		return handleResponse<SubmissionResponse>(response);
	},

	chat: async (messages: { role: string; content: string }[]): Promise<{ response: string }> => {
		const response = await fetchWithAuth("/api/ai/chat", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ messages }),
		});
		return handleResponse<{ response: string }>(response);
	},

	getAllCodeSubmissions: async (limit?: number): Promise<CodeSubmission[]> => {
		const url = limit ? `/api/submissions?limit=${limit}` : "/api/submissions";
		const response = await fetchWithAuth(url);
		const data = await handleResponse<CodeSubmission[] | null>(response);
		return data || [];
	},

	getAllQuizSubmissions: async (limit?: number): Promise<QuizSubmission[]> => {
		const url = limit ? `/api/quiz-submissions?limit=${limit}` : "/api/quiz-submissions";
		const response = await fetchWithAuth(url);
		const data = await handleResponse<QuizSubmission[] | null>(response);
		return data || [];
	},
};

