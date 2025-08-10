<script lang="ts">
	import type { Round } from '$lib/types/tournament';
	import { tournamentForm } from '$lib/stores/tournamentForm';

	let rounds: Round[] = $state([
		{
			name: 'Round 1',
			groupCount: 1,
			playersPerGroup: 128,
			matchesPerGroup: 1,
			advancingPlayersPerGroup: 1,
			concurrentGroups: 1
		}
	]);
	const totalRounds = $derived(rounds.length || 0);

	const firstRoundGroupCount = $derived(
		Math.ceil($tournamentForm.playerCount / rounds[0]?.playersPerGroup || 1)
	);

	$effect(() => {
		if (rounds[0]) {
			rounds[0].groupCount = firstRoundGroupCount;
		}
	});

	function handleAddRound(): void {
		const lastRound = rounds[rounds.length - 1];

		rounds = [
			...rounds,
			{
				name: `Round ${rounds.length + 1}`,
				groupCount: 1,
				playersPerGroup: lastRound.groupCount * lastRound.advancingPlayersPerGroup,
				matchesPerGroup: 1,
				advancingPlayersPerGroup: 1,
				concurrentGroups: 1
			}
		];
	}

	function handleRemoveRound(index: number): void {
		rounds = rounds.filter((_, i) => i !== index);
	}

	function updateRoundField(index: number, field: keyof Round, value: number | null): void {
		if (index < 0 || index >= rounds.length) return;

		const updatedRounds = [...rounds];
		updatedRounds[index] = {
			...updatedRounds[index],
			[field]: value
		};

		if (field === 'playersPerGroup') {
			const groupCount = Math.ceil(
				$tournamentForm.playerCount / updatedRounds[index].playersPerGroup
			);

			console.log(groupCount === Infinity, groupCount);

			if (
				!isNaN(groupCount) &&
				groupCount !== Infinity &&
				updatedRounds[index].playersPerGroup > 0
			) {
				updatedRounds[index].groupCount = groupCount;
			}
		}
		rounds = updatedRounds;
	}
</script>

