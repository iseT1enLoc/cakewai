package repository

import (
	"cakewai/cakewai.com/domain"
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userId int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	if user.GoogleId != "" {
		res, err := tx.Exec(`INSERT INTO users (email, google_id, name, profile_picture) VALUES (?, ?, ?,?)`, user.Email, user.GoogleId, user.Name, user.ProfilePicture)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		user.Id = int(id)
		return user, nil
	}

	res, err := tx.Exec(`INSERT INTO users (email, password) VALUES (?, ?) `, user.Email, user.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user.Id = int(id)

	return user, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, userId int) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetUserByEmail implements UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	stm := `SELECT id, google_id, profile_picture, name, password, phone, created_at, updated_at 
	        FROM users WHERE email = $1` // Use $1 for positional parameter

	// Execute the query and scan the result into the user struct
	err := u.db.QueryRowContext(ctx, stm, email).Scan(
		&user.Id,
		&user.GoogleId,
		&user.ProfilePicture,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No user found, return nil
	} else if err != nil {
		return nil, err // Return the error if any other
	}

	return &user, nil
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	user := domain.User{}

	// Use QueryRowContext to support context
	err := u.db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, id).Scan(&user.Id) // Adjust based on your User fields
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no user is found
		}
		return nil, err // Return other errors
	}

	return &user, nil
}

// GetUsers implements UserRepository.
func (u *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	stm := `SELECT id, google_id, profile_picture, name, email, phone, created_at, updated_at FROM users`

	// Use QueryContext for context support
	records, err := u.db.Query(stm)
	if err != nil {
		return nil, err // Return the error if the query fails
	}
	defer records.Close() // Ensure the rows are closed

	for records.Next() {
		var user domain.User

		// Scan the values into the user struct
		if err := records.Scan(
			&user.Id,
			&user.GoogleId,
			&user.ProfilePicture,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err // Return the error if scanning fails
		}

		// Append the user to the slice
		users = append(users, &user)
	}

	// Check for errors from iterating over rows.
	if err := records.Err(); err != nil {
		return nil, err // Return the error if there was an error during iteration
	}

	return users, nil // Return the slice of users
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	if user.Password != "" {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		user.Password = string(encryptedPassword)
	}

	fieldsQuery := ""
	if user.Email != "" {
		fieldsQuery += "email = :email,"
	}
	if user.Name != "" {
		fieldsQuery += "name = :name,"
	}
	if user.Password != "" {
		fieldsQuery += "password = :password,"
	}
	if user.Phone != "" {
		fieldsQuery += "phone = :phone,"
	}
	fieldsQuery = fieldsQuery[:len(fieldsQuery)-1]

	_, err = tx.Exec("UPDATE users SET "+fieldsQuery+" WHERE id = ?", user.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
