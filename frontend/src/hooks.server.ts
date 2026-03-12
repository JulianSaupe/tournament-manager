import type { Handle } from '@sveltejs/kit';
import { ApiClient } from '$lib/adapters/api-client';

const apiClient = new ApiClient('http://localhost:3000/api');

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.tournamentProvider = apiClient;
	return resolve(event);
};
