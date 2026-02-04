<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { fade, scale } from 'svelte/transition';
    import Button from '$lib/components/Button.svelte';
    
    export let isOpen: boolean = false;
    export let title: string;
    export let placeholder: string = "Enter name...";
    
    let inputValue: string = "";
    const dispatch = createEventDispatcher();
    
    const handleCreate = () => {
        if (inputValue.trim()) {
            dispatch('create', { name: inputValue });
            inputValue = "";
            isOpen = false;
        }
    };
    
    const handleCancel = () => {
        inputValue = "";
        isOpen = false;
    };
    
    const handleKeydown = (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleCreate();
        } else if (e.key === 'Escape') {
            handleCancel();
        }
    };
</script>

{#if isOpen}
    <div 
        class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4"
        on:click={handleCancel}
        on:keydown={handleKeydown}
        role="button"
        tabindex="0"
        transition:fade={{ duration: 200 }}
    >
        <div 
            class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6"
            on:click|stopPropagation
            role="dialog"
            aria-modal="true"
            transition:scale={{ duration: 200, start: 0.95 }}
        >
            <h2 class="text-2xl font-semibold text-gray-900 mb-4">{title}</h2>
            
            <input 
                type="text" 
                bind:value={inputValue}
                {placeholder}
                class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors"
                autofocus
                on:keydown={handleKeydown}
            />
            
            <div class="flex gap-3 mt-6 justify-end">
                <button 
                    on:click={handleCancel}
                    class="px-6 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-semibold hover:bg-gray-50 transition-colors"
                >
                    Cancel
                </button>
                <Button size="sm" on:click={handleCreate} class="px-6">
                    Create
                </Button>
            </div>
        </div>
    </div>
{/if}