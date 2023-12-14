use super::user_repository;
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

pub enum ErrorContext {
    ClientError,
    ServerError,
}

pub enum UserOrError {
    User(User),
    Error(ErrorContext),
}

pub struct UserService {
    user_repository: user_repository::UserRepository,
}

impl UserService {
    pub fn new(user_repository: user_repository::UserRepository) -> Self {
        Self { user_repository }
    }

    pub async fn sign_up(&self, name: String, email: String) -> UserOrError {
        let user = self.user_repository.sign_up(name, email).await;

        match user {
            Ok(user) => match user {
                UserOrError::User(User::SigningUp { name, email }) => {
                    UserOrError::User(User::SigningUp { name, email })
                }
                _ => UserOrError::Error(ErrorContext::ClientError),
            },
            Err(_) => UserOrError::Error(ErrorContext::ServerError),
        }
    }

    pub async fn sign_in(&self, email: String) -> &str {
        self.user_repository.sign_in(email).await
    }
}
