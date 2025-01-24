package option

import (
	"errors"
	"fmt"
	"io"
	"mygo/internal/migrate"
	"mygo/pkg/database"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
)

var (
	errNotPrepared = errors.New("options not prepared")
)

type Options struct {
	ConfigFile string `toml:"-"`
	DBDir      string `toml:"-"`
	LogLevel   string `toml:"-"`
	Http       struct {
		Port int `toml:"port"`
	} `toml:"http"`
	clients struct {
		database *sqlx.DB
	}
}

func (opt *Options) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, opt)
}

func (opt *Options) Prepare() error {
	if opt.DBDir == "" {
		return errors.New("database dir required")
	}

	if err := migrate.Execute(migrate.FS, "scripts", fmt.Sprintf("sqlite3://%s", opt.getDBDSN())); err != nil {
		return fmt.Errorf("database migrate, %w", err)
	}

	db, err := database.NewDB(database.Option{
		Driver:       "sqlite3",
		DSN:          opt.getDBDSN(),
		MaxIdleConns: 3,
		MaxIdleTime:  10 * time.Minute,
	})

	if err != nil {
		return fmt.Errorf("connect database, %w", err)
	}

	opt.clients.database = db.Unsafe()
	return nil
}

func (opt *Options) getDBDSN() string {
	file := filepath.Join(opt.DBDir, "main.db")
	return fmt.Sprintf("%s?mode=rwc&_timeout=5000&_journal=WAL&_sync=NORMAL&_encoding=UTF-8&_fk=1", filepath.ToSlash(file))
}
