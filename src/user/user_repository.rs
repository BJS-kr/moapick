use crate::db::conn::Pool;
use anyhow::Result;

use super::user_state::{User, UserOrFail};

pub struct UserRepository {
    pool: Pool,
}

impl UserRepository {
    pub fn new(pool: Pool) -> Self {
        Self { pool }
    }

    pub async fn sign_up(&self, name: String, email: String) -> Result<UserOrFail> {
        // check if user already exists
        let conn = self.pool.get().unwrap();

        // 이미 가입된 유저이므로 Ok(ClientError)

        // 가입되지 않은 유저이므로 가입처리

        Ok(UserOrFail::User(User::SignedUp { name, email }))
    }

    pub async fn sign_in(&self, email: String) -> Result<UserOrFail> {
        todo!()
    }
}
