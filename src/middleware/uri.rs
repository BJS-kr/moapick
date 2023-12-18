use axum::{
    body::Body,
    http::{Request, Uri},
    middleware::Next,
    response::Response,
};

#[derive(Debug, Clone)]
pub struct RequestUri(pub Uri);

pub async fn uri_middleware(request: Request<Body>, next: Next) -> Response {
    let uri = request.uri().clone();
    let mut response = next.run(request).await;

    response.extensions_mut().insert(RequestUri(uri));

    response
}
