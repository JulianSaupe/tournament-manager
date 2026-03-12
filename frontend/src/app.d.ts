// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { TournamentProvider } from '$lib/ports/tournament-provider';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			tournamentProvider: TournamentProvider;
		}

		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
