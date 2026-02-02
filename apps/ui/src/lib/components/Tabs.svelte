<script lang="ts">
	export let tabs: string[] = [];
	export let active: string = "";

	const handleTabClick = (tab: string): void => {
		active = tab;
	};
</script>

<div class="flex border-b border-gray-200">
	{#each tabs as tab}
		<button
			on:click={() => handleTabClick(tab)}
			class="px-8 py-4 text-base font-medium transition-colors duration-200 relative {active === tab 
				? 'text-gray-900 font-semibold' 
				: 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'}"
		>
			{tab}
			{#if active === tab}
				<span class="absolute bottom-0 left-0 right-0 h-0.5 bg-cyan-500"></span>
			{/if}
		</button>
	{/each}
</div>

<style>
	:global([data-tab-content]) {
		animation: fadeIn 0.2s ease-in-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
</style>