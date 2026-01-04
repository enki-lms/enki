<script lang="ts">
    export let removable = false;
    export let onRemove: (() => void) | undefined = undefined;
</script>

<span class="inline-flex items-center gap-2 px-3 py-1 text-sm bg-slate-300 text-slate-900 rounded-full {$$props.class}">
    <slot />
    {#if removable}
        <button
            on:click={onRemove}
            class="hover:opacity-70 font-bold ml-1"
            aria-label="Remove tag"
        >
            ×
        </button>
    {/if}
</span>
