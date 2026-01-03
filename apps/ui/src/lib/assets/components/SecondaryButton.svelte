<script lang="ts">
    export let size: "sm" | "md" = "md";

    const sizes = new Map<string, string>([
        ["md", "px-6 py-3"],
        ["sm", "px-4 py-2"]
    ]);
</script>

<button
    class="rounded-lg font-semibold border-2 border-slate-300 text-slate-900 hover:bg-slate-100 transition-colors {sizes.get(size)} {$$props.class}"
    on:click
>
    <slot />
</button>
