import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';
import { env } from '$env/dynamic/private';

const PROTO_DIR = path.resolve(process.cwd(), '../backend/shared/proto');

const authDef = protoLoader.loadSync(path.join(PROTO_DIR, 'authentication.proto'), {
	keepCase: true,
	longs: String,
	enums: String,
	defaults: true,
	oneofs: true
});

const userDef = protoLoader.loadSync(path.join(PROTO_DIR, 'user.proto'), {
	keepCase: true,
	longs: String,
	enums: String,
	defaults: true,
	oneofs: true
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const authProto = grpc.loadPackageDefinition(authDef) as any;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const userProto = grpc.loadPackageDefinition(userDef) as any;

// Lazy singletons — created on first use so env vars are guaranteed to be loaded
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let _authClient: any = null;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let _userClient: any = null;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getAuthClient(): any {
	if (!_authClient) {
		_authClient = new authProto.authentication.AuthenticationService(
			env.GRPC_IDENTITY_URL ?? '[::1]:5000',
			grpc.credentials.createInsecure()
		);
	}
	return _authClient;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getUserClient(): any {
	if (!_userClient) {
		_userClient = new userProto.user.UserService(
			env.GRPC_IDENTITY_URL ?? '[::1]:5000',
			grpc.credentials.createInsecure()
		);
	}
	return _userClient;
}

function makeMetadata(): grpc.Metadata {
	const meta = new grpc.Metadata();
	meta.add('authorization', `Bearer ${env.GRPC_AUTH_TOKEN ?? ''}`);
	return meta;
}

export interface LoginResult {
	success: boolean;
	session_id: string;
	expires_at: { seconds: string; nanos: number } | null;
	message: string;
}

export interface ValidateSessionResult {
	valid: boolean;
	user_id: string;
	expires_at: { seconds: string; nanos: number } | null;
	message: string;
}

export interface LogoutResult {
	success: boolean;
	message: string;
}

export interface CreateUserResult {
	user_id: string;
	success: boolean;
}

export function grpcLogin(
	email: string,
	password: string,
	ipAddress: string,
	userAgent: string
): Promise<LoginResult> {
	return new Promise((resolve, reject) => {
		getAuthClient().Login(
			{ email, password, ip_address: ipAddress, user_agent: userAgent },
			makeMetadata(),
			(err: grpc.ServiceError | null, res: LoginResult) => {
				if (err) reject(err);
				else resolve(res);
			}
		);
	});
}

export function grpcValidateSession(sessionId: string): Promise<ValidateSessionResult> {
	return new Promise((resolve, reject) => {
		getAuthClient().ValidateSession(
			{ session_id: sessionId },
			makeMetadata(),
			(err: grpc.ServiceError | null, res: ValidateSessionResult) => {
				if (err) reject(err);
				else resolve(res);
			}
		);
	});
}

export function grpcLogout(sessionId: string): Promise<LogoutResult> {
	return new Promise((resolve, reject) => {
		getAuthClient().Logout(
			{ session_id: sessionId },
			makeMetadata(),
			(err: grpc.ServiceError | null, res: LogoutResult) => {
				if (err) reject(err);
				else resolve(res);
			}
		);
	});
}

export function grpcCreateUser(
	username: string,
	email: string,
	password: string
): Promise<CreateUserResult> {
	return new Promise((resolve, reject) => {
		getUserClient().Create(
			{ username, email, password },
			makeMetadata(),
			(err: grpc.ServiceError | null, res: CreateUserResult) => {
				if (err) reject(err);
				else resolve(res);
			}
		);
	});
}
