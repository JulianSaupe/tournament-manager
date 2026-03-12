import { transform, z } from 'zod';
import { BaseApiSchema } from '$lib/adapters/schemas/base';
import { type Tournament, TournamentStatus } from '$lib/types/tournament/tournament';

const TournamentSchema: z.ZodType<Tournament> = z
	.object({
		id: z.string(),
		name: z.string(),
		description: z.string(),
		start_date: z.string(),
		end_date: z.string(),
		status: z.enum(TournamentStatus),
		player_count: z.number().optional()
	})
	.transform((data): Tournament => {
		return {
			id: data.id,
			name: data.name,
			description: data.description,
			startDate: data.start_date,
			endDate: data.end_date,
			status: data.status,
			playerCount: data.player_count
		};
	});

export const TournamentApiSchema: z.ZodType<Tournament> = BaseApiSchema(TournamentSchema).transform(
	(response) => {
		return response.data;
	}
);

export const TournamentListApiSchema: z.ZodType<Tournament[]> = BaseApiSchema(
	z.array(TournamentSchema)
).transform((response) => response.data);
