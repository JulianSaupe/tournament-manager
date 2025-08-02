<script lang="ts">
    import {formatDate, getTournamentsByStatus, mockTournaments, type Tournament} from '$lib/mockData';
    import {goto} from '$app/navigation';
    import {CircleCheckBig, Clipboard, SquarePen} from "lucide-svelte";

    // Group tournaments by status
    const tournamentsByStatus = getTournamentsByStatus(mockTournaments);

    // Get active tournaments
    const activeTournaments = tournamentsByStatus['active'] || [];

    // Calculate summary statistics
    const totalTournaments = mockTournaments.length;
    const activeTournamentCount = activeTournaments.length;
    const draftTournamentCount = (tournamentsByStatus['draft'] || []).length;
    const completedTournamentCount = (tournamentsByStatus['completed'] || []).length;

    // Define your types
    type TournamentStatus = 'draft' | 'active' | 'completed' | 'cancelled';

    type StatusConfig = {
        name: string;
        color: string;
        bgColor: string;
        icon: string;
    };

    const statusConfig: Record<TournamentStatus, StatusConfig> = {
        draft: {
            name: 'Draft',
            color: 'badge-neutral',
            bgColor: 'bg-neutral/10',
            icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />'
        },
        active: {
            name: 'Active',
            color: 'badge-success',
            bgColor: 'bg-success/10',
            icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />'
        },
        completed: {
            name: 'Completed',
            color: 'badge-info',
            bgColor: 'bg-info/10',
            icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />'
        },
        cancelled: {
            name: 'Cancelled',
            color: 'badge-error',
            bgColor: 'bg-error/10',
            icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />'
        }
    };

    // Helper function to get status config safely
    function getStatusConfig(status: string): StatusConfig {
        return statusConfig[status as TournamentStatus] || {
            name: status.charAt(0).toUpperCase() + status.slice(1),
            color: 'badge-ghost',
            bgColor: 'bg-base-200',
            icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />'
        };
    }

    // Sorting for table
    let sortField: keyof Tournament = 'name';
    let sortDirection: 'asc' | 'desc' = 'asc';

    // Sort tournaments
    function sortTournaments(tournaments: Tournament[]): Tournament[] {
        return [...tournaments].sort((a, b) => {
            let aValue = a[sortField];
            let bValue = b[sortField];

            // Handle date fields
            if (sortField === 'startDate' || sortField === 'endDate') {
                aValue = new Date(aValue as string).getTime();
                bValue = new Date(bValue as string).getTime();
            }

            if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1;
            if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1;
            return 0;
        });
    }

    // Set sort field and direction
    function setSorting(field: keyof Tournament) {
        if (sortField === field) {
            sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            sortField = field;
            sortDirection = 'asc';
        }
    }

    // Function to navigate to tournament creation page
    function navigateToCreate() {
        goto('/tournaments/create');
    }

    // Get sorted tournaments
    $: sortedTournaments = sortTournaments(mockTournaments);
</script>

