<script lang="ts">
  import { onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import PlusIcon from "$lib/components/icons/PlusIcon.svelte";
  import EmptyDocumentIcon from "$lib/components/icons/EmptyDocumentIcon.svelte";
  import ItemCard from "$lib/components/ItemCard.svelte";
  import FormModal, {
    type FieldConfig,
  } from "$lib/components/FormModal.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import Select from "$lib/components/Select.svelte";
  import ProblemGroupDetail from "$lib/components/ProblemGroupDetail.svelte";
  import { api, type Course, type ProblemGroup } from "$lib/api";

  let courses: Course[] = [];
  let selectedCourseId: number | "" = "";
  let problemGroups: ProblemGroup[] = [];
  let isLoading = true;
  let isLoadingGroups = false;
  let error: string | null = null;

  // Modal states
  let isFormModalOpen = false;
  let isDeleteDialogOpen = false;
  let isSubmitting = false;
  let editingGroup: ProblemGroup | null = null;
  let deletingGroup: ProblemGroup | null = null;

  // Detail view
  let selectedGroup: ProblemGroup | null = null;

  const formFields: FieldConfig[] = [
    {
      name: "name",
      label: "Name",
      type: "text",
      placeholder: "Enter problem group name...",
      required: true,
    },
    {
      name: "description",
      label: "Description",
      type: "textarea",
      placeholder: "Describe this problem group...",
      required: false,
    },
    {
      name: "type",
      label: "Type",
      type: "select",
      required: true,
      options: [
        { value: "practice", label: "Practice" },
        { value: "exam", label: "Exam" },
      ],
    },
  ];

  async function fetchCourses() {
    try {
      isLoading = true;
      error = null;
      courses = await api.getCourses();
      if (courses.length > 0 && !selectedCourseId) {
        selectedCourseId = courses[0].id;
      }
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load courses";
    } finally {
      isLoading = false;
    }
  }

  async function fetchProblemGroups() {
    if (!selectedCourseId) {
      problemGroups = [];
      return;
    }
    try {
      isLoadingGroups = true;
      error = null;
      problemGroups = await api.getCourseProblemGroups(
        selectedCourseId as number,
      );
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load problem groups";
    } finally {
      isLoadingGroups = false;
    }
  }

  onMount(() => {
    fetchCourses();
  });

  $: if (selectedCourseId) {
    fetchProblemGroups();
  }

  function handleNewGroup() {
    editingGroup = null;
    isFormModalOpen = true;
  }

  function handleEditGroup(group: ProblemGroup) {
    editingGroup = group;
    isFormModalOpen = true;
  }

  function handleDeleteClick(group: ProblemGroup) {
    deletingGroup = group;
    isDeleteDialogOpen = true;
  }

  async function handleFormSubmit(
    event: CustomEvent<Record<string, string | number>>,
  ) {
    const { name, description, type } = event.detail;
    isSubmitting = true;

    try {
      const input = {
        name: String(name),
        description: String(description || ""),
        type: type as "exam" | "practice",
      };

      if (editingGroup) {
        const updated = await api.updateProblemGroup(editingGroup.id, input);
        problemGroups = problemGroups.map((g) =>
          g.id === updated.id ? updated : g,
        );
      } else {
        const created = await api.createProblemGroup(
          selectedCourseId as number,
          input,
        );
        problemGroups = [created, ...problemGroups];
      }
      isFormModalOpen = false;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to save problem group";
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteConfirm() {
    if (!deletingGroup) return;
    isSubmitting = true;

    try {
      await api.deleteProblemGroup(deletingGroup.id);
      problemGroups = problemGroups.filter((g) => g.id !== deletingGroup!.id);
      isDeleteDialogOpen = false;
      deletingGroup = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to delete problem group";
    } finally {
      isSubmitting = false;
    }
  }

  function handleCardClick(group: ProblemGroup) {
    selectedGroup = group;
  }

  function handleBack() {
    selectedGroup = null;
    fetchProblemGroups();
  }
</script>

{#if selectedGroup}
  <ProblemGroupDetail group={selectedGroup} on:back={handleBack} />
{:else}
  <div class="w-full">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-semibold text-gray-900">Problem Sets</h2>
      <Button size="md" on:click={handleNewGroup} disabled={!selectedCourseId}>
        <span slot="icon" class="text-white">
          <PlusIcon />
        </span>
        New Problem Set
      </Button>
    </div>

    <!-- Course Selector -->
    <div class="mb-6">
      <label
        for="course-select"
        class="block text-sm font-medium text-gray-700 mb-2"
        >Select Course</label
      >
      {#if isLoading}
        <div class="h-12 bg-gray-100 rounded-lg animate-pulse"></div>
      {:else}
        <Select
          id="course-select"
          bind:value={selectedCourseId}
          placeholder="Select a course..."
          options={courses.map((c) => ({ value: c.id, label: c.name }))}
        />
      {/if}
    </div>

    {#if error}
      <div
        class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
      >
        {error}
        <button class="ml-2 underline" on:click={fetchProblemGroups}
          >Retry</button
        >
      </div>
    {/if}

    {#if !selectedCourseId}
      <div class="py-16 text-center text-gray-500">
        <p class="text-lg font-medium">Select a course first</p>
        <p class="text-sm mt-1">Choose a course to view its problem sets</p>
      </div>
    {:else if isLoadingGroups}
      <div class="py-16 text-center text-gray-500">
        <div
          class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-4">Loading problem sets...</p>
      </div>
    {:else if problemGroups.length === 0}
      <div class="py-16 text-center text-gray-500">
        <EmptyDocumentIcon />
        <p class="text-lg font-medium">No problem sets yet</p>
        <p class="text-sm mt-1">Create your first problem set to get started</p>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each problemGroups as group (group.id)}
          <div class="relative group/card">
            <ItemCard
              title={group.name}
              description={group.description}
              date={group.type === "exam" ? "📝 Exam" : "📚 Practice"}
              type="problem-set"
              on:click={() => handleCardClick(group)}
            />
            <div
              class="absolute top-2 right-2 opacity-0 group-hover/card:opacity-100 transition-opacity flex gap-1"
            >
              <button
                class="p-2 bg-white rounded-lg shadow-md hover:bg-gray-50 text-gray-600"
                on:click|stopPropagation={() => handleEditGroup(group)}
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
                on:click|stopPropagation={() => handleDeleteClick(group)}
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
  title={editingGroup ? "Edit Problem Set" : "Create New Problem Set"}
  submitText={editingGroup ? "Save" : "Create"}
  fields={formFields}
  initialValues={editingGroup
    ? {
        name: editingGroup.name,
        description: editingGroup.description,
        type: editingGroup.type,
      }
    : { type: "practice" }}
  isLoading={isSubmitting}
  on:submit={handleFormSubmit}
  on:cancel={() => (isFormModalOpen = false)}
/>

<ConfirmDialog
  bind:isOpen={isDeleteDialogOpen}
  title="Delete Problem Set"
  message="Are you sure you want to delete '{deletingGroup?.name}'? This will also delete all problems and test cases in this group."
  isLoading={isSubmitting}
  on:confirm={handleDeleteConfirm}
  on:cancel={() => {
    isDeleteDialogOpen = false;
    deletingGroup = null;
  }}
/>
