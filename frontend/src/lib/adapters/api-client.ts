import type { TournamentProvider } from '$lib/ports/tournament-provider';
import type { Tournament } from '$lib/types/tournament/tournament';
import { TournamentApiSchema, TournamentListApiSchema } from '$lib/adapters/schemas/tournament';
import type { Qualifying } from '$lib/types/tournament/qualifying';
import { QualifyingApiSchema } from '$lib/adapters/schemas/qualifying';
import { ApiErrorResult, type ApiResult, ApiSuccessResult } from '$lib/types/api-result';
import type { z } from 'zod';

export class ApiClient implements TournamentProvider {
	constructor(private readonly baseUrl: string = 'http://localhost:3000/api') {}

	async loadTournament(id: string): Promise<ApiResult<Tournament>> {
		const response = await fetch(`${this.baseUrl}/tournament/${id}`);
		return this.handleApiResponse(response, TournamentApiSchema, 'Failed to load tournament');
	}

	async loadQualifying(tournamentId: string): Promise<ApiResult<Qualifying>> {
		const response = await fetch(`${this.baseUrl}/tournament/${tournamentId}/qualifying`);
		return this.handleApiResponse(response, QualifyingApiSchema, 'Failed to load qualifying');
	}

	async listTournaments(): Promise<ApiResult<Tournament[]>> {
		const response = await fetch(`${this.baseUrl}/tournaments`);
		return this.handleApiResponse(response, TournamentListApiSchema, 'Failed to load tournaments');
	}

	private async handleApiResponse<T>(
		response: Response,
		schema: z.ZodSchema<T>,
		errorMessage: string
	): Promise<ApiResult<T>> {
		if (!response.ok) {
			return ApiErrorResult(new Error(errorMessage));
		}

		const result = schema.safeParse(await response.json());

		if (result.success) {
			return ApiSuccessResult(result.data);
		}

		return ApiErrorResult(new Error(errorMessage));
	}
}
