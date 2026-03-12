import type { TournamentProvider } from '$lib/ports/tournament-provider';
import type { Tournament } from '$lib/models/tournament/tournament';
import { TournamentApiSchema, TournamentListApiSchema } from '$lib/adapters/schemas/tournament';
import type { Qualifying } from '$lib/models/tournament/qualifying';
import { QualifyingApiSchema } from '$lib/adapters/schemas/qualifying';

export class ApiClient implements TournamentProvider {
	constructor(private readonly baseUrl: string = 'http://localhost:3000/api') {}

	async loadTournament(id: string): Promise<Tournament> {
		const response = await fetch(`${this.baseUrl}/tournament/${id}`);

		if (!response.ok) {
			throw new Error('Failed to load tournament');
		}

		const result = TournamentApiSchema.safeParse(await response.json());

		if (result.success) {
			return result.data;
		}

		throw new Error('Failed to load tournament');
	}

	async loadQualifying(tournamentId: string): Promise<Qualifying> {
		const response = await fetch(`${this.baseUrl}/tournament/${tournamentId}/qualifying`);

		if (!response.ok) {
			throw new Error('Failed to load qualifying');
		}

		const result = QualifyingApiSchema.safeParse(await response.json());

		if (result.success) {
			return result.data;
		}

		throw new Error('Failed to load qualifying');
	}

	async listTournaments(): Promise<Tournament[]> {
		const response = await fetch(`${this.baseUrl}/tournaments`);

		if (!response.ok) {
			throw new Error('Failed to load tournaments');
		}

		const result = TournamentListApiSchema.safeParse(await response.json());

		if (result.success) {
			return result.data;
		}

		throw new Error('Failed to load tournaments');
	}
}
