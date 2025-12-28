use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};
use tonic::Status;

const JWT_SECRET: &[u8] = b"secret";
const TOKEN_EXPIRATION: i64 = 10_000_000_000;

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Claims {
    sub: String,
    exp: i64,
}

pub fn generate_token(username: &str) -> Result<String, Status> {
    let claims = Claims {
        sub: username.to_owned(),
        exp: TOKEN_EXPIRATION,
    };

    encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(JWT_SECRET),
    )
    .map_err(|_| Status::internal("Failed to create token"))
}

pub fn verify_token(token: &str) -> Result<bool, Status> {
    let token_data = decode::<Claims>(
        token,
        &DecodingKey::from_secret(JWT_SECRET),
        &Validation::default(),
    )
    .map_err(|_| Status::unauthenticated(""));

    Ok(token_data.ok().is_some())
}
