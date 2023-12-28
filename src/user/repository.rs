use crate::db::schema::users::{self, dsl};
use crate::{
    db::{conn::Pool, models::user::*},
    fault::Fault,
};
use anyhow::{anyhow, Ok, Result};
use diesel::prelude::*;
use diesel::r2d2::{ConnectionManager, PooledConnection};

use super::state::{User, UserOrFail};

pub struct UserRepository {
    pool: Pool,
}

impl UserRepository {
    pub fn new(pool: Pool) -> Self {
        Self { pool }
    }
    fn get_conn(&self) -> PooledConnection<ConnectionManager<PgConnection>> {
        self.pool
            .get()
            .map_err(|e| anyhow!("cannot get connection out of the pool: {}", e))
            .expect("cannot get connection out of the pool")
    }

    fn find_user_by_email(
        &self,
        email: &str,
        conn: &mut PooledConnection<ConnectionManager<PgConnection>>,
    ) -> Vec<UserModel> {
        dsl::users
            .filter(dsl::email.eq(email))
            .limit(1)
            .select(UserModel::as_select())
            .load(conn)
            .expect("error loading existing user")
    }

    pub async fn sign_up(&self, name: String, email: String) -> Result<UserOrFail> {
        let mut conn = self.get_conn();
        // check if user already exists
        let existing_user: Vec<UserModel> = self.find_user_by_email(&email, &mut conn);

        // 이미 가입된 유저이므로 Ok(ClientFault)
        // 연산의 성공은 Ok이다
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

    pub async fn get_signed_in_user(&self, email: String) -> Result<UserOrFail> {
        let mut conn = self.get_conn();
        let user_vec = self.find_user_by_email(&email, &mut conn);

        if user_vec.is_empty() {
            return Ok(UserOrFail::Fail(Fault::Client));
        }

        let user = user_vec[0].to_owned();

        Ok(UserOrFail::User(User::SignedIn {
            id: user.id,
            name: user.name,
            email: user.email,
        }))
    }
}
