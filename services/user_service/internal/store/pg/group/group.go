package group

import (
	"github.com/Masterminds/squirrel"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"
)

// TODO ADD CONNECTION FROM TEACHER TO GROUP

func (r *Repository) CreateGroup(groupData models.Group) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("groups").
		Columns("name", "created_at").
		Values(groupData.Name, groupData.CreatedAt).
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

func (r *Repository) GetGroupByUserID(userID string) (models.Group, error) {
	var u models.Group

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("id", "name", "created_at").
		From("groups").
		Join("groups_users ON groups.id = groups_users.group_id").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return u, err
	}

	// executing query
	row := r.db.QueryRow(query, args...)

	err = row.Scan(&u.ID, &u.Name, &u.CreatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *Repository) GetGroupsByUserID(userID string) ([]models.Group, error) {
	result := make([]models.Group, 0)

	// build query
	query := `
	SELECT
			g.*
	FROM group_user gu
	JOIN groups g
			ON g.id = gu.group_id
	WHERE gu.user_id = %s;
	`

	// executing query
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.Group{}, err
	}

	for rows.Next() {
		var g models.Group
		err = rows.Scan(&g.ID, &g.Name, &g.CreatedAt)
		if err != nil {
			return []models.Group{}, err
		}

		result = append(result, g)
	}

	return result, nil
}

func (r *Repository) AssignUserToGroup(userID, groupID string) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("groups_users").
		Columns("user_id", "group_id").
		Values(userID, groupID).
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

func (r *Repository) DeleteGroup(ID string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Delete("groups").
		Where(squirrel.Eq{"id": ID}).
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

func (r *Repository) GetGroupsByCourseID(courseID string) ([]models.Group, error) {
	result := make([]models.Group, 0)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("id", "name", "created_at").
		From("groups").
		Join("groups_courses ON groups.id = groups_courses.group_id").
		Where(squirrel.Eq{"course_id": courseID}).
		ToSql()
	if err != nil {
		return []models.Group{}, err
	}

	// executing query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.Group{}, err
	}

	for rows.Next() {
		var g models.Group
		err = rows.Scan(&g.ID, &g.Name, &g.CreatedAt)
		if err != nil {
			return []models.Group{}, err
		}

		result = append(result, g)
	}

	return result, nil
}

func (r *Repository) AttachGroupToCourse(courseID, groupID string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("groups_courses").
		Columns("course_id", "group_id").
		Values(courseID, groupID).
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

func (r *Repository) DetachGroupToCourse(groupID string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Delete("groups_courses").
		Where(squirrel.Eq{"group_id": groupID}).
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

func (r *Repository) GetGroups() ([]models.Group, error) {
	result := make([]models.Group, 0)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("id", "name", "created_at").
		From("groups").
		Offset(15).
		ToSql()
	if err != nil {
		return []models.Group{}, err
	}

	// executing query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.Group{}, err
	}

	for rows.Next() {
		var g models.Group
		err = rows.Scan(&g.ID, &g.Name, &g.CreatedAt)
		if err != nil {
			return []models.Group{}, err
		}

		result = append(result, g)
	}

	return result, nil
}
