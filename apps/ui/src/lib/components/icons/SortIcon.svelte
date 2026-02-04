<script lang="ts">
    export let direction: 'asc' | 'desc' = 'asc';
</script>

<svg class="w-4 h-4 transform {direction === 'desc' ? 'rotate-180' : ''}" viewBox="0 0 16 16" fill="currentColor">
    <path d="M8 3l4 5H4l4-5z"/>
</svg>
