<script lang="ts">
    import RightArrow from "$lib/components/icons/RightArrowIcon.svelte";

    export let color: "blue" = "blue";
    export let size: "sm" | "md" = "md";

    const colors = new Map<string, string>([
        ["blue", ""]
    ]);

    const sizes = new Map<string, string>([
        ["md", "w-36 h-14"],
        ["sm", "w-24 h-14"]
    ]);
</script>

<button class="rounded-full bg-sky-400 hover:bg-sky-500 font-semibold font-sans text-lg {sizes.get(size)} text-white {$$props.class}" on:click>
    <span class="flex flex-row gap-1 items-center justify-center">
        <span class="text-black">
            <slot/>
        </span>
        <slot name="icon"/>
    </span>
</button>