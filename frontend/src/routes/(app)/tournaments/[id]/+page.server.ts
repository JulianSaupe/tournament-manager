import type { PageServerLoad } from './$types';
import type { Tournament } from '$lib/types/tournament/tournament';
import { TournamentStatus } from '$lib/types/tournament/tournament';
import { type Qualifying, type QualifyingPlayer } from '$lib/types/tournament/qualifying';

export const load: PageServerLoad = async ({ locals, params }) => {
	try {
		const tournamentId = params.id;

		const [tournament, qualifying] = await Promise.all([
			locals.tournamentProvider.loadTournament(tournamentId),
			locals.tournamentProvider.loadQualifying(tournamentId)
		]);

		if (tournament.success && qualifying.success) {
			return {
				tournament: tournament.data,
				qualifying: qualifying.data
			};
		}

		const err = tournament.error || qualifying.error;

		return {
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	} catch (err) {
		return {
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
