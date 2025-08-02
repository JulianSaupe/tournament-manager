/**
 * Types and interfaces for tournament structure
 */

/**
 * Represents a round in the tournament
 */
export interface Round {
    name: string;
    groupCount: number; // Will be calculated based on playerCount and playersPerGroup
    playersPerGroup: number;
    matchesPerGroup: number;
    advancingPlayersPerGroup: number;
    concurrentGroups: number; // Number of groups that can play concurrently
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