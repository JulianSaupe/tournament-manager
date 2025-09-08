<script lang="ts">
    import '../../app.css';

    let {children} = $props();
</script>

<svelte:head>
    <title>Tournament Manager</title>
    <meta name="description" content="Professional tournament management system"/>
</svelte:head>

<div class="min-h-screen bg-base-100">
    <main class="w-full p-">
        {@render children?.()}
    </main>
</div>