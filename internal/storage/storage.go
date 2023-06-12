package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type Configer interface {
	GetAddressBaseData() string
}

type NodeStorage struct {
	conf Configer
	db   *sql.DB
}

func New(ctx context.Context, conf Configer) (*NodeStorage, error) {
	node := &NodeStorage{conf: conf}
	err := node.openDB()
	if err != nil {
		return nil, err
	}
	err = node.createTables(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r *NodeStorage) openDB() error {
	db, err := sql.Open("pgx", r.conf.GetAddressBaseData())
	if err != nil {
		logger.Log.Error("failed to open the database", zap.Error(err))
		return err
	}
	r.db = db
	return nil
}

func (r *NodeStorage) createTables(ctx context.Context) error {
	con, err := r.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer con.Close()
	err = createTableUsers(ctx, con)
	if err != nil {
		return err
	}
	err = createTableOrders(ctx, con)
	if err != nil {
		return err
	}
	return nil
}

func createTableUsers(ctx context.Context, con *sql.Conn) error {
	_, err := con.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users 
(
	id serial PRIMARY KEY,
	userID VARCHAR(45)  NOT NULL,
	login VARCHAR(250) NOT NULL,
    password VARCHAR(250) NOT NULL,
    balance INT DEFAULT 0
);`)
	return err
}

func createTableOrders(ctx context.Context, con *sql.Conn) error {
	_, err := con.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS orders 
(
	id serial PRIMARY KEY,
	userID VARCHAR(45)  NOT NULL,
	uploaded_at timestamp  NOT NULL,
	number VARCHAR(250) NOT NULL,
    status VARCHAR(250) NOT NULL,
    accural INT DEFAULT 0
);`)
	return err
}

func (r *NodeStorage) CreateNewUser(ctx context.Context, userID, login, password string) error {
	con, err := r.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return err
	}
	defer con.Close()
	row := con.QueryRowContext(ctx, "SELECT COUNT(id) FROM users WHERE login = $1", login)
	var countLogin int
	err = row.Scan(&countLogin)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("when scanning a request for a login", zap.Error(err))
		return err
	}
	if countLogin != 0 {
		return errorsapp.ErrLoginAlreadyBorrowed
	}
	_, err = con.ExecContext(ctx, "INSERT INTO users(userID, login, password) VALUES ($1, $2, $3);",
		userID, login, password)
	if err != nil {
		logger.Log.Error("error when adding a record to the database", zap.Error(err))
		return err
	}

	return nil
}

func (r *NodeStorage) AuthenticationUser(ctx context.Context, login, password string) (string, error) {
	con, err := r.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return "", err
	}
	defer con.Close()
	row := con.QueryRowContext(ctx, "SELECT userID FROM users WHERE login = $1 AND password = $2", login, password)
	var userID string
	err = row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errorsapp.ErrNoFoundUser
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("when scanning a request for a login", zap.Error(err))
		return "", err
	}
	return userID, nil
}

func (r *NodeStorage) Close() error {
	r.db.Close()
	return nil
}

func (r *NodeStorage) CreateOrder(ctx context.Context, number, userID string) error {

	con, err := r.db.Begin()
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return err
	}
	defer con.Rollback()
	row := con.QueryRowContext(ctx, "SELECT COUNT(id) FROM orders WHERE number = $1", number)
	var countNumber int
	row.Scan(&countNumber)
	if countNumber != 0 {

		row2 := con.QueryRowContext(ctx, "SELECT COUNT(id) FROM orders WHERE number = $1 AND userID = $2", number, userID)
		var countNumberUser int
		row2.Scan(&countNumberUser)
		if countNumberUser != 0 {
			return errorsapp.ErrOrderAlreadyThisUser
		}
		return errorsapp.ErrOrderAlreadyAnotherUser
	}

	_, err = con.ExecContext(ctx, "INSERT INTO orders (userID, number, uploaded_at, status) VALUES ($1, $2, $3, $4);",
		userID, number, time.Now(), "NOW")
	if err != nil {
		return err
	}

	con.Commit()
	return nil
}
