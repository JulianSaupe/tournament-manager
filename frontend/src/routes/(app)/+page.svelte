<script lang="ts">
    import {goto} from '$app/navigation';
    import {CircleCheckBig, Clipboard, SquarePen} from 'lucide-svelte';
    import type {IndexTournament} from './+page.server';
    import type {PageData} from './$types';

    export let data: PageData;

    type TournamentStatus = 'draft' | 'active' | 'completed' | 'cancelled';

    type Tournament = {
        id: string;
        name: string;
        description: string;
        startDate: string;
        endDate: string;
        status: TournamentStatus;
        playerCount?: number;
    };

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
            icon:
                '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />'
        },
        active: {
            name: 'Active',
            color: 'badge-success',
            bgColor: 'bg-success/10',
            icon:
                '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />'
        },
        completed: {
            name: 'Completed',
            color: 'badge-info',
            bgColor: 'bg-info/10',
            icon:
                '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />'
        },
        cancelled: {
            name: 'Cancelled',
            color: 'badge-error',
            bgColor: 'bg-error/10',
            icon:
                '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />'
        }
    };

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });
    }

    function getTournamentsByStatus(
        tournaments: Tournament[]
    ): Record<string, Tournament[]> {
        return tournaments.reduce(
            (acc, tournament) => {
                if (!acc[tournament.status]) {
                    acc[tournament.status] = [];
                }
                acc[tournament.status].push(tournament);
                return acc;
            },
            {} as Record<string, Tournament[]>
        );
    }

    function normalizeStatus(status: string): TournamentStatus {
        const s = status.toLowerCase();
        if (s === 'draft' || s === 'active' || s === 'completed' || s === 'cancelled') return s;
        return 'draft';
    }

    const tournaments: Tournament[] = (Array.isArray(data.tournaments) ? data.tournaments : []).map((t: IndexTournament) => ({
        ...t,
        status: normalizeStatus(t.status)
    }));

    const tournamentsByStatus = getTournamentsByStatus(tournaments);
    const activeTournaments = tournamentsByStatus['active'] || [];
    const totalTournaments = tournaments.length;
    const activeTournamentCount = activeTournaments.length;
    const draftTournamentCount = (tournamentsByStatus['draft'] || []).length;
    const completedTournamentCount = (tournamentsByStatus['completed'] || []).length;

    let sortField: keyof Tournament = 'name';
    let sortDirection: 'asc' | 'desc' = 'asc';
    let statusFilter: 'all' | 'draft' | 'active' | 'completed' | 'cancelled' = 'all';
    let searchTerm: string = '';

    function sortTournaments(tournaments: Tournament[]): Tournament[] {
        return [...tournaments].sort((a, b) => {
            let aValue: any = a[sortField];
            let bValue: any = b[sortField];

            if (sortField === 'startDate' || sortField === 'endDate') {
                aValue = new Date(aValue as string).getTime();
                bValue = new Date(bValue as string).getTime();
            }
            if (sortField === 'playerCount') {
                aValue = a.playerCount ?? 0;
                bValue = b.playerCount ?? 0;
            }

            if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1;
            if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1;
            return 0;
        });
    }

    function setSorting(field: keyof Tournament) {
        if (sortField === field) {
            sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            sortField = field;
            sortDirection = 'asc';
        }
    }

    function filterTournaments(tournaments: Tournament[]): Tournament[] {
        let result = tournaments;
        if (statusFilter !== 'all') {
            result = result.filter((t) => t.status === statusFilter);
        }
        const q = searchTerm.trim().toLowerCase();
        if (q) {
            result = result.filter((t) =>
                t.name.toLowerCase().includes(q) || (t.description || '').toLowerCase().includes(q)
            );
        }
        return result;
    }

    function navigateToCreate() {
        goto('/tournaments/create');
    }

    $: filteredTournaments = filterTournaments(tournaments);
    $: sortedTournaments = sortTournaments(filteredTournaments);
    const error = data.error;
</script>

