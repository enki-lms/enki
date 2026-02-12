<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import {
    api,
    type QuizProblem,
    type QuizSubmission,
    type GradeInput,
  } from "$lib/api";

  export let problem: QuizProblem;

  const dispatch = createEventDispatcher<{ back: void }>();

  let submissions: QuizSubmission[] = [];
  let isLoading = true;
  let error: string | null = null;

  let gradingSubmission: QuizSubmission | null = null;
  let gradeScore = 0;
  let gradeFeedback = "";
  let gradeIsCorrect: boolean | null = null;
  let isGrading = false;

  async function fetchSubmissions() {
    try {
      isLoading = true;
      error = null;
      submissions = await api.getQuizProblemSubmissions(problem.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load submissions";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchSubmissions();
  });

  function handleBack() {
    dispatch("back");
  }

  function startGrading(submission: QuizSubmission) {
    gradingSubmission = submission;
    gradeScore = submission.score ?? 0;
    gradeFeedback = submission.feedback ?? "";
    gradeIsCorrect = submission.is_correct;
  }

  function cancelGrading() {
    gradingSubmission = null;
    gradeScore = 0;
    gradeFeedback = "";
    gradeIsCorrect = null;
  }

  async function submitGrade() {
    if (!gradingSubmission) return;
    isGrading = true;

    try {
      const input: GradeInput = {
        score: gradeScore,
        feedback: gradeFeedback || undefined,
        is_correct: gradeIsCorrect ?? undefined,
      };
      const updated = await api.gradeQuizSubmission(
        gradingSubmission.id,
        input,
      );
      submissions = submissions.map((s) => (s.id === updated.id ? updated : s));
      cancelGrading();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save grade";
    } finally {
      isGrading = false;
    }
  }

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }
</script>

<div class="w-full">
  <div class="flex items-center gap-4 mb-6">
    <BackButton onclick={handleBack} />
    <div class="flex-1">
      <h2 class="text-2xl font-semibold text-gray-900">Submissions</h2>
      <p class="text-gray-500 text-sm">{problem.name}</p>
    </div>
  </div>

  {#if error}
    <div
      class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
    >
      {error}
      <button class="ml-2 underline" on:click={fetchSubmissions}>Retry</button>
    </div>
  {/if}

  {#if isLoading}
    <div class="py-12 text-center text-gray-500">
      <div
        class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
      ></div>
      <p class="mt-4">Loading submissions...</p>
    </div>
  {:else if submissions.length === 0}
    <div class="py-12 text-center text-gray-500">
      <p class="text-lg font-medium">No submissions yet</p>
      <p class="text-sm mt-1">
        Students haven't submitted answers for this question
      </p>
    </div>
  {:else}
    <div class="space-y-4">
      {#each submissions as submission (submission.id)}
        <div class="bg-white border border-gray-200 rounded-xl p-5">
          <div class="flex items-start justify-between mb-3">
            <div>
              <p class="font-medium text-gray-900">
                {submission.user_name || `User ${submission.user_id}`}
              </p>
              <p class="text-sm text-gray-500">
                {formatDate(submission.created_at)}
              </p>
            </div>
            <div class="flex items-center gap-2">
              {#if submission.is_correct === true}
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  ✓ Correct
                </span>
              {:else if submission.is_correct === false}
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"
                >
                  ✗ Incorrect
                </span>
              {:else}
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
                >
                  Pending
                </span>
              {/if}
              <span class="text-sm font-medium text-gray-600">
                {submission.score}/{submission.max_score} pts
              </span>
            </div>
          </div>

          <!-- Answer -->
          <div class="mb-3">
            <p class="text-xs font-medium text-gray-500 mb-1">Answer</p>
            {#if submission.answer_text}
              <p class="bg-gray-50 rounded-lg p-3 text-gray-800">
                {submission.answer_text}
              </p>
            {:else if submission.selected_options?.length > 0}
              <div class="flex flex-wrap gap-2">
                {#each submission.selected_options as optionId}
                  {@const option = problem.options?.find(
                    (o) => o.id === optionId,
                  )}
                  <span
                    class="inline-flex items-center px-3 py-1 rounded-lg bg-gray-100 text-gray-800 text-sm"
                  >
                    {option?.option_text || `Option ${optionId}`}
                  </span>
                {/each}
              </div>
            {:else}
              <p class="text-gray-400 italic">No answer provided</p>
            {/if}
          </div>

          {#if submission.feedback}
            <div class="mb-3">
              <p class="text-xs font-medium text-gray-500 mb-1">Feedback</p>
              <p class="text-gray-700">{submission.feedback}</p>
            </div>
          {/if}

          <!-- Grade Button -->
          {#if gradingSubmission?.id === submission.id}
            <div class="border-t pt-4 mt-4 space-y-3">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label
                    for="grade-score"
                    class="block text-sm font-medium text-gray-700 mb-1"
                    >Score</label
                  >
                  <input
                    id="grade-score"
                    type="number"
                    bind:value={gradeScore}
                    min="0"
                    max={problem.points}
                    class="w-full px-3 py-2 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1"
                    >Correct?</label
                  >
                  <div class="flex gap-2">
                    <button
                      type="button"
                      class="flex-1 py-2 rounded-lg border-2 transition-colors
												{gradeIsCorrect === true
                        ? 'border-green-500 bg-green-50 text-green-700'
                        : 'border-gray-200 hover:border-gray-300'}"
                      on:click={() => (gradeIsCorrect = true)}
                    >
                      ✓ Yes
                    </button>
                    <button
                      type="button"
                      class="flex-1 py-2 rounded-lg border-2 transition-colors
												{gradeIsCorrect === false
                        ? 'border-red-500 bg-red-50 text-red-700'
                        : 'border-gray-200 hover:border-gray-300'}"
                      on:click={() => (gradeIsCorrect = false)}
                    >
                      ✗ No
                    </button>
                  </div>
                </div>
              </div>
              <div>
                <label
                  for="grade-feedback"
                  class="block text-sm font-medium text-gray-700 mb-1"
                  >Feedback</label
                >
                <textarea
                  id="grade-feedback"
                  bind:value={gradeFeedback}
                  placeholder="Optional feedback for the student..."
                  rows="2"
                  class="w-full px-3 py-2 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 resize-none"
                ></textarea>
              </div>
              <div class="flex justify-end gap-2">
                <button
                  type="button"
                  class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg"
                  on:click={cancelGrading}
                  disabled={isGrading}
                >
                  Cancel
                </button>
                <Button size="sm" on:click={submitGrade} disabled={isGrading}>
                  {#if isGrading}
                    <span
                      class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
                    ></span>
                  {/if}
                  Save Grade
                </Button>
              </div>
            </div>
          {:else}
            <button
              class="text-sm text-sky-500 hover:text-sky-600 font-medium"
              on:click={() => startGrading(submission)}
            >
              Grade Submission
            </button>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
