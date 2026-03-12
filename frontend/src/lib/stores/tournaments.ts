import { writable } from 'svelte/store';
import type { Tournament } from '$lib/types/tournament/tournament';

export const tournaments = writable<Tournament[]>([]);

export function addTournament(tournament: Tournament) {
	tournaments.update((list) => [tournament, ...list]);
}

export function updateTournament(id: string, updates: Partial<Tournament>) {
	tournaments.update((list) => list.map((t) => (t.id === id ? { ...t, ...updates } : t)));
}

export function removeTournament(id: string) {
	tournaments.update((list) => list.filter((t) => t.id !== id));
}

export function setTournaments(list: Tournament[]) {
	tournaments.set(list);
}