<div class="card bg-base-100 shadow-sm">
	<div class="card-body">
		<div class="mb-4 card-title flex items-center justify-between">
			<h2 class="card-title">Tournament Structure</h2>
			<span class="badge badge-primary">{totalRounds} Rounds</span>
		</div>

		<!-- Rounds Configuration -->
		<div class="space-y-6">
			{#each rounds as round, index}
				<div class="rounded-lg border border-base-300 bg-base-100 p-4">
					<div class="mb-4 flex items-center justify-between">
						<h4 class="font-medium">{round.name}</h4>
						{#if index > 0}
							<button
								type="button"
								class="btn btn-outline btn-sm btn-error"
								onclick={() => handleRemoveRound(index)}
							>
								Remove
							</button>
						{/if}
					</div>

					<!-- Round configuration -->
					<div class="flex flex-col space-y-4">
						<!-- Group information display -->
						<div class="flex items-center justify-between rounded-lg bg-base-200 p-3">
							<div class="flex items-center">
								<div class="mr-2 badge badge-primary">{round.groupCount}</div>
								<span class="text-sm font-medium">Groups</span>
							</div>
							<div class="text-sm text-base-content/70">
								Total: {round.groupCount * (round.playersPerGroup || 0) || 0} players
							</div>
						</div>

						<div class="flex flex-wrap gap-3">
							<!-- Players per Group -->
							<div class="form-control min-w-[150px] flex-1">
								<label class="label" for="playersPerGroup">
									<span class="label-text">Players per Group</span>
								</label>
								<input
									type="number"
									class="input-bordered input w-full"
									min="2"
									value={round.playersPerGroup || ''}
									oninput={(e) => {
										const value = e.currentTarget.value;
										if (value === '') {
											// Allow empty value during editing
											updateRoundField(index, 'playersPerGroup', null);
										} else {
											updateRoundField(index, 'playersPerGroup', parseInt(value));
										}
									}}
									onblur={(e) => {
										// Apply fallback only when field loses focus
										const value = parseInt(e.currentTarget.value);
										if (isNaN(value) || value < 2) {
											updateRoundField(index, 'playersPerGroup', 2);
											e.currentTarget.value = '2';
										}
									}}
								/>
							</div>

							<!-- Matches per Group -->
							<div class="form-control min-w-[150px] flex-1">
								<label class="label" for="matchesPerGroup">
									<span class="label-text">Matches per Group</span>
								</label>
								<input
									type="number"
									class="input-bordered input w-full"
									min="1"
									value={round.matchesPerGroup || ''}
									oninput={(e) => {
										const value = e.currentTarget.value;
										if (value === '') {
											// Allow empty value during editing
											updateRoundField(index, 'matchesPerGroup', null);
										} else {
											updateRoundField(index, 'matchesPerGroup', parseInt(value));
										}
									}}
									onblur={(e) => {
										// Apply fallback only when field loses focus
										const value = parseInt(e.currentTarget.value);
										if (isNaN(value) || value < 1) {
											updateRoundField(index, 'matchesPerGroup', 1);
											e.currentTarget.value = '1';
										}
									}}
								/>
							</div>

							<!-- Advancing Players per Group -->
							<div class="form-control min-w-[150px] flex-1">
								<label class="label" for="advancingPlayersPerGroup">
									<span class="label-text">Advancing Players</span>
								</label>
								<input
									type="number"
									class="input-bordered input w-full"
									min="1"
									value={round.advancingPlayersPerGroup || ''}
									oninput={(e) => {
										const value = e.currentTarget.value;
										if (value === '') {
											// Allow empty value during editing
											updateRoundField(index, 'advancingPlayersPerGroup', null);
										} else {
											updateRoundField(index, 'advancingPlayersPerGroup', parseInt(value));
										}
									}}
									onblur={(e) => {
										// Apply fallback only when field loses focus
										const value = parseInt(e.currentTarget.value);
										if (isNaN(value) || value < 1) {
											updateRoundField(index, 'advancingPlayersPerGroup', 1);
											e.currentTarget.value = '1';
										} else if (round.playersPerGroup && value >= round.playersPerGroup) {
											// Ensure advancing players is less than total players
											const maxAdvancing = Math.max(1, round.playersPerGroup - 1);
											updateRoundField(index, 'advancingPlayersPerGroup', maxAdvancing);
											e.currentTarget.value = maxAdvancing.toString();
										}
									}}
								/>
							</div>

							<!-- Concurrent Groups -->
							<div class="form-control min-w-[150px] flex-1">
								<label class="label" for="concurrentGroups">
									<span class="label-text">Concurrent Groups</span>
								</label>
								<input
									type="number"
									class="input-bordered input w-full"
									min="1"
									max={round.groupCount}
									value={round.concurrentGroups || ''}
									oninput={(e) => {
										const value = e.currentTarget.value;
										if (value === '') {
											// Allow empty value during editing
											updateRoundField(index, 'concurrentGroups', null);
										} else {
											updateRoundField(index, 'concurrentGroups', parseInt(value));
										}
									}}
									onblur={(e) => {
										// Apply fallback only when field loses focus
										const value = parseInt(e.currentTarget.value);
										if (isNaN(value) || value < 1) {
											updateRoundField(index, 'concurrentGroups', 1);
											e.currentTarget.value = '1';
										} else if (value > round.groupCount) {
											// Ensure concurrent groups is not more than total groups
											updateRoundField(index, 'concurrentGroups', round.groupCount);
											e.currentTarget.value = round.groupCount.toString();
										}
									}}
								/>
								<label class="label" for="concurrentGroups">
									<span class="label-text-alt text-base-content/70"
										>Groups that can play at the same time</span
									>
								</label>
							</div>
						</div>
					</div>
				</div>
			{/each}

			<!-- Add Round Button -->
			<div class="mt-4 flex justify-center">
				<button type="button" class="btn btn-outline btn-primary" onclick={handleAddRound}>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
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
					Add Round
				</button>
			</div>
		</div>
	</div>
</div>
