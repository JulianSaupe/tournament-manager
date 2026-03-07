<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Tournament Manager</title>
	<meta name="description" content="Professional tournament management system" />
</svelte:head>

{@render children?.()}
