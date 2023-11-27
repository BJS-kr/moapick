use std::sync::Arc;

use axum::{
    routing::{get, post},
    Json, Router,
};

use super::user_service::{self, User};

pub struct UserController {
    user_service: Arc<user_service::UserService>,
}

impl UserController {
    pub fn new(user_service: user_service::UserService) -> Self {
        Self {
            user_service: Arc::new(user_service),
        }
    }

    pub fn routes(self) -> Router {
        let user_service = self.user_service.clone();

        Router::new()
            .route(
                "/sign_up",
                post(move |Json(signing_up): Json<User>| {
                    let user_service = user_service.clone();
                    async move {
                        if let User::SigningUp { name, email } = signing_up {
                            user_service.sign_up(name, email).await;
                        } else {
                            panic!("signing_up is not User::SigningUp")
                        }
                    }
                }),
            )
            .route(
                "/sign_in",
                post(move |Json(signing_in): Json<User>| {
                    let user_service = user_service.clone();
                    async move {
                        if let User::SigningIn { email } = signing_in {
                            user_service.sign_in(email).await;
                        } else {
                            panic!("signing_in is not User::SigningIn")
                        }
                    }
                }),
            )
    }
}
