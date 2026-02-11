<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import QuizProblemModal from "$lib/components/QuizProblemModal.svelte";
  import QuizProblemSubmissions from "$lib/components/QuizProblemSubmissions.svelte";
  import {
    api,
    type QuizGroup,
    type QuizProblem,
    type QuizProblemInput,
  } from "$lib/api";

  export let group: QuizGroup;

  const dispatch = createEventDispatcher<{ back: void }>();

  let problems: QuizProblem[] = [];
  let isLoading = true;
  let error: string | null = null;

  let isFormModalOpen = false;
  let isDeleteDialogOpen = false;
  let isSubmitting = false;
  let editingProblem: QuizProblem | null = null;
  let deletingProblem: QuizProblem | null = null;

  let viewingSubmissionsProblem: QuizProblem | null = null;

  const problemTypeLabels: Record<string, string> = {
    open_ended: "📝 Open Ended",
    true_false: "✅ True/False",
    mcq_single: "🔘 Single Choice",
    mcq_multi: "☑️ Multiple Choice",
    fill_blank: "📄 Fill in Blank",
  };

  async function fetchProblems() {
    try {
      isLoading = true;
      error = null;
      problems = await api.getQuizGroupProblems(group.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load problems";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchProblems();
  });

  function handleBack() {
    dispatch("back");
  }

  function handleNewProblem() {
    editingProblem = null;
    isFormModalOpen = true;
  }

  function handleEditProblem(problem: QuizProblem) {
    editingProblem = problem;
    isFormModalOpen = true;
  }

  function handleDeleteClick(problem: QuizProblem) {
    deletingProblem = problem;
    isDeleteDialogOpen = true;
  }

  function handleViewSubmissions(problem: QuizProblem) {
    viewingSubmissionsProblem = problem;
  }

  async function handleFormSubmit(event: CustomEvent<QuizProblemInput>) {
    const input = event.detail;
    isSubmitting = true;

    try {
      if (editingProblem) {
        const updated = await api.updateQuizProblem(editingProblem.id, input);
        problems = problems.map((p) => (p.id === updated.id ? updated : p));
      } else {
        const created = await api.createQuizProblem(group.id, input);
        problems = [...problems, created];
      }
      isFormModalOpen = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save problem";
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteConfirm() {
    if (!deletingProblem) return;
    isSubmitting = true;

    try {
      await api.deleteQuizProblem(deletingProblem.id);
      problems = problems.filter((p) => p.id !== deletingProblem!.id);
      isDeleteDialogOpen = false;
      deletingProblem = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete problem";
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if viewingSubmissionsProblem}
  <QuizProblemSubmissions
    problem={viewingSubmissionsProblem}
    on:back={() => (viewingSubmissionsProblem = null)}
  />
{:else}
  <div class="w-full">
    <div class="flex items-center gap-4 mb-6">
      <BackButton onclick={handleBack} />
      <div class="flex-1">
        <h2 class="text-2xl font-semibold text-gray-900">{group.name}</h2>
        <p class="text-gray-500 text-sm">
          {group.type === "exam" ? "📝 Exam" : "📚 Practice"}
        </p>
      </div>
      <Button size="md" on:click={handleNewProblem}>New Question</Button>
    </div>

    {#if group.description}
      <p class="text-gray-600 mb-6">{group.description}</p>
    {/if}

    {#if error}
      <div
        class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
      >
        {error}
        <button class="ml-2 underline" on:click={fetchProblems}>Retry</button>
      </div>
    {/if}

    {#if isLoading}
      <div class="py-12 text-center text-gray-500">
        <div
          class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-4">Loading questions...</p>
      </div>
    {:else if problems.length === 0}
      <div class="py-12 text-center text-gray-500">
        <p class="text-lg font-medium">No questions yet</p>
        <p class="text-sm mt-1">Add quiz questions to this group</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each problems as problem, idx (problem.id)}
          <div
            class="bg-white border border-gray-200 rounded-xl p-4 group relative hover:border-sky-300 transition-colors"
          >
            <div class="flex items-start gap-4">
              <div
                class="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center text-purple-600 font-bold flex-shrink-0"
              >
                {idx + 1}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-semibold text-gray-900">{problem.name}</h3>
                  <span
                    class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600"
                  >
                    {problemTypeLabels[problem.problem_type] ||
                      problem.problem_type}
                  </span>
                  <span class="text-sm text-green-600 font-medium"
                    >{problem.points} pts</span
                  >
                </div>
                {#if problem.description}
                  <p class="text-sm text-gray-500 line-clamp-1">
                    {problem.description}
                  </p>
                {/if}
                <p class="text-sm text-gray-700 mt-1 line-clamp-2">
                  {problem.problem_text}
                </p>
              </div>
            </div>
            <div
              class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1"
            >
              <button
                class="p-2 bg-gray-100 rounded-lg hover:bg-purple-100 text-purple-600"
                on:click={() => handleViewSubmissions(problem)}
                title="View Submissions"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
                  />
                </svg>
              </button>
              <button
                class="p-2 bg-gray-100 rounded-lg hover:bg-gray-200 text-gray-600"
                on:click={() => handleEditProblem(problem)}
                title="Edit"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
              </button>
              <button
                class="p-2 bg-gray-100 rounded-lg hover:bg-red-100 text-red-500"
                on:click={() => handleDeleteClick(problem)}
                title="Delete"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
{/if}

<QuizProblemModal
  bind:isOpen={isFormModalOpen}
  {editingProblem}
  isLoading={isSubmitting}
  on:submit={handleFormSubmit}
  on:cancel={() => (isFormModalOpen = false)}
/>

<ConfirmDialog
  bind:isOpen={isDeleteDialogOpen}
  title="Delete Question"
  message="Are you sure you want to delete '{deletingProblem?.name}'?"
  isLoading={isSubmitting}
  on:confirm={handleDeleteConfirm}
  on:cancel={() => {
    isDeleteDialogOpen = false;
    deletingProblem = null;
  }}
/>
