<script lang="ts">
    import BackButton from "$lib/components/BackButton.svelte";
    import { goto } from "$app/navigation";

    export let href: string;
</script>

<div class="flex items-center justify-between w-full relative {$$props.class}">
    <BackButton class="absolute" on:click={() => goto(href)} />
    <div class="flex-1 flex items-center justify-center">
        <h1 class="text-white text-2xl font-bold">
            <slot />
        </h1>
    </div>
</div>