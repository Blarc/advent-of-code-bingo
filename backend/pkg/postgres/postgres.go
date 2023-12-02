package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func New(host, port, user, pass, name string) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s TimeZone=Europe/Ljubljana",
		host,
		port,
		user,
		pass,
		name,
	)

	// TODO: add logger via gorm config
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}

func (p *Postgres) Close() error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
