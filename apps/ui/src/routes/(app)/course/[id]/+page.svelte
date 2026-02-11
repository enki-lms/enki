<script lang="ts">
  import { page } from "$app/stores";
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    api,
    type Course,
    type ProblemGroup,
    type QuizGroup,
  } from "$lib/api";

  export let data;

  const courseId = parseInt($page.params.id);

  let course: Course | null = null;
  let problemGroups: ProblemGroup[] = [];
  let quizGroups: QuizGroup[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const [courseData, problemGroupsData, quizGroupsData] = await Promise.all(
        [
          api.getCourse(courseId),
          api.getCourseProblemGroups(courseId),
          api.getCourseQuizGroups(courseId),
        ],
      );
      course = courseData;
      problemGroups = problemGroupsData;
      quizGroups = quizGroupsData;
    } catch (e) {
      console.error("Error loading course data:", e);
      error = e instanceof Error ? e.message : "Failed to load course";
    } finally {
      loading = false;
    }
  });

  function navigateToProblemGroup(groupId: number) {
    goto(`/problem-group/${groupId}`);
  }

  function navigateToQuizGroup(groupId: number) {
    goto(`/quiz-group/${groupId}`);
  }
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex items-center gap-4">
    <div class="p-2">
      <BackButton onclick={() => goto("/home")} />
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
      {#if loading}
        <div class="bg-white rounded-2xl p-12 shadow-sm flex justify-center">
          <div class="flex flex-col items-center gap-3">
            <div
              class="w-8 h-8 border-3 border-sky-400 border-t-transparent rounded-full animate-spin"
            ></div>
            <p class="text-gray-500">Loading course...</p>
          </div>
        </div>
      {:else if error}
        <div class="bg-white rounded-2xl p-12 shadow-sm flex justify-center">
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
      {:else if course}
        <div
          class="bg-white rounded-2xl p-6 md:p-8 shadow-sm border border-gray-100 mb-6"
        >
          <div class="flex items-start gap-4">
            <div
              class="w-16 h-16 rounded-xl bg-gradient-to-br from-sky-400 to-emerald-400 flex items-center justify-center flex-shrink-0"
            >
              <svg
                class="w-8 h-8 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
            </div>
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-gray-900">
                {course.name}
              </h1>
              <p class="mt-1 text-gray-500">{course.subject}</p>
            </div>
          </div>
        </div>

        {#if problemGroups.length > 0}
          <div class="mb-6">
            <h2
              class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2"
            >
              <svg
                class="w-5 h-5 text-sky-500"
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
              Programming Exercises
            </h2>
            <div
              class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden divide-y divide-gray-100"
            >
              {#each problemGroups as group, index}
                <button
                  class="w-full p-5 text-left hover:bg-sky-50 transition-colors flex items-center justify-between group"
                  style="animation: slideIn 0.3s ease-out {index *
                    0.05}s backwards"
                  on:click={() => navigateToProblemGroup(group.id)}
                >
                  <div class="flex items-center gap-4">
                    <div
                      class="w-10 h-10 rounded-lg bg-sky-100 flex items-center justify-center text-sky-600 font-semibold"
                    >
                      {index + 1}
                    </div>
                    <div>
                      <h3 class="font-semibold text-gray-900">{group.name}</h3>
                      {#if group.description}
                        <p class="text-sm text-gray-500 line-clamp-1">
                          {group.description}
                        </p>
                      {/if}
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <span
                      class="text-xs px-2 py-1 rounded-full bg-sky-100 text-sky-700 capitalize"
                    >
                      {group.type}
                    </span>
                    <svg
                      class="w-5 h-5 text-gray-400 group-hover:text-sky-500 group-hover:translate-x-1 transition-all"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 5l7 7-7 7"
                      />
                    </svg>
                  </div>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        {#if quizGroups.length > 0}
          <div class="mb-6">
            <h2
              class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2"
            >
              <svg
                class="w-5 h-5 text-emerald-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
                />
              </svg>
              Quizzes
            </h2>
            <div
              class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden divide-y divide-gray-100"
            >
              {#each quizGroups as group, index}
                <button
                  class="w-full p-5 text-left hover:bg-emerald-50 transition-colors flex items-center justify-between group"
                  style="animation: slideIn 0.3s ease-out {index *
                    0.05}s backwards"
                  on:click={() => navigateToQuizGroup(group.id)}
                >
                  <div class="flex items-center gap-4">
                    <div
                      class="w-10 h-10 rounded-lg bg-emerald-100 flex items-center justify-center text-emerald-600 font-semibold"
                    >
                      {index + 1}
                    </div>
                    <div>
                      <h3 class="font-semibold text-gray-900">{group.name}</h3>
                      {#if group.description}
                        <p class="text-sm text-gray-500 line-clamp-1">
                          {group.description}
                        </p>
                      {/if}
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <span
                      class="text-xs px-2 py-1 rounded-full bg-emerald-100 text-emerald-700 capitalize"
                    >
                      {group.type}
                    </span>
                    <svg
                      class="w-5 h-5 text-gray-400 group-hover:text-emerald-500 group-hover:translate-x-1 transition-all"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 5l7 7-7 7"
                      />
                    </svg>
                  </div>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        {#if problemGroups.length === 0 && quizGroups.length === 0}
          <div
            class="bg-white rounded-2xl p-12 shadow-sm flex flex-col items-center text-center"
          >
            <div
              class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center mb-4"
            >
              <svg
                class="w-8 h-8 text-gray-400"
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
            <p class="text-gray-500 text-lg">No content available yet.</p>
            <p class="text-gray-400 text-sm mt-1">
              Your teacher hasn't added any exercises or quizzes to this course
              yet.
            </p>
          </div>
        {/if}
      {/if}
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
