package public

import (
	"context"
	"fmt"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/user/pubuser"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PublicServer) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetMyProfileResponse, error) {
	userData, err := s.store.User().GetUserByID(req.UserID)
	if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.GetMyProfileResponse{
		ProfileData: &pb.UserData{
			FullName: fmt.Sprintf("%s %s", userData.Name, userData.Surname),
			Email:    userData.Email,
		},
	}, nil
}

func (s *PublicServer) GetMyCourses(ctx context.Context, req *pb.GetMyCoursesRequest) (*pb.GetMyCoursesResponse, error) {
	coursesData, err := s.store.Course().GetCoursesByUserID(req.UserID)

	if err != nil {
		s.logger.Errorf("Failed to get courses data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	courses := make([]*pb.CourseData, 0)
	for _, course := range coursesData {
		c := &pb.CourseData{
			CourseID:   course.ID,
			CourseName: course.Name,
		}
		courses = append(courses, c)
	}

	return &pb.GetMyCoursesResponse{
		Courses: courses,
	}, nil
}

func (s *PublicServer) StudentGetCourseTeachers(ctx context.Context, req *pb.StudentGetCourseTeachersRequest) (*pb.StudentGetCourseTeachersResponse, error) {
	teacherData, err := s.store.User().GetTeachersByCourseID(req.CourseID)

	if err != nil {
		s.logger.Errorf("Failed to get courses data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	teachers := make([]*pb.UserData, 0)
	for _, teacher := range teacherData {
		t := &pb.UserData{
			FullName: fmt.Sprintf("%s %s", teacher.Name, teacher.Surname),
			Email:    teacher.Email,
		}

		teachers = append(teachers, t)
	}

	return &pb.StudentGetCourseTeachersResponse{
		TeachersData: teachers,
	}, nil
}

func (s *PublicServer) TeacherGetCourseGroups(ctx context.Context, req *pb.TeacherGetCourseGroupsRequest) (*pb.TeacherGetCourseGroupsResponse, error) {
	coursesData, err := s.store.Group().GetGroupsByCourseID(req.CourseID)

	if err != nil {
		s.logger.Errorf("Failed to get courses data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	groups := make([]*pb.GroupData, 0)
	for _, group := range coursesData {
		g := &pb.GroupData{
			GroupID:   group.ID,
			GroupName: group.Name,
		}

		groups = append(groups, g)
	}

	return &pb.TeacherGetCourseGroupsResponse{
		Groups: groups,
	}, nil
}

func (s *PublicServer) TeacherGetGroupStudents(ctx context.Context, req *pb.TeacherGetGroupStudentsRequest) (*pb.TeacherGetGroupStudentsResponse, error) {
	studentsData, err := s.store.User().GetUsersByGroupID(req.GroupID)

	if err != nil {
		s.logger.Errorf("Failed to get courses data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	students := make([]*pb.UserData, 0)
	for _, student := range studentsData {
		s := &pb.UserData{
			FullName: fmt.Sprintf("%s %s", student.Name, student.Surname),
			Email:    student.Email,
		}

		students = append(students, s)
	}

	return &pb.TeacherGetGroupStudentsResponse{
		StudentsData: students,
	}, nil
}
