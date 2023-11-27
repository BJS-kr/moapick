use super::user_repository;
use axum::{extract::Path, response::IntoResponse, Json};
use serde::{Deserialize, Serialize};
#[derive(Deserialize, Serialize)]
pub enum User {
    SigningUp {
        name: String,
        email: String,
    },
    SigningIn {
        email: String,
    },
    SignedIn {
        id: i32,
        name: String,
        email: String,
    },
}

pub struct UserService {
    user_repository: user_repository::UserRepository,
}

impl UserService {
    pub fn new(user_repository: user_repository::UserRepository) -> Self {
        Self { user_repository }
    }

    pub async fn sign_up(&self, name: String, email: String) -> &str {
        self.user_repository.sign_up(name, email).await
    }

    pub async fn sign_in(&self, email: String) -> &str {
        self.user_repository.sign_in(email).await
    }
}
