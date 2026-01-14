<script lang="ts">

    import DownArrowIcon from "$lib/components/icons/ArrowDownIcon.svelte";

	let isOpen = false;
	let selectedLanguage = "Python";

	const languages = ["Python"];

	const toggleDropdown = () => {
		isOpen = !isOpen;
	};

	const selectLanguage = (lang: string) => {
		selectedLanguage = lang;
		isOpen = false;
	};
</script>

<div class="relative">
	<button
		on:click={toggleDropdown}
		class="flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
	>
		<span class="font-semibold text-gray-800">{selectedLanguage}</span>
		<span class="text-gray-400 transition-transform duration-200 {isOpen ? 'rotate-180' : ''}">
            <DownArrowIcon />
        </span>
	</button>

	{#if isOpen}
		<div class="absolute top-full mt-1 w-40 bg-white border border-gray-300 rounded-lg shadow-lg z-10 animate-fadeIn">
			{#each languages as lang}
				<button
					on:click={() => selectLanguage(lang)}
					class="w-full text-left px-4 py-2 hover:bg-gray-100"
				>
					{lang}
				</button>
			{/each}
		</div>
	{/if}
</div>

<style>
	:global(.animate-fadeIn) {
		animation: fadeIn 0.2s ease-in-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
