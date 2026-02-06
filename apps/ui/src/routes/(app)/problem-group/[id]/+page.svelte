<script lang="ts">
  import { page } from "$app/stores";
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import ProblemSidebarItem from "$lib/components/ProblemSidebarItem.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import { onMount } from "svelte";
  import { api, type Problem } from "$lib/api";

  let {
    data,
  }: { data: { user?: { fullName?: string; email?: string; role?: string } } } =
    $props();

  // Get the problem group ID from the URL
  const groupId = $page.params.id ?? "";

  let problems: Problem[] = $state([]);
  let loading = $state(true);
  let error: string | null = $state(null);
  let groupTitle = $state("Problem List");

  // Mock completion data - in a real app, this would come from the API
  function isCompleted(problemId: number): boolean {
    // Simulate some problems being completed
    const completed = [1, 3]; // IDs of completed problems
    return completed.includes(problemId);
  }

  onMount(async () => {
    try {
      problems = await api.getGroupProblems(groupId);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load problems";
    } finally {
      loading = false;
    }
  });

  // Calculate overall progress
  const completedCount = $derived(
    problems.filter((p) => isCompleted(p.id)).length,
  )
  const overallProgress = $derived(
    problems.length > 0 ? (completedCount / problems.length) * 100 : 0,
  );
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex items-center gap-4">
    <div class="p-2">
      <BackButton onclick={() => history.back()} />
    </div>
    <div class="p-2">
      <ProfileMenu
        width="60"
        height="60"
        name={data.user?.fullName ?? ""}
        email={data.user?.email ?? ""}
        role={data.user?.role ?? ""}
      />
    </div>
    <div class="p-3"><LogoPlaceHolder /></div>
  </header>

  <main class="p-4 md:p-8">
    <div class="max-w-4xl mx-auto">
      <!-- Header Card -->
      <div
        class="bg-white rounded-2xl p-6 md:p-8 shadow-sm border border-gray-100 mb-6"
      >
        <div
          class="flex flex-col md:flex-row md:items-center md:justify-between gap-4"
        >
          <div>
            <h1 class="text-2xl md:text-3xl font-bold text-gray-900">
              {groupTitle}
            </h1>
            <p class="mt-1 text-gray-500">
              {problems.length} problem{problems.length !== 1 ? "s" : ""} in this
              group
            </p>
          </div>

          <!-- Overall Progress -->
          {#if !loading && !error && problems.length > 0}
            <div
              class="flex items-center gap-4 bg-gradient-to-r from-sky-50 to-emerald-50 rounded-xl px-5 py-3 border border-sky-100"
            >
              <div class="relative">
                <svg class="w-12 h-12 transform -rotate-90">
                  <circle
                    cx="24"
                    cy="24"
                    r="20"
                    stroke="currentColor"
                    stroke-width="4"
                    fill="none"
                    class="text-gray-200"
                  />
                  <circle
                    cx="24"
                    cy="24"
                    r="20"
                    stroke="currentColor"
                    stroke-width="4"
                    fill="none"
                    class="text-sky-500 transition-all duration-700"
                    stroke-dasharray={2 * Math.PI * 20}
                    stroke-dashoffset={2 *
                      Math.PI *
                      20 *
                      (1 - overallProgress / 100)}
                    stroke-linecap="round"
                  />
                </svg>
                <span
                  class="absolute inset-0 flex items-center justify-center text-xs font-bold text-gray-700"
                >
                  {Math.round(overallProgress)}%
                </span>
              </div>
              <div>
                <p class="text-sm font-semibold text-gray-700">
                  Overall Progress
                </p>
                <p class="text-xs text-gray-500">
                  {completedCount} of {problems.length} completed
                </p>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Problems Sidebar -->
      <div
        class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden"
      >
        {#if loading}
          <div class="flex justify-center p-12">
            <div class="flex flex-col items-center gap-3">
              <div
                class="w-8 h-8 border-3 border-sky-400 border-t-transparent rounded-full animate-spin"
              ></div>
              <p class="text-gray-500">Loading problems...</p>
            </div>
          </div>
        {:else if error}
          <div class="flex justify-center p-12">
            <div class="flex flex-col items-center gap-3 text-center">
              <div
                class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center"
              >
                <svg
                  class="w-6 h-6 text-red-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
              </div>
              <p class="text-red-500 font-medium">{error}</p>
            </div>
          </div>
        {:else if problems.length === 0}
          <div class="flex justify-center p-12">
            <div class="flex flex-col items-center gap-3 text-center">
              <div
                class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center"
              >
                <svg
                  class="w-6 h-6 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              </div>
              <p class="text-gray-500">No problems found in this group.</p>
            </div>
          </div>
        {:else}
          <div class="divide-y divide-gray-100">
            {#each problems as problem, index}
              <div
                class="p-2 first:pt-3 last:pb-3 transition-all duration-200"
                style="animation: slideIn 0.3s ease-out {index *
                  0.05}s backwards"
              >
                <ProblemSidebarItem
                  id={problem.id.toString()}
                  title={problem.name}
                  description={problem.description}
                  completed={isCompleted(problem.id)}
                />
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </main>
</div>

<style>
  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateX(-10px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
</style>
