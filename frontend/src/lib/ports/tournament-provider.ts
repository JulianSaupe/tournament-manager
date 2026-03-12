import type { Tournament } from '$lib/models/tournament/tournament';
import type { Qualifying } from '$lib/models/tournament/qualifying';

export interface TournamentProvider {
	loadTournament(id: string): Promise<Tournament>;

	loadQualifying(tournamentId: string): Promise<Qualifying>;

	listTournaments(): Promise<Tournament[]>;
}
