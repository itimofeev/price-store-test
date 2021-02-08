package pg

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func doMigrationIfNeeded(db *pg.DB) error {
	col := getMigrations()
	_, _, _ = col.Run(db, "init") // nolint:dogsled
	_, _, err := col.Run(db, "up")
	return err
}

func getMigrations() *migrations.Collection {
	sqls := []string{`
		CREATE TABLE products
		(
			id           VARCHAR(64) NOT NULL PRIMARY KEY,
			name         TEXT        NOT NULL,
			price        BIGINT      NOT NULL,
			last_update  TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
			update_count BIGINT      NOT NULL DEFAULT 0,
			CHECK (price > 0)
		);
		
		CREATE UNIQUE INDEX products__name__index ON products (name);`,
	}

	migs := make([]*migrations.Migration, 0, len(sqls))

	for i, sql := range sqls {
		sql := sql
		migs = append(migs, &migrations.Migration{
			Version: int64(i + 1),
			UpTx:    true,
			Up: func(db migrations.DB) error {
				_, err := db.Exec(sql)
				return err
			},
		})
	}

	return migrations.NewCollection(migs...)
}
