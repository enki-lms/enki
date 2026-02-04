<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import PlusIcon from '$lib/components/icons/PlusIcon.svelte';
    import EmptyDocumentIcon from '$lib/components/icons/EmptyDocumentIcon.svelte';
    import ItemCard from '$lib/components/ItemCard.svelte';
    import Modal from '$lib/components/Modal.svelte';

    interface ProblemSet {
        id: string;
        title: string;
        description: string;
        date: string;
    }

    // za backend
    let problemSets: ProblemSet[] = [];
    let isModalOpen: boolean = false;

    const handleNewProblemSet = () => {
        isModalOpen = true;
    };

    const handleCreate = (event: CustomEvent<{ name: string }>) => {
        const newProblemSet: ProblemSet = {
            id: Date.now().toString(),
            title: event.detail.name,
            description: "",
            date: new Date().toLocaleDateString('en-US', { 
                month: 'long', 
                day: 'numeric', 
                year: 'numeric' 
            })
        };
        problemSets = [newProblemSet, ...problemSets];
        // za backend
    };

    const handleCardClick = (id: string) => {
        // za backend
        console.log('Problem set clicked:', id);
    };
</script>

<div class="w-full">
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-semibold text-gray-900">Problem Sets</h2>
        <Button size="md" on:click={handleNewProblemSet}>
            <span slot="icon" class="text-white">
                <PlusIcon />
            </span>
            New Problem Set
        </Button>
    </div>

    {#if problemSets.length === 0}
        <div class="py-16 text-center text-gray-500">
            <EmptyDocumentIcon />
            <p class="text-lg font-medium">No problem sets yet</p>
            <p class="text-sm mt-1">Create your first problem set to get started</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each problemSets as problemSet (problemSet.id)}
                <ItemCard
                    title={problemSet.title}
                    description={problemSet.description}
                    date={problemSet.date}
                    type="problem-set"
                    on:click={() => handleCardClick(problemSet.id)}
                />
            {/each}
        </div>
    {/if}
</div>

<Modal 
    bind:isOpen={isModalOpen}
    title="Create New Problem Set"
    placeholder="Problem set name..."
    on:create={handleCreate}
/>
