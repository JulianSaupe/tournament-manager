import type { Tournament } from '$lib/domain/tournament';
import type { Qualifying } from '$lib/domain/qualifying';
import type { ApiResult } from '$lib/domain/api-result';

export interface TournamentProvider {
	loadTournament(id: string): Promise<ApiResult<Tournament>>;

	loadQualifying(tournamentId: string): Promise<ApiResult<Qualifying>>;

	listTournaments(): Promise<ApiResult<Tournament[]>>;
}
