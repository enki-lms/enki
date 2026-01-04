<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	function handleKeyDown(event: KeyboardEvent) {
		if (event.shiftKey && event.key === "D") {
			goto("/component-demos");
		}
	}

	onMount(() => {
		goto("/auth");
	});
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="w-full h-screen bg-slate-900 flex items-center justify-center">
	<p class="text-slate-400 text-sm">Press <span class="font-mono font-semibold text-blue-50">Shift + D</span> to view components</p>
</div>
