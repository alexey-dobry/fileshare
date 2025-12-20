package course

import (
	"github.com/Masterminds/squirrel"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"
)

func (r *Repository) CreateCourse(courseData models.Course) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("course").
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
	query, args, err := squirrel.Select("c.id", "c.name", "c.created_at").
		From("user u").
		Join("user_group gu ON gu.user_id = u.id").
		Join("group_course gc ON gc.group_id = gu.group_id").
		Join("course c ON c.id = gc.course_id").
		Where(squirrel.Eq{"u.uuid": userID}).
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

func (r *Repository) AssignTeacherToCourse(teacherID, courseID string) error {
	// making query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// build query
	query, args, err := psql.Insert("teacher_course").
		Columns("user_id", "course_id").
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
	query, args, err := psql.Delete("course").
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
	query, args, err := psql.
		Delete("teacher_course").
		Where(squirrel.Eq{"user_id": teacherID}).
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
		From("course").
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
