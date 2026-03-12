import type { PageServerLoad } from './$types';
import type { Tournament } from '$lib/models/tournament/tournament';
import { TournamentStatus } from '$lib/models/tournament/tournament';
import { type Qualifying, type QualifyingPlayer } from '$lib/models/tournament/qualifying';

export const load: PageServerLoad = async ({ locals, params }) => {
	try {
		const tournamentId = params.id;

		const [tournament, qualifying] = await Promise.all([
			locals.tournamentProvider.loadTournament(tournamentId),
			locals.tournamentProvider.loadQualifying(tournamentId)
		]);

		return {
			tournament: tournament,
			qualifying: qualifying
		};
	} catch (err) {
		return {
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
