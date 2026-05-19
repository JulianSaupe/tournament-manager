// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { TournamentProvider } from '$lib/ports/tournament-repository';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			tournamentProvider: TournamentProvider;
			user: {
				id: string;
				sessionId: string;
				// populated for future RBAC — fetch via AuthorizationService when needed
				permissions: string[];
				roles: string[];
			} | null;
		}

		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
