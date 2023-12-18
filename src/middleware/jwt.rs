use anyhow::Result;
use axum::{
    body::Body,
    http::{Request, Uri},
    middleware::Next,
    response::Response,
};
use jsonwebtoken::{decode, Algorithm, DecodingKey, Validation};

use crate::user::user_service::AccessTokenClaims;

pub enum JWTError {
    NoAuthorizationHeader,
    NotBearerAuthorization,
}

const SECRET: String = std::env::var("JWT_SECRET").unwrap();
const BEARER: &str = "Bearer";
const AUTHORIZATION: &str = "Authorization";

pub fn jwt_middleware(request: Request<Body>, next: Next) -> Result<(), JWTError> {
    let headers = request.headers();
    let auth_header = match headers.get(AUTHORIZATION) {
        Some(h) => h,
        None => return Err(JWTError::NoAuthorizationHeader),
    };
    let str_auth_header = match std::str::from_utf8(auth_header.as_bytes()) {
        Ok(h) => h,
        Err(_) => return Err(JWTError::NoAuthorizationHeader),
    };

    if !str_auth_header.starts_with(BEARER) {
        return Err(JWTError::NotBearerAuthorization);
    }

    let access_token = str_auth_header.trim_start_matches(BEARER).to_owned();

    let decoded = decode::<AccessTokenClaims>(
        &access_token,
        &DecodingKey::from_secret(SECRET.as_ref()),
        &Validation::new(Algorithm::HS256),
    );

    todo!()
}
