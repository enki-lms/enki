<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import ConfirmDialog from "$lib/components/ConfirmDialog.svelte";
  import { api, type ExamSession, type ExamStudent } from "$lib/api";

  export let session: ExamSession;

  const dispatch = createEventDispatcher<{ back: void; updated: void }>();

  let students: ExamStudent[] = [];
  let isLoading = true;
  let error: string | null = null;
  let isActionLoading = false;

  let isDiscontinueDialogOpen = false;
  let discontinuingStudent: ExamStudent | null = null;

  const studentStatusColors: Record<string, string> = {
    active: "bg-green-100 text-green-800",
    submitted: "bg-blue-100 text-blue-800",
    discontinued: "bg-red-100 text-red-800",
  };

  const studentStatusLabels: Record<string, string> = {
    active: "🟢 Active",
    submitted: "✓ Submitted",
    discontinued: "✗ Discontinued",
  };

  async function fetchStudents() {
    try {
      isLoading = true;
      error = null;
      students = await api.getExamSessionStudents(session.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load students";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchStudents();
  });

  function handleBack() {
    dispatch("back");
  }

  async function handleStartSession() {
    isActionLoading = true;
    try {
      const updated = await api.startExamSession(session.id);
      session = updated;
      dispatch("updated");
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to start session";
    } finally {
      isActionLoading = false;
    }
  }

  async function handleEndSession() {
    isActionLoading = true;
    try {
      const updated = await api.endExamSession(session.id);
      session = updated;
      dispatch("updated");
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to end session";
    } finally {
      isActionLoading = false;
    }
  }

  function handleDiscontinueClick(student: ExamStudent) {
    discontinuingStudent = student;
    isDiscontinueDialogOpen = true;
  }

  async function handleDiscontinueConfirm() {
    if (!discontinuingStudent) return;
    isActionLoading = true;

    try {
      await api.discontinueStudent(session.id, discontinuingStudent.student_id);
      await fetchStudents();
      isDiscontinueDialogOpen = false;
      discontinuingStudent = null;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to discontinue student";
    } finally {
      isActionLoading = false;
    }
  }

  function formatDate(dateStr: string | undefined): string {
    if (!dateStr) return "-";
    const date = new Date(dateStr);
    return date.toLocaleString();
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
    <div class="flex-1">
      <h2 class="text-2xl font-semibold text-gray-900">
        Exam Session #{session.id}
      </h2>
      <p class="text-gray-500 text-sm">
        {session.problem_group_type === "comp_sci"
          ? "💻 CS Problems"
          : "📝 Quiz"} • {session.duration_minutes} minutes
      </p>
    </div>
    <div class="flex gap-2">
      {#if session.status === "pending"}
        <Button
          size="md"
          on:click={handleStartSession}
          disabled={isActionLoading}
        >
          {#if isActionLoading}
            <span
              class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
            ></span>
          {/if}
          ▶ Start Exam
        </Button>
      {:else if session.status === "active"}
        <button
          class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-full font-medium transition-colors disabled:opacity-50"
          on:click={handleEndSession}
          disabled={isActionLoading}
        >
          {#if isActionLoading}
            <span
              class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
            ></span>
          {/if}
          ⏹ End Exam
        </button>
      {/if}
    </div>
  </div>

  {#if error}
    <div
      class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
    >
      {error}
      <button class="ml-2 underline" on:click={fetchStudents}>Retry</button>
    </div>
  {/if}

  <div class="grid grid-cols-4 gap-4 mb-6 bg-gray-50 rounded-xl p-4">
    <div>
      <p class="text-sm text-gray-500">Status</p>
      <span
        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
				{session.status === 'pending' ? 'bg-yellow-100 text-yellow-800' : ''}
				{session.status === 'active' ? 'bg-green-100 text-green-800' : ''}
				{session.status === 'ended' ? 'bg-gray-100 text-gray-600' : ''}"
      >
        {session.status === "pending"
          ? "⏳ Pending"
          : session.status === "active"
            ? "🟢 Active"
            : "✓ Ended"}
      </span>
    </div>
    <div>
      <p class="text-sm text-gray-500">Created</p>
      <p class="font-medium text-gray-900 text-sm">
        {formatDate(session.created_at)}
      </p>
    </div>
    <div>
      <p class="text-sm text-gray-500">Started</p>
      <p class="font-medium text-gray-900 text-sm">
        {formatDate(session.started_at)}
      </p>
    </div>
    <div>
      <p class="text-sm text-gray-500">Ended</p>
      <p class="font-medium text-gray-900 text-sm">
        {formatDate(session.ended_at)}
      </p>
    </div>
  </div>

  <div class="bg-gray-50 rounded-xl p-6">
    <h3 class="text-lg font-semibold text-gray-800 mb-4">
      Students ({students.length})
    </h3>

    {#if isLoading}
      <div class="py-8 text-center text-gray-500">
        <div
          class="inline-block w-6 h-6 border-3 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-2 text-sm">Loading...</p>
      </div>
    {:else if students.length === 0}
      <div class="py-8 text-center text-gray-500">
        <p>No students in this session</p>
      </div>
    {:else}
      <div class="space-y-2">
        {#each students as student (student.id)}
          <div
            class="flex items-center justify-between bg-white rounded-lg p-4 border border-gray-200"
          >
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold text-sm"
              >
                {getInitials(student.student?.full_name || "?")}
              </div>
              <div>
                <p class="font-medium text-gray-900">
                  {student.student?.full_name ||
                    `Student ${student.student_id}`}
                </p>
                <p class="text-sm text-gray-500">
                  {student.student?.email || ""}
                </p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {studentStatusColors[
                  student.status
                ]}"
              >
                {studentStatusLabels[student.status] || student.status}
              </span>
              {#if session.status === "active" && student.status === "active"}
                <button
                  class="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                  on:click={() => handleDiscontinueClick(student)}
                  title="Discontinue student"
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
                      d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                    />
                  </svg>
                </button>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<ConfirmDialog
  bind:isOpen={isDiscontinueDialogOpen}
  title="Discontinue Student"
  message="Are you sure you want to discontinue '{discontinuingStudent?.student
    ?.full_name ||
    'this student'}' from the exam? They will not be able to rejoin."
  confirmText="Discontinue"
  isLoading={isActionLoading}
  on:confirm={handleDiscontinueConfirm}
  on:cancel={() => {
    isDiscontinueDialogOpen = false;
    discontinuingStudent = null;
  }}
/>
