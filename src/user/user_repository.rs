use crate::{
    db::{conn::Pool, models::user::*},
    fault::Fault,
};
use anyhow::{Ok, Result};
use diesel::prelude::*;

use super::user_state::{User, UserOrFail};

pub struct UserRepository {
    pool: Pool,
}

impl UserRepository {
    pub fn new(pool: Pool) -> Self {
        Self { pool }
    }

    pub async fn sign_up(&self, name: String, email: String) -> Result<UserOrFail> {
        use crate::db::schema::users;
        use crate::db::schema::users::dsl;
        let mut conn = self
            .pool
            .get()
            .expect("cannot get connection out of the pool");
        // check if user already exists
        let existing_user: Vec<UserModel> = dsl::users
            .filter(dsl::email.eq(&email))
            .limit(1)
            .select(UserModel::as_select())
            .load(&mut conn)
            .expect("error loading existing user");

        // 이미 가입된 유저이므로 Ok(ClientFault)
        if !(existing_user.is_empty()) {
            return Ok(UserOrFail::Fail(Fault::Client));
        }

        // 가입되지 않은 유저이므로 가입처리
        let new_user = NewUser { name, email };

        let inserted = diesel::insert_into(users::table)
            .values(&new_user)
            .returning(UserModel::as_returning())
            .get_result(&mut conn)
            .expect("failed to insert signing up user");

        Ok(UserOrFail::User(User::SignedUp {
            name: inserted.name,
            email: inserted.email,
        }))
    }

    pub async fn sign_in(&self, email: String) -> Result<UserOrFail> {
        todo!()
    }
}
