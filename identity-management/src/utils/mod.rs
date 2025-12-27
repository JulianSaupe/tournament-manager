use argon2::{
    password_hash::{
        rand_core::OsRng,
        PasswordHash, PasswordHasher, PasswordVerifier,
        SaltString
    },
    Argon2
};

pub fn hash_string(input: &str) -> Result<String, argon2::password_hash::Error> {
    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();
    let password_hash = argon2.hash_password(input.as_bytes(), &salt)?.to_string();

    Ok(password_hash)
}

pub fn verify_hash(input: &str, hash: &str) -> bool {
    PasswordHash::new(hash)
        .ok()
        .map(|h| Argon2::default().verify_password(input.as_bytes(), &h).is_ok())
        .unwrap_or(false)
}