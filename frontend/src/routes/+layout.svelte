<script lang="ts">
    import '../app.css';
    import favicon from '$lib/assets/favicon.svg';

    let {children} = $props();
</script>

<svelte:head>
    <link href={favicon} rel="icon"/>
    <title>TournamentProvider Manager</title>
    <meta content="Professional tournament management system" name="description"/>
</svelte:head>

{@render children?.()}
