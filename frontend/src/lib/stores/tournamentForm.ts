import { writable } from 'svelte/store';

export const tournamentForm = writable({
	name: '',
	playerCount: 0,
	startDate: '',
	endDate: '',
	errors: {
		name: false,
		playerCount: 0,
		startDate: false,
		endDate: false
	}
});
