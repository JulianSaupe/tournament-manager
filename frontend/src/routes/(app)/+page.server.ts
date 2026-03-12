import type { PageServerLoad } from './$types';
import { type Tournament, TournamentStatus } from '$lib/models/tournament/tournament';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const tournaments = await locals.tournamentProvider.listTournaments();

		return {
			tournaments: tournaments,
			error: null as string | null
		};
	} catch (err) {
		return {
			tournaments: [] as Tournament[],
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
