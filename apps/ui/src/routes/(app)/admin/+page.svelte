<script lang="ts">
	import Tabs from '$lib/components/Tabs.svelte';
	import Card from '$lib/components/Card.svelte';
	import LogoPlaceHolder from '$lib/components/LogoPlaceHolder.svelte';
	import Avatar from '$lib/components/Avatar.svelte';
	import FileManager from '$lib/components/FileManager.svelte';
	import ProblemSets from '$lib/components/ProblemSets.svelte';
	import Lectures from '$lib/components/Lectures.svelte';
	import Courses from '$lib/components/Courses.svelte';
	import CourseMaterials from '$lib/components/CourseMaterials.svelte';
	import Students from '$lib/components/Students.svelte';

	type TabType = 'Problem Sets' | 'Lectures' | 'Courses' | 'Students' | 'Course Materials';

	const tabsList: TabType[] = ['Problem Sets', 'Lectures', 'Courses', 'Students', 'Course Materials'];
	let activeTab: TabType = 'Problem Sets';
</script>

<div class="min-h-screen p-8 bg-gray-100">
	<div class="max-w-6xl mx-auto">
		<!-- Header with Avatar and Logo -->
		<div class="mb-6 flex items-center gap-6">
			<Avatar width="80" height="80" src="src/lib/assets/images/temp-pfp.png" />
			<LogoPlaceHolder />
		</div>

		<!-- Main Card -->
		<Card padding="p-0" shadow={false} class="bg-white rounded-xl border-gray-200">
			<!-- Tabs -->
			<div class="px-6 pt-6">
				<Tabs tabs={tabsList} bind:active={activeTab} />
			</div>

			<!-- Tab Content Area -->
			<div class="px-6 pb-6 pt-8 min-h-[600px]">
				{#if activeTab === 'Problem Sets'}
					<ProblemSets />
				{:else if activeTab === 'Lectures'}
					<Lectures />
				{:else if activeTab === 'Courses'}
					<Courses />
				{:else if activeTab === 'Students'}
					<Students />
				{:else if activeTab === 'Course Materials'}
					<CourseMaterials />
				{/if}
			</div>
		</Card>
	</div>
</div>