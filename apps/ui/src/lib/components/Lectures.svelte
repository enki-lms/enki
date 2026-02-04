<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import PlusIcon from '$lib/components/icons/PlusIcon.svelte';
    import EmptyVideoIcon from '$lib/components/icons/EmptyVideoIcon.svelte';
    import ItemCard from '$lib/components/ItemCard.svelte';
    import Modal from '$lib/components/Modal.svelte';

    interface Lecture {
        id: string;
        title: string;
        description: string;
        date: string;
    }

    // za backend
    let lectures: Lecture[] = [];
    let isModalOpen: boolean = false;

    const handleNewLecture = () => {
        isModalOpen = true;
    };

    const handleCreate = (event: CustomEvent<{ name: string }>) => {
        const newLecture: Lecture = {
            id: Date.now().toString(),
            title: event.detail.name,
            description: "",
            date: new Date().toLocaleDateString('en-US', { 
                month: 'long', 
                day: 'numeric', 
                year: 'numeric' 
            })
        };
        lectures = [newLecture, ...lectures];
        // za backend
    };

    const handleCardClick = (id: string) => {
        // za backend
        console.log('Lecture clicked:', id);
    };
</script>

<div class="w-full">
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-semibold text-gray-900">Lectures</h2>
        <Button size="md" on:click={handleNewLecture}>
            <span slot="icon" class="text-white">
                <PlusIcon />
            </span>
            New Lecture
        </Button>
    </div>

    {#if lectures.length === 0}
        <div class="py-16 text-center text-gray-500">
            <EmptyVideoIcon />
            <p class="text-lg font-medium">No lectures yet</p>
            <p class="text-sm mt-1">Create your first lecture to get started</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each lectures as lecture (lecture.id)}
                <ItemCard
                    title={lecture.title}
                    description={lecture.description}
                    date={lecture.date}
                    type="lecture"
                    on:click={() => handleCardClick(lecture.id)}
                />
            {/each}
        </div>
    {/if}
</div>

<Modal 
    bind:isOpen={isModalOpen}
    title="Create New Lecture"
    placeholder="Lecture name..."
    on:create={handleCreate}
/>
