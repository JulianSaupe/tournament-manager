/**
 * Types and interfaces for tournament structure
 */

/**
 * Represents a round in the tournament
 */
export interface Round {
    name: string;
    groupCount: number;
    playersPerGroup: number;
    matchesPerGroup: number;
    advancingPlayersPerGroup: number;
    concurrentGroups: number;
}

/**
 * Form data for tournament creation
 */
export interface TournamentFormData {
    name: string;
    startDate: string;
    endDate: string;
    playerCount: number;
    groupPhase: boolean;
    allowPartiallyFilledGroups: boolean;
    groupSize: number;
    rounds: Round[];
}

/**
 * Form validation errors
 */
export interface TournamentFormErrors {
    name: string;
    startDate: string;
    endDate: string;
    playerCount: string;
}