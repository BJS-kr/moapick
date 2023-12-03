use crate::user::user_controller::user_routes;
use axum::{
    http::{Method, StatusCode},
    response::IntoResponse,
    Extension, Json, Router,
};
use serde::{Deserialize, Serialize};
use std::{net::SocketAddr, sync::Arc};
use tower_http::cors::{Any, CorsLayer};

pub mod db;
pub mod user;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST])
        .allow_origin(Any);

    let db_client = db::conn::connect().await.unwrap();

    // user domain
    let user_repository = user::user_repository::UserRepository::new(db_client);
    let user_service = Arc::new(user::user_service::UserService::new(user_repository));

    // main app
    let app = Router::new()
        .nest("/user", user_routes())
        .layer(Extension(user_service))
        .layer(cors);

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));

    tracing::debug!("listening on {}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
