use crate::db::conn::Pool;
use anyhow::Result;

pub struct UserRepository {
    pool: Pool,
}

impl UserRepository {
    pub fn new(pool: Pool) -> Self {
        Self { pool }
    }

    pub async fn sign_up(&self, name: String, email: String) -> Result<()> {
        // 이미 가입된 유저이므로 에러처리

        // 가입되지 않은 유저이므로 가입처리

        Ok(())
    }

    pub async fn sign_in(&self, email: String) -> &str {
        "signed in"
    }
}
