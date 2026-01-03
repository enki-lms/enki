<script lang="ts">
    export let padding = "p-6";
    export let shadow = true;
</script>

<div class="bg-blue-50 rounded-lg {padding} {shadow ? 'shadow-md' : ''} border border-slate-300 {$$props.class}">
    <slot />
</div>
