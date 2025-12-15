package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"
)

func (r *Repository) CreateUser(userData models.User) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("users").
		Columns("uuid", "name", "surname", "email", "role", "created_at").
		Values(userData.ID, userData.Name, userData.Surname, userData.Email, userData.Role, userData.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	// executing query
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByID(ID string) (models.User, error) {
	var u models.User

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("name", "surname", "email").
		From("users").
		Where(squirrel.Eq{"uuid": ID}).
		ToSql()
	if err != nil {
		return u, err
	}

	// executing query
	row := r.db.QueryRow(query, args...)

	err = row.Scan(&u.Name, &u.Surname, &u.Email)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *Repository) DeleteUser(Email string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Delete("users").
		Where(squirrel.Eq{"email": Email}).
		ToSql()
	if err != nil {
		return err
	}

	// executing query
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUsersByGroupID(groupID string) ([]models.User, error) {
	result := make([]models.User, 0)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("uuid", "name", "surname", "email").
		From("groups_users").
		Join("users ON groups_users.user_id = users.uuid").
		Where(squirrel.Eq{"group_id": groupID}).
		ToSql()
	if err != nil {
		return []models.User{}, err
	}

	// executing query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.User{}, err
	}

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.CreatedAt)
		if err != nil {
			return []models.User{}, err
		}

		result = append(result, u)
	}

	return result, nil
}

func (r *Repository) GetTeachersByCourseID(courseID string) ([]models.User, error) {
	result := make([]models.User, 0)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("uuid", "name", "surname", "email").
		From("users").
		Join("teacher_courses ON teacher_courses.user_id = users.uuid").
		Where(squirrel.Eq{"course_id": courseID}).
		ToSql()
	if err != nil {
		return []models.User{}, err
	}

	// executing query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.User{}, err
	}

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.CreatedAt)
		if err != nil {
			return []models.User{}, err
		}

		result = append(result, u)
	}

	return result, nil
}
