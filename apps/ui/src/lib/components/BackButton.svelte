<script lang="ts">
  import ArrowLeftIcon from "$lib/components/icons/ArrowLeftIcon.svelte";

  let {
    onclick = () => history.back(),
    class: className = "",
  }: { onclick?: () => void; class?: string } = $props();
</script>

<button
  {onclick}
  class="w-20 h-20 border border-gray-300 rounded-2xl bg-white flex items-center justify-center hover:bg-gray-50 {className}"
>
  <ArrowLeftIcon />
</button>
