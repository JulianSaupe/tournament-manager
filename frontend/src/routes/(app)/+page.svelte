<script lang="ts">
    import {goto} from '$app/navigation';
    import {Check, CircleCheckBig, Clipboard, Funnel, Info, SquarePen} from 'lucide-svelte';
    import moment from 'moment';
    import type {Tournament} from "$lib/types/tournament/tournament";
    import {TournamentStatus} from "$lib/types/tournament/tournament";
    import {statusConfig} from "$lib/types/tournament/statusConfig";

    let {data} = $props();

    const formatDate = (dateString: string): string => moment(dateString).format('MMM D, YYYY');

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

    const tournaments: Tournament[] = Array.isArray(data.tournaments) ? data.tournaments : [];
    const tournamentsByStatus = getTournamentsByStatus(tournaments);
    const activeTournaments = tournamentsByStatus[TournamentStatus.ACTIVE] || [];
    const totalTournaments = tournaments.length;
    const activeTournamentCount = activeTournaments.length;
    const draftTournamentCount = (tournamentsByStatus[TournamentStatus.DRAFT] || []).length;
    const completedTournamentCount = (tournamentsByStatus[TournamentStatus.COMPLETED] || []).length;

    let sortField: keyof Tournament = $state('name');
    let sortDirection: 'asc' | 'desc' = $state('asc');
    let statusFilter: TournamentStatus | null = $state(null);
    let searchTerm: string = $state('');

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

        if (statusFilter !== null) {
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

    let filteredTournaments = $derived(filterTournaments(tournaments));
    let sortedTournaments = $derived(sortTournaments(filteredTournaments));
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
            <button class="btn btn-sm btn-primary" onclick={navigateToCreate}>
                <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="mr-1 h-5 w-5"
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
                <Info/>
                <span>No active tournaments found.</span>
            </div>
        {/if}
    </div>

    <!-- All Tournaments Table -->
    <div class="mb-8">
        <div class="mb-6 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
            <h2 class="text-xl font-semibold sm:text-2xl">All Tournaments</h2>
            <div class="flex gap-2">
                <input class="input input-sm input-bordered" placeholder="Search..." bind:value={searchTerm}/>
                <div class="dropdown dropdown-end">
                    <div tabindex="0" role="button" class="btn btn-outline btn-sm">
                        <Funnel/>
                        <span class="hidden sm:inline">Filter</span>
                    </div>
                    <ul class="dropdown-content menu z-[1] w-52 rounded-box bg-base-100 p-2 shadow">
                        <li>
                            <button onclick={() => (statusFilter = null)} class="justify-between">
                                All

                                {#if statusFilter === null}
                                    <Check/>
                                {/if}
                            </button>
                        </li>
                        <li>
                            <button onclick={() => (statusFilter = TournamentStatus.ACTIVE)} class="justify-between">
                                Active

                                {#if statusFilter === TournamentStatus.ACTIVE}
                                    <Check/>
                                {/if}
                            </button>
                        </li>
                        <li>
                            <button onclick={() => (statusFilter = TournamentStatus.DRAFT)} class="justify-between">
                                Draft

                                {#if statusFilter === TournamentStatus.DRAFT}
                                    <Check/>
                                {/if}
                            </button>
                        </li>
                        <li>
                            <button onclick={() => (statusFilter = TournamentStatus.COMPLETED)} class="justify-between">
                                Completed

                                {#if statusFilter === TournamentStatus.COMPLETED}
                                    <Check/>
                                {/if}
                            </button>
                        </li>
                        <li>
                            <button onclick={() => (statusFilter = TournamentStatus.CANCELLED)} class="justify-between">
                                Cancelled

                                {#if statusFilter === TournamentStatus.CANCELLED}
                                    <Check/>
                                {/if}
                            </button>
                        </li>
                    </ul>
                </div>
            </div>
        </div>

        <!-- Table for medium and larger screens -->
        <div class="hidden overflow-x-auto rounded-box bg-base-100 shadow sm:block">
            <table class="table w-full table-zebra">
                <thead>
                <tr>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('name')}>
                        Name
                        {#if sortField === 'name'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('status')}>
                        Status
                        {#if sortField === 'status'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="hidden cursor-pointer md:table-cell select-none" onclick={() => setSorting('startDate')}>
                        Start Date
                        {#if sortField === 'startDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="hidden cursor-pointer md:table-cell select-none" onclick={() => setSorting('endDate')}>
                        End Date
                        {#if sortField === 'endDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('playerCount')}>
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
                                <button class="btn btn-outline btn-xs"
                                        onclick={() => goto('/tournaments/' + tournament.id)}>View
                                </button>
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
                            <button class="btn btn-outline btn-xs"
                                    onclick={() => goto('/tournaments/' + tournament.id)}>View
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>
