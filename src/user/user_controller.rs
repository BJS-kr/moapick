use axum::{
    http::StatusCode, middleware::from_fn, response::IntoResponse, routing::post, Extension, Json,
    Router,
};
use serde_json::json;
use std::sync::Arc;

use crate::middleware::jwt::jwt_middleware;
use crate::{
    fault::Fault,
    user::{
        user_service::UserService,
        user_state::{TokensOrFail, User, UserOrFail},
    },
};

pub fn user_routes() -> Router {
    let router = Router::new()
        .route("/sign-up", post(sign_up))
        .route("/sign-in", post(sign_in))
        .layer(from_fn(jwt_middleware));

    async fn sign_in(
        Extension(user_service): Extension<Arc<UserService>>,
        Json(user): Json<User>,
    ) -> impl IntoResponse {
        if let User::SigningIn { email } = user {
            let tokens_or_fail = user_service.sign_in(email).await;

            match tokens_or_fail {
                TokensOrFail::Tokens(tokens) => (StatusCode::CREATED, Json(json!(tokens))),
                TokensOrFail::Fail(Fault::Client) => (
                    StatusCode::NOT_FOUND,
                    Json(json!({
                        "error": "user not found"
                    })),
                ),
                _ => (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(json!({
                        "error": "internal server error"
                    })),
                ),
            }
        } else {
            (
                StatusCode::BAD_REQUEST,
                Json(json!({"error": "bad request"})),
            )
        }
    }

    async fn sign_up(
        Extension(user_service): Extension<Arc<UserService>>,
        Json(user): Json<User>,
    ) -> impl IntoResponse {
        if let User::SigningUp { name, email } = user {
            let sign_up_result = user_service.sign_up(name, email).await;

            match sign_up_result {
                UserOrFail::User(User::SignedUp { .. }) => StatusCode::CREATED,
                UserOrFail::Fail(Fault::Client) => StatusCode::CONFLICT,
                _ => StatusCode::INTERNAL_SERVER_ERROR,
            }
        } else {
            StatusCode::BAD_REQUEST
        }
    }

    router
}
