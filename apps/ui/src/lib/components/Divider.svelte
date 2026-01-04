<script lang="ts">
    export let text: string = "";
</script>

{#if text}
    <div class="flex items-center gap-4 my-4">
        <hr class="flex-1 border-t border-slate-300" />
        <span class="text-slate-600 text-sm font-medium px-2">{text}</span>
        <hr class="flex-1 border-t border-slate-300" />
    </div>
{:else}
    <hr class="my-4 border-t border-slate-300" />
{/if}
