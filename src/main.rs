use crate::user::user_controller::user_routes;
use axum::{http::Method, Extension, Router};
use dotenvy::dotenv;
use std::{net::SocketAddr, sync::Arc};
use tower_http::cors::{Any, CorsLayer};

pub mod db;
pub mod user;

#[tokio::main]
async fn main() {
    dotenv().ok();
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST])
        .allow_origin(Any);
    let db_url =
        std::env::var("DATABASE_URL").expect("DATABASE_URL must be provided as a env variable");
    let pool = db::conn::establish_connection(db_url);
    // user domain
    let user_repository = user::user_repository::UserRepository::new(pool);
    let user_service = Arc::new(user::user_service::UserService::new(user_repository));

    // main app
    let app = Router::new()
        .nest("/user", user_routes())
        .layer(Extension(user_service))
        .layer(cors);

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));

    tracing::debug!("listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
