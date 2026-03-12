import type { PageServerLoad } from './$types';
import { type Tournament, TournamentStatus } from '$lib/types/tournament/tournament';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const result = await locals.tournamentProvider.listTournaments();

		if (result.success) {
			return {
				tournaments: result.data,
				error: null as string | null
			};
		}

		return {
			tournaments: [] as Tournament[],
			error: result.error.message
		};
	} catch (err) {
		return {
			tournaments: [] as Tournament[],
			error: err instanceof Error ? err.message : 'Unknown error'
		};
	}
};
