<script lang="ts">
    export let disabled = false;
</script>

<button {disabled} on:click class:opacity-50={disabled} class="bg-slate-300 rounded-full w-7 h-7 flex justify-center items-center">
    <slot />
</button>