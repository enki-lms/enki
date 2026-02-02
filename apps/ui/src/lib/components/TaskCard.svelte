<script lang="ts">
  import { goto } from "$app/navigation";

  export let title: string;
  export let id: string;

  function navigateToProblemGroup() {
    goto(`/problem-group/${id}`);
  }
</script>

<button
  class="w-full bg-[#EFF9FB] rounded-xl border border-gray-300 p-8 cursor-pointer hover:shadow-md transition-shadow text-left"
  onclick={navigateToProblemGroup}
>
  <h2 class="text-2xl font-bold text-gray-900">{title}</h2>
</button>
