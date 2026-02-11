<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import { api, type Course, type Enrollment, type Student } from "$lib/api";

  export let course: Course;

  const dispatch = createEventDispatcher<{ back: void }>();

  let enrollments: Enrollment[] = [];
  let allStudents: Student[] = [];
  let isLoading = true;
  let error: string | null = null;

  let isAddStudentModalOpen = false;
  let isRemoveDialogOpen = false;
  let isSubmitting = false;
  let removingEnrollment: Enrollment | null = null;
  let selectedStudentId: number | null = null;

  async function fetchData() {
    try {
      isLoading = true;
      error = null;
      const [enrollmentsData, studentsData] = await Promise.all([
        api.getCourseEnrollments(course.id),
        api.getSchoolStudents(),
      ]);
      enrollments = enrollmentsData;
      allStudents = studentsData;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load data";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchData();
  });

  $: availableStudents = allStudents.filter(
    (student) => !enrollments.some((e) => e.user_id === student.id),
  );

  function handleBack() {
    dispatch("back");
  }

  function handleAddStudent() {
    selectedStudentId = null;
    isAddStudentModalOpen = true;
  }

  function handleRemoveClick(enrollment: Enrollment) {
    removingEnrollment = enrollment;
    isRemoveDialogOpen = true;
  }

  async function handleEnrollStudent() {
    if (!selectedStudentId) return;
    isSubmitting = true;

    try {
      const enrollment = await api.enrollStudent(course.id, selectedStudentId);
      const student = allStudents.find((s) => s.id === selectedStudentId);
      enrollment.user = student;
      enrollments = [...enrollments, enrollment];
      isAddStudentModalOpen = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to enroll student";
    } finally {
      isSubmitting = false;
    }
  }

  async function handleRemoveConfirm() {
    if (!removingEnrollment) return;
    isSubmitting = true;

    try {
      await api.unenrollStudent(course.id, removingEnrollment.user_id);
      enrollments = enrollments.filter((e) => e.id !== removingEnrollment!.id);
      isRemoveDialogOpen = false;
      removingEnrollment = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to remove student";
    } finally {
      isSubmitting = false;
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

<div class="w-full">
  <div class="flex items-center gap-4 mb-6">
    <BackButton onclick={handleBack} />
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">{course.name}</h2>
      <p class="text-gray-500">{course.subject}</p>
    </div>
  </div>

  {#if error}
    <div
      class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
    >
      {error}
      <button class="ml-2 underline" on:click={fetchData}>Retry</button>
    </div>
  {/if}

  <div class="bg-gray-50 rounded-xl p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800">Enrolled Students</h3>
      <Button size="sm" on:click={handleAddStudent}>Add Student</Button>
    </div>

    {#if isLoading}
      <div class="py-8 text-center text-gray-500">
        <div
          class="inline-block w-6 h-6 border-3 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-2 text-sm">Loading...</p>
      </div>
    {:else if enrollments.length === 0}
      <div class="py-8 text-center text-gray-500">
        <p>No students enrolled yet</p>
        <p class="text-sm mt-1">Add students to this course</p>
      </div>
    {:else}
      <div class="space-y-2">
        {#each enrollments as enrollment (enrollment.id)}
          <div
            class="flex items-center justify-between bg-white rounded-lg p-4 border border-gray-200"
          >
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold text-sm"
              >
                {getInitials(enrollment.user?.full_name || "?")}
              </div>
              <div>
                <p class="font-medium text-gray-900">
                  {enrollment.user?.full_name || `User ${enrollment.user_id}`}
                </p>
                <p class="text-sm text-gray-500">
                  {enrollment.user?.email || ""}
                </p>
              </div>
            </div>
            <button
              class="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
              on:click={() => handleRemoveClick(enrollment)}
              title="Remove from course"
            >
              <svg
                class="w-5 h-5"
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
        {/each}
      </div>
    {/if}
  </div>
</div>

{#if isAddStudentModalOpen}
  <div
    class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4"
    on:click={() => (isAddStudentModalOpen = false)}
    on:keydown={(e) => e.key === "Escape" && (isAddStudentModalOpen = false)}
    role="button"
    tabindex="0"
  >
    <div
      class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 max-h-[80vh] overflow-y-auto"
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
    >
      <h2 class="text-xl font-semibold text-gray-900 mb-4">
        Add Student to Course
      </h2>

      {#if availableStudents.length === 0}
        <p class="text-gray-500 py-4">
          All students are already enrolled in this course.
        </p>
      {:else}
        <div class="space-y-2 max-h-64 overflow-y-auto">
          {#each availableStudents as student (student.id)}
            <button
              class="w-full flex items-center gap-3 p-3 rounded-lg border-2 transition-colors
								{selectedStudentId === student.id
                ? 'border-sky-400 bg-sky-50'
                : 'border-gray-200 hover:border-gray-300'}"
              on:click={() => (selectedStudentId = student.id)}
            >
              <div
                class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold text-sm"
              >
                {getInitials(student.full_name)}
              </div>
              <div class="text-left">
                <p class="font-medium text-gray-900">{student.full_name}</p>
                <p class="text-sm text-gray-500">{student.email}</p>
              </div>
            </button>
          {/each}
        </div>
      {/if}

      <div class="flex gap-3 mt-6 justify-end">
        <button
          class="px-5 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-medium hover:bg-gray-50"
          on:click={() => (isAddStudentModalOpen = false)}
        >
          Cancel
        </button>
        {#if availableStudents.length > 0}
          <Button
            size="sm"
            on:click={handleEnrollStudent}
            disabled={!selectedStudentId || isSubmitting}
          >
            {#if isSubmitting}
              <span
                class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
              ></span>
            {/if}
            Add
          </Button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<ConfirmDialog
  bind:isOpen={isRemoveDialogOpen}
  title="Remove Student"
  message="Are you sure you want to remove '{removingEnrollment?.user
    ?.full_name || 'this student'}' from {course.name}?"
  confirmText="Remove"
  isLoading={isSubmitting}
  on:confirm={handleRemoveConfirm}
  on:cancel={() => {
    isRemoveDialogOpen = false;
    removingEnrollment = null;
  }}
/>
