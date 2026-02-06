<script lang="ts">
  import { page } from "$app/stores";
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import SubmissionCard from "$lib/components/SubmissionCard.svelte";
  import Tabs from "$lib/components/Tabs.svelte";
  import { onMount } from "svelte";
  import {
    api,
    type CodeSubmission,
    type QuizSubmission,
    type Problem,
  } from "$lib/api";

  let {
    data,
  }: { data: { user?: { fullName?: string; email?: string; role?: string } } } =
    $props();

  let activeTab = $state("Code");
  const tabs = ["Code", "Quiz"];

  let codeSubmissions: CodeSubmission[] = $state([]);
  let quizSubmissions: QuizSubmission[] = $state([]);
  let problemNames: Map<number, string> = $state(new Map());
  let loading = $state(true);
  let error: string | null = $state(null);

  onMount(async () => {
    try {
      const [codeData, quizData] = await Promise.all([
        api.getAllCodeSubmissions(),
        api.getAllQuizSubmissions(),
      ]);
      codeSubmissions = codeData;
      quizSubmissions = quizData;

      // Fetch problem names for all submissions
      const problemIds = new Set([
        ...codeData.map((s) => s.problem_id),
        ...quizData.map((s) => s.problem_id),
      ]);

      // For now, we'll just use the problem_id as the name
      // A proper implementation would fetch problem names from the API
      problemIds.forEach((id) => {
        problemNames.set(id, `Problem #${id}`);
      });
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load submissions";
    } finally {
      loading = false;
    }
  });

  const currentSubmissions = $derived(
    activeTab === "Code" ? codeSubmissions : quizSubmissions,
  );

  const submissionType = $derived(activeTab === "Code" ? "code" : "quiz") as
    | "code"
    | "quiz";
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
              Submission History
            </h1>
            <p class="mt-1 text-gray-500">View all your past submissions</p>
          </div>

          <!-- Stats -->
          {#if !loading && !error}
            <div class="flex gap-4">
              <div
                class="flex items-center gap-2 bg-gradient-to-r from-sky-50 to-sky-100 rounded-xl px-4 py-2 border border-sky-200"
              >
                <svg
                  class="w-5 h-5 text-sky-600"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                  />
                </svg>
                <span class="font-bold text-sky-700"
                  >{codeSubmissions.length}</span
                >
                <span class="text-sm text-sky-600">Code</span>
              </div>
              <div
                class="flex items-center gap-2 bg-gradient-to-r from-purple-50 to-purple-100 rounded-xl px-4 py-2 border border-purple-200"
              >
                <svg
                  class="w-5 h-5 text-purple-600"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <span class="font-bold text-purple-700"
                  >{quizSubmissions.length}</span
                >
                <span class="text-sm text-purple-600">Quiz</span>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Tabs and Content -->
      <div
        class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden"
      >
        <Tabs {tabs} bind:active={activeTab} />

        {#if loading}
          <div class="flex justify-center p-12">
            <div class="flex flex-col items-center gap-3">
              <div
                class="w-8 h-8 border-3 border-sky-400 border-t-transparent rounded-full animate-spin"
              ></div>
              <p class="text-gray-500">Loading submissions...</p>
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
        {:else if currentSubmissions.length === 0}
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
              <p class="text-gray-500">
                No {activeTab.toLowerCase()} submissions yet.
              </p>
            </div>
          </div>
        {:else}
          <div class="p-4 space-y-3">
            {#each currentSubmissions as submission, index}
              <div
                style="animation: slideIn 0.3s ease-out {index *
                  0.05}s backwards"
              >
                <SubmissionCard
                  {submission}
                  type={submissionType}
                  problemName={problemNames.get(submission.problem_id) ??
                    `Problem #${submission.problem_id}`}
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
