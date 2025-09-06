import type {Round} from "$lib/types/tournament/tournament";

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