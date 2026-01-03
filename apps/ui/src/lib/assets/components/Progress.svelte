<script lang="ts">
    export let value: number = 0; // 0-100
    export let showLabel = true;
</script>

<div class="w-full {$$props.class}">
    <div class="relative w-full h-3 bg-slate-300 rounded-full overflow-hidden">
        <div
            class="h-full bg-sky-400 rounded-full transition-all duration-300"
            style="width: {Math.min(100, Math.max(0, value))}%"
        />
    </div>
    {#if showLabel}
        <p class="text-sm text-slate-600 mt-1">{Math.round(value)}%</p>
    {/if}
</div>
