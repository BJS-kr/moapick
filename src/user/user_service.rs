use super::{
    user_repository,
    user_state::{AccessTokenAndRefreshToken, TokensOrFail, User, UserOrFail},
};
use crate::fault::Fault;
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct AccessTokenClaims<'a> {
    id: i32,
    name: &'a str,
    email: &'a str,
    exp: usize,
}

#[derive(Debug, Serialize, Deserialize)]
struct RefreshTokenClaims<'a> {
    id: i32,
    name: &'a str,
    email: &'a str,
    exp: usize,
}

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
            &AccessTokenClaims {
                id,
                name,
                email,
                exp,
            },
            // as_ref는 reference들끼리 형변환 하게 해주는 메서드. into랑 비슷하다고 보면 된다.
            &EncodingKey::from_secret(secret.as_ref()),
        )
        .unwrap();

        jwt_token
    }
}
