<script lang="ts">
    import CourseCard from '$lib/components/CourseCard.svelte';
    import FileManager from '$lib/components/FileManager.svelte';

    interface Course {
        id: string;
        name: string;
        description: string;
        color: string;
    }

    // Placeholder courses - za backend
    const courses: Course[] = [
        
    ];

    let selectedCourse: Course | null = null;

    const handleCourseClick = (course: Course) => {
        selectedCourse = course;
        // za backend
        console.log('Selected course:', course.id);
    };

    const handleBackToCourses = () => {
        selectedCourse = null;
    };
</script>

<div class="w-full">
    {#if !selectedCourse}

        <div class="mb-6">
            <h2 class="text-2xl font-semibold text-gray-900">Course Materials</h2>
            <p class="text-sm text-gray-600 mt-1">Select a course to view its materials</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each courses as course (course.id)}
                <CourseCard
                    name={course.name}
                    description={course.description}
                    color={course.color}
                    on:click={() => handleCourseClick(course)}
                />
            {/each}
        </div>
    {:else}

        <div class="mb-6 flex items-center gap-4">
            <button
                on:click={handleBackToCourses}
                class="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
            >
                <svg class="w-5 h-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
                </svg>
                <span class="font-medium">Back to Courses</span>
            </button>
            <div class="flex-1">
                <h2 class="text-xl font-semibold text-gray-900">{selectedCourse.name}</h2>
                <p class="text-sm text-gray-600">{selectedCourse.description}</p>
            </div>
        </div>

        <FileManager courseId={selectedCourse.id} />
    {/if}
</div>
