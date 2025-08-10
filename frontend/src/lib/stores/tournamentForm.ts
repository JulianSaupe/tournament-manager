import { writable } from 'svelte/store';
import type { TournamentFormData } from '$lib/validation/tournamentSchema';
import { tournamentSchema } from '$lib/validation/tournamentSchema';
import { z } from 'zod';

// Initial form data with all required fields
const initialFormData: TournamentFormData = {
	name: '',
	description: '',
	startDate: '',
	endDate: '',
	allowUnderfilledGroups: false,
	playerCount: 0,
	rounds: []
};

// Store for form data
export const tournamentForm = writable<TournamentFormData>(initialFormData);

// Store for validation errors
export const tournamentFormErrors = writable<Record<string, string>>({});

// Store for form validation state
export const tournamentFormValid = writable<boolean>(false);

// Validation function
export function validateTournamentForm(data: TournamentFormData): { isValid: boolean; errors: Record<string, string> } {
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

// Function to reset form to initial state
export function resetTournamentForm() {
	tournamentForm.set(initialFormData);
	tournamentFormErrors.set({});
	tournamentFormValid.set(false);
}
