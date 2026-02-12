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
      sessions = await api.getActiveExamSessions();
    } catch (e) {
      error = (e as Error).message;
    } finally {
      isLoading = false;
    }
  });
</script>

<div class="min-h-screen bg-gray-50 p-8">
  <div class="max-w-7xl mx-auto space-y-8">
    <div>
      <h1 class="text-3xl font-bold text-gray-900">My Exams</h1>
      <p class="text-lg text-gray-600 mt-2">
        View and join your scheduled exams.
      </p>
    </div>

    {#if isLoading}
      <div class="text-center py-20">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"
        ></div>
        <p class="mt-4 text-gray-500">Loading your exams...</p>
      </div>
    {:else if error}
      <div class="bg-red-50 text-red-700 p-6 rounded-xl border border-red-100">
        {error}
      </div>
    {:else if sessions.length === 0}
      <Card>
        <div class="text-center py-20">
          <div class="mx-auto h-16 w-16 text-gray-300 mb-4">
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-900">No active exams</h3>
          <p class="text-gray-500 mt-2">
            You don't have any exams scheduled at the moment.
          </p>
        </div>
      </Card>
    {:else}
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {#each sessions as session}
          <div
            class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 flex flex-col justify-between h-full hover:shadow-md transition-shadow"
          >
            <div class="space-y-4">
              <div class="flex justify-between items-start">
                <span
                  class={`px-3 py-1 text-xs font-semibold rounded-full 
                    ${session.status === "active" ? "bg-green-100 text-green-700" : "bg-yellow-100 text-yellow-700"}`}
                >
                  {session.status.toUpperCase()}
                </span>
                <span class="text-xs text-mono text-gray-400"
                  >#{session.id}</span
                >
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 mb-1">
                  {session.problem_group_type === "comp_sci"
                    ? "Computer Science"
                    : "Quiz"}
                </h3>
                <p class="text-sm text-gray-500 flex items-center gap-2">
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  {session.duration_minutes} minutes
                </p>
              </div>
            </div>

            <div class="mt-6 pt-6 border-t border-gray-100">
              {#if session.status === "active"}
                <Button
                  href={`/student/exams/${session.id}/take`}
                  variant="primary"
                  class="w-full justify-center"
                >
                  Join Exam
                </Button>
              {:else}
                <Button disabled class="w-full justify-center opacity-70"
                  >Waiting to Start</Button
                >
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
