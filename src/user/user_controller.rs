use axum::{
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Extension, Json, Router,
};
use std::sync::Arc;

use crate::user::{user_service::UserService, user_state::User};

pub fn user_routes() -> Router {
    let router = Router::new()
        .route("/hello", get(|| async { "Hello, World!" }))
        .route("/sign-up", post(sign_up))
        .route("/sign-in", post(sign_in));

    async fn sign_in(
        Extension(user_service): Extension<Arc<UserService>>,
        Json(user): Json<User>,
    ) -> impl IntoResponse {
        if let User::SigningIn { email } = user {
            user_service.sign_in(email).await;
            StatusCode::CREATED
        } else {
            StatusCode::BAD_REQUEST
        }
    }

    async fn sign_up(
        Extension(user_service): Extension<Arc<UserService>>,
        Json(user): Json<User>,
    ) -> impl IntoResponse {
        if let User::SigningUp { name, email } = user {
            user_service.sign_up(name, email).await;
            StatusCode::CREATED
        } else {
            StatusCode::BAD_REQUEST
        }
    }

    router
}
