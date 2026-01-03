<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import button from '$lib/assets/components/Button.svelte';
	

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
