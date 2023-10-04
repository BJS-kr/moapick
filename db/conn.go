package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func InitDB() {
	var err error

	/**
	TODO default DB말고 특정한 DB쓰자! ex)dbname='...' <- postgres 키고 DB만들어줘야함. 볼륨 지정해두면 계속 쓸 수 있음
	FIXME 연결 정보 환경 변수로 옮기기
	*/
	dsn := "host=localhost user=postgres password=test port=5432 sslmode=disable"
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})


	if err != nil {
		panic(err)
	}
}

