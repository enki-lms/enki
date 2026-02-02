<script lang="ts">
    export let padding: string = "p-6";
    export let shadow: boolean = true;
</script>

<div class="bg-white rounded-xl {padding} {shadow ? 'shadow-sm' : ''} border border-gray-200 {$$props.class || ''}">
    <slot />
</div>