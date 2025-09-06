<script lang="ts">
    import {Calendar, Check, Funnel, Info, User} from "lucide-svelte";
    import {type Tournament} from "$lib/types/tournament/tournament";
    import {statusConfig} from "$lib/types/tournament/statusConfig";
    import moment from "moment/moment";
    import {type Qualifying, QualifyingFilter} from "$lib/types/tournament/qualifying";

    let {data} = $props();

    let tournament: Tournament = $derived(data.tournament);
    let qualifying: Qualifying = $derived(data.qualifying);

    let sortField: 'position' | 'name' | 'signupDate' | 'time' = $state('name');
    let sortDirection: 'asc' | 'desc' = $state('asc');
    let statusFilter: QualifyingFilter = $state(QualifyingFilter.ALL);

    function setSorting(field: 'position' | 'name' | 'signupDate' | 'time') {
        if (sortField === field) {
            sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            sortField = field;
            sortDirection = 'asc';
        }
    }
</script>

<!-- Header -->
<div class="mb-8">
    <h1 class="mb-2 text-3xl font-bold">{tournament.name}</h1>
    <p class="text-base-content/70">Manage your tournament here</p>
</div>

<!-- Stats Overview -->
<div class="mb-8 grid grid-cols-2 gap-3 sm:grid-cols-2 md:gap-4 lg:grid-cols-4">
    <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
        <div class="stat-figure hidden text-primary sm:block">
            <Info class="h-6 w-6 md:h-8 md:w-8"/>
        </div>
        <div class="stat-title text-xs sm:text-sm">Status</div>
        <div class="stat-value text-lg sm:text-2xl md:text-3xl">{statusConfig[tournament.status].name}</div>
    </div>

    <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
        <div class="stat-figure hidden text-primary sm:block">
            <User class="h-6 w-6 md:h-8 md:w-8"/>
        </div>
        <div class="stat-title text-xs sm:text-sm">Players</div>
        <div class="stat-value text-lg text-primary sm:text-2xl md:text-3xl">
            {tournament.playerCount}
        </div>
    </div>

    <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
        <div class="stat-figure hidden text-neutral sm:block">
            <Calendar class="h-6 w-6 md:h-8 md:w-8"/>
        </div>
        <div class="stat-title text-xs sm:text-sm">Start</div>
        <div class="stat-value text-lg text-neutral sm:text-2xl md:text-3xl">
            {moment(tournament.startDate).format('MMM D, YYYY')}
        </div>
    </div>

    <div class="stat rounded-box bg-base-100 p-3 shadow md:p-6">
        <div class="stat-figure hidden text-primary sm:block">
            <Calendar class="h-6 w-6 md:h-8 md:w-8"/>
        </div>
        <div class="stat-title text-xs sm:text-sm">End</div>
        <div class="stat-value text-lg text-primary sm:text-2xl md:text-3xl">
            {moment(tournament.endDate).format('MMM D, YYYY')}
        </div>
    </div>
</div>

<div class="mb-8">
    <div class="mb-6 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
        <h2 class="flex items-center text-xl font-semibold sm:text-2xl">
            Qualifying
        </h2>
        <div class="flex gap-2 mb-2">
            <div class="dropdown dropdown-end">
                <div tabindex="0" role="button" class="btn btn-outline btn-sm">
                    <Funnel/>
                    <span class="hidden sm:inline">Filter</span>
                </div>
                <ul class="dropdown-content menu z-[1] w-52 rounded-box bg-base-100 p-2 shadow">
                    <li>
                        <button onclick={() => (statusFilter = QualifyingFilter.ALL)} class="justify-between">
                            All

                            {#if statusFilter === QualifyingFilter.ALL}
                                <Check/>
                            {/if}
                        </button>
                    </li>
                    <li>
                        <button onclick={() => (statusFilter = QualifyingFilter.QUALIFIED)} class="justify-between">
                            Qualified

                            {#if statusFilter === QualifyingFilter.QUALIFIED}
                                <Check/>
                            {/if}
                        </button>
                    </li>
                    <li>
                        <button onclick={() => (statusFilter = QualifyingFilter.UNQUALIFIED)} class="justify-between">
                            Unqualified

                            {#if statusFilter === QualifyingFilter.UNQUALIFIED}
                                <Check/>
                            {/if}
                        </button>
                    </li>
                </ul>
            </div>
        </div>
    </div>

    {#if qualifying?.players.length > 0}
        <!-- Table -->
        <div class="hidden overflow-x-auto rounded-box bg-base-100 shadow sm:block">
            <table class="table w-full table-zebra">
                <thead>
                <tr>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('position')}>
                        Position
                        {#if sortField === 'position'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('name')}>
                        Name
                        {#if sortField === 'name'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="cursor-pointer select-none" onclick={() => setSorting('time')}>
                        Best Time
                        {#if sortField === 'time'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th class="hidden cursor-pointer md:table-cell select-none"
                        onclick={() => setSorting('signupDate')}>
                        Signup Date
                        {#if sortField === 'signupDate'}
                            <span>{sortDirection === 'asc' ? '↑' : '↓'}</span>
                        {/if}
                    </th>
                    <th>Actions</th>
                </tr>
                </thead>
                <tbody>
                {#each qualifying.players as player}
                    <tr>
                        <td>{player.position}</td>
                        <td>{player.name}</td>
                        <td class="hidden md:table-cell">{player.time === -1 ? '-' : player.time}</td>
                        <td class="hidden md:table-cell">{moment(player.signupDate).format('MMM D, YYYY')}</td>
                        <td>
                            <div class="flex gap-2">
                                <button class="btn btn-outline btn-xs">Edit</button>
                            </div>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {:else}
        <div class="alert alert-info">
            <Info/>
            <span>No qualifying players found.</span>
        </div>
    {/if}
</div>

<h2 class="flex items-center text-xl font-semibold sm:text-2xl">
    Players
</h2>