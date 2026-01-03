<script lang="ts">
    import { goto } from "$app/navigation";
    import BackArrowIcon from "$lib/assets/components/icons/BackArrowIcon.svelte";
    import ArrowButton from "$lib/assets/components/ArrowButton.svelte";
    import ShareIcon from "$lib/assets/components/icons/ShareIcon.svelte";

    export let backPath: string;
</script>

<nav class="w-full font-open-sans text-slate-900 pt-8 h-24 bg-blue-50 rounded-b-3xl flex items-center p-6 justify-between {$$props.class}">
    <ArrowButton disabled={false} on:click={async () => await goto(backPath)}>
        <BackArrowIcon/>
    </ArrowButton>
    <span class="font-semibold">
        <slot name="first" />
    </span>
    <slot name="second" />
    <slot name="icon" />
    <!--        <span class="text-blue">15:30 - 15:36</span>-->
    <!--        <ShareIcon/>-->
</nav>
