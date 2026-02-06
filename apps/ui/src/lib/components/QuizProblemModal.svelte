<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { fade, scale } from "svelte/transition";
  import Button from "$lib/components/Button.svelte";
  import type {
    QuizProblem,
    QuizProblemInput,
    QuizProblemType,
    QuizOption,
  } from "$lib/api";

  export let isOpen = false;
  export let editingProblem: QuizProblem | null = null;
  export let isLoading = false;

  const dispatch = createEventDispatcher<{
    submit: QuizProblemInput;
    cancel: void;
  }>();

  // Form values
  let problemType: QuizProblemType = "open_ended";
  let name = "";
  let description = "";
  let problemText = "";
  let points = 5;
  let correctAnswer = "";
  let options: {
    option_text: string;
    is_correct: boolean;
    display_order: number;
  }[] = [];

  const problemTypes: { value: QuizProblemType; label: string }[] = [
    { value: "open_ended", label: "📝 Open Ended" },
    { value: "true_false", label: "✅ True/False" },
    { value: "mcq_single", label: "🔘 Single Choice" },
    { value: "mcq_multi", label: "☑️ Multiple Choice" },
    { value: "fill_blank", label: "📄 Fill in Blank" },
  ];

  $: if (isOpen) {
    if (editingProblem) {
      problemType = editingProblem.problem_type;
      name = editingProblem.name;
      description = editingProblem.description;
      problemText = editingProblem.problem_text;
      points = editingProblem.points;
      correctAnswer = editingProblem.correct_answer || "";
      options =
        editingProblem.options?.map((o) => ({
          option_text: o.option_text,
          is_correct: o.is_correct,
          display_order: o.display_order,
        })) || [];
    } else {
      resetForm();
    }
  }

  $: if (problemType === "true_false" && options.length !== 2) {
    options = [
      { option_text: "True", is_correct: true, display_order: 1 },
      { option_text: "False", is_correct: false, display_order: 2 },
    ];
  }

  $: needsOptions = ["true_false", "mcq_single", "mcq_multi"].includes(
    problemType,
  );
  $: needsCorrectAnswer = problemType === "fill_blank";

  function resetForm() {
    problemType = "open_ended";
    name = "";
    description = "";
    problemText = "";
    points = 5;
    correctAnswer = "";
    options = [];
  }

  function addOption() {
    options = [
      ...options,
      {
        option_text: "",
        is_correct: false,
        display_order: options.length + 1,
      },
    ];
  }

  function removeOption(index: number) {
    options = options
      .filter((_, i) => i !== index)
      .map((o, i) => ({
        ...o,
        display_order: i + 1,
      }));
  }

  function toggleCorrect(index: number) {
    if (problemType === "mcq_single" || problemType === "true_false") {
      // Only one can be correct
      options = options.map((o, i) => ({
        ...o,
        is_correct: i === index,
      }));
    } else {
      // Multiple can be correct
      options = options.map((o, i) =>
        i === index ? { ...o, is_correct: !o.is_correct } : o,
      );
    }
  }

  function handleSubmit() {
    if (!name.trim() || !problemText.trim()) return;

    const input: QuizProblemInput = {
      problem_type: problemType,
      name: name.trim(),
      description: description.trim(),
      problem_text: problemText.trim(),
      points,
    };

    if (needsCorrectAnswer) {
      input.correct_answer = correctAnswer.trim();
    }

    if (needsOptions && options.length > 0) {
      input.options = options
        .filter((o) => o.option_text.trim())
        .map((o) => ({
          option_text: o.option_text.trim(),
          is_correct: o.is_correct,
          display_order: o.display_order,
        }));
    }

    dispatch("submit", input);
  }

  function handleCancel() {
    dispatch("cancel");
    isOpen = false;
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      handleCancel();
    }
  }
</script>

