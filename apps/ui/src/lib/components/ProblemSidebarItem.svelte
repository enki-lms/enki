<script lang="ts">
  import { goto } from "$app/navigation";

  let {
    id,
    title,
    description = "",
    completed = false,
    active = false,
  }: {
    id: string;
    title: string;
    description?: string;
    completed?: boolean;
    active?: boolean;
  } = $props();

  function navigateToProblem() {
    goto(`/problem/${id}`);
  }
</script>

<button
  class="group w-full text-left p-4 rounded-xl border-2 transition-all duration-200 ease-out
    {active
    ? 'bg-sky-50 border-sky-400 shadow-md'
    : 'bg-white border-transparent hover:bg-[#EFF9FB] hover:border-sky-200 hover:shadow-sm'}"
  onclick={navigateToProblem}
>
  <div class="flex items-center gap-3">
    <!-- Status Icon -->
    <div
      class="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-medium
        {completed
        ? 'bg-emerald-500'
        : 'bg-slate-300'} transition-transform duration-200 group-hover:scale-110"
    >
      {completed ? "✓" : "○"}
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <h3
        class="font-semibold text-gray-900 truncate group-hover:text-sky-700 transition-colors duration-200"
      >
        {title}
      </h3>
      {#if description}
        <p class="mt-1 text-sm text-gray-500 line-clamp-2">
          {description}
        </p>
      {/if}
    </div>

    <!-- Arrow -->
    <div
      class="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
    >
      <svg class="w-5 h-5 text-sky-500" viewBox="0 0 20 20" fill="currentColor">
        <path
          fill-rule="evenodd"
          d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
  </div>
</button>
