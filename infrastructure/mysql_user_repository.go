package infrastructure

import (
	"database/sql"
	"ddd-bottomup/domain/entity"
	"ddd-bottomup/domain/repository"
	"ddd-bottomup/domain/valueobject"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) repository.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) FindByID(id *entity.UserID) (*entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, is_premium 
		FROM users 
		WHERE id = ?
	`

	var userID, firstName, lastName, email string
	var isPremium bool
	err := r.db.QueryRow(query, id.Value()).Scan(&userID, &firstName, &lastName, &email, &isPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// エンティティを再構成
	reconstructedID, _ := entity.ReconstructUserID(userID)
	fullName, _ := valueobject.NewFullName(firstName, lastName)
	emailValue, _ := valueobject.NewEmail(email)
	user := entity.ReconstructUser(reconstructedID, fullName, emailValue, isPremium)

	return user, nil
}

func (r *MySQLUserRepository) FindByName(name *valueobject.FullName) (*entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, is_premium 
		FROM users 
		WHERE first_name = ? AND last_name = ?
	`

	var userID, firstName, lastName, email string
	var isPremium bool
	err := r.db.QueryRow(query, name.FirstName(), name.LastName()).Scan(&userID, &firstName, &lastName, &email, &isPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// エンティティを再構成
	reconstructedID, _ := entity.ReconstructUserID(userID)
	fullName, _ := valueobject.NewFullName(firstName, lastName)
	emailValue, _ := valueobject.NewEmail(email)
	user := entity.ReconstructUser(reconstructedID, fullName, emailValue, isPremium)

	return user, nil
}

func (r *MySQLUserRepository) Save(user *entity.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, email, is_premium, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE 
		first_name = VALUES(first_name), 
		last_name = VALUES(last_name), 
		email = VALUES(email),
		is_premium = VALUES(is_premium),
		updated_at = NOW()
	`

	_, err := r.db.Exec(query,
		user.ID().Value(),
		user.Name().FirstName(),
		user.Name().LastName(),
		user.Email().Value(),
		user.IsPremium(),
	)

	return err
}

func (r *MySQLUserRepository) Delete(id *entity.UserID) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id.Value())
	return err
}

// テーブル作成用SQL（参考）
/*
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
*/