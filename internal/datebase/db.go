package datebase

import (
	"context"
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(ctx context.Context) (*Repository, error) {

	dsn := "root:ComplexPassw0rd!@tcp(127.0.0.1:3307)/domain_checker"
	db, err := sql.Open("mysql", dsn)

	// connect to db
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	// check connection to db
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки подключения к базе данных", err)
	}

	go func() {
		defer db.Close()
		<-ctx.Done()
	}()

	// Передача объекта базы данных в пакет restapi
	return &Repository{db: db}, nil
}
