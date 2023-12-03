use sea_query::{Expr, Iden, PostgresQueryBuilder, Query};
use tokio_postgres::{Client, Error};

#[derive(Iden)]
enum User {
    Table,
    Id,
    Name,
    Email,
}

pub struct UserRepository {
    client: Client,
}

impl UserRepository {
    pub fn new(client: Client) -> Self {
        Self { client }
    }

    pub async fn sign_up(&self, name: String, email: String) -> Result<(), Error> {
        let q = Query::select()
            .from(User::Table)
            .and_where(Expr::col(User::Email).eq(&email))
            .to_string(PostgresQueryBuilder);

        let user = self.client.simple_query(&q).await?;

        if user.len() > 0 {
            // 이미 가입된 유저이므로 에러처리
        }

        // 가입되지 않은 유저이므로 가입처리

        Ok(())
    }

    pub async fn sign_in(&self, email: String) -> &str {
        "signed in"
    }
}
