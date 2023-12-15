use crate::fault::Fault;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub enum User {
    SigningUp {
        name: String,
        email: String,
    },
    SignedUp {
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

pub enum UserOrFail {
    User(User),
    Fail(Fault),
}
