import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { grpcLogout } from '$lib/server/grpc/identity-client';

export const actions: Actions = {
	default: async ({ cookies }) => {
		const sessionId = cookies.get('session_id');
		if (sessionId) {
			try {
				await grpcLogout(sessionId);
			} catch {
				// Ignore errors — clear the cookie regardless
			}
		}
		cookies.delete('session_id', { path: '/' });
		redirect(302, '/login');
	}
};
