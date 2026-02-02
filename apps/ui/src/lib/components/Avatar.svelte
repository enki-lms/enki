<script lang="ts">
    export let src: string;
    export let width: string = "80";
    export let height: string = "80";
</script>

<div class="rounded-full overflow-hidden border border-slate-300 p-1">
    <img {src} alt="Profile Avatar" class="size-full object-cover rounded-full" style="max-width: {width}px; max-height: {height}px; min-width: {width}px; min-height: {height}px;"/>
</div>