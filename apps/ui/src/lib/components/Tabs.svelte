<script lang="ts">
	export let tabs: string[] = [];
	export let active: string = "";

	const handleTabClick = (tab: string) => {
		active = tab;
	};
</script>

<div class="flex border-b border-gray-300">
	{#each tabs as tab}
		<button
			on:click={() => handleTabClick(tab)}
			class="px-4 py-3 font-semibold text-gray-800 pb-3 transition-colors duration-200 {active === tab ? 'border-b-2 border-[#77B6EA] text-black' : 'border-b-2 border-transparent text-gray-600'}"
		>
			{tab}
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
