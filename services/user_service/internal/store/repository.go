package store

import "github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"

type UserRepository interface {
	CreateUser(userData models.User) error

	GetUserByID(ID string) (models.User, error)
	GetUsersByGroupID(groupID string) ([]models.User, error)
	GetTeachersByCourseID(courseID string) ([]models.User, error)

	DeleteUser(Email string) error
}

type GroupRepository interface {
	CreateGroup(groupData models.Group) error

	GetGroupByUserID(userID string) (models.Group, error)
	GetGroupsByUserID(userID string) ([]models.Group, error)
	GetGroupsByCourseID(courseID string) ([]models.Group, error)
	GetGroups() ([]models.Group, error)

	AttachGroupToCourse(courseID, groupID string) error
	DetachGroupToCourse(groupID string) error

	AssignUserToGroup(userID, groupID string) error

	DeleteGroup(ID string) error
}

type CourseRepository interface {
	CreateCourse(courseData models.Course) error
	GetCoursesByUserID(userID string) ([]models.Course, error)
	GetCourses() ([]models.Course, error)

	AssignTeacherToCourse(teacherID, courseID string) error
	DetachTeacherToCourse(teacherID string) error

	DeleteCourse(ID string) error
}
