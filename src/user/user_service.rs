use super::{
    user_repository,
    user_state::{User, UserOrFail},
};
use crate::fault::Fault;

pub struct UserService {
    user_repository: user_repository::UserRepository,
}

impl UserService {
    pub fn new(user_repository: user_repository::UserRepository) -> Self {
        Self { user_repository }
    }

    pub async fn sign_up(&self, name: String, email: String) -> UserOrFail {
        let user_or_fail = self.user_repository.sign_up(name, email).await;

        if let Ok(user_or_fail) = user_or_fail {
            match user_or_fail {
                signed_up_user @ UserOrFail::User(User::SignedUp { .. }) => signed_up_user,
                client_fault @ UserOrFail::Fail(Fault::Client) => client_fault,
                _ => UserOrFail::Fail(Fault::Server),
            }
        } else {
            UserOrFail::Fail(Fault::Server)
        }
    }

    pub async fn sign_in(&self, email: String) -> UserOrFail {
        self.user_repository.sign_in(email).await;
        todo!();
    }
}
