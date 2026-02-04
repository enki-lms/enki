<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import PlusIcon from '$lib/components/icons/PlusIcon.svelte';
    import EmptyBookIcon from '$lib/components/icons/EmptyBookIcon.svelte';
    import ItemCard from '$lib/components/ItemCard.svelte';
    import Modal from '$lib/components/Modal.svelte';

    interface Course {
        id: string;
        title: string;
        description: string;
        date: string;
    }

    // za backend
    let courses: Course[] = [];
    let isModalOpen: boolean = false;

    const handleNewCourse = () => {
        isModalOpen = true;
    };

    const handleCreate = (event: CustomEvent<{ name: string }>) => {
        const newCourse: Course = {
            id: Date.now().toString(),
            title: event.detail.name,
            description: "",
            date: new Date().toLocaleDateString('en-US', { 
                month: 'long', 
                day: 'numeric', 
                year: 'numeric' 
            })
        };
        courses = [newCourse, ...courses];
        // za backend
    };

    const handleCardClick = (id: string) => {
        // za backend
        console.log('Course clicked:', id);
    };
</script>

<div class="w-full">

    <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-semibold text-gray-900">Courses</h2>
        <Button size="md" on:click={handleNewCourse}>
            <span slot="icon" class="text-white">
                <PlusIcon />
            </span>
            New Course
        </Button>
    </div>

    <!-- Courses Grid -->
    {#if courses.length === 0}
        <div class="py-16 text-center text-gray-500">
            <EmptyBookIcon />
            <p class="text-lg font-medium">No courses yet</p>
            <p class="text-sm mt-1">Create your first course to get started</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each courses as course (course.id)}
                <ItemCard
                    title={course.title}
                    description={course.description}
                    date={course.date}
                    type="course"
                    on:click={() => handleCardClick(course.id)}
                />
            {/each}
        </div>
    {/if}
</div>

<Modal 
    bind:isOpen={isModalOpen}
    title="Create New Course"
    placeholder="Course name..."
    on:create={handleCreate}
/>
