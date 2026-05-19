import type { PageServerLoad } from './$types';
import type { Tournament } from '$lib/domain/tournament';
import { TournamentStatus } from '$lib/domain/tournament';
import { type Qualifying, type QualifyingPlayer } from '$lib/domain/qualifying';

export const load: PageServerLoad = async ({ locals, params }) => {
	try {
		const tournamentId = params.id;

		const [tournamentResult, qualifyingResult] = await Promise.all([
			locals.tournamentProvider.loadTournament(tournamentId),
			locals.tournamentProvider.loadQualifying(tournamentId)
		]);

		if (tournamentResult.success && qualifyingResult.success) {
			return {
				tournament: tournamentResult.data,
				qualifying: qualifyingResult.data
			};
		}

		const err = tournamentResult.error || qualifyingResult.error;

		return {
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	} catch (err) {
		return {
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
