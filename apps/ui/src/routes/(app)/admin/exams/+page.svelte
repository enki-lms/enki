<script lang="ts">
  import { onMount } from "svelte";
  import { api, type ExamSession } from "$lib/api";
  import Button from "$lib/components/Button.svelte";
  import Card from "$lib/components/Card.svelte";
  import { formatDate } from "$lib/utils/date";

  let sessions: ExamSession[] = [];
  let isLoading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      sessions = await api.getExamSessions();
    } catch (e) {
      error = (e as Error).message;
    } finally {
      isLoading = false;
    }
  });

  async function handleStartSession(id: number) {
    if (
      !confirm(
        "Are you sure you want to start this exam session? Students will be able to join.",
      )
    )
      return;
    try {
      await api.startExamSession(id);
      sessions = await api.getExamSessions(); // Refresh
    } catch (e) {
      alert((e as Error).message);
    }
  }

  async function handleEndSession(id: number) {
    if (
      !confirm(
        "Are you sure you want to end this exam session? This will forcibly submit all active students.",
      )
    )
      return;
    try {
      await api.endExamSession(id);
      sessions = await api.getExamSessions(); // Refresh
    } catch (e) {
      alert((e as Error).message);
    }
  }
</script>

<div class="min-h-screen bg-gray-50 p-8">
  <div class="max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Exam Sessions</h1>
        <p class="text-gray-600 mt-1">
          Manage your exam sessions and invigilation.
        </p>
      </div>
      <Button href="/admin/exams/new" variant="primary">Create New Exam</Button>
    </div>

    {#if isLoading}
      <div class="text-center py-12">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"
        ></div>
      </div>
    {:else if error}
      <div class="bg-red-50 text-red-700 p-4 rounded-lg shadow-sm">
        {error}
      </div>
    {:else if sessions.length === 0}
      <Card>
        <div class="text-center py-16">
          <div class="mx-auto h-12 w-12 text-gray-400">
            <!-- Icon placeholder -->
            <svg
              class="h-12 w-12"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
              />
            </svg>
          </div>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No exams</h3>
          <p class="mt-1 text-sm text-gray-500 mb-6">
            Get started by creating a new exam session.
          </p>
          <Button href="/admin/exams/new" variant="secondary"
            >Create your first exam</Button
          >
        </div>
      </Card>
    {:else}
      <div
        class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden"
      >
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >ID</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >Type</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >Group ID</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >Status</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >Created At</th
              >
              <th
                class="px-6 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >Actions</th
              >
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each sessions as session}
              <tr class="hover:bg-gray-50 transition-colors">
                <td
                  class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900"
                  >#{session.id}</td
                >
                <td
                  class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize"
                  >{session.problem_group_type === "comp_sci"
                    ? "Computer Science"
                    : "Quiz"}</td
                >
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                  >{session.problem_group_id}</td
                >
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class={`px-2.5 py-0.5 inline-flex text-xs leading-5 font-medium rounded-full 
                    ${
                      session.status === "active"
                        ? "bg-green-100 text-green-800"
                        : session.status === "ended"
                          ? "bg-gray-100 text-gray-800"
                          : "bg-yellow-100 text-yellow-800"
                    }`}
                  >
                    {session.status.charAt(0).toUpperCase() +
                      session.status.slice(1)}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                  >{formatDate(session.created_at)}</td
                >
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-3"
                >
                  {#if session.status === "pending"}
                    <button
                      class="text-indigo-600 hover:text-indigo-900 font-medium transition-colors"
                      on:click={() => handleStartSession(session.id)}
                      >Start</button
                    >
                  {/if}
                  {#if session.status === "active"}
                    <button
                      class="text-red-600 hover:text-red-900 font-medium transition-colors"
                      on:click={() => handleEndSession(session.id)}>End</button
                    >
                    <!-- Future: <a href={`/admin/exams/${session.id}/monitor`} class="text-blue-600 hover:text-blue-900">Monitor</a> -->
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>
