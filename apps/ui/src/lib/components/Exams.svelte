<script lang="ts">
  import { onMount } from "svelte";
  import Button from "$lib/components/Button.svelte";
  import PlusIcon from "$lib/components/icons/PlusIcon.svelte";
  import ExamSessionDetail from "$lib/components/ExamSessionDetail.svelte";
  import ExamSessionModal from "$lib/components/ExamSessionModal.svelte";
  import { api, type ExamSession } from "$lib/api";

  let sessions: ExamSession[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Modal states
  let isCreateModalOpen = false;
  let isSubmitting = false;

  // Detail view
  let selectedSession: ExamSession | null = null;

  const statusColors: Record<string, string> = {
    pending: "bg-yellow-100 text-yellow-800",
    active: "bg-green-100 text-green-800",
    ended: "bg-gray-100 text-gray-600",
  };

  const statusLabels: Record<string, string> = {
    pending: "⏳ Pending",
    active: "🟢 Active",
    ended: "✓ Ended",
  };

  async function fetchSessions() {
    try {
      isLoading = true;
      error = null;
      sessions = await api.getExamSessions();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load exam sessions";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchSessions();
  });

  function handleNewSession() {
    isCreateModalOpen = true;
  }

  function handleSessionClick(session: ExamSession) {
    selectedSession = session;
  }

  function handleBack() {
    selectedSession = null;
    fetchSessions();
  }

  function handleCreated() {
    isCreateModalOpen = false;
    fetchSessions();
  }

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }
</script>

{#if selectedSession}
  <ExamSessionDetail
    session={selectedSession}
    on:back={handleBack}
    on:updated={fetchSessions}
  />
{:else}
  <div class="w-full">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-semibold text-gray-900">Exam Sessions</h2>
      <Button size="md" on:click={handleNewSession}>
        <span slot="icon" class="text-white">
          <PlusIcon />
        </span>
        New Exam Session
      </Button>
    </div>

    {#if error}
      <div
        class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700"
      >
        {error}
        <button class="ml-2 underline" on:click={fetchSessions}>Retry</button>
      </div>
    {/if}

    {#if isLoading}
      <div class="py-16 text-center text-gray-500">
        <div
          class="inline-block w-8 h-8 border-4 border-gray-300 border-t-sky-500 rounded-full animate-spin"
        ></div>
        <p class="mt-4">Loading exam sessions...</p>
      </div>
    {:else if sessions.length === 0}
      <div class="py-16 text-center text-gray-500">
        <p class="text-lg font-medium">No exam sessions</p>
        <p class="text-sm mt-1">Create an exam session to start an exam</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each sessions as session (session.id)}
          <button
            class="w-full bg-white border border-gray-200 rounded-xl p-5 hover:border-sky-300 hover:shadow-sm transition-all text-left"
            on:click={() => handleSessionClick(session)}
          >
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-3">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {statusColors[
                    session.status
                  ]}"
                >
                  {statusLabels[session.status] || session.status}
                </span>
                <span class="text-sm text-gray-500">
                  {session.problem_group_type === "comp_sci"
                    ? "💻 CS Problems"
                    : "📝 Quiz"}
                </span>
              </div>
              <span class="text-sm text-gray-500">ID: {session.id}</span>
            </div>
            <div class="grid grid-cols-3 gap-4 text-sm">
              <div>
                <p class="text-gray-500">Duration</p>
                <p class="font-medium text-gray-900">
                  {session.duration_minutes} minutes
                </p>
              </div>
              <div>
                <p class="text-gray-500">Created</p>
                <p class="font-medium text-gray-900">
                  {formatDate(session.created_at)}
                </p>
              </div>
              <div>
                {#if session.started_at}
                  <p class="text-gray-500">Started</p>
                  <p class="font-medium text-gray-900">
                    {formatDate(session.started_at)}
                  </p>
                {:else if session.status === "pending"}
                  <p class="text-gray-500">Status</p>
                  <p class="font-medium text-yellow-600">Waiting to start</p>
                {/if}
              </div>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
{/if}

<ExamSessionModal
  bind:isOpen={isCreateModalOpen}
  on:created={handleCreated}
  on:cancel={() => (isCreateModalOpen = false)}
/>
