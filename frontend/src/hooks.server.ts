import type { Handle } from '@sveltejs/kit';
import { TournamentHttpAdapter } from '$lib/adapters/http/tournament-http-adapter';
import { grpcValidateSession } from '$lib/server/grpc/identity-client';

const tournamentAdapter = new TournamentHttpAdapter('http://localhost:3000/api');

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.tournamentProvider = tournamentAdapter;
	event.locals.user = null;

	const sessionId = event.cookies.get('session_id');
	if (sessionId) {
		try {
			const result = await grpcValidateSession(sessionId);
			if (result.valid) {
				event.locals.user = {
					id: result.user_id,
					sessionId,
					permissions: [],
					roles: []
				};
			}
		} catch {
			// Identity service unavailable or session invalid — remain unauthenticated
		}
	}

	return resolve(event);
};
