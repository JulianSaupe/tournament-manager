import {z} from 'zod';

/**
 * Validation schema for tournament round creation
 * Matches the backend CreateTournamentRoundRequest structure
 */
export const tournamentRoundSchema = z.object({
    name: z.string().min(3, 'Round name must be at least 3 characters').max(255, 'Round name must be at most 255 characters'),
    matchCount: z.number().min(1, 'Match count must be at least 1'),
    playerAdvancementCount: z.number().min(0, 'Player advancement count must be at least 0'),
    groupSize: z.number().min(2, 'Group size must be at least 2'),
    groupCount: z.number().min(1, 'Group count must be at least 1'),
    concurrentGroupCount: z.number().min(1, 'Concurrent group count must be at least 1')
});

/**
 * Validation schema for tournament creation
 * Matches the backend CreateTournamentRequest structure
 */
export const tournamentSchema = z.object({
    name: z.string().min(3, 'Tournament name must be at least 3 characters').max(255, 'Tournament name must be at most 255 characters'),
    description: z.string().min(3, 'Description must be at least 3 characters').max(255, 'Description must be at most 255 characters'),
    startDate: z.string().min(1, 'Start date is required'),
    endDate: z.string().min(1, 'End date is required'),
    allowUnderfilledGroups: z.boolean(),
    playerCount: z.number().min(1, 'Player count must be at least 1'),
    rounds: z.array(tournamentRoundSchema).min(1, 'Tournament must have at least one round')
}).refine((data) => {
    // Validate that end date is after start date
    const startDate = new Date(data.startDate);
    const endDate = new Date(data.endDate);
    return endDate >= startDate;
}, {
    message: 'End date must be after start date',
    path: ['endDate']
}).refine((data) => {
    console.log('Validating tournament structure:', data);

    // Validate tournament structure rules (matching backend validation)
    if (!data.rounds || data.rounds.length === 0) {
        return false;
    }

    // Last round must have exactly one group
    const lastRound = data.rounds[data.rounds.length - 1];
    if (lastRound.groupCount !== 1) {
        return false;
    }

    // Validate round progression
    let previousRound: typeof data.rounds[0] | null = null;

    for (const round of data.rounds) {
        if (!data.allowUnderfilledGroups) {
            const playersInRound = round.groupCount * round.groupSize;

            if (previousRound === null) {
                // First round: players in round must equal total tournament players
                if (playersInRound !== data.playerCount) {
                    return false;
                }
            } else {
                // Subsequent rounds: players must equal advancing players from previous round
                const advancingFromPrevious = previousRound.playerAdvancementCount * previousRound.groupCount;
                if (playersInRound !== advancingFromPrevious) {
                    return false;
                }
            }
        }

        // Player advancement count cannot exceed group size
        if (round.playerAdvancementCount > round.groupSize) {
            return false;
        }

        previousRound = round;
    }

    return true;
}, {
    message: 'Invalid tournament structure: check round progression and player counts',
    path: ['rounds']
});

export type TournamentFormData = z.infer<typeof tournamentSchema>;
export type TournamentRoundData = z.infer<typeof tournamentRoundSchema>;