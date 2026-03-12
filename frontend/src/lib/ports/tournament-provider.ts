import type { Tournament } from '$lib/types/tournament/tournament';
import type { Qualifying } from '$lib/types/tournament/qualifying';
import type { ApiResult } from '$lib/types/api-result';

export interface TournamentProvider {
	loadTournament(id: string): Promise<ApiResult<Tournament>>;

	loadQualifying(tournamentId: string): Promise<ApiResult<Qualifying>>;

	listTournaments(): Promise<ApiResult<Tournament[]>>;
}
