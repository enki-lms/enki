<script lang="ts">
  import Avatar from "$lib/components/Avatar.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import TaskCard from "$lib/components/TaskCard.svelte";
  import { onMount } from "svelte";
  import { api, type ProblemGroup } from "$lib/api";

  let problemGroups: ProblemGroup[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const courses = await api.getCourses();
      console.log("Fetched courses:", courses);
      if (!courses) {
        console.warn("No courses returned from API");
        problemGroups = [];
        return;
      }

      const allGroupsPromises = courses.map((course) => {
        console.log("Processing course:", course, "ID:", course.id);
        return api.getCourseProblemGroups(course.id);
      });
      const groupsArrays = await Promise.all(allGroupsPromises);
      problemGroups = groupsArrays.flat();
    } catch (e) {
      console.error("Error loading courses:", e);
      error = e instanceof Error ? e.message : "Failed to load courses";
    } finally {
      loading = false;
    }
  });
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex">
    <div class="p-2">
      <Avatar width="80" height="80" src="src/lib/assets/images/temp-pfp.png" />
    </div>
    <div class="p-3"><LogoPlaceHolder /></div>
  </header>
  <main class="p-8">
    <div class="bg-white rounded-xl p-12 shadow-sm">
      {#if loading}
        <div class="flex justify-center p-8">
          <p class="text-gray-500">Loading courses...</p>
        </div>
      {:else if error}
        <div class="flex justify-center p-8">
          <p class="text-red-500">{error}</p>
        </div>
      {:else}
        <div class="grid grid-cols-2 gap-6">
          {#each problemGroups as group}
            <TaskCard id={group.id.toString()} title={group.name} />
          {/each}
          {#if problemGroups.length === 0}
            <div class="col-span-2 text-center text-gray-500">
              No problem groups found.
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </main>
</div>
