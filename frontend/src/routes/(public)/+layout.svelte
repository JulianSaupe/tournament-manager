<script lang="ts">
    import '../../app.css';

    let {children} = $props();
</script>

<svelte:head>
    <title>TournamentProvider Manager</title>
    <meta content="Professional tournament management system" name="description"/>
</svelte:head>

<div class="min-h-screen bg-base-100">
    <main class="p- w-full">
        {@render children?.()}
    </main>
</div>
