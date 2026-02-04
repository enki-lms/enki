<script lang="ts">
    import DocumentIcon from '$lib/components/icons/DocumentIcon.svelte';
    import VideoIcon from '$lib/components/icons/VideoIcon.svelte';
    import BookIcon from '$lib/components/icons/BookIcon.svelte';

    export let title: string;
    export let description: string = "";
    export let date: string;
    export let type: 'problem-set' | 'lecture' | 'course' = 'problem-set';
    
    const typeColors = {
        'problem-set': 'bg-blue-50 border-blue-200 hover:bg-blue-100',
        'lecture': 'bg-purple-50 border-purple-200 hover:bg-purple-100',
        'course': 'bg-green-50 border-green-200 hover:bg-green-100'
    };
</script>

<button 
    class="w-full p-5 rounded-xl border-2 transition-all duration-200 text-left {typeColors[type]} hover:shadow-md"
    on:click
>
    <div class="flex items-start gap-4">
        <div class="flex-shrink-0 mt-1">
            {#if type === 'problem-set'}
                <DocumentIcon />
            {:else if type === 'lecture'}
                <VideoIcon />
            {:else if type === 'course'}
                <BookIcon />
            {/if}
        </div>
        <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-gray-900 mb-1 truncate">{title}</h3>
            {#if description}
                <p class="text-sm text-gray-600 mb-2 line-clamp-2">{description}</p>
            {/if}
            <p class="text-xs text-gray-500">{date}</p>
        </div>
    </div>
</button>