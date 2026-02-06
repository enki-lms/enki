<script lang="ts">
  import { onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import PlusIcon from "$lib/components/icons/PlusIcon.svelte";
  import EmptyBookIcon from "$lib/components/icons/EmptyBookIcon.svelte";
  import ItemCard from "$lib/components/ItemCard.svelte";
  import FormModal, {
    type FieldConfig,
  } from "$lib/components/FormModal.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import CourseDetail from "$lib/components/CourseDetail.svelte";
  import { api, type Course } from "$lib/api";

  let courses: Course[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Modal states
  let isFormModalOpen = false;
  let isDeleteDialogOpen = false;
  let isSubmitting = false;
  let editingCourse: Course | null = null;
  let deletingCourse: Course | null = null;

  // Detail view
  let selectedCourse: Course | null = null;

  const formFields: FieldConfig[] = [
    {
      name: "name",
      label: "Course Name",
      type: "text",
      placeholder: "Enter course name...",
      required: true,
    },
    {
      name: "subject",
      label: "Subject",
      type: "text",
      placeholder: "e.g., Computer Science, Mathematics...",
      required: true,
    },
  ];

  async function fetchCourses() {
    try {
      isLoading = true;
      error = null;
      courses = await api.getCourses();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load courses";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchCourses();
  });

  function handleNewCourse() {
    editingCourse = null;
    isFormModalOpen = true;
  }

  function handleEditCourse(course: Course) {
    editingCourse = course;
    isFormModalOpen = true;
  }

  function handleDeleteClick(course: Course) {
    deletingCourse = course;
    isDeleteDialogOpen = true;
  }

  async function handleFormSubmit(
    event: CustomEvent<Record<string, string | number>>,
  ) {
    const { name, subject } = event.detail;
    isSubmitting = true;

    try {
      if (editingCourse) {
        const updated = await api.updateCourse(editingCourse.id, {
          name: String(name),
          subject: String(subject),
        });
        courses = courses.map((c) => (c.id === updated.id ? updated : c));
      } else {
        const created = await api.createCourse({
          name: String(name),
          subject: String(subject),
        });
        courses = [created, ...courses];
      }
      isFormModalOpen = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save course";
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteConfirm() {
    if (!deletingCourse) return;
    isSubmitting = true;

    try {
      await api.deleteCourse(deletingCourse.id);
      courses = courses.filter((c) => c.id !== deletingCourse!.id);
      isDeleteDialogOpen = false;
      deletingCourse = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete course";
    } finally {
      isSubmitting = false;
    }
  }

  function handleCardClick(course: Course) {
    selectedCourse = course;
  }

  function handleBack() {
    selectedCourse = null;
    fetchCourses(); // Refresh in case enrollments changed
  }
</script>

{#if selectedCourse}
  <CourseDetail course={selectedCourse} on:back={handleBack} />
{:else}
  <div class="w-full">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-semibold text-gray-900">Courses</h2>
      <Button size="md" on:click={handleNewCourse}>
        <span slot="icon" class="text-white">
          <PlusIcon />
        </span>
        New Course
      </Button>
    </div>

    {#if error}
      <div
        class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
      >
        {error}
        <button class="ml-2 underline" on:click={fetchCourses}>Retry</button>
      </div>
    {/if}

    {#if isLoading}
      <div class="py-16 text-center text-gray-500">
        <div
          class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-4">Loading courses...</p>
      </div>
    {:else if courses.length === 0}
      <div class="py-16 text-center text-gray-500">
        <EmptyBookIcon />
        <p class="text-lg font-medium">No courses yet</p>
        <p class="text-sm mt-1">Create your first course to get started</p>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each courses as course (course.id)}
          <div class="relative group">
            <ItemCard
              title={course.name}
              description={course.subject}
              date=""
              type="course"
              on:click={() => handleCardClick(course)}
            />
            <div
              class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1"
            >
              <button
                class="p-2 bg-white rounded-lg shadow-md hover:bg-gray-50 text-gray-600"
                on:click|stopPropagation={() => handleEditCourse(course)}
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
                class="p-2 bg-white rounded-lg shadow-md hover:bg-red-50 text-red-500"
                on:click|stopPropagation={() => handleDeleteClick(course)}
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
  title={editingCourse ? "Edit Course" : "Create New Course"}
  submitText={editingCourse ? "Save" : "Create"}
  fields={formFields}
  initialValues={editingCourse
    ? { name: editingCourse.name, subject: editingCourse.subject }
    : {}}
  isLoading={isSubmitting}
  on:submit={handleFormSubmit}
  on:cancel={() => (isFormModalOpen = false)}
/>

<ConfirmDialog
  bind:isOpen={isDeleteDialogOpen}
  title="Delete Course"
  message="Are you sure you want to delete '{deletingCourse?.name}'? This will also remove all enrollments, problem groups, and quiz groups associated with this course."
  isLoading={isSubmitting}
  on:confirm={handleDeleteConfirm}
  on:cancel={() => {
    isDeleteDialogOpen = false;
    deletingCourse = null;
  }}
/>
