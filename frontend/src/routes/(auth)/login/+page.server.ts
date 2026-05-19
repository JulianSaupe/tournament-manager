import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { grpcLogin } from '$lib/server/grpc/identity-client';

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.user) {
		redirect(302, '/');
	}
};

export const actions: Actions = {
	default: async ({ request, cookies, getClientAddress }) => {
		const data = await request.formData();
		const email = String(data.get('email') ?? '');
		const password = String(data.get('password') ?? '');

		if (!email || !password) {
			return fail(400, { error: 'Email and password are required', email });
		}

		try {
			const result = await grpcLogin(
				email,
				password,
				getClientAddress(),
				request.headers.get('user-agent') ?? ''
			);

			if (!result.success) {
				return fail(401, { error: result.message || 'Invalid credentials', email });
			}

			const maxAge = result.expires_at
				? parseInt(result.expires_at.seconds) - Math.floor(Date.now() / 1000)
				: 60 * 60 * 24 * 7;

			cookies.set('session_id', result.session_id, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				maxAge
			});
		} catch (e) {
			console.error('Authentication error:', e);
			return fail(503, { error: 'Authentication service unavailable', email });
		}

		redirect(302, '/');
	}
};
