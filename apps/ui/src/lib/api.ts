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

export interface Course {
	id: number;
	name: string;
	institution: string;
}

export interface ProblemGroup {
	id: number;
	course_id?: number; // Backend might not return this directly if nested, but useful to have
	name: string;
	description: string;
	type: 'exam' | 'practice';
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
	getCourses: async (): Promise<Course[]> => {
		const response = await fetchWithAuth('/api/courses');
		const data = await handleResponse<Course[] | null>(response);
		return data || [];
	},

	getCourseProblemGroups: async (courseId: number): Promise<ProblemGroup[]> => {
		const response = await fetchWithAuth(`/api/courses/${courseId}/problem-groups`);
		return handleResponse<ProblemGroup[]>(response);
	},

	getGroupProblems: async (groupId: string): Promise<Problem[]> => {
		const response = await fetchWithAuth(`/api/problem-groups/${groupId}/problems`);
		return handleResponse<Problem[]>(response);
	},

	getProblem: async (id: string): Promise<Problem> => {
		const response = await fetchWithAuth(`/api/problems/${id}`);
		return handleResponse<Problem>(response);
	},
	
	submitSolution: async (problemId: string, code: string): Promise<SubmissionResponse> => {
		const response = await fetchWithAuth(`/api/problems/${problemId}/submit`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ code }),
		});
		return handleResponse<SubmissionResponse>(response);
	},

	chat: async (messages: { role: string; content: string }[]): Promise<{ response: string }> => {
		const response = await fetchWithAuth('/api/ai/chat', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ messages }),
		});
		return handleResponse<{ response: string }>(response);
	}
};
