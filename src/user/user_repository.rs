use tokio_postgres::Client;

pub struct UserRepository {
    client: Client,
}

impl UserRepository {
    pub fn new(client: Client) -> Self {
        Self { client }
    }

    pub async fn sign_up(&self, name: String, email: String) -> &str {
        "signed up"
    }

    pub async fn sign_in(&self, email: String) -> &str {
        "signed in"
    }
}
