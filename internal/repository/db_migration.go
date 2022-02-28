package repository

import "github.com/golang-migrate/migrate"

func RunMigration(databaseDSN string) (bool, error) {
	m, err := migrate.New("file://internal/repository/migration", databaseDSN)
	if err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	return true, nil
}
