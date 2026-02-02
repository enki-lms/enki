<script lang="ts">
  import { page } from "$app/stores";
  import Avatar from "$lib/components/Avatar.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import ProblemCard from "$lib/components/ProblemCard.svelte";
  import BackButton from "$lib/components/BackButton.svelte";
  import { onMount } from "svelte";
  import { api, type Problem } from "$lib/api";

  // Get the problem group ID from the URL
  const groupId = $page.params.id ?? "";

  let problems: Problem[] = [];
  let loading = true;
  let error: string | null = null;
  // Note: We don't have an endpoint to get just the group details easily without iterating courses,
  // so we'll just show the problems for now, or we could fetch group details if needed.
  // For the demo, "Problem List" is sufficient title if we can't get the group name easily.
  let groupTitle = "Problem List";

  onMount(async () => {
    try {
      problems = await api.getGroupProblems(groupId);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load problems";
    } finally {
      loading = false;
    }
  });
</script>

<div class="min-h-screen bg-gray-200">
  <header class="p-2 flex items-center gap-4">
    <div class="p-2">
      <BackButton onclick={() => history.back()} />
    </div>
    <div class="p-2">
      <Avatar
        width="60"
        height="60"
        src="/src/lib/assets/images/temp-pfp.png"
      />
    </div>
    <div class="p-3"><LogoPlaceHolder /></div>
  </header>
  <main class="p-8">
    <div class="bg-white rounded-xl p-12 shadow-sm">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">
        {groupTitle}
      </h1>

      {#if loading}
        <div class="flex justify-center p-8">
          <p class="text-gray-500">Loading problems...</p>
        </div>
      {:else if error}
        <div class="flex justify-center p-8">
          <p class="text-red-500">{error}</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {#each problems as problem}
            <ProblemCard
              id={problem.id.toString()}
              title={problem.name}
              description={problem.description}
            />
          {/each}
          {#if problems.length === 0}
            <div class="col-span-2 text-center text-gray-500">
              No problems found in this group.
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </main>
</div>
