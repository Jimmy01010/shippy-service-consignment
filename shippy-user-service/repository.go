package main

import (
	"context"
	pb "github.com/Jimmy01010/protocol/shippy-user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       string `sql:"id"`
	Name     string `sql:"name"`
	Email    string `sql:"email"`
	Company  string `sql:"company"`
	Password string `sql:"password"`
}

type Repository interface {
	GetAll(ctx context.Context) ([]*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (p PostgresRepository) Get(ctx context.Context, id string) (*User, error) {
	var user *User
	if err := p.db.GetContext(ctx, &user, "select * from users where id = $1", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (p PostgresRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)
	if err := p.db.GetContext(ctx, users, "select * from users"); err != nil {
		return nil, err
	}
	return nil, nil
}

func (p PostgresRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()
	insertSql := "insert into users (id, name, email, company, password) values ($1, $2, $3, $4, $5)"
	if _, err := p.db.
		ExecContext(ctx, insertSql, user.ID, user.Name, user.Email, user.Company, user.Password); err != nil {
		return nil
	}
	return nil
}

func (p PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user *User
	if err := p.db.GetContext(ctx, &user, "select * from users where email = $1", email); err != nil {
		return nil, err
	}

	return user, nil
}

// MarshalUser 将pb结构解析为DB的user结构
func MarshalUser(user *pb.User) *User {
	return &User{
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Company:  user.Company,
		Password: user.Password,
	}
}

// UnmarshalUser 将DB的user结构解析为pb结构
func UnmarshalUser(user *User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Company:  user.Company,
		Password: user.Password,
	}
}

// UnmarshalUserCollection 将DB的users结构解析为pb结构
func UnmarshalUserCollection(users []*User) []*pb.User {
	u := make([]*pb.User, len(users))
	for _, val := range users {
		u = append(u, UnmarshalUser(val))
	}
	return u
}
