use diesel::{
    r2d2::{self, ConnectionManager},
    PgConnection,
};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};

pub type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;
pub type DB = diesel::pg::Pg;

const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations/");
// macos 환경 설정 내역
// brew install libpq && brew link --force libpq && PQ_LIB_DIR="$(brew --prefix libpq)/lib"
// cargo clean && cargo run

// windows 환경 설정 내역
// window에서는 choco에 libpq가 없어서 postgres를 직접 설치함. VS build tools에도 없음
// libpq가 없다는 에러 자체는 사용자 환경변수에 PQ_LIB_DIR 변수를 postgres설치된 곳에 lib경로를 넣어주니 해결되었음. 다만 cargo clean해야 함
// STATUS_DLL_NOTFOUND에러가 windows에서는 추가되는데, 이는 환경변수의 PATH 변수에 pg lib경로를 추가해주라는 조언이 있었음. 어느정도 해결되었지만 찾을 수 없는 dll 세개가 발견됨. libcrypto-3, libssl-3, libintl-9
// openssl 최신버전으로 재설치해봄 컴퓨터에 1버전이 깔려있었고 최신버전은 3버전이었음 -> libssl-3-x64.dll 과 libcrypto-3-x64.dll 에러가 해결됨. 남은건 libintl-9.dll
// libintl-9.dll 수동으로 다운받아서 딱히 넣을 곳이 없어서 그냥 postgres lib에 넣음. PATH에 포함된 어느 경로에건 있으면 될 듯
// 해결 완료
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
