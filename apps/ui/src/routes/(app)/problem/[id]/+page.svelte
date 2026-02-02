<script lang="ts">
  import { page } from "$app/stores";
  import BackButton from "$lib/components/BackButton.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import Tabs from "$lib/components/Tabs.svelte";
  import LanguageSelector from "$lib/components/LanguageSelector.svelte";
  import Chat from "$lib/components/Chat.svelte";
  import CodeEditor from "$lib/components/CodeEditor.svelte";

  import { onMount } from "svelte";
  import { api, type Problem } from "$lib/api";

  // Get problem ID from URL
  const problemId = $page.params.id ?? "";

  let activeDescTab = "Description";
  const descTabs = ["Description", "Scratchpad", "Assistant"];

  let activeSolutionTab = "Solution";
  const solutionTabs = ["Solution"];

  let activeTestTab = "Test Result";
  const testTabs = ["Test Result"];

  let solutionText = "";
  let testResultText = "";
  let scratchpadText = "";
  let messages: Array<{
    id: number;
    text: string;
    sender: "user" | "assistant";
  }> = [];

  let problem: Problem | null = null;
  let error: string | null = null;
  let loading = true;

  onMount(async () => {
    try {
      problem = await api.getProblem(problemId);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load problem";
    } finally {
      loading = false;
    }
  });

  async function handleSubmit(code: string) {
    console.log("Submitted code for problem", problemId, ":", code);
    testResultText = "Submitting solution...";

    try {
      const response = await api.submitSolution(problemId, code);
      const result = response.result;

      let output = `Score: ${result.score}/${result.max_score} (${result.passed}/${result.total_test_cases} passed)\n\n`;

      result.results.forEach((tc, i) => {
        output += `Test Case ${i + 1}: ${tc.passed ? "PASSED" : "FAILED"}\n`;
        if (!tc.passed) {
          if (tc.error) output += `Error: ${tc.error}\n`;
          if (tc.expected) output += `Expected: ${tc.expected}\n`;
          if (tc.actual) output += `Actual:   ${tc.actual}\n`;
        }
        output += "\n";
      });

      testResultText = output;
      // Switch to test result tab
      activeTestTab = "Test Result";
    } catch (e) {
      testResultText = `Error submitting solution: ${e instanceof Error ? e.message : "Unknown error"}`;
    }
  }
</script>

<div class="h-screen bg-[#E8EEF2] flex flex-col">
  <!-- Header -->
  <div
    class="bg-white border-b border-gray-300 px-4 md:px-6 py-3 md:py-4 flex items-center gap-2 md:gap-4 flex-wrap md:flex-nowrap"
  >
    <BackButton onclick={() => history.back()} />
    <div class="hidden md:block">
      <LogoPlaceHolder />
    </div>
    <div class="flex-1 md:flex-none">
      <LanguageSelector />
    </div>
  </div>

  <!-- Main Content -->
  <div
    class="flex flex-1 gap-3 md:gap-6 p-3 md:p-6 overflow-hidden flex-col lg:flex-row"
  >
    <!-- Left Panel - Problem Description -->
    <div
      class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white"
    >
      <Tabs tabs={descTabs} bind:active={activeDescTab} />
      <div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
        {#if activeDescTab === "Description"}
          <div class="overflow-y-auto h-full p-3 md:p-6">
            {#if loading}
              <div class="flex items-center justify-center h-full">
                <p class="text-gray-500">Loading problem...</p>
              </div>
            {:else if error}
              <div class="flex items-center justify-center h-full">
                <p class="text-red-500">{error}</p>
              </div>
            {:else if problem}
              <h2 class="text-xl md:text-2xl font-bold mb-4">
                {problem.name}
              </h2>
              <div
                class="mb-4 text-sm md:text-base text-gray-700 whitespace-pre-wrap"
              >
                {problem.problem_text}
              </div>

              {#if problem.time_limit_ms || problem.memory_limit_mb}
                <p
                  class="mb-3 font-semibold text-sm md:text-base text-gray-800"
                >
                  Constraints:
                </p>
                <ul
                  class="list-disc list-inside space-y-2 ml-2 text-sm md:text-base text-gray-700"
                >
                  {#if problem.time_limit_ms}
                    <li>Time Limit: {problem.time_limit_ms}ms</li>
                  {/if}
                  {#if problem.memory_limit_mb}
                    <li>Memory Limit: {problem.memory_limit_mb}MB</li>
                  {/if}
                </ul>
              {/if}
            {:else}
              <div class="flex items-center justify-center h-full">
                <p class="text-gray-500">Problem not found</p>
              </div>
            {/if}
          </div>
        {:else if activeDescTab === "Scratchpad"}
          <textarea
            bind:value={scratchpadText}
            class="w-full h-full p-3 md:p-4 border-none resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA]"
            placeholder="Write your scratchpad notes here..."
          ></textarea>
        {:else if activeDescTab === "Assistant"}
          <Chat
            bind:messages
            additionalContext={`Problem: ${problem?.name}\n\nDescription:\n${problem?.problem_text}\n\nCurrent Code:\n${solutionText}`}
          />
        {/if}
      </div>
    </div>

    <!-- Right Panel - Solution and Test Result as separate frames -->
    <div class="flex-1 lg:w-1/2 flex flex-col gap-3 md:gap-6 min-w-0">
      <!-- Solution Frame -->
      <div
        class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white"
      >
        <Tabs tabs={solutionTabs} bind:active={activeSolutionTab} />
        <div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
          <CodeEditor onSubmit={handleSubmit} bind:code={solutionText} />
        </div>
      </div>

      <!-- Test Result Frame -->
      <div
        class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white"
      >
        <Tabs tabs={testTabs} bind:active={activeTestTab} />
        <div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
          <textarea
            readonly
            bind:value={testResultText}
            class="w-full h-full p-3 md:p-6 resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA]"
            placeholder="Test results will appear here..."
          ></textarea>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  :global(.animate-fadeIn) {
    animation: fadeIn 0.2s ease-in-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @media (max-width: 1024px) {
    :global(.animate-fadeIn) {
      animation: fadeIn 0.2s ease-in-out;
    }
  }
</style>
