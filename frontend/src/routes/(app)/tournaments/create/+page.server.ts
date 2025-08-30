import { fail, redirect, type Actions } from '@sveltejs/kit';
import { tournamentSchema, type TournamentFormData } from '$lib/validation/tournamentSchema';
import { z } from 'zod';

function zodErrorsToRecord(error: z.ZodError): Record<string, string> {
	const map: Record<string, string> = {};
	for (const issue of error.issues) {
		const path = issue.path.join('.');
		// If multiple errors for same path, keep the first message
		if (!map[path]) {
			map[path] = issue.message;
		}
	}
	return map;
}

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const payload = formData.get('payload');

		if (!payload || typeof payload !== 'string') {
			return fail(400, {
				message: 'Invalid form submission',
				errors: { form: 'Missing or invalid payload' }
			});
		}

		let data: TournamentFormData;
		try {
			data = JSON.parse(payload) as TournamentFormData;
		} catch (e) {
			return fail(400, {
				message: 'Invalid payload: unable to parse JSON',
				errors: { form: 'Malformed payload JSON' }
			});
		}

		// Validate with Zod
		const parsed = tournamentSchema.safeParse(data);
		if (!parsed.success) {
			const errors = zodErrorsToRecord(parsed.error);
			return fail(400, {
				message: 'Please fix the validation errors before submitting.',
				errors,
				values: data
			});
		}

		// Forward to Go backend
		try {
			const response = await fetch('http://localhost:3000/api/tournament', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(parsed.data)
			});

			if (!response.ok) {
				let backendError = '';
				try {
					backendError = await response.text();
				} catch {
					backendError = response.statusText || 'Unknown backend error';
				}
				return fail(response.status, {
					message: 'Failed to create tournament',
					backendError,
					errors: {},
					values: data
				});
			}
		} catch (err) {
			return fail(500, {
				message: 'Could not contact backend service',
				backendError: err instanceof Error ? err.message : 'Unknown error'
			});
		}

		throw redirect(303, '/');
	}
};
