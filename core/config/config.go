package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"

	// psql driver.
	_ "github.com/lib/pq"
)

type Interface struct {
	str interface{}
}

type Config struct {
	ConfigFile string
	EnvPrefix  string
	DSN        string
}

// Simple constructor.
func New(str interface{}) Interface {
	return Interface{str: str}
}

// Method wraps discrete methods.
func (s Interface) Combine(c Config) error {
	if c.ConfigFile != "" {
		fmt.Printf("try to apply config from file %s...\n", c.ConfigFile)
		if err := s.SetFromFile(c.ConfigFile); err != nil {
			return fmt.Errorf("can't apply config from file: %w", err)
		}
	}
	if c.EnvPrefix != "" {
		fmt.Printf("try to apply config from environment...\n")
		if err := s.SetFromEnv(c.EnvPrefix); err != nil {
			return fmt.Errorf("can't apply envvars to config:%w", err)
		}
	}
	if c.DSN != "" {
		fmt.Printf("try to apply config from DSN %s...\n", c.DSN)
		db, dbname, err := DialDSN(c.DSN)
		if err != nil {
			return fmt.Errorf("can't dial DB:%w", err)
		}
		if err := s.SetFromDB(db, dbname); err != nil {
			return fmt.Errorf("can't apply db lines to config:%w", err)
		}
	}
	return nil
}

// Method adds and replace config fields from file.
func (s Interface) SetFromFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("can't open config file: %w", err)
	}
	defer f.Close()
	l, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("can't read content of the config file : %w", err)
	}
	_, err = toml.Decode(string(l), s.str)
	if err != nil {
		return fmt.Errorf("can't parce config file : %w", err)
	}
	return nil
}

// Method adds and replace config fields from env.
func (s Interface) SetFromEnv(prefix string) error {
	return getEnvVar(reflect.ValueOf(s.str), reflect.TypeOf(s.str), -1, prefix)
}

func DialDSN(dsn string) (db *sql.DB, dbname string, err error) {
	m := strings.FieldsFunc(dsn, func(r rune) bool { return r == ':' || r == '@' || r == '/' })
	dbName := m[len(m)-1]
	if dbName == "" {
		return nil, "", fmt.Errorf("DSN not contains database name: %s", dsn)
	}

	var driver string
	switch {
	case strings.HasPrefix(dsn, "postgres://"):
		driver = "postgres"
		dsn = strings.TrimLeft(dsn, "postgres://")
	case strings.HasPrefix(dsn, "mysql://"):
		driver = "mysql"
		dsn = strings.TrimLeft(dsn, "mysql://")
	default:
		return nil, "", fmt.Errorf("can't use unknown SQL dialect")
	}

	db, err = sql.Open(driver, dsn)
	if err != nil {
		return nil, "", fmt.Errorf("can't connect to DB: %w", err)
	}
	return db, dbName, nil
}

// Method adds and replace config fields from db.
func (s Interface) SetFromDB(db *sql.DB, dbname string) error {
	defer db.Close()
	res := make(map[string]string)
	var key, val string

	//TODO: Перенести это в параметры.
	table := "config"
	q := "SELECT " + table + ".key, " + table + ".value FROM " + table
	results, err := db.Query(q)
	if err != nil || results.Err() != nil {
		return fmt.Errorf("can't get key-value pairs from DB: %w", err)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&key, &val)
		if err != nil {
			return fmt.Errorf("can't parse key-value into vars: %w", err)
		}
		res[key] = val
	}
	if err = parseToStruct(reflect.ValueOf(s.str), reflect.TypeOf(s.str), -1, "", res); err != nil {
		return fmt.Errorf("can't parse into struct: %w", err)
	}
	return nil
}
