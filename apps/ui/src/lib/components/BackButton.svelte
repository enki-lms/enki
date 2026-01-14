<script lang="ts">
    import ArrowLeftIcon from "$lib/components/icons/ArrowLeftIcon.svelte";
</script>

<button on:click class="w-20 h-20 border border-gray-300 rounded-2xl bg-white flex items-center justify-center hover:bg-gray-50 {$$props.class}">
    <ArrowLeftIcon/>
</button>
