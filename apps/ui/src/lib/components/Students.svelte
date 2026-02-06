<script lang="ts">
  import { onMount } from "svelte";
  import StudentCard from "./StudentCard.svelte";
  import { api, type Student } from "$lib/api";

  let students: Student[] = [];
  let isLoading = true;
  let error: string | null = null;

  async function fetchStudents(): Promise<void> {
    try {
      isLoading = true;
      error = null;
      students = await api.getSchoolStudents();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load students";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchStudents();
  });
</script>

<div class="w-full">
  <div class="mb-6">
    <h2 class="text-2xl font-bold text-gray-800">Students</h2>
    <p class="text-gray-500 text-sm mt-1">All students in your school</p>
  </div>

  {#if error}
    <div
      class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
    >
      {error}
      <button class="ml-2 underline" on:click={fetchStudents}>Retry</button>
    </div>
  {/if}

  {#if isLoading}
    <div class="py-16 text-center text-gray-500">
      <div
        class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
      ></div>
      <p class="mt-4">Loading students...</p>
    </div>
  {:else if students.length > 0}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each students as student (student.id)}
        <StudentCard
          student={{
            id: String(student.id),
            name: student.full_name,
            email: student.email,
          }}
        />
      {/each}
    </div>
  {:else}
    <div class="text-center py-12 text-gray-500">
      <p class="text-lg">No students found</p>
    </div>
  {/if}
</div>
