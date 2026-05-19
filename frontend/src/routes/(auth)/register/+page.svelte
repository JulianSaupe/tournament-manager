<script lang="ts">
    import type { ActionData } from './$types';

    let { form }: { form: ActionData } = $props();
</script>

<div class="card bg-base-100 shadow-xl">
    <div class="card-body">
        <h2 class="card-title text-2xl mb-2">Create account</h2>

        {#if form?.error}
            <div class="alert alert-error mb-4">
                <span>{form.error}</span>
            </div>
        {/if}

        <form method="POST" class="flex flex-col gap-4">
            <label class="form-control">
                <div class="label"><span class="label-text">Username</span></div>
                <input
                    type="text"
                    name="username"
                    class="input input-bordered"
                    value={form?.username ?? ''}
                    required
                    autocomplete="username"
                />
            </label>

            <label class="form-control">
                <div class="label"><span class="label-text">Email</span></div>
                <input
                    type="email"
                    name="email"
                    class="input input-bordered"
                    value={form?.email ?? ''}
                    required
                    autocomplete="email"
                />
            </label>

            <label class="form-control">
                <div class="label"><span class="label-text">Password</span></div>
                <input
                    type="password"
                    name="password"
                    class="input input-bordered"
                    required
                    autocomplete="new-password"
                />
            </label>

            <button type="submit" class="btn btn-primary mt-2">Create account</button>
        </form>

        <div class="divider">or</div>

        <p class="text-center text-sm">
            Already have an account?
            <a href="/login" class="link link-primary">Sign in</a>
        </p>
    </div>
</div>
