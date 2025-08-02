<script lang="ts">
    import {goto} from '$app/navigation';
    import {mockTournaments} from '$lib/mockData';

    // Form data
    let formData = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: 0,
        groupPhase: false,
        allowPartiallyFilledGroups: false,
        groupSize: 4,
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
                status: 'draft',
                groupPhase: formData.groupPhase,
                groupSize: formData.groupPhase ? formData.groupSize : null,
                allowPartiallyFilledGroups: formData.groupPhase ? formData.allowPartiallyFilledGroups : false
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

<div class="container mx-auto">
    <h1 class="text-2xl mb-6">Create New Tournament</h1>

    <form onsubmit={handleSubmit} class="space-y-6">
        <div class="card bg-base-100 shadow-sm">
            <div class="card-body">
                <div class="card-title">
                    <h2 class="card-title">Tournament Details</h2>
                </div>
                <div class="grid grid-cols-2 gap-6">
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
                            <label class="label" for="name">
                                <span class="label-text-alt text-error">{errors.name}</span>
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
                            <label class="label" for="playerCount">
                                <span class="label-text-alt text-error">{errors.playerCount}</span>
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
                            <label class="label" for="startDate">
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
                            <label class="label" for="endDate">
                                <span class="label-text-alt text-error">{errors.endDate}</span>
                            </label>
                        {/if}
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-sm">
            <div class="card-body">
                <div class="card-title mb-4">
                    <h2 class="card-title">Rounds</h2>
                </div>
                
                <!-- Group Phase Toggle -->
                <div class="form-control bg-base-200 p-4 rounded-lg mb-4">
                    <div class="flex items-left mb-1">
                        <span class="label-text font-medium text-base">Enable Group Phase</span>
                        <input type="checkbox" class="toggle toggle-primary ms-3" bind:checked={formData.groupPhase}/>
                    </div>
                    <p class="text-xs text-base-content/70">Organize players into groups for the initial tournament
                        phase</p>
                </div>
                
                <!-- Group Phase Settings (only visible when group phase is enabled) -->
                {#if formData.groupPhase}
                    <div class="bg-base-100 border border-base-300 rounded-lg p-4 ml-2 sm:ml-4">
                        <h3 class="font-medium mb-4">Group Phase Settings</h3>
                        
                        <div class="grid gap-4">
                            <!-- Group Size -->
                            <div class="form-control w-full">
                                <label for="groupSize" class="label">
                                    <span class="label-text">Group Size</span>
                                </label>
                                <input 
                                    type="number" 
                                    id="groupSize" 
                                    class="input input-bordered w-full" 
                                    placeholder="Enter group size"
                                    min="2"
                                    bind:value={formData.groupSize}
                                />
                                <label class="label" for="groupSize">
                                    <span class="label-text-alt text-base-content/70">Recommended: 4-6 players per group</span>
                                </label>
                            </div>
                            
                            <!-- Allow Partially Filled Groups -->
                            <div class="form-control">
                                <label class="label cursor-pointer justify-between flex-wrap sm:flex-nowrap">
                                    <span class="label-text w-full sm:w-auto mb-2 sm:mb-0">Allow Partially Filled Groups</span>
                                    <input 
                                        type="checkbox" 
                                        class="checkbox" 
                                        bind:checked={formData.allowPartiallyFilledGroups}
                                    />
                                </label>
                                <label class="label mt-[-8px]" for="allowPartiallyFilledGroups">
                                    <span class="label-text-alt text-base-content/70">If disabled, all groups will have the same number of players</span>
                                </label>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
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