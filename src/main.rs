use crate::user::user_controller::user_routes;
use axum::{
    http::{Method, StatusCode},
    middleware::from_fn,
    response::{IntoResponse, Response},
    routing::get,
    Extension, Router,
};
use dotenvy::dotenv;
use middleware::jwt::jwt_middleware;
use middleware::uri::uri_middleware;
use std::{net::SocketAddr, sync::Arc, time::Duration};
use tower_http::{
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};
use tracing::Span;

pub mod db;
pub mod fault;
pub mod middleware;
pub mod user;

async fn handler_404() -> impl IntoResponse {
    (
        StatusCode::NOT_FOUND,
        "The requested resource was not found",
    )
}

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
    let trace_layer = TraceLayer::new_for_http().on_response(
        |response: &Response, latency: Duration, _span: &Span| {
            println!(
                "{:?} {} {}ms",
                response
                    .extensions()
                    .get::<middleware::uri::RequestUri>()
                    .map(|r| &r.0)
                    .unwrap(),
                response.status(),
                latency.as_millis()
            )
        },
    );

    let app = Router::new()
        .layer(cors)
        .route("/hello", get(|| async { "Hello, World!" }))
        .nest("/user", user_routes())
        .layer(from_fn(jwt_middleware))
        .layer(Extension(user_service))
        .fallback(handler_404)
        .layer(from_fn(uri_middleware))
        .layer(trace_layer);

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));

    tracing::debug!("listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
