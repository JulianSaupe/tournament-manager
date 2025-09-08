<script lang="ts">
    import {CircleCheckBig} from 'lucide-svelte';

    let {data} = $props();

    let name = $state('');
    let email = $state('');
    let errors: { name?: string; email?: string } = $state({});
    let submitted = $state(false);

    function validate(): boolean {
        errors = {};
        if (!name.trim()) {
            errors.name = 'Name is required';
        }
        if (!email.trim()) {
            errors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
            errors.email = 'Enter a valid email address';
        }
        return Object.keys(errors).length === 0;
    }

    function handleSubmit(event: Event) {
        event.preventDefault();
        if (!validate()) return;
        submitted = true;
    }

    function resetForm() {
        name = '';
        email = '';
        errors = {};
        submitted = false;
    }
</script>

<div class="min-h-[calc(100vh-0rem)]">
    <div class="container mx-auto max-w-xl px-4 py-12 flex min-h-[calc(100vh-0rem)] items-center justify-center">
        <div class="w-full">
            <div class="mb-6 text-center">
                <h1 class="text-3xl font-semibold">Tournament Signup</h1>
            </div>

            {#if submitted}
                <div class="card bg-base-100 shadow-lg">
                    <div class="card-body">
                        <div class="flex items-center gap-3">
                            <CircleCheckBig class="h-6 w-6 text-success"/>
                            <h2 class="card-title">You're signed up!</h2>
                        </div>
                        <p class="text-base-content/70">
                            Thanks, {name}. You have entered tournament
                            <span class="badge badge-outline">{data.id}</span>.
                        </p>
                        <div class="mt-4 card-actions justify-center">
                            <button class="btn btn-ghost" onclick={resetForm}>Signup another player</button>
                        </div>
                    </div>
                </div>
            {:else}
                <form onsubmit={handleSubmit} class="space-y-6">
                    <div class="card bg-base-100 shadow-lg">
                        <div class="card-body">
                            <h2 class="card-title">Player Information</h2>
                            <div class="form-control w-full">
                                <label for="name" class="label">
                                    <span class="label-text">Name</span>
                                </label>
                                <input
                                        id="name"
                                        type="text"
                                        value={name}
                                        oninput={(e) => (name = e.currentTarget.value)}
                                        class={`input-bordered input w-full ${errors.name ? 'input-error' : ''}`}
                                        placeholder="Enter your name"
                                />
                                {#if errors.name}
                                    <label class="label" for="name">
                                        <span class="label-text-alt text-error">{errors.name}</span>
                                    </label>
                                {/if}
                            </div>

                            <div class="form-control w-full">
                                <label for="email" class="label">
                                    <span class="label-text">Email</span>
                                </label>
                                <input
                                        id="email"
                                        type="email"
                                        value={email}
                                        oninput={(e) => (email = e.currentTarget.value)}
                                        class={`input-bordered input w-full ${errors.email ? 'input-error' : ''}`}
                                        placeholder="Enter your email"
                                />
                                {#if errors.email}
                                    <label class="label" for="email">
                                        <span class="label-text-alt text-error">{errors.email}</span>
                                    </label>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <div class="card-actions justify-center">
                        <button type="submit" class="btn btn-primary w-full">Sign Up</button>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</div>
