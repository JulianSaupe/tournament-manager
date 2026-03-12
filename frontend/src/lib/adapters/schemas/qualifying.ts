import { z, type ZodType } from 'zod';
import type { Qualifying, QualifyingPlayer } from '$lib/types/tournament/qualifying';
import { BaseApiSchema } from '$lib/adapters/schemas/base';

const PlayerSchema: ZodType<QualifyingPlayer> = z
	.object({
		name: z.string(),
		position: z.number(),
		signup_date: z.string(),
		time: z.number()
	})
	.transform((data) => ({
		name: data.name,
		position: data.position,
		signupDate: data.signup_date,
		time: data.time
	}));

const QualifyingSchema: ZodType<Qualifying> = z
	.object({
		players: z.array(PlayerSchema)
	})
	.transform((data) => ({
		status: 'test',
		players: data.players
	}));

export const QualifyingApiSchema: ZodType<Qualifying> = BaseApiSchema(QualifyingSchema).transform(
	(response) => response.data
);
