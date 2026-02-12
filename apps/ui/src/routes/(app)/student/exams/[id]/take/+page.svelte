<script lang="ts">
  import { page } from "$app/stores";
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { api, type Problem, type QuizProblem } from "$lib/api";
  import Button from "$lib/components/Button.svelte";
  import Card from "$lib/components/Card.svelte";
  import CodeEditor from "$lib/components/CodeEditor.svelte";
  import { browser } from "$app/environment";

  const sessionId = $page.params.id;
  let socket: WebSocket;
  let isConnected = false;
  let error: string | null = null;
  let status: string = "connecting"; // connecting, active, submitted, ended

  let timeRemaining = 0; // in seconds
  let endTime: number | null = null; // unix timestamp
  let timerInterval: any;

  // Exam Data
  let problemGroupType: "comp_sci" | "quiz" | null = null;
  let problemGroupId: number | null = null;

  let csProblems: Problem[] = [];
  let quizProblems: QuizProblem[] = [];

  let currentProblemIndex = 0;

  // Student Answers (map problemId -> answer)
  let answers: Record<number, string> = {}; // For Quiz: JSON string or value? For CS: Code
  let autoSaveTimeout: any;

  function formatTime(seconds: number): string {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h.toString().padStart(2, "0")}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
  }

  function getCookie(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
  }

  onMount(() => {
    if (!browser) return;

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const host = window.location.host;
    const token = getCookie("token");
    const wsUrl = `${protocol}//${host}/ws/exam/${sessionId}?token=${token}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      isConnected = true;
      status = "active";
    };

    socket.onmessage = async (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "joined") {
        const data = msg.payload;
        problemGroupType = data.problem_group_type;
        problemGroupId = data.problem_group_id;
        endTime = data.ends_at;
        timeRemaining = data.remaining_seconds;

        startTimer();
        loadProblems();
      } else if (msg.type === "timer_sync") {
        timeRemaining = msg.payload.remaining_seconds;
      } else if (msg.type === "end_session") {
        // Teacher ended or time expired handled by server closing?
        // Handle explicit end message if any (server sends "teacher_ended" reason on close)
        status = "ended";
        alert("The exam session has ended.");
        goto("/student/exams");
      }
    };

    socket.onclose = (event) => {
      isConnected = false;
      if (status !== "submitted" && status !== "ended") {
        if (event.code !== 1000) {
          // Normal closure
          error = "Connection lost. Reconnecting...";
          // Implement reconnection logic here if needed
        }
      }
    };

    socket.onerror = (e) => {
      console.error("WebSocket error", e);
      error = "Connection error";
    };
  });

  onDestroy(() => {
    if (socket) socket.close();
    clearInterval(timerInterval);
    clearTimeout(autoSaveTimeout);
  });

  function startTimer() {
    if (timerInterval) clearInterval(timerInterval);
    timerInterval = setInterval(() => {
      if (timeRemaining > 0) {
        timeRemaining--;
      } else {
        clearInterval(timerInterval);
        status = "ended";
        alert("Time is up! Submitting exam...");
        handleSubmitExam(); // Auto-submit on frontend too?
      }
    }, 1000);
  }

  async function loadProblems() {
    try {
      if (problemGroupType === "comp_sci" && problemGroupId) {
        csProblems = await api.getGroupProblems(problemGroupId);
      } else if (problemGroupType === "quiz" && problemGroupId) {
        quizProblems = await api.getQuizGroupProblems(problemGroupId);
      }
    } catch (e) {
      error = "Failed to load problems";
    }
  }

  function handleAnswerChange(problemId: number, value: string) {
    answers[problemId] = value;
    // Debounce auto-save
    if (autoSaveTimeout) clearTimeout(autoSaveTimeout);
    autoSaveTimeout = setTimeout(() => {
      saveProgress(problemId, value);
    }, 2000);
  }

  function saveProgress(problemId: number, value: string) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: "save_progress",
          payload: {
            problem_id: problemId,
            problem_type: problemGroupType,
            current_answer: value,
          },
        }),
      );
    }
  }

  async function handleSubmitExam() {
    if (
      !confirm(
        "Are you sure you want to finish the exam? You cannot change your answers after submitting.",
      )
    )
      return;

    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: "submit_exam",
          payload: {},
        }),
      );
      status = "submitted";
      goto("/student/exams");
    }
  }

  function nextProblem() {
    const max =
      problemGroupType === "comp_sci" ? csProblems.length : quizProblems.length;
    if (currentProblemIndex < max - 1) currentProblemIndex++;
  }

  function prevProblem() {
    if (currentProblemIndex > 0) currentProblemIndex--;
  }
</script>

<div class="min-h-screen bg-gray-50 flex flex-col">
  <!-- Header -->
  <header
    class="bg-white shadow-sm border-b px-6 py-4 flex justify-between items-center sticky top-0 z-10"
  >
    <div>
      <h1 class="text-xl font-bold text-gray-900">
        {problemGroupType === "comp_sci"
          ? "Computer Science Exam"
          : "Quiz Exam"}
      </h1>
      <p class="text-sm text-gray-500">Session #{sessionId}</p>
    </div>

    <div class="flex items-center gap-6">
      <div class="text-center">
        <p class="text-xs text-gray-500 uppercase tracking-wide">
          Time Remaining
        </p>
        <p
          class={`text-2xl font-mono font-bold ${timeRemaining < 300 ? "text-red-600 animate-pulse" : "text-gray-900"}`}
        >
          {formatTime(timeRemaining)}
        </p>
      </div>
      <Button on:click={handleSubmitExam} variant="primary">Submit Exam</Button>
    </div>
  </header>

  <div class="flex-1 flex overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-64 bg-white border-r overflow-y-auto">
      <div class="p-4 border-b">
        <h3 class="font-medium text-gray-700">Problems</h3>
      </div>
      <div class="p-2 space-y-1">
        {#if problemGroupType === "comp_sci"}
          {#each csProblems as problem, i}
            <button
              class={`w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors
                         ${i === currentProblemIndex ? "bg-indigo-50 text-indigo-700" : "text-gray-600 hover:bg-gray-50"}`}
              on:click={() => (currentProblemIndex = i)}
            >
              {i + 1}. {problem.name}
              {#if answers[problem.id]}
                <span class="float-right text-green-500">✓</span>
              {/if}
            </button>
          {/each}
        {:else if problemGroupType === "quiz"}
          {#each quizProblems as problem, i}
            <button
              class={`w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors
                         ${i === currentProblemIndex ? "bg-indigo-50 text-indigo-700" : "text-gray-600 hover:bg-gray-50"}`}
              on:click={() => (currentProblemIndex = i)}
            >
              {i + 1}. {problem.name}
              {#if answers[problem.id]}
                <span class="float-right text-green-500">✓</span>
              {/if}
            </button>
          {/each}
        {/if}
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto p-8">
      {#if error}
        <div class="bg-red-50 text-red-700 p-4 rounded-lg mb-4">{error}</div>
      {/if}

      {#if problemGroupType === "comp_sci" && csProblems.length > 0}
        {@const problem = csProblems[currentProblemIndex]}
        <div class="max-w-4xl mx-auto space-y-6">
          <div>
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              {problem.name}
            </h2>
            <div
              class="prose max-w-none bg-white p-6 rounded-lg shadow-sm border"
            >
              {@html problem.problem_text}
              <!-- Assumption: problem_text is HTML or markdown rendered -->
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm border overflow-hidden">
            <div
              class="bg-gray-50 px-4 py-2 border-b flex justify-between items-center"
            >
              <span class="font-medium text-gray-700">Code Editor</span>
              <span class="text-xs text-gray-500"
                >Changes saved automatically</span
              >
            </div>
            <div class="h-[500px]">
              <CodeEditor
                value={answers[problem.id] || ""}
                on:change={(e) => handleAnswerChange(problem.id, e.detail)}
                language="python"
              />
            </div>
          </div>
        </div>
      {:else if problemGroupType === "quiz" && quizProblems.length > 0}
        {@const problem = quizProblems[currentProblemIndex]}
        <div class="max-w-3xl mx-auto space-y-8">
          <div>
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              {problem.name}
            </h2>
            <div
              class="prose max-w-none bg-white p-6 rounded-lg shadow-sm border"
            >
              <p>{problem.problem_text}</p>
            </div>
          </div>

          <Card>
            <div class="space-y-4">
              <h3 class="font-medium text-gray-900">Your Answer</h3>

              {#if problem.problem_type === "open_ended"}
                <textarea
                  class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                  rows="6"
                  value={answers[problem.id] || ""}
                  on:input={(e) =>
                    handleAnswerChange(problem.id, e.currentTarget.value)}
                  placeholder="Type your answer here..."
                ></textarea>
              {:else if problem.problem_type === "true_false"}
                <div class="space-y-2">
                  <label class="flex items-center space-x-3">
                    <input
                      type="radio"
                      name={`problem-${problem.id}`}
                      value="true"
                      checked={answers[problem.id] === "true"}
                      on:change={() => handleAnswerChange(problem.id, "true")}
                      class="h-4 w-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                    />
                    <span class="text-gray-900">True</span>
                  </label>
                  <label class="flex items-center space-x-3">
                    <input
                      type="radio"
                      name={`problem-${problem.id}`}
                      value="false"
                      checked={answers[problem.id] === "false"}
                      on:change={() => handleAnswerChange(problem.id, "false")}
                      class="h-4 w-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                    />
                    <span class="text-gray-900">False</span>
                  </label>
                </div>
              {:else if problem.problem_type === "mcq_single"}
                <div class="space-y-2">
                  {#each problem.options || [] as option}
                    <label class="flex items-center space-x-3">
                      <input
                        type="radio"
                        name={`problem-${problem.id}`}
                        value={option.id.toString()}
                        checked={answers[problem.id] === option.id.toString()}
                        on:change={() =>
                          handleAnswerChange(problem.id, option.id.toString())}
                        class="h-4 w-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                      />
                      <span class="text-gray-900">{option.option_text}</span>
                    </label>
                  {/each}
                </div>
              {/if}
            </div>
          </Card>
        </div>
      {/if}

      <!-- Navigation Buttons -->
      <div class="max-w-4xl mx-auto mt-8 flex justify-between">
        <Button
          variant="secondary"
          on:click={prevProblem}
          disabled={currentProblemIndex === 0}
        >
          Previous
        </Button>
        <Button
          variant="secondary"
          on:click={nextProblem}
          disabled={problemGroupType === "comp_sci"
            ? currentProblemIndex === csProblems.length - 1
            : currentProblemIndex === quizProblems.length - 1}
        >
          Next
        </Button>
      </div>
    </main>
  </div>
</div>

<style>
  /* Hide standard layout headers for exam mode if possible */
  /* This requires layout changes, but for now we just overlay */
</style>
