package course

import (
	"github.com/Masterminds/squirrel"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"
)

func (r *Repository) CreateCourse(courseData models.Course) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("courses").
		Columns("name", "created_at").
		Values(courseData.Name, courseData.CreatedAt).
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

func (r *Repository) GetCoursesByUserID(userID string) ([]models.Course, error) {
	result := make([]models.Course, 0)

	// build query
	query := `
	SELECT DISTINCT
			c.*
	FROM users u
	JOIN groups_users gu
			ON gu.user_id = u.id
	JOIN courses_groups cg
			ON cg.group_id = gu.group_id
	JOIN courses c
			ON c.id = cg.course_id
	WHERE u.uuid = %s;
	`

	// executing query
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.Course{}, err
	}

	for rows.Next() {
		var c models.Course
		err = rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return []models.Course{}, err
		}

		result = append(result, c)
	}

	return result, nil
}

func (r *Repository) AssignTeacherToCourse(teacherID, courseID string) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("groups_courses").
		Columns("teacher_id", "course_id").
		Values(teacherID, courseID).
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

func (r *Repository) DeleteCourse(ID string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Delete("courses").
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

func (r *Repository) DetachTeacherToCourse(teacherID string) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Delete("groups_courses").
		Where(squirrel.Eq{"teacher_id": teacherID}).
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

func (r *Repository) GetCourses() ([]models.Course, error) {
	result := make([]models.Course, 0)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Select("id", "name", "created_at").
		From("courses").
		Offset(15).
		ToSql()
	if err != nil {
		return []models.Course{}, err
	}

	// executing query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []models.Course{}, err
	}

	for rows.Next() {
		var c models.Course
		err = rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return []models.Course{}, err
		}

		result = append(result, c)
	}

	return result, nil
}
