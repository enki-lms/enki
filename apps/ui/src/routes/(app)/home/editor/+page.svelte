<script lang="ts">
	import BackButton from "$lib/components/BackButton.svelte";
	import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
	import Tabs from "$lib/components/Tabs.svelte";
	import LanguageSelector from "$lib/components/LanguageSelector.svelte";
	import Chat from "$lib/components/Chat.svelte";

	let activeDescTab = "Description";
	const descTabs = ["Description", "Scratchpad", "Assistant"];
	
	let activeSolutionTab = "Solution";
	const solutionTabs = ["Solution"];
	
	let activeTestTab = "Test Result";
	const testTabs = ["Test Result"];
	
	let solutionText = "";
	let testResultText = "";
	let scratchpadText = "";
	let messages: Array<{ id: number; text: string; sender: "user" | "assistant" }> = [];

	const problemDescription = {
		title: "",
		description: "",
		rules: []
	};
</script>

<div class="h-screen bg-[#E8EEF2] flex flex-col">
	<!-- Header -->
	<div class="bg-white border-b border-gray-300 px-4 md:px-6 py-3 md:py-4 flex items-center gap-2 md:gap-4 flex-wrap md:flex-nowrap">
		<BackButton />
		<div class="hidden md:block">
			<LogoPlaceHolder />
		</div>
		<div class="flex-1 md:flex-none">
			<LanguageSelector />
		</div>
	</div>

	<!-- Main Content -->
	<div class="flex flex-1 gap-3 md:gap-6 p-3 md:p-6 overflow-hidden flex-col lg:flex-row">
		<!-- Left Panel - Problem Description -->
		<div class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white">
			<Tabs tabs={descTabs} bind:active={activeDescTab} />
			<div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
				{#if activeDescTab === "Description"}
					<div class="overflow-y-auto h-full p-3 md:p-6">
						{#if problemDescription.title}
							<h2 class="text-xl md:text-2xl font-bold mb-4">{problemDescription.title}</h2>
						{/if}
						{#if problemDescription.description}
							<p class="mb-4 text-sm md:text-base text-gray-700">{problemDescription.description}</p>
						{/if}
						{#if problemDescription.rules && problemDescription.rules.length > 0}
							<p class="mb-3 font-semibold text-sm md:text-base text-gray-800">You can move according to these rules:</p>
							<ul class="list-disc list-inside space-y-2 ml-2 text-sm md:text-base text-gray-700">
								{#each problemDescription.rules as rule}
									<li>{rule}</li>
								{/each}
							</ul>
						{/if}
					</div>
				{:else if activeDescTab === "Scratchpad"}
					<textarea
						bind:value={scratchpadText}
						class="w-full h-full p-3 md:p-4 border-none resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA]"
						placeholder="Write your scratchpad notes here..."
					/>
				{:else if activeDescTab === "Assistant"}
					<Chat bind:messages />
				{/if}
			</div>
		</div>

		<!-- Right Panel - Solution and Test Result as separate frames -->
		<div class="flex-1 lg:w-1/2 flex flex-col gap-3 md:gap-6 min-w-0">
			<!-- Solution Frame -->
			<div class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white">
				<Tabs tabs={solutionTabs} bind:active={activeSolutionTab} />
				<div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
					<textarea
						bind:value={solutionText}
						class="w-full h-full p-3 md:p-6 resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA]"
						placeholder="Write your solution here..."
					/>
				</div>
			</div>

			<!-- Test Result Frame -->
			<div class="flex-1 flex flex-col min-w-0 border border-gray-300 rounded-lg bg-white">
				<Tabs tabs={testTabs} bind:active={activeTestTab} />
				<div class="flex-1 overflow-hidden animate-fadeIn" data-tab-content>
					<textarea
						bind:value={testResultText}
						class="w-full h-full p-3 md:p-6 resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA]"
						placeholder="Test results will appear here..."
					/>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	:global(.animate-fadeIn) {
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

	@media (max-width: 1024px) {
		:global(.animate-fadeIn) {
			animation: fadeIn 0.2s ease-in-out;
		}
	}
</style>
