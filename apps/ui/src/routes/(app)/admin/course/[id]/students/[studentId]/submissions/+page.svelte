<script lang="ts">
  import { page } from "$app/stores";
  import BackButton from "$lib/components/BackButton.svelte";
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import Card from "$lib/components/Card.svelte";
  import { onMount } from "svelte";
  import {
    api,
    type Course,
    type SubmissionHistoryItem,
    type Student,
  } from "$lib/api";
  import { goto } from "$app/navigation";

  export let data;

  const courseId = parseInt($page.params.id);
  const studentId = parseInt($page.params.studentId);

  let course: Course | null = null;
  let student: Student | null = null;
  let submissions: SubmissionHistoryItem[] = [];
  let loading = true;
  let error: string | null = null;
  let showModal = false;
  let selectedSubmission: SubmissionHistoryItem | null = null;

  async function loadData() {
    try {
      loading = true;
      const [courseData, enrollments, submissionsData, studentsData] =
        await Promise.all([
          api.getCourse(courseId),
          api.getCourseEnrollments(courseId),
          api.getStudentSubmissions(courseId, studentId),
          api.getSchoolStudents(),
        ]);
      course = courseData;
      submissions = submissionsData;

      const enrollment = enrollments.find((e) => e.user_id === studentId);
      const studentUser = studentsData.find((s) => s.id === studentId);

      if (enrollment && studentUser) {
        student = studentUser;
      } else {
        throw new Error("Student not found in this course");
      }
    } catch (e) {
      console.error("Error loading data:", e);
      error = e instanceof Error ? e.message : "Failed to load data";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  function openModal(submission: SubmissionHistoryItem) {
    selectedSubmission = submission;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    selectedSubmission = null;
  }

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleString();
  }

  function getScoreColor(score: number, max: number) {
    const percentage = max > 0 ? (score / max) * 100 : 0;
    if (percentage >= 80) return "text-green-600 bg-green-50";
    if (percentage >= 50) return "text-yellow-600 bg-yellow-50";
    return "text-red-600 bg-red-50";
  }
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex items-center gap-4">
    <div class="p-2">
      <BackButton onclick={() => goto(`/course/${courseId}`)} />
    </div>
    <div class="p-2">
      <ProfileMenu
        width="60"
        height="60"
        name={data.user?.fullName ?? ""}
        email={data.user?.email ?? ""}
        role={data.user?.role ?? ""}
      />
    </div>
    <div class="p-3"><LogoPlaceHolder /></div>
  </header>

  <main class="p-4 md:p-8">
    <div class="max-w-5xl mx-auto">
      {#if loading}
        <div class="bg-white rounded-2xl p-12 shadow-sm flex justify-center">
          <div class="flex flex-col items-center gap-3">
            <div
              class="w-8 h-8 border-3 border-sky-400 border-t-transparent rounded-full animate-spin"
            ></div>
            <p class="text-gray-500">Loading submissions...</p>
          </div>
        </div>
      {:else if error}
        <div class="bg-white rounded-2xl p-12 shadow-sm flex justify-center">
          <div class="flex flex-col items-center gap-3 text-center">
            <div
              class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center"
            >
              <svg
                class="w-6 h-6 text-red-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
            </div>
            <p class="text-red-500 font-medium">{error}</p>
            <button
              class="text-sky-600 hover:underline text-sm"
              on:click={loadData}>Retry</button
            >
          </div>
        </div>
      {:else if course && student}
        <div
          class="bg-white rounded-2xl p-6 md:p-8 shadow-sm border border-gray-100 mb-6"
        >
          <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-4"
          >
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-gray-900">
                {student.full_name}
              </h1>
              <p class="text-gray-500 mt-1">
                Submissions for <span class="font-medium text-gray-700"
                  >{course.name}</span
                >
              </p>
            </div>
            <div class="flex flex-col items-end gap-1">
              <span class="text-sm text-gray-400 font-mono"
                >{student.email}</span
              >
            </div>
          </div>
        </div>

        <div
          class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden"
        >
          {#if submissions.length === 0}
            <div class="p-12 text-center text-gray-500">
              <p>No submissions found for this student.</p>
            </div>
          {:else}
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr
                    class="bg-gray-50 border-b border-gray-100 text-xs text-gray-500 uppercase tracking-wider"
                  >
                    <th class="p-4 font-medium">Type</th>
                    <th class="p-4 font-medium">Problem</th>
                    <th class="p-4 font-medium">Date</th>
                    <th class="p-4 font-medium">Score</th>
                    <th class="p-4 font-medium text-right">Details</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100">
                  {#each submissions as sub (sub.id)}
                    <tr class="hover:bg-gray-50 transition-colors">
                      <td class="p-4">
                        <span
                          class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${sub.type === "comp_sci" ? "bg-indigo-100 text-indigo-800" : "bg-emerald-100 text-emerald-800"}`}
                        >
                          {sub.type === "comp_sci" ? "Coding" : "Quiz"}
                        </span>
                      </td>
                      <td class="p-4 text-sm text-gray-900">
                        Problem #{sub.problem_id}
                        <!-- TODO: Would be nice to have problem name, but API doesn't return it yet (or we'd need to fetch) -->
                      </td>
                      <td class="p-4 text-sm text-gray-500">
                        {formatDate(sub.created_at)}
                      </td>
                      <td class="p-4">
                        <span
                          class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getScoreColor(sub.score, sub.max_score)}`}
                        >
                          {sub.score} / {sub.max_score}
                        </span>
                      </td>
                      <td class="p-4 text-right">
                        <button
                          class="text-sky-600 hover:text-sky-800 text-sm font-medium hover:underline"
                          on:click={() => openModal(sub)}
                        >
                          View Details
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </main>

  {#if showModal && selectedSubmission}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      on:click|self={closeModal}
    >
      <div
        class="bg-white rounded-2xl w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-2xl flex flex-col"
      >
        <div
          class="p-6 border-b border-gray-100 flex items-center justify-between sticky top-0 bg-white z-10"
        >
          <h3 class="text-xl font-bold text-gray-900">Submission Details</h3>
          <button
            on:click={closeModal}
            class="text-gray-400 hover:text-gray-600"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        <div class="p-6 flex-1 overflow-y-auto">
          <div class="grid grid-cols-2 gap-4 mb-6 text-sm">
            <div class="p-3 bg-gray-50 rounded-lg">
              <span class="block text-gray-500 text-xs uppercase mb-1"
                >Date</span
              >
              <span class="font-medium"
                >{formatDate(selectedSubmission.created_at)}</span
              >
            </div>
            <div class="p-3 bg-gray-50 rounded-lg">
              <span class="block text-gray-500 text-xs uppercase mb-1"
                >Score</span
              >
              <span
                class={`font-medium ${getScoreColor(selectedSubmission.score, selectedSubmission.max_score).split(" ")[0]}`}
              >
                {selectedSubmission.score} / {selectedSubmission.max_score}
              </span>
            </div>
          </div>

          {#if selectedSubmission.type === "comp_sci"}
            <div class="mb-6">
              <h4 class="font-semibold text-gray-900 mb-2">Code</h4>
              <div
                class="bg-gray-900 rounded-lg p-4 overflow-x-auto text-sm font-mono text-gray-100"
              >
                <pre>{selectedSubmission.code}</pre>
              </div>
            </div>
            <div>
              <h4 class="font-semibold text-gray-900 mb-2">Test Results</h4>
              <div class="space-y-2">
                <div class="flex items-center gap-2">
                  <span class="text-sm text-gray-600">Passed:</span>
                  <span class="font-mono font-medium"
                    >{selectedSubmission.passed_tests} / {selectedSubmission.total_tests}</span
                  >
                </div>
                {#if selectedSubmission.results_json}
                  <div
                    class="mt-2 text-xs text-gray-500 bg-gray-50 p-3 rounded border border-gray-100 font-mono whitespace-pre-wrap"
                  >
                    {selectedSubmission.results_json}
                    <!-- Render raw JSON for now, can be improved -->
                  </div>
                {/if}
              </div>
            </div>
          {:else if selectedSubmission.type === "quiz"}
            {#if selectedSubmission.answer_text}
              <div class="mb-6">
                <h4 class="font-semibold text-gray-900 mb-2">Answer Text</h4>
                <div
                  class="bg-gray-50 p-4 rounded-lg text-gray-800 border border-gray-200"
                >
                  {selectedSubmission.answer_text}
                </div>
              </div>
            {/if}
            {#if selectedSubmission.selected_options && selectedSubmission.selected_options.length > 0}
              <div class="mb-6">
                <h4 class="font-semibold text-gray-900 mb-2">
                  Selected Options
                </h4>
                <div
                  class="bg-gray-50 p-4 rounded-lg text-gray-800 border border-gray-200 font-mono"
                >
                  IDs: {selectedSubmission.selected_options.join(", ")}
                </div>
              </div>
            {/if}
            {#if selectedSubmission.feedback}
              <div>
                <h4 class="font-semibold text-gray-900 mb-2">Feedback</h4>
                <div
                  class="bg-blue-50 p-4 rounded-lg text-blue-800 border border-blue-200"
                >
                  {selectedSubmission.feedback}
                </div>
              </div>
            {/if}
          {/if}
        </div>
        <div
          class="p-4 border-t border-gray-100 bg-gray-50 rounded-b-2xl flex justify-end"
        >
          <button
            class="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg shadow-sm hover:bg-gray-50 font-medium transition-colors"
            on:click={closeModal}
          >
            Close
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
