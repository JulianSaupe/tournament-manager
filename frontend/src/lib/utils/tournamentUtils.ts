import type { Round, TournamentFormData, TournamentFormErrors, PhaseVisualizationData } from '$lib/types/tournament';

/**
 * Calculate default number of matches for a group
 * Round-robin tournament: each player plays against every other player once
 */
export function calculateDefaultMatches(players: number): number {
    return players > 1 ? Math.floor((players * (players - 1)) / 2) : 0;
}

/**
 * Add a new round to the tournament
 */
export function addRound(formData: TournamentFormData): TournamentFormData {
    const lastRound = formData.rounds[formData.rounds.length - 1];
    const newRoundNumber = formData.rounds.length + 1;

    // Get values from last round with fallbacks for null/undefined values
    const lastAdvancingPlayers = lastRound.advancingPlayersPerGroup || 2; // Default to 2 if null
    const lastConcurrentGroups = lastRound.concurrentGroups || 1; // Default to 1 if null

    // Calculate default values based on previous round
    const newPlayersPerGroup = lastAdvancingPlayers * 2; // Double the advancing players

    const newRound: Round = {
        name: `Round ${newRoundNumber}`,
        groupCount: 1, // Will be calculated automatically by reactive statement
        playersPerGroup: newPlayersPerGroup,
        matchesPerGroup: calculateDefaultMatches(newPlayersPerGroup),
        advancingPlayersPerGroup: Math.max(1, Math.floor(lastAdvancingPlayers / 2)), // Half the advancing players
        concurrentGroups: Math.max(1, Math.floor(lastConcurrentGroups / 2)) // Half the concurrent groups, minimum 1
    };

    return {
        ...formData,
        rounds: [...formData.rounds, newRound]
    };
}

/**
 * Remove a round from the tournament
 */
export function removeRound(formData: TournamentFormData, index: number): TournamentFormData {
    if (formData.rounds.length > 1) {
        return {
            ...formData,
            rounds: formData.rounds.filter((_, i) => i !== index)
        };
    }
    return formData;
}

/**
 * Calculate total number of rounds in the tournament
 */
export function calculateTotalRounds(formData: TournamentFormData): number {
    return formData.rounds.length;
}

/**
 * Calculate matches for a group
 * Default calculation: each player plays against every other player once
 */
function calculateMatchesForGroup(players: number): number {
    return players > 1 ? (players * (players - 1)) / 2 : 0;
}

/**
 * Generate visualization data for the tournament
 */
export function generateVisualizationData(formData: TournamentFormData): PhaseVisualizationData[] {
    const data: PhaseVisualizationData[] = [];
    let totalPlayers = formData.playerCount;

    // If group phase is enabled, use it as the starting point
    if (formData.groupPhase && formData.playerCount > 0) {
        const groupCount = Math.ceil(totalPlayers / formData.groupSize);
        const playersPerGroup = formData.groupSize;

        // Calculate how many players advance to the first round
        const advancingPlayers = groupCount * formData.rounds[0].advancingPlayersPerGroup;

        data.push({
            name: 'Group Phase',
            groupCount,
            playersPerGroup,
            totalPlayers,
            advancingPlayers,
            matchesPerGroup: calculateMatchesForGroup(playersPerGroup),
            advancingPlayersPerGroup: formData.rounds[0].advancingPlayersPerGroup,
            concurrentGroups: formData.groupSize // For group phase, default to group size as concurrent groups
        });

        totalPlayers = advancingPlayers;
    }

    // Add data for each configured round
    formData.rounds.forEach((round, index) => {
        const advancingPlayers = index < formData.rounds.length - 1
            ? formData.rounds[index + 1].groupCount * formData.rounds[index + 1].playersPerGroup
            : 1; // Final winner

        data.push({
            name: round.name,
            groupCount: round.groupCount,
            playersPerGroup: round.playersPerGroup,
            totalPlayers,
            advancingPlayers,
            matchesPerGroup: round.matchesPerGroup,
            advancingPlayersPerGroup: round.advancingPlayersPerGroup,
            concurrentGroups: round.concurrentGroups
        });

        totalPlayers = advancingPlayers;
    });

    return data;
}

/**
 * Update group counts for each round based on player count and players per group
 */
export function updateGroupCounts(formData: TournamentFormData): TournamentFormData {
    if (formData.playerCount <= 0) {
        return formData;
    }

    const updatedRounds = [...formData.rounds];
    let availablePlayers = formData.playerCount;

    // If group phase is enabled, calculate players for the first round
    if (formData.groupPhase) {
        // Calculate number of groups in the group phase
        const groupPhaseGroupCount = Math.ceil(formData.playerCount / formData.groupSize);

        // Calculate players advancing from group phase to first round
        availablePlayers = groupPhaseGroupCount * formData.rounds[0].advancingPlayersPerGroup;
    }

    // Update each round's group count
    updatedRounds.forEach((round, index) => {
        if (index > 0) {
            // For rounds after the first, available players come from previous round
            const prevRound = updatedRounds[index - 1];
            availablePlayers = prevRound.groupCount * prevRound.advancingPlayersPerGroup;
        }

        // Calculate group count based on available players and players per group
        // Ensure we always have at least 1 group
        round.groupCount = Math.max(1, Math.ceil(availablePlayers / round.playersPerGroup));
    });

    return {
        ...formData,
        rounds: updatedRounds
    };
}

/**
 * Validate the tournament form
 * Returns an object with validation result and errors
 */
export function validateForm(formData: TournamentFormData): { isValid: boolean; errors: TournamentFormErrors } {
    let isValid = true;

    // Reset errors
    const errors: TournamentFormErrors = {
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

    // Validate rounds configuration
    if (formData.rounds.length === 0) {
        alert('At least one round is required');
        isValid = false;
    }

    // Check if each round has valid settings
    for (let i = 0; i < formData.rounds.length; i++) {
        const round = formData.rounds[i];

        if (round.groupCount <= 0) {
            alert(`Round ${i + 1}: Number of groups must be greater than 0`);
            isValid = false;
        }

        if (round.playersPerGroup <= 1) {
            alert(`Round ${i + 1}: Players per group must be at least 2`);
            isValid = false;
        }

        if (round.matchesPerGroup <= 0) {
            alert(`Round ${i + 1}: Matches per group must be greater than 0`);
            isValid = false;
        }

        if (round.advancingPlayersPerGroup <= 0 || round.advancingPlayersPerGroup >= round.playersPerGroup) {
            alert(`Round ${i + 1}: Advancing players must be between 1 and ${round.playersPerGroup - 1}`);
            isValid = false;
        }

        if (round.concurrentGroups <= 0) {
            alert(`Round ${i + 1}: Concurrent groups must be at least 1`);
            isValid = false;
        }

        if (round.concurrentGroups > round.groupCount) {
            alert(`Round ${i + 1}: Concurrent groups cannot exceed the total number of groups (${round.groupCount})`);
            isValid = false;
        }
    }

    return { isValid, errors };
}