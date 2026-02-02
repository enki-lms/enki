<script lang="ts">
    export let type: 'document' | 'image' | 'zip' = 'document';
</script>

<svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    {#if type === 'document'}
        <path d="M7 3C5.89543 3 5 3.89543 5 5V19C5 20.1046 5.89543 21 7 21H17C18.1046 21 19 20.1046 19 19V9L13 3H7Z" fill="#8B95A5"/>
        <path d="M13 3V7C13 8.10457 13.8954 9 15 9H19" fill="#6B7280"/>
    {:else if type === 'image'}
        <path d="M7 3C5.89543 3 5 3.89543 5 5V19C5 20.1046 5.89543 21 7 21H17C18.1046 21 19 20.1046 19 19V9L13 3H7Z" fill="#34D399"/>
        <path d="M13 3V7C13 8.10457 13.8954 9 15 9H19" fill="#10B981"/>
    {:else if type === 'zip'}
        <path d="M7 3C5.89543 3 5 3.89543 5 5V19C5 20.1046 5.89543 21 7 21H17C18.1046 21 19 20.1046 19 19V9L13 3H7Z" fill="#FFB340"/>
        <path d="M13 3V7C13 8.10457 13.8954 9 15 9H19" fill="#FF9800"/>
    {/if}
</svg>
