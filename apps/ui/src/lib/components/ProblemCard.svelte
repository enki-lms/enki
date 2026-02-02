<script lang="ts">
  import { goto } from "$app/navigation";

  let {
    title,
    id,
    description = "",
  }: { title: string; id: string; description?: string } = $props();

  function navigateToProblem() {
    goto(`/problem/${id}`);
  }
</script>

<button
  class="w-full bg-[#EFF9FB] rounded-xl border border-gray-300 p-6 cursor-pointer hover:shadow-md transition-shadow text-left"
  onclick={navigateToProblem}
>
  <h3 class="text-xl font-bold text-gray-900">{title}</h3>
  {#if description}
    <p class="mt-2 text-gray-600 text-sm">{description}</p>
  {/if}
</button>
