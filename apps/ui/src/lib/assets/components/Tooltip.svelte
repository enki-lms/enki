<script lang="ts">
    export let text: string;
    export let position: "top" | "bottom" | "left" | "right" = "top";

    let showTooltip = false;

    const positionClasses = {
        top: "bottom-full mb-2 left-1/2 -translate-x-1/2",
        bottom: "top-full mt-2 left-1/2 -translate-x-1/2",
        left: "right-full mr-2 top-1/2 -translate-y-1/2",
        right: "left-full ml-2 top-1/2 -translate-y-1/2"
    };
</script>

<div class="relative inline-block">
    <div
        on:mouseenter={() => (showTooltip = true)}
        on:mouseleave={() => (showTooltip = false)}
    >
        <slot />
    </div>

    {#if showTooltip}
        <div
            class="absolute z-50 px-2 py-1 text-xs font-semibold bg-slate-900 text-blue-50 rounded whitespace-nowrap {positionClasses[position]}"
        >
            {text}
        </div>
    {/if}
</div>
