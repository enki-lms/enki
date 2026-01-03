<script lang="ts">
    import { goto } from "$app/navigation";

    export let href = "";
</script>

<div class="mt-[-1px] border-t border-b border-lightgray p-4">
    <button on:click class="flex border-none justify-between items-center w-full">
        <slot/>
        <slot name="icon" />
    </button>
</div>