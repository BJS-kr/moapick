use axum::{http::StatusCode, response::IntoResponse, Json, Router};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
use user::user_controller::UserController;

pub mod db;
pub mod user;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let db_client = db::conn::connect().await.unwrap();
    let user_repository = user::user_repository::UserRepository::new(db_client);
    let user_service = user::user_service::UserService::new(user_repository);
    let app = Router::new().nest("/user", UserController::new(user_service).routes());
    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));

    tracing::debug!("listening on {}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