<div class="w-full">
    <!-- Dashboard Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-bold mb-2">Dashboard</h1>
        <p class="text-base-content/70">Welcome to your tournament management dashboard</p>
    </div>

    <!-- Stats Overview -->
    <div class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4 mb-8">
        <div class="stat bg-base-100 shadow rounded-box p-3 md:p-6">
            <div class="stat-figure text-primary hidden sm:block">
                <Clipboard class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Total</div>
            <div class="stat-value text-lg sm:text-2xl md:text-3xl">{totalTournaments}</div>
        </div>

        <div class="stat bg-base-100 shadow rounded-box p-3 md:p-6">
            <div class="stat-figure text-primary hidden sm:block">
                <CircleCheckBig class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Active</div>
            <div class="stat-value text-lg sm:text-2xl md:text-3xl text-primary">{activeTournamentCount}</div>
        </div>

        <div class="stat bg-base-100 shadow rounded-box p-3 md:p-6">
            <div class="stat-figure text-neutral hidden sm:block">
                <SquarePen class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Draft</div>
            <div class="stat-value text-lg sm:text-2xl md:text-3xl text-neutral">{draftTournamentCount}</div>
        </div>

        <div class="stat bg-base-100 shadow rounded-box p-3 md:p-6">
            <div class="stat-figure text-primary hidden sm:block">
                <Clipboard class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Completed</div>
            <div class="stat-value text-lg sm:text-2xl md:text-3xl text-primary">{completedTournamentCount}</div>
        </div>
    </div>

    <!-- Active Tournaments Cards -->
    <div class="mb-10">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 mb-6">
            <h2 class="text-xl sm:text-2xl font-semibold flex items-center">
                <span class="badge badge-success mr-2">Active</span>
                Tournaments
            </h2>
            <button class="btn btn-primary btn-sm" on:click={navigateToCreate}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24"
                     stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
                <span class="hidden sm:inline">Create New</span>
                <span class="sm:hidden">New</span>
            </button>
        </div>

        {#if activeTournaments.length > 0}
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
                {#each activeTournaments as tournament}
                    <div class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow duration-300">
                        <div class="card-body p-4 md:p-6">
                            <div class="flex justify-between items-start">
                                <h3 class="card-title text-lg sm:text-xl">{tournament.name}</h3>
                                <div class="badge badge-success">Active</div>
                            </div>

                            <div class="divider my-2"></div>

                            <div class="mt-2 space-y-2 text-sm">
                                <div class="flex justify-between">
                                    <span class="font-medium">Start:</span>
                                    <span>{formatDate(tournament.startDate)}</span>
                                </div>
                                <div class="flex justify-between">
                                    <span class="font-medium">End:</span>
                                    <span>{formatDate(tournament.endDate)}</span>
                                </div>
                                <div class="flex justify-between">
                                    <span class="font-medium">Players:</span>
                                    <span>{tournament.playerCount}</span>
                                </div>
                            </div>

                            <div class="card-actions justify-end mt-4">
                                <button class="btn btn-xs sm:btn-sm btn-outline">View Details</button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="alert alert-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                     class="stroke-current shrink-0 w-5 h-5 sm:w-6 sm:h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span>No active tournaments found.</span>
            </div>
        {/if}
    </div>

    <!-- All Tournaments Table -->
    <div class="mb-8">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 mb-6">
            <h2 class="text-xl sm:text-2xl font-semibold">All Tournaments</h2>
            <div class="flex gap-2">
                <div class="dropdown dropdown-end">
                    <div tabindex="0" role="button" class="btn btn-sm btn-outline">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 sm:h-5 sm:w-5 mr-1" fill="none"
                             viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/>
                        </svg>
                        <span class="hidden sm:inline">Filter</span>
                    </div>
                    <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                        <li><a>All Statuses</a></li>
                        <li><a>Active Only</a></li>
                        <li><a>Draft Only</a></li>
                        <li><a>Completed Only</a></li>
                    </ul>
                </div>
                <button class="btn btn-sm btn-outline">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 sm:h-5 sm:w-5 mr-1" fill="none"
                         viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
                    </svg>
                    <span class="hidden sm:inline">Export</span>
                </button>
            </div>
        </div>

        <!-- Table for medium and larger screens -->
        <div class="hidden sm:block overflow-x-auto bg-base-100 rounded-box shadow">
            <table class="table table-zebra w-full">
                <thead>
                <tr>
                    <th class="cursor-pointer" on:click={() => setSorting('name')}>
                        Name
                        {#if sortField === 'name'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer" on:click={() => setSorting('status')}>
                        Status
                        {#if sortField === 'status'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer md:table-cell hidden" on:click={() => setSorting('startDate')}>
                        Start Date
                        {#if sortField === 'startDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer md:table-cell hidden" on:click={() => setSorting('endDate')}>
                        End Date
                        {#if sortField === 'endDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer" on:click={() => setSorting('playerCount')}>
                        Players
                        {#if sortField === 'playerCount'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th>Actions</th>
                </tr>
                </thead>
                <tbody>
                {#each sortedTournaments as tournament}
                    <tr>
                        <td>{tournament.name}</td>
                        <td>
                            <div class="badge {statusConfig[tournament.status].color}">
                                {statusConfig[tournament.status].name}
                            </div>
                        </td>
                        <td class="md:table-cell hidden">{formatDate(tournament.startDate)}</td>
                        <td class="md:table-cell hidden">{formatDate(tournament.endDate)}</td>
                        <td>{tournament.playerCount}</td>
                        <td>
                            <div class="flex gap-2">
                                <button class="btn btn-xs btn-outline">View</button>
                                {#if tournament.status === 'draft'}
                                    <button class="btn btn-xs btn-primary">Edit</button>
                                {/if}
                            </div>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>

        <!-- Card-based layout for small screens -->
        <div class="sm:hidden space-y-4">
            {#each sortedTournaments as tournament}
                <div class="card bg-base-100 shadow-md">
                    <div class="card-body p-4">
                        <div class="flex justify-between items-start">
                            <h3 class="card-title text-base">{tournament.name}</h3>
                            <div class="badge {statusConfig[tournament.status].color}">
                                {statusConfig[tournament.status].name}
                            </div>
                        </div>

                        <div class="text-sm space-y-1 mt-2">
                            <div class="flex justify-between">
                                <span class="font-medium">Start:</span>
                                <span>{formatDate(tournament.startDate)}</span>
                            </div>
                            <div class="flex justify-between">
                                <span class="font-medium">End:</span>
                                <span>{formatDate(tournament.endDate)}</span>
                            </div>
                            <div class="flex justify-between">
                                <span class="font-medium">Players:</span>
                                <span>{tournament.playerCount}</span>
                            </div>
                        </div>

                        <div class="card-actions justify-end mt-3">
                            <button class="btn btn-xs btn-outline">View</button>
                            {#if tournament.status === 'draft'}
                                <button class="btn btn-xs btn-primary">Edit</button>
                            {/if}
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>