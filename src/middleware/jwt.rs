use axum::{
    body::Body,
    http::{Request, StatusCode},
    middleware::Next,
    response::Response,
};
use jsonwebtoken::{decode, errors::ErrorKind, Algorithm, DecodingKey, Validation};

use crate::user::service::TokenClaims;

const BEARER: &str = "Bearer";
const AUTHORIZATION: &str = "Authorization";

fn build_response(status: StatusCode, message: &str) -> Response<Body> {
    Response::builder()
        .status(status)
        .body(Body::from(message.to_owned()))
        .unwrap()
}

pub async fn jwt_middleware(mut request: Request<Body>, next: Next) -> Response<Body> {
    let secret = std::env::var("JWT_SECRET").expect("JWT_SECRET not set");
    let headers = request.headers();
    let auth_header = match headers.get(AUTHORIZATION) {
        Some(h) => h,
        None => return build_response(StatusCode::BAD_REQUEST, "No Authorization Header"),
    };

    let str_auth_header = match std::str::from_utf8(auth_header.as_bytes()) {
        Ok(h) => h,
        Err(_) => return build_response(StatusCode::BAD_REQUEST, "No Authorization Header"),
    };

    if !str_auth_header.starts_with(BEARER) {
        return build_response(StatusCode::BAD_REQUEST, "Not Bearer Authentication");
    }

    let access_token = str_auth_header.trim_start_matches(BEARER).to_owned();

    let decoded = decode::<TokenClaims>(
        &access_token,
        &DecodingKey::from_secret(secret.as_ref()),
        &Validation::new(Algorithm::HS256),
    );

    let token_claims = match decoded {
        Ok(data) => data,
        Err(e) => match e.kind() {
            ErrorKind::ExpiredSignature => {
                return build_response(StatusCode::BAD_REQUEST, "expired access token")
            }
            _ => return build_response(StatusCode::BAD_REQUEST, "invalid token"),
        },
    };

    request.extensions_mut().insert(token_claims);

    next.run(request).await
}
