import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { grpcCreateUser, grpcLogin } from '$lib/server/grpc/identity-client';

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.user) {
		redirect(302, '/');
	}
};

export const actions: Actions = {
	default: async ({ request, cookies, getClientAddress }) => {
		const data = await request.formData();
		const username = String(data.get('username') ?? '');
		const email = String(data.get('email') ?? '');
		const password = String(data.get('password') ?? '');

		if (!username || !email || !password) {
			return fail(400, { error: 'All fields are required', username, email });
		}

		try {
			const created = await grpcCreateUser(username, email, password);
			if (!created.success) {
				return fail(422, { error: 'Registration failed', username, email });
			}
		} catch {
			return fail(503, { error: 'Registration service unavailable', username, email });
		}

		// Auto-login after registration
		try {
			const session = await grpcLogin(
				email,
				password,
				getClientAddress(),
				request.headers.get('user-agent') ?? ''
			);

			if (session.success) {
				const maxAge = session.expires_at
					? parseInt(session.expires_at.seconds) - Math.floor(Date.now() / 1000)
					: 60 * 60 * 24 * 7;

				cookies.set('session_id', session.session_id, {
					path: '/',
					httpOnly: true,
					sameSite: 'lax',
					maxAge
				});
			}
		} catch {
			// Registration succeeded but auto-login failed — redirect to login
			redirect(302, '/login');
		}

		redirect(302, '/');
	}
};