<div class="w-full">
    <!-- Dashboard Header -->
    <div class="mb-8">
        <h1 class="mb-2 text-3xl font-bold">Dashboard</h1>
        <p class="text-base-content/70">Welcome to your tournament management dashboard</p>
    </div>

    <!-- Stats Overview -->
    <div class="mb-8 grid grid-cols-2 gap-3 sm:grid-cols-2 md:gap-4 lg:grid-cols-4">
        <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
            <div class="stat-figure hidden text-primary sm:block">
                <Clipboard class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Total</div>
            <div class="stat-value text-lg sm:text-2xl md:text-3xl">{totalTournaments}</div>
        </div>

        <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
            <div class="stat-figure hidden text-primary sm:block">
                <CircleCheckBig class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Active</div>
            <div class="stat-value text-lg text-primary sm:text-2xl md:text-3xl">
                {activeTournamentCount}
            </div>
        </div>

        <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
            <div class="stat-figure hidden text-neutral sm:block">
                <SquarePen class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Draft</div>
            <div class="stat-value text-lg text-neutral sm:text-2xl md:text-3xl">
                {draftTournamentCount}
            </div>
        </div>

        <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
            <div class="stat-figure hidden text-primary sm:block">
                <Clipboard class="h-6 w-6 md:h-8 md:w-8"/>
            </div>
            <div class="stat-title text-xs sm:text-sm">Completed</div>
            <div class="stat-value text-lg text-primary sm:text-2xl md:text-3xl">
                {completedTournamentCount}
            </div>
        </div>
    </div>

    <!-- Active Tournaments Cards -->
    <div class="mb-10">
        <div class="mb-6 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
            <h2 class="flex items-center text-xl font-semibold sm:text-2xl">
                <span class="mr-2 badge badge-success">Active</span>
                Tournaments
            </h2>
            <button class="btn btn-sm btn-primary" on:click={navigateToCreate}>
                <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="mr-1 h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                >
                    <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 4v16m8-8H4"
                    />
                </svg>
                <span class="hidden sm:inline">Create New</span>
                <span class="sm:hidden">New</span>
            </button>
        </div>

        {#if activeTournaments.length > 0}
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6 lg:grid-cols-3">
                {#each activeTournaments as tournament}
                    <div class="card bg-base-100 shadow-xl transition-shadow duration-300 hover:shadow-2xl">
                        <div class="card-body p-4 md:p-6">
                            <div class="flex items-start justify-between">
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
                                    <span>{tournament.playerCount ?? 0}</span>
                                </div>
                            </div>

                            <div class="mt-4 card-actions justify-end">
                                <button class="btn btn-outline btn-xs sm:btn-sm">View Details</button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="alert alert-info">
                <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        class="h-5 w-5 shrink-0 stroke-current sm:h-6 sm:w-6"
                >
                    <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                </svg>
                <span>No active tournaments found.</span>
            </div>
        {/if}
    </div>

    <!-- All Tournaments Table -->
    <div class="mb-8">
        <div class="mb-6 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
            <h2 class="text-xl font-semibold sm:text-2xl">All Tournaments</h2>
            <div class="flex gap-2">
                <input class="input input-sm input-bordered" placeholder="Search..." bind:value={searchTerm} />
                <div class="dropdown dropdown-end">
                    <div tabindex="0" role="button" class="btn btn-outline btn-sm">
                        <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="mr-1 h-4 w-4 sm:h-5 sm:w-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                        >
                            <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
                            />
                        </svg>
                        <span class="hidden sm:inline">Filter</span>
                    </div>
                    <ul
                            tabindex="0"
                            class="dropdown-content menu z-[1] w-52 rounded-box bg-base-100 p-2 shadow"
                    >
                        <li><a href="#" on:click|preventDefault={() => (statusFilter = 'all')}>All Statuses</a></li>
                        <li><a href="#" on:click|preventDefault={() => (statusFilter = 'active')}>Active Only</a></li>
                        <li><a href="#" on:click|preventDefault={() => (statusFilter = 'draft')}>Draft Only</a></li>
                        <li><a href="#" on:click|preventDefault={() => (statusFilter = 'completed')}>Completed Only</a></li>
                        <li><a href="#" on:click|preventDefault={() => (statusFilter = 'cancelled')}>Cancelled Only</a></li>
                    </ul>
                </div>
                <button class="btn btn-outline btn-sm">
                    <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="mr-1 h-4 w-4 sm:h-5 sm:w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                    >
                        <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                        />
                    </svg>
                    <span class="hidden sm:inline">Export</span>
                </button>
            </div>
        </div>

        <!-- Table for medium and larger screens -->
        <div class="hidden overflow-x-auto rounded-box bg-base-100 shadow sm:block">
            <table class="table w-full table-zebra">
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
                    <th
                            class="hidden cursor-pointer md:table-cell"
                            on:click={() => setSorting('startDate')}
                    >
                        Start Date
                        {#if sortField === 'startDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="hidden cursor-pointer md:table-cell" on:click={() => setSorting('endDate')}>
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
                        <td class="hidden md:table-cell">{formatDate(tournament.startDate)}</td>
                        <td class="hidden md:table-cell">{formatDate(tournament.endDate)}</td>
                        <td>{tournament.playerCount ?? 0}</td>
                        <td>
                            <div class="flex gap-2">
                                <button class="btn btn-outline btn-xs">View</button>
                            </div>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>

        <!-- Card-based layout for small screens -->
        <div class="space-y-4 sm:hidden">
            {#each sortedTournaments as tournament}
                <div class="card bg-base-100 shadow-md">
                    <div class="card-body p-4">
                        <div class="flex items-start justify-between">
                            <h3 class="card-title text-base">{tournament.name}</h3>
                            <div class="badge {statusConfig[tournament.status].color}">
                                {statusConfig[tournament.status].name}
                            </div>
                        </div>

                        <div class="mt-2 space-y-1 text-sm">
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
                                <span>{tournament.playerCount ?? 0}</span>
                            </div>
                        </div>

                        <div class="mt-3 card-actions justify-end">
                            <button class="btn btn-outline btn-xs">View</button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>
