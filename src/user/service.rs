use super::{
    repository,
    state::{AccessTokenAndRefreshToken, TokensOrFail, User, UserOrFail},
};
use crate::fault::Fault;
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct TokenClaims {
    pub id: i32,
    pub name: String,
    pub email: String,
    pub exp: usize,
}

pub struct UserService {
    user_repository: repository::UserRepository,
}

impl UserService {
    pub fn new(user_repository: repository::UserRepository) -> Self {
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

    pub async fn sign_in(&self, email: String) -> TokensOrFail {
        let signed_in_user = self.user_repository.get_signed_in_user(email).await;
        let secret =
            std::env::var("JWT_SECRET").expect("JWT_SECRET must be provided as a env variable");
        if let Ok(signed_in_user) = signed_in_user {
            match signed_in_user {
                UserOrFail::User(User::SignedIn { id, name, email }) => {
                    let access_token = self.generate_jwt_token(&secret, id, &name, &email, 1000000);
                    let refresh_token =
                        self.generate_jwt_token(&secret, id, &name, &email, 10000000);

                    TokensOrFail::Tokens(AccessTokenAndRefreshToken {
                        access_token,
                        refresh_token,
                    })
                }
                UserOrFail::Fail(Fault::Client) => TokensOrFail::Fail(Fault::Client),
                _ => TokensOrFail::Fail(Fault::Server),
            }
        } else {
            TokensOrFail::Fail(Fault::Server)
        }
    }

    fn generate_jwt_token(
        &self,
        secret: &str,
        id: i32,
        name: &str,
        email: &str,
        exp: usize,
    ) -> String {
        let jwt_token = encode(
            &Header::default(),
            &TokenClaims {
                id,
                name: name.into(),
                email: email.into(),
                exp,
            },
            // as_ref는 reference들끼리 형변환 하게 해주는 메서드. into랑 비슷하다고 보면 된다.
            &EncodingKey::from_secret(secret.as_ref()),
        )
        .unwrap();

        jwt_token
    }
}
