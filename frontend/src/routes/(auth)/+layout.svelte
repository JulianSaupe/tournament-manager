<script lang="ts">
    import '../../app.css';

    let { children } = $props();
</script>

<svelte:head>
    <title>Tournament Manager</title>
</svelte:head>

<div class="min-h-screen bg-base-200 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
        <div class="text-center mb-8">
            <a href="/" class="text-2xl font-bold text-primary">Tournament Manager</a>
        </div>
        {@render children()}
    </div>
</div>
