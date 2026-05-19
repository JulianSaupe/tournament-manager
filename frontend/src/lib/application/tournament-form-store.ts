import { writable } from 'svelte/store';
import type { TournamentFormData } from '$lib/application/tournament-validator';
import { tournamentSchema } from '$lib/application/tournament-validator';
import { z } from 'zod';

const initialFormData: TournamentFormData = {
	name: '',
	description: '',
	startDate: '',
	endDate: '',
	allowUnderfilledGroups: false,
	playerCount: 0,
	rounds: []
};

export const tournamentForm = writable<TournamentFormData>(initialFormData);

export const tournamentFormErrors = writable<Record<string, string>>({});

export const tournamentFormValid = writable<boolean>(false);

export function validateTournamentForm(data: TournamentFormData): {
	isValid: boolean;
	errors: Record<string, string>;
} {
	try {
		tournamentSchema.parse(data);
		return { isValid: true, errors: {} };
	} catch (error) {
		const errors: Record<string, string> = {};

		if (error instanceof z.ZodError) {
			error.issues.forEach((issue) => {
				const path = issue.path.join('.');
				errors[path] = issue.message;
			});
		}

		return { isValid: false, errors };
	}
}

export function resetTournamentForm() {
	tournamentForm.set(initialFormData);
	tournamentFormErrors.set({});
	tournamentFormValid.set(false);
}
