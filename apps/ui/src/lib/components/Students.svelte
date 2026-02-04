<script lang="ts">
	import { onMount } from 'svelte';
	import StudentCard from './StudentCard.svelte';

	export interface Student {
		id: string;
		name: string;
		email: string;
	}

	let students: Student[] = [];

	async function fetchStudents(): Promise<void> {
		// za backend
		students = [
	]
	}

	onMount(() => {
		fetchStudents();
	});
</script>

<div class="students-container">
	<div class="mb-6">
		<h2 class="text-2xl font-bold text-gray-800">Students</h2>
	</div>

	{#if students.length > 0}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each students as student (student.id)}
				<StudentCard {student} />
			{/each}
		</div>
	{:else}
		<div class="text-center py-12 text-gray-500">
			<p class="text-lg">No students found</p>
		</div>
	{/if}
</div>

<style>
	.students-container {
		width: 100%;
	}
</style>