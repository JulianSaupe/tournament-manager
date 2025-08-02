<script lang="ts">
    import {goto} from '$app/navigation';
    import {mockTournaments} from '$lib/mockData';

    // Form data
    let formData = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: 0
    };

    // Form validation
    let errors = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: ''
    };

    // Form validation function
    function validateForm() {
        let isValid = true;

        // Reset errors
        errors = {
            name: '',
            startDate: '',
            endDate: '',
            playerCount: ''
        };

        // Validate name
        if (!formData.name.trim()) {
            errors.name = 'Tournament name is required';
            isValid = false;
        }

        // Validate start date
        if (!formData.startDate) {
            errors.startDate = 'Start date is required';
            isValid = false;
        }

        // Validate end date
        if (!formData.endDate) {
            errors.endDate = 'End date is required';
            isValid = false;
        } else if (formData.startDate && formData.endDate && new Date(formData.endDate) < new Date(formData.startDate)) {
            errors.endDate = 'End date must be after start date';
            isValid = false;
        }

        // Validate player count
        if (!formData.playerCount || formData.playerCount <= 0) {
            errors.playerCount = 'Number of players must be greater than 0';
            isValid = false;
        }

        return isValid;
    }

    // Form submission
    function handleSubmit(event: Event) {
        event.preventDefault();
        if (validateForm()) {
            // In a real application, this would send data to the backend
            // For now, we'll just log it and navigate back to the home page

            // Create a new tournament object
            const newTournament = {
                id: (mockTournaments.length + 1).toString(),
                name: formData.name,
                startDate: new Date(formData.startDate).toISOString(),
                endDate: new Date(formData.endDate).toISOString(),
                playerCount: formData.playerCount,
                status: 'draft'
            };

            console.log('New tournament created:', newTournament);

            // Navigate back to the home page
            goto('/');
        }
    }

    // Cancel and go back to home page
    function handleCancel() {
        goto('/');
    }
</script>

<div class="container mx-auto max-w-3xl">
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h1 class="card-title text-2xl mb-6">Create New Tournament</h1>

            <form onsubmit={handleSubmit} class="space-y-6">
                <!-- Tournament Name -->
                <div class="form-control w-full">
                    <label for="name" class="label">
                        <span class="label-text">Tournament Name</span>
                    </label>
                    <input
                            type="text"
                            id="name"
                            bind:value={formData.name}
                            class="input input-bordered w-full {errors.name ? 'input-error' : ''}"
                            placeholder="Enter tournament name"
                    />
                    {#if errors.name}
                        <label class="label">
                            <span class="label-text-alt text-error">{errors.name}</span>
                        </label>
                    {/if}
                </div>

                <!-- Start Date -->
                <div class="form-control w-full">
                    <label for="startDate" class="label">
                        <span class="label-text">Start Date</span>
                    </label>
                    <input
                            type="date"
                            id="startDate"
                            bind:value={formData.startDate}
                            class="input input-bordered w-full {errors.startDate ? 'input-error' : ''}"
                    />
                    {#if errors.startDate}
                        <label class="label">
                            <span class="label-text-alt text-error">{errors.startDate}</span>
                        </label>
                    {/if}
                </div>

                <!-- End Date -->
                <div class="form-control w-full">
                    <label for="endDate" class="label">
                        <span class="label-text">End Date</span>
                    </label>
                    <input
                            type="date"
                            id="endDate"
                            bind:value={formData.endDate}
                            class="input input-bordered w-full {errors.endDate ? 'input-error' : ''}"
                    />
                    {#if errors.endDate}
                        <label class="label">
                            <span class="label-text-alt text-error">{errors.endDate}</span>
                        </label>
                    {/if}
                </div>

                <!-- Number of Players -->
                <div class="form-control w-full">
                    <label for="playerCount" class="label">
                        <span class="label-text">Number of Players</span>
                    </label>
                    <input
                            type="number"
                            id="playerCount"
                            bind:value={formData.playerCount}
                            min="1"
                            class="input input-bordered w-full {errors.playerCount ? 'input-error' : ''}"
                            placeholder="Enter number of players"
                    />
                    {#if errors.playerCount}
                        <label class="label">
                            <span class="label-text-alt text-error">{errors.playerCount}</span>
                        </label>
                    {/if}
                </div>

                <!-- Form Actions -->
                <div class="card-actions justify-end mt-8">
                    <button type="button" class="btn btn-ghost" onclick={handleCancel}>
                        Cancel
                    </button>
                    <button type="submit" class="btn btn-primary">
                        Create Tournament
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>