{#if isOpen}
  <div
    class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4"
    on:click={handleCancel}
    on:keydown={handleKeydown}
    role="button"
    tabindex="0"
    transition:fade={{ duration: 200 }}
  >
    <div
      class="bg-white rounded-2xl shadow-xl max-w-2xl w-full p-6 max-h-[90vh] overflow-y-auto"
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <h2 class="text-2xl font-semibold text-gray-900 mb-6">
        {editingProblem ? "Edit Question" : "Create Question"}
      </h2>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <!-- Problem Type -->
        <div>
          <label
            for="type"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Question Type <span class="text-red-500">*</span>
          </label>
          <select
            id="type"
            bind:value={problemType}
            class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 bg-white"
            disabled={isLoading}
          >
            {#each problemTypes as pt}
              <option value={pt.value}>{pt.label}</option>
            {/each}
          </select>
        </div>

        <!-- Name -->
        <div>
          <label
            for="name"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Name <span class="text-red-500">*</span>
          </label>
          <input
            id="name"
            type="text"
            bind:value={name}
            placeholder="Question name..."
            class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
            disabled={isLoading}
            required
          />
        </div>

        <!-- Description -->
        <div>
          <label
            for="description"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Description
          </label>
          <input
            id="description"
            type="text"
            bind:value={description}
            placeholder="Brief description..."
            class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
            disabled={isLoading}
          />
        </div>

        <!-- Problem Text -->
        <div>
          <label
            for="problem_text"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Question Text <span class="text-red-500">*</span>
          </label>
          <textarea
            id="problem_text"
            bind:value={problemText}
            placeholder="The actual question students will see..."
            rows="4"
            class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 resize-none"
            disabled={isLoading}
            required
          ></textarea>
        </div>

        <!-- Points -->
        <div>
          <label
            for="points"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Points <span class="text-red-500">*</span>
          </label>
          <input
            id="points"
            type="number"
            bind:value={points}
            min="0"
            max="100"
            class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
            disabled={isLoading}
            required
          />
        </div>

        <!-- Correct Answer (for fill_blank) -->
        {#if needsCorrectAnswer}
          <div>
            <label
              for="correct_answer"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Correct Answer <span class="text-red-500">*</span>
            </label>
            <input
              id="correct_answer"
              type="text"
              bind:value={correctAnswer}
              placeholder="The correct answer..."
              class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
              disabled={isLoading}
              required
            />
          </div>
        {/if}

        <!-- Options (for MCQ/TrueFalse) -->
        {#if needsOptions}
          <div>
            <div class="flex items-center justify-between mb-2">
              <label class="block text-sm font-medium text-gray-700">
                Options <span class="text-red-500">*</span>
              </label>
              {#if problemType !== "true_false"}
                <button
                  type="button"
                  class="text-sm text-sky-500 hover:text-sky-600"
                  on:click={addOption}
                  disabled={isLoading}
                >
                  + Add Option
                </button>
              {/if}
            </div>
            <div class="space-y-2">
              {#each options as option, idx}
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    class="w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors
											{option.is_correct
                      ? 'border-green-500 bg-green-500 text-white'
                      : 'border-gray-300 hover:border-sky-400'}"
                    on:click={() => toggleCorrect(idx)}
                    disabled={isLoading}
                    title={option.is_correct ? "Correct" : "Mark as correct"}
                  >
                    {#if option.is_correct}
                      <svg
                        class="w-4 h-4"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                          clip-rule="evenodd"
                        />
                      </svg>
                    {/if}
                  </button>
                  <input
                    type="text"
                    bind:value={option.option_text}
                    placeholder="Option {idx + 1}..."
                    class="flex-1 px-3 py-2 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
                    disabled={isLoading || problemType === "true_false"}
                  />
                  {#if problemType !== "true_false"}
                    <button
                      type="button"
                      class="p-2 text-red-500 hover:bg-red-50 rounded"
                      on:click={() => removeOption(idx)}
                      disabled={isLoading}
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
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                    </button>
                  {/if}
                </div>
              {/each}
            </div>
            {#if problemType === "mcq_single"}
              <p class="text-xs text-gray-500 mt-1">
                Click the circle to mark the correct answer
              </p>
            {:else if problemType === "mcq_multi"}
              <p class="text-xs text-gray-500 mt-1">
                Click circles to mark all correct answers
              </p>
            {/if}
          </div>
        {/if}

        <div class="flex gap-3 mt-6 justify-end pt-4">
          <button
            type="button"
            on:click={handleCancel}
            class="px-6 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-semibold hover:bg-gray-50 transition-colors"
            disabled={isLoading}
          >
            Cancel
          </button>
          <Button size="sm" type="submit" disabled={isLoading}>
            {#if isLoading}
              <span
                class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
              ></span>
            {/if}
            {editingProblem ? "Save" : "Create"}
          </Button>
        </div>
      </form>
    </div>
  </div>
{/if}
