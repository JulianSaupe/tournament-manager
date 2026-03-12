import type { PageServerLoad } from './$types';
import { type Tournament, TournamentStatus } from '$lib/types/tournament/tournament';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const tournaments = await locals.tournamentProvider.listTournaments();

		if (tournaments.success) {
			return {
				tournaments: tournaments,
				error: null as string | null
			};
		}

		return {
			tournaments: [] as Tournament[],
			error: tournaments.error.message
		};
	} catch (err) {
		return {
			tournaments: [] as Tournament[],
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
