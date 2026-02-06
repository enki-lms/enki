<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { fade, scale } from "svelte/transition";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import {
    api,
    type Course,
    type ProblemGroup,
    type QuizGroup,
    type Student,
    type Enrollment,
    type ExamSessionInput,
  } from "$lib/api";

  export let isOpen = false;

  const dispatch = createEventDispatcher<{
    created: void;
    cancel: void;
  }>();

  let courses: Course[] = [];
  let problemGroups: ProblemGroup[] = [];
  let quizGroups: QuizGroup[] = [];
  let courseEnrollments: Enrollment[] = [];
  let isLoading = true;
  let isSubmitting = false;
  let error: string | null = null;

  // Form values
  let selectedCourseId: number | "" = "";
  let problemGroupType: "comp_sci" | "quiz" = "comp_sci";
  let selectedGroupId: number | "" = "";
  let durationMinutes = 60;
  let selectedStudentIds: Set<number> = new Set();

  $: if (isOpen) {
    fetchInitialData();
  }

  $: if (selectedCourseId) {
    fetchCourseData();
  }

  $: availableGroups =
    problemGroupType === "comp_sci" ? problemGroups : quizGroups;

  async function fetchInitialData() {
    try {
      isLoading = true;
      error = null;
      courses = await api.getCourses();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load data";
    } finally {
      isLoading = false;
    }
  }

  async function fetchCourseData() {
    if (!selectedCourseId) return;
    try {
      error = null;
      const [pg, qg, enrollments] = await Promise.all([
        api.getCourseProblemGroups(selectedCourseId as number),
        api.getCourseQuizGroups(selectedCourseId as number),
        api.getCourseEnrollments(selectedCourseId as number),
      ]);
      problemGroups = pg;
      quizGroups = qg;
      courseEnrollments = enrollments;
      selectedGroupId = "";
      selectedStudentIds = new Set();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load course data";
    }
  }

  function toggleStudent(studentId: number) {
    if (selectedStudentIds.has(studentId)) {
      selectedStudentIds.delete(studentId);
    } else {
      selectedStudentIds.add(studentId);
    }
    selectedStudentIds = new Set(selectedStudentIds);
  }

  function selectAllStudents() {
    selectedStudentIds = new Set(courseEnrollments.map((e) => e.user_id));
  }

  function deselectAllStudents() {
    selectedStudentIds = new Set();
  }

  async function handleSubmit() {
    if (!selectedGroupId || selectedStudentIds.size === 0) return;
    isSubmitting = true;

    try {
      const input: ExamSessionInput = {
        problem_group_type: problemGroupType,
        problem_group_id: selectedGroupId as number,
        duration_minutes: durationMinutes,
        student_ids: Array.from(selectedStudentIds),
      };
      await api.createExamSession(input);
      dispatch("created");
      resetForm();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to create exam session";
    } finally {
      isSubmitting = false;
    }
  }

  function resetForm() {
    selectedCourseId = "";
    problemGroupType = "comp_sci";
    selectedGroupId = "";
    durationMinutes = 60;
    selectedStudentIds = new Set();
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

  function getInitials(name: string): string {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
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
        Create Exam Session
      </h2>

      {#if error}
        <div
          class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm"
        >
          {error}
        </div>
      {/if}

      {#if isLoading}
        <div class="py-12 text-center text-gray-500">
          <div
            class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
          ></div>
          <p class="mt-4">Loading...</p>
        </div>
      {:else}
        <form on:submit|preventDefault={handleSubmit} class="space-y-5">
          <!-- Course Selection -->
          <div>
            <label
              for="course"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Course <span class="text-red-500">*</span>
            </label>
            <Select
              id="course"
              bind:value={selectedCourseId}
              placeholder="Select a course..."
              options={courses.map((c) => ({ value: c.id, label: c.name }))}
            />
          </div>

          <!-- Problem Group Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Exam Type <span class="text-red-500">*</span>
            </label>
            <div class="flex gap-2">
              <button
                type="button"
                class="flex-1 py-3 px-4 rounded-lg border-2 transition-colors
									{problemGroupType === 'comp_sci'
                  ? 'border-sky-400 bg-sky-50 text-sky-700'
                  : 'border-gray-200 hover:border-gray-300'}"
                on:click={() => (problemGroupType = "comp_sci")}
              >
                💻 CS Problems
              </button>
              <button
                type="button"
                class="flex-1 py-3 px-4 rounded-lg border-2 transition-colors
									{problemGroupType === 'quiz'
                  ? 'border-sky-400 bg-sky-50 text-sky-700'
                  : 'border-gray-200 hover:border-gray-300'}"
                on:click={() => (problemGroupType = "quiz")}
              >
                📝 Quiz
              </button>
            </div>
          </div>

          <!-- Problem Group Selection -->
          {#if selectedCourseId}
            <div>
              <label
                for="group"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                Problem Group <span class="text-red-500">*</span>
              </label>
              {#if availableGroups.length === 0}
                <p class="text-gray-500 text-sm py-2">
                  No {problemGroupType === "comp_sci"
                    ? "problem groups"
                    : "quiz groups"} found for this course.
                </p>
              {:else}
                <Select
                  id="group"
                  bind:value={selectedGroupId}
                  placeholder="Select a {problemGroupType === 'comp_sci'
                    ? 'problem group'
                    : 'quiz group'}..."
                  options={availableGroups.map((g) => ({
                    value: g.id,
                    label: `${g.name} (${g.type})`,
                  }))}
                />
              {/if}
            </div>
          {/if}

          <!-- Duration -->
          <div>
            <label
              for="duration"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Duration (minutes) <span class="text-red-500">*</span>
            </label>
            <input
              id="duration"
              type="number"
              bind:value={durationMinutes}
              min="5"
              max="480"
              class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400"
            />
          </div>

          <!-- Student Selection -->
          {#if selectedCourseId && courseEnrollments.length > 0}
            <div>
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">
                  Students <span class="text-red-500">*</span>
                  <span class="text-gray-500 font-normal"
                    >({selectedStudentIds.size} selected)</span
                  >
                </label>
                <div class="flex gap-2">
                  <button
                    type="button"
                    class="text-xs text-sky-500 hover:text-sky-600"
                    on:click={selectAllStudents}
                  >
                    Select All
                  </button>
                  <button
                    type="button"
                    class="text-xs text-gray-500 hover:text-gray-600"
                    on:click={deselectAllStudents}
                  >
                    Clear
                  </button>
                </div>
              </div>
              <div
                class="max-h-48 overflow-y-auto border-2 border-gray-200 rounded-lg p-2 space-y-1"
              >
                {#each courseEnrollments as enrollment (enrollment.id)}
                  <button
                    type="button"
                    class="w-full flex items-center gap-3 p-2 rounded-lg transition-colors
											{selectedStudentIds.has(enrollment.user_id)
                      ? 'bg-sky-50 border border-sky-200'
                      : 'hover:bg-gray-50'}"
                    on:click={() => toggleStudent(enrollment.user_id)}
                  >
                    <div
                      class="w-5 h-5 rounded border-2 flex items-center justify-center transition-colors
												{selectedStudentIds.has(enrollment.user_id)
                        ? 'border-sky-500 bg-sky-500 text-white'
                        : 'border-gray-300'}"
                    >
                      {#if selectedStudentIds.has(enrollment.user_id)}
                        <svg
                          class="w-3 h-3"
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
                    </div>
                    <div
                      class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xs font-medium"
                    >
                      {getInitials(enrollment.user?.full_name || "?")}
                    </div>
                    <span class="text-sm text-gray-900"
                      >{enrollment.user?.full_name ||
                        `User ${enrollment.user_id}`}</span
                    >
                  </button>
                {/each}
              </div>
            </div>
          {:else if selectedCourseId}
            <p class="text-gray-500 text-sm">
              No students enrolled in this course.
            </p>
          {/if}

          <div class="flex gap-3 mt-6 justify-end pt-4">
            <button
              type="button"
              on:click={handleCancel}
              class="px-6 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-semibold hover:bg-gray-50 transition-colors"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <Button
              size="sm"
              type="submit"
              disabled={isSubmitting ||
                !selectedGroupId ||
                selectedStudentIds.size === 0}
            >
              {#if isSubmitting}
                <span
                  class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
                ></span>
              {/if}
              Create Session
            </Button>
          </div>
        </form>
      {/if}
    </div>
  </div>
{/if}
