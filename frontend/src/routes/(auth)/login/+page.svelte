<script lang="ts">
    import type {ActionData} from './$types';

    let {form}: { form: ActionData } = $props();
</script>

<div class="card bg-base-100 shadow-xl">
    <div class="card-body">
        <h2 class="card-title text-2xl mb-2">Sign in</h2>

        {#if form?.error}
            <div class="alert alert-error mb-4">
                <span>{form.error}</span>
            </div>
        {/if}

        <form class="flex flex-col gap-4" method="POST">
            <label class="form-control">
                <div class="label"><span class="label-text">Email</span></div>
                <input
                        autocomplete="email"
                        class="input input-bordered"
                        name="email"
                        required
                        type="email"
                        value={form?.email ?? ''}
                />
            </label>

            <label class="form-control">
                <div class="label"><span class="label-text">Password</span></div>
                <input
                        autocomplete="current-password"
                        class="input input-bordered"
                        name="password"
                        required
                        type="password"
                />
            </label>

            <button class="btn btn-primary mt-2" type="submit">Sign in</button>
        </form>

        <div class="divider">or</div>

        <p class="text-center text-sm">
            Don't have an account?
            <a class="link link-primary" href="/register">Register</a>
        </p>
    </div>
</div>
