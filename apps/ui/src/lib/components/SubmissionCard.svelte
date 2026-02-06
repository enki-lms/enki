<script lang="ts">
  import type { CodeSubmission, QuizSubmission } from "$lib/api";

  type SubmissionType = "code" | "quiz";

  interface Props {
    submission: CodeSubmission | QuizSubmission;
    type: SubmissionType;
    problemName?: string;
  }

  let { submission, type, problemName = "Problem" }: Props = $props();

  let expanded = $state(false);

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  }

  function getScorePercent(): number {
    return submission.max_score > 0
      ? (submission.score / submission.max_score) * 100
      : 0;
  }

  function getScoreColor(): string {
    const percent = getScorePercent();
    if (percent >= 80) return "text-emerald-600 bg-emerald-50";
    if (percent >= 50) return "text-amber-600 bg-amber-50";
    return "text-red-600 bg-red-50";
  }

  function getProgressColor(): string {
    const percent = getScorePercent();
    if (percent >= 80) return "bg-emerald-500";
    if (percent >= 50) return "bg-amber-500";
    return "bg-red-500";
  }

  function isCodeSubmission(
    s: CodeSubmission | QuizSubmission,
  ): s is CodeSubmission {
    return "code" in s;
  }

  function isQuizSubmission(
    s: CodeSubmission | QuizSubmission,
  ): s is QuizSubmission {
    return "answer_text" in s || "selected_options" in s;
  }
</script>

<div
  class="bg-white rounded-xl border border-gray-200 overflow-hidden transition-all duration-200 hover:shadow-md hover:border-gray-300"
>
  <!-- Header - Always visible -->
  <button
    type="button"
    class="w-full p-4 text-left cursor-pointer"
    onclick={() => (expanded = !expanded)}
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <!-- Type Badge -->
          <span
            class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium {type ===
            'code'
              ? 'bg-sky-100 text-sky-700'
              : 'bg-purple-100 text-purple-700'}"
          >
            {#if type === "code"}
              <svg
                class="w-3 h-3 mr-1"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                />
              </svg>
              Code
            {:else}
              <svg
                class="w-3 h-3 mr-1"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Quiz
            {/if}
          </span>
          <span class="text-xs text-gray-500"
            >{formatDate(submission.created_at)}</span
          >
        </div>
        <h3 class="font-semibold text-gray-900 truncate">{problemName}</h3>
      </div>

      <!-- Score -->
      <div class="flex items-center gap-3">
        <div class="text-right">
          <span
            class="inline-flex items-center px-2.5 py-1 rounded-lg text-sm font-bold {getScoreColor()}"
          >
            {submission.score}/{submission.max_score}
          </span>
        </div>
        <!-- Expand Arrow -->
        <svg
          class="w-5 h-5 text-gray-400 transition-transform duration-200 {expanded
            ? 'rotate-180'
            : ''}"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="mt-3 h-1.5 bg-gray-100 rounded-full overflow-hidden">
      <div
        class="h-full rounded-full transition-all duration-500 {getProgressColor()}"
        style="width: {getScorePercent()}%"
      ></div>
    </div>
  </button>

  <!-- Expanded Details -->
  {#if expanded}
    <div class="px-4 pb-4 pt-0 border-t border-gray-100 animate-slideDown">
      {#if type === "code" && isCodeSubmission(submission)}
        <div class="mt-3 space-y-3">
          <div class="flex items-center gap-4 text-sm text-gray-600">
            <span class="flex items-center gap-1">
              <svg
                class="w-4 h-4 text-emerald-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M5 13l4 4L19 7"
                />
              </svg>
              {submission.passed_tests} passed
            </span>
            <span class="flex items-center gap-1">
              <svg
                class="w-4 h-4 text-red-500"
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
              {submission.total_tests - submission.passed_tests} failed
            </span>
          </div>
          <div class="bg-gray-900 rounded-lg p-3 overflow-x-auto">
            <pre
              class="text-sm text-gray-100 font-mono whitespace-pre-wrap">{submission.code}</pre>
          </div>
        </div>
      {:else if type === "quiz" && isQuizSubmission(submission)}
        <div class="mt-3 space-y-3">
          {#if submission.answer_text}
            <div>
              <p class="text-sm font-medium text-gray-700 mb-1">Your Answer:</p>
              <p class="text-sm text-gray-600 bg-gray-50 rounded-lg p-3">
                {submission.answer_text}
              </p>
            </div>
          {/if}
          {#if submission.is_correct !== null}
            <div class="flex items-center gap-2">
              {#if submission.is_correct}
                <span
                  class="inline-flex items-center gap-1 text-sm font-medium text-emerald-600"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  Correct
                </span>
              {:else}
                <span
                  class="inline-flex items-center gap-1 text-sm font-medium text-red-600"
                >
                  <svg
                    class="w-4 h-4"
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
                  Incorrect
                </span>
              {/if}
            </div>
          {/if}
          {#if submission.feedback}
            <div>
              <p class="text-sm font-medium text-gray-700 mb-1">Feedback:</p>
              <p class="text-sm text-gray-600 bg-blue-50 rounded-lg p-3">
                {submission.feedback}
              </p>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-slideDown {
    animation: slideDown 0.2s ease-out;
  }
</style>
