use std::sync::Arc;

use diesel::{
    r2d2::{self, ConnectionManager},
    PgConnection,
};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};

pub type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;
pub type DB = diesel::pg::Pg;

const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations/");

// brew install libpq && brew link --force libpq && PQ_LIB_DIR="$(brew --prefix libpq)/lib"
// cargo clean && cargo run
pub fn establish_connection(db_url: String) -> Pool {
    let manager = ConnectionManager::<PgConnection>::new(db_url);
    let pool = r2d2::Pool::builder()
        .build(manager)
        .expect("failed to create pool");

    let mut conn = pool.get().expect("failed to apply db migrations");

    // same with: <PgConnection as MigrationHarness<DB>>::run_pending_migrations(&mut conn, MIGRATIONS);
    run_migrations(&mut conn);

    pool
}

pub fn run_migrations(conn: &mut impl MigrationHarness<DB>) {
    conn.run_pending_migrations(MIGRATIONS)
        .expect("failed to run migrations");
}
