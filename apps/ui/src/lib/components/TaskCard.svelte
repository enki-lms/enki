<script lang="ts">
	let isOpen = false;
	let isClosing = false;

	export let title: string;
	export let tasks: string[] = [];

	function openOverlay() {
		isOpen = true;
		isClosing = false;
	}

	function closeOverlay() {
		isClosing = true;
		setTimeout(() => {
			isOpen = false;
			isClosing = false;
		}, 300);
	}
</script>

<div class="bg-[#EFF9FB] rounded-xl border border-gray-300 p-8 cursor-pointer hover:shadow-md transition-shadow" on:click={openOverlay}>
	<h2 class="text-2xl font-bold text-gray-900">{title}</h2>
</div>

{#if isOpen}
    <div class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 {isClosing ? 'animate-fadeOut' : 'animate-fadeIn'}" role="presentation" on:click={closeOverlay} on:keydown={(e) => e.key === 'Escape' && closeOverlay()}>
        <div class="bg-white rounded-xl shadow-lg p-8 max-w-md w-full {isClosing ? 'animate-fadeOutScale' : 'animate-fadeInScale'}" role="dialog" aria-modal="true" on:click|stopPropagation>

            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-bold text-gray-900">{title}</h2>
                <button on:click={closeOverlay} class="text-gray-500 hover:text-gray-700 text-2xl">×</button>
            </div>

            <div class="mb-6 max-h-96 overflow-y-auto">
                {#each tasks as task}
                    <div class="py-2 px-3 bg-gray-50 rounded mb-2 text-gray-700">{task}</div>
                {/each}
                {#if tasks.length === 0}
                    <p class="text-gray-500 text-center py-4">No tasks available</p>
                {/if}
            </div>

            <button on:click={closeOverlay} class="w-full bg-[#77B6EA] hover:bg-[#77B6EA]/80 text-white font-semibold py-2 rounded-lg transition-colors">
                Continue
            </button>
        </div>
    </div>
{/if}

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

    @keyframes fadeOut {
		from {
			opacity: 1;
		}
		to {
			opacity: 0;
		}
	}

	@keyframes fadeInScale {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	@keyframes fadeOutScale {
		from {
			opacity: 1;
			transform: scale(1);
		}
		to {
			opacity: 0;
			transform: scale(0.95);
		}
	}

	:global(.animate-fadeIn) {
		animation: fadeIn 0.3s ease-out;
	}

	:global(.animate-fadeInScale) {
		animation: fadeInScale 0.3s ease-out;
	}

    :global(.animate-fadeOut) {
        animation: fadeOut 0.3s ease-out;
    }

    :global(.animate-fadeOutScale) {
        animation: fadeOutScale 0.3s ease-out;
    }
</style>
