export interface Tournament {
    id: string;
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    status: TournamentStatus;
    playerCount?: number;
}

export interface Round {
	name: string;
	groupCount: number;
	playersPerGroup: number;
	matchesPerGroup: number;
	advancingPlayersPerGroup: number;
	concurrentGroups: number;
}

export enum TournamentStatus {
    DRAFT = 'draft',
    ACTIVE = 'active',
    COMPLETED = 'completed',
    CANCELLED = 'cancelled'
}