<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import FormModal, {
    type FieldConfig,
  } from "$lib/components/FormModal.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import ProblemDetail from "$lib/components/ProblemDetail.svelte";
  import { api, type ProblemGroup, type Problem } from "$lib/api";

  export let group: ProblemGroup;

  const dispatch = createEventDispatcher<{ back: void }>();

  let problems: Problem[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Modal states
  let isFormModalOpen = false;
  let isDeleteDialogOpen = false;
  let isSubmitting = false;
  let editingProblem: Problem | null = null;
  let deletingProblem: Problem | null = null;

  // Detail view
  let selectedProblem: Problem | null = null;

  const formFields: FieldConfig[] = [
    {
      name: "name",
      label: "Problem Name",
      type: "text",
      placeholder: "Enter problem name...",
      required: true,
    },
    {
      name: "description",
      label: "Description",
      type: "textarea",
      placeholder: "Brief description...",
      required: false,
    },
    {
      name: "problem_text",
      label: "Problem Text",
      type: "textarea",
      placeholder: "Full problem statement...",
      required: true,
    },
    {
      name: "time_limit_ms",
      label: "Time Limit (ms)",
      type: "number",
      placeholder: "1000",
      min: 100,
      max: 60000,
    },
    {
      name: "memory_limit_mb",
      label: "Memory Limit (MB)",
      type: "number",
      placeholder: "128",
      min: 16,
      max: 1024,
    },
  ];

  async function fetchProblems() {
    try {
      isLoading = true;
      error = null;
      problems = await api.getGroupProblems(group.id);
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

  function handleEditProblem(problem: Problem) {
    editingProblem = problem;
    isFormModalOpen = true;
  }

  function handleDeleteClick(problem: Problem) {
    deletingProblem = problem;
    isDeleteDialogOpen = true;
  }

  async function handleFormSubmit(
    event: CustomEvent<Record<string, string | number>>,
  ) {
    const { name, description, problem_text, time_limit_ms, memory_limit_mb } =
      event.detail;
    isSubmitting = true;

    try {
      const input = {
        name: String(name),
        description: String(description || ""),
        problem_text: String(problem_text),
        time_limit_ms: time_limit_ms ? Number(time_limit_ms) : undefined,
        memory_limit_mb: memory_limit_mb ? Number(memory_limit_mb) : undefined,
      };

      if (editingProblem) {
        const updated = await api.updateProblem(editingProblem.id, input);
        problems = problems.map((p) => (p.id === updated.id ? updated : p));
      } else {
        const created = await api.createProblem(group.id, input);
        problems = [created, ...problems];
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
      await api.deleteProblem(deletingProblem.id);
      problems = problems.filter((p) => p.id !== deletingProblem!.id);
      isDeleteDialogOpen = false;
      deletingProblem = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete problem";
    } finally {
      isSubmitting = false;
    }
  }

  function handleProblemClick(problem: Problem) {
    selectedProblem = problem;
  }

  function handleProblemBack() {
    selectedProblem = null;
    fetchProblems();
  }
</script>

{#if selectedProblem}
  <ProblemDetail problem={selectedProblem} on:back={handleProblemBack} />
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
      <Button size="md" on:click={handleNewProblem}>New Problem</Button>
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
        <p class="mt-4">Loading problems...</p>
      </div>
    {:else if problems.length === 0}
      <div class="py-12 text-center text-gray-500">
        <p class="text-lg font-medium">No problems yet</p>
        <p class="text-sm mt-1">Create your first problem to get started</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each problems as problem, idx (problem.id)}
          <div class="relative group">
            <button
              class="w-full flex items-center gap-4 bg-white border border-gray-200 rounded-xl p-4 hover:border-sky-300 hover:shadow-sm transition-all text-left"
              on:click={() => handleProblemClick(problem)}
            >
              <div
                class="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center text-blue-600 font-bold"
              >
                {idx + 1}
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-semibold text-gray-900 truncate">
                  {problem.name}
                </h3>
                {#if problem.description}
                  <p class="text-sm text-gray-500 truncate">
                    {problem.description}
                  </p>
                {/if}
              </div>
              <div class="text-xs text-gray-400 flex gap-3">
                {#if problem.time_limit_ms}
                  <span>⏱ {problem.time_limit_ms}ms</span>
                {/if}
                {#if problem.memory_limit_mb}
                  <span>💾 {problem.memory_limit_mb}MB</span>
                {/if}
              </div>
            </button>
            <div
              class="absolute top-1/2 -translate-y-1/2 right-4 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1"
            >
              <button
                class="p-2 bg-gray-100 rounded-lg hover:bg-gray-200 text-gray-600"
                on:click|stopPropagation={() => handleEditProblem(problem)}
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
                on:click|stopPropagation={() => handleDeleteClick(problem)}
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

<FormModal
  bind:isOpen={isFormModalOpen}
  title={editingProblem ? "Edit Problem" : "Create New Problem"}
  submitText={editingProblem ? "Save" : "Create"}
  fields={formFields}
  initialValues={editingProblem
    ? {
        name: editingProblem.name,
        description: editingProblem.description,
        problem_text: editingProblem.problem_text,
        time_limit_ms: editingProblem.time_limit_ms || 0,
        memory_limit_mb: editingProblem.memory_limit_mb || 0,
      }
    : { time_limit_ms: 1000, memory_limit_mb: 128 }}
  isLoading={isSubmitting}
  on:submit={handleFormSubmit}
  on:cancel={() => (isFormModalOpen = false)}
/>

<ConfirmDialog
  bind:isOpen={isDeleteDialogOpen}
  title="Delete Problem"
  message="Are you sure you want to delete '{deletingProblem?.name}'? This will also delete all test cases for this problem."
  isLoading={isSubmitting}
  on:confirm={handleDeleteConfirm}
  on:cancel={() => {
    isDeleteDialogOpen = false;
    deletingProblem = null;
  }}
/>
