<script lang="ts">
	import Tabs from '$lib/components/Tabs.svelte';
	import Card from '$lib/components/Card.svelte';
	import LogoPlaceHolder from '$lib/components/LogoPlaceHolder.svelte';
	import Avatar from '$lib/components/Avatar.svelte';
	import FileManager from '$lib/components/FileManager.svelte';

	type TabType = 'Problem Sets' | 'Lectures' | 'Courses' | 'Students' | 'Course Materials';

	const tabsList: TabType[] = ['Problem Sets', 'Lectures', 'Courses', 'Students', 'Course Materials'];
	let activeTab: TabType = 'Problem Sets';

    
	const renderContent = (tab: TabType): string => {
		switch (tab) {
			case 'Problem Sets':
				return 'Problem Sets content';
			case 'Lectures':
				return 'Lectures content';
			case 'Courses':
				return 'Courses content';
			case 'Students':
				return 'Students content';
			default:
				return '';
		}
	};
</script>





<div class="min-h-screen p-8 bg-gray-100">
	<div class="max-w-6xl mx-auto">
		<div class="mb-6 flex items-center gap-6">
			<Avatar width="80" height="80" src="src/lib/assets/images/temp-pfp.png" />
			<LogoPlaceHolder />
		</div>


		<Card padding="p-0" shadow={false} class="bg-white rounded-xl border-gray-200">
			<div class="px-6 pt-6">
				<Tabs tabs={tabsList} bind:active={activeTab} />
			</div>



			<div class="px-6 pb-6 pt-8 min-h-[600px]">
				{#key activeTab}
					{#if activeTab === 'Course Materials'}
						<FileManager />
					{:else}
						<div class="text-lg text-gray-700" data-tab-content>
							{renderContent(activeTab)}
						</div>
					{/if}
				{/key}
			</div>
		</Card>
	</div>
</div>