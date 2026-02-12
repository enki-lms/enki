<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import FormModal, {
    type FieldConfig,
  } from "$lib/components/FormModal.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import { api, type Problem, type TestCase } from "$lib/api";

  export let problem: Problem;

  const dispatch = createEventDispatcher<{ back: void }>();

  let testCases: TestCase[] = [];
  let isLoading = true;
  let error: string | null = null;

  let isFormModalOpen = false;
  let isDeleteDialogOpen = false;
  let isSubmitting = false;
  let editingTestCase: TestCase | null = null;
  let deletingTestCase: TestCase | null = null;

  const formFields: FieldConfig[] = [
    {
      name: "input",
      label: "Input",
      type: "textarea",
      placeholder: "Test case input...",
      required: true,
    },
    {
      name: "output",
      label: "Expected Output",
      type: "textarea",
      placeholder: "Expected output...",
      required: true,
    },
    {
      name: "correct_points",
      label: "Points",
      type: "number",
      placeholder: "10",
      min: 0,
      max: 100,
      required: true,
    },
  ];

  /* 
     Dynamically generate fields based on problem type.
     Ideally we'd use a derived store or reactive declaration for fields, 
     but `formFields` is a const. We should make it reactive or compute it.
  */
  let dyanmicFields = formFields;
  $: {
    if (problem.type === "turtle") {
      dyanmicFields = [
        ...formFields.filter((f) => f.name !== "output"), // Remove output for turtle? Or keep it as empty/ignored?
        // Turtle usually doesn't need text output match if we do image match.
        // But schema requires `output` column NOT NULL. We can set a default.
        // Let's keep input (turtle code snippet? no, input is stdin).
        // For turtle, input might be empty.
        {
          name: "image_file",
          label: "Ideal Image",
          type: "file",
          accept: "image/*",
          required: false, // Optional on edit if already has one
        },
      ];
    } else {
      dyanmicFields = formFields;
    }
  }

  async function fetchTestCases() {
    try {
      isLoading = true;
      error = null;
      testCases = await api.getProblemTestCases(problem.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load test cases";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchTestCases();
  });

  function handleBack() {
    dispatch("back");
  }

  function handleNewTestCase() {
    editingTestCase = null;
    isFormModalOpen = true;
  }

  function handleEditTestCase(tc: TestCase) {
    editingTestCase = tc;
    isFormModalOpen = true;
  }

  function handleDeleteClick(tc: TestCase) {
    deletingTestCase = tc;
    isDeleteDialogOpen = true;
  }

  async function handleFormSubmit(event: CustomEvent<Record<string, any>>) {
    const { input, output, correct_points, image_file } = event.detail;
    isSubmitting = true;

    try {
      const inputData = {
        input: String(input || ""),
        output: String(output || ""), // Output can be empty for turtle
        correct_points: Number(correct_points),
      };

      if (image_file && image_file instanceof File) {
        const uploadRes = await api.uploadFile(image_file);
        inputData.image_url = uploadRes.url;
      } else if (editingTestCase?.image_url) {
        // Keep existing image if not replacing
        inputData.image_url = editingTestCase.image_url;
      }

      if (editingTestCase) {
        const updated = await api.updateTestCase(editingTestCase.id, inputData);
        testCases = testCases.map((tc) =>
          tc.id === updated.id ? updated : tc,
        );
      } else {
        const created = await api.createTestCase(problem.id, inputData);
        testCases = [...testCases, created];
      }
      isFormModalOpen = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save test case";
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteConfirm() {
    if (!deletingTestCase) return;
    isSubmitting = true;

    try {
      await api.deleteTestCase(deletingTestCase.id);
      testCases = testCases.filter((tc) => tc.id !== deletingTestCase!.id);
      isDeleteDialogOpen = false;
      deletingTestCase = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete test case";
    } finally {
      isSubmitting = false;
    }
  }

  $: totalPoints = testCases.reduce((sum, tc) => sum + tc.correct_points, 0);
</script>

<div class="w-full">
  <div class="flex items-center gap-4 mb-6">
    <BackButton onclick={handleBack} />
    <div class="flex-1">
      <h2 class="text-2xl font-semibold text-gray-900">{problem.name}</h2>
      <div class="flex gap-3 text-sm text-gray-500">
        {#if problem.time_limit_ms}
          <span>⏱ {problem.time_limit_ms}ms</span>
        {/if}
        {#if problem.memory_limit_mb}
          <span>💾 {problem.memory_limit_mb}MB</span>
        {/if}
      </div>
    </div>
  </div>

  {#if problem.description}
    <div class="mb-4 p-4 bg-gray-50 rounded-lg">
      <h3 class="text-sm font-medium text-gray-700 mb-1">Description</h3>
      <p class="text-gray-600">{problem.description}</p>
    </div>
  {/if}

  <div class="mb-6 p-4 bg-blue-50 rounded-lg">
    <h3 class="text-sm font-medium text-gray-700 mb-1">Problem Text</h3>
    <pre
      class="text-gray-800 whitespace-pre-wrap text-sm">{problem.problem_text}</pre>
  </div>

  {#if error}
    <div
      class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
    >
      {error}
      <button class="ml-2 underline" on:click={fetchTestCases}>Retry</button>
    </div>
  {/if}

  <div class="bg-gray-50 rounded-xl p-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h3 class="text-lg font-semibold text-gray-800">Test Cases</h3>
        <p class="text-sm text-gray-500">Total points: {totalPoints}</p>
      </div>
      <Button size="sm" on:click={handleNewTestCase}>Add Test Case</Button>
    </div>

    {#if isLoading}
      <div class="py-8 text-center text-gray-500">
        <div
          class="inline-block w-6 h-6 border-3 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-2 text-sm">Loading...</p>
      </div>
    {:else if testCases.length === 0}
      <div class="py-8 text-center text-gray-500">
        <p>No test cases yet</p>
        <p class="text-sm mt-1">Add test cases to verify student submissions</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each testCases as tc, idx (tc.id)}
          <div
            class="bg-white rounded-lg border border-gray-200 p-4 group relative"
          >
            <div class="flex items-start justify-between mb-3">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                Test #{idx + 1}
              </span>
              <span class="text-sm font-medium text-green-600"
                >{tc.correct_points} pts</span
              >
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-gray-500 mb-1"
                  >Input</label
                >
                <pre
                  class="text-sm bg-gray-50 p-2 rounded overflow-x-auto max-h-24">{tc.input}</pre>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 mb-1"
                  >Expected Output</label
                >
                <pre
                  class="text-sm bg-gray-50 p-2 rounded overflow-x-auto max-h-24">{tc.output}</pre>
              </div>
              {#if tc.image_url}
                <div class="col-span-2">
                  <label class="block text-xs font-medium text-gray-500 mb-1"
                    >Ideal Image</label
                  >
                  <img
                    src={tc.image_url}
                    alt="Ideal Solution"
                    class="h-32 w-auto object-contain border rounded bg-gray-50"
                  />
                </div>
              {/if}
            </div>
            <div
              class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1"
            >
              <button
                class="p-1.5 bg-gray-100 rounded hover:bg-gray-200 text-gray-600"
                on:click={() => handleEditTestCase(tc)}
                title="Edit"
              >
                <svg
                  class="w-3.5 h-3.5"
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
                class="p-1.5 bg-gray-100 rounded hover:bg-red-100 text-red-500"
                on:click={() => handleDeleteClick(tc)}
                title="Delete"
              >
                <svg
                  class="w-3.5 h-3.5"
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
</div>

<FormModal
  bind:isOpen={isFormModalOpen}
  title={editingTestCase ? "Edit Test Case" : "Add Test Case"}
  submitText={editingTestCase ? "Save" : "Add"}
  fields={dyanmicFields}
  initialValues={editingTestCase
    ? {
        input: editingTestCase.input,
        output: editingTestCase.output,
        correct_points: editingTestCase.correct_points,
        image_file: editingTestCase.image_url, // Pass existing URL for preview
      }
    : { correct_points: 10 }}
  isLoading={isSubmitting}
  on:submit={handleFormSubmit}
  on:cancel={() => (isFormModalOpen = false)}
/>

<ConfirmDialog
  bind:isOpen={isDeleteDialogOpen}
  title="Delete Test Case"
  message="Are you sure you want to delete this test case?"
  isLoading={isSubmitting}
  on:confirm={handleDeleteConfirm}
  on:cancel={() => {
    isDeleteDialogOpen = false;
    deletingTestCase = null;
  }}
/>
