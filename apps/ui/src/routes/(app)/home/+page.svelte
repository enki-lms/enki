<script lang="ts">
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { api, type Course } from "$lib/api";

  export let data;

  let courses: Course[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      courses = await api.getEnrolledCourses();
    } catch (e) {
      console.error("Error loading courses:", e);
      error = e instanceof Error ? e.message : "Failed to load courses";
    } finally {
      loading = false;
    }
  });

  function navigateToCourse(courseId: number) {
    goto(`/course/${courseId}`);
  }
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex">
    <div class="p-2">
      <ProfileMenu
        width="80"
        height="80"
        name={data.user?.fullName ?? ""}
        email={data.user?.email ?? ""}
        role={data.user?.role ?? ""}
      />
    </div>
    <div class="p-3"><LogoPlaceHolder /></div>
  </header>
  <main class="p-8">
    <div class="bg-white rounded-xl p-12 shadow-sm">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">My Courses</h1>
      {#if loading}
        <div class="flex justify-center p-8">
          <div class="flex flex-col items-center gap-3">
            <div
              class="w-8 h-8 border-3 border-sky-400 border-t-transparent rounded-full animate-spin"
            ></div>
            <p class="text-gray-500">Loading courses...</p>
          </div>
        </div>
      {:else if error}
        <div class="flex justify-center p-8">
          <p class="text-red-500">{error}</p>
        </div>
      {:else if courses.length === 0}
        <div class="flex flex-col items-center justify-center p-12 text-center">
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
                d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
              />
            </svg>
          </div>
          <p class="text-gray-500 text-lg">
            You are not enrolled in any courses yet.
          </p>
          <p class="text-gray-400 text-sm mt-1">
            Contact your teacher to get enrolled.
          </p>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {#each courses as course}
            <button
              class="w-full bg-gradient-to-br from-sky-50 to-emerald-50 rounded-xl border-2 border-sky-100 p-8 cursor-pointer hover:shadow-lg hover:border-sky-200 transition-all text-left group"
              on:click={() => navigateToCourse(course.id)}
            >
              <div class="flex items-start gap-4">
                <div
                  class="w-12 h-12 rounded-xl bg-gradient-to-br from-sky-400 to-emerald-400 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform"
                >
                  <svg
                    class="w-6 h-6 text-white"
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
                <div class="flex-1 min-w-0">
                  <h2 class="text-xl font-bold text-gray-900 mb-1 truncate">
                    {course.name}
                  </h2>
                  <p class="text-sm text-gray-600">{course.subject}</p>
                </div>
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
      {/if}
    </div>
  </main>
</div>
