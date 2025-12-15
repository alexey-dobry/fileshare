package internal

import (
	"context"
	"time"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/user/intuser"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/domain/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *InternalServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	u := models.User{
		ID:        req.Id,
		Name:      req.Name,
		Surname:   req.Surname,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	err := s.store.User().CreateUser(u)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.store.User().DeleteUser(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) CreateGroup(ctx context.Context, req *pb.GroupRequest) (*pb.CreateGroupResponse, error) {
	groupID := uuid.New().String()
	g := models.Group{
		ID:        groupID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	err := s.store.Group().CreateGroup(g)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.CreateGroupResponse{
		ID: groupID,
	}, nil
}

func (s *InternalServer) DeleteGroup(ctx context.Context, req *pb.GroupRequest) (*emptypb.Empty, error) {
	err := s.store.Group().DeleteGroup(req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) CreateCourse(ctx context.Context, req *pb.CourseRequest) (*pb.CreateCourseResponse, error) {
	courseID := uuid.New().String()
	c := models.Course{
		ID:        courseID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	err := s.store.Course().CreateCourse(c)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.CreateCourseResponse{
		ID: courseID,
	}, nil
}

func (s *InternalServer) DeleteCourse(ctx context.Context, req *pb.CourseRequest) (*emptypb.Empty, error) {
	err := s.store.Course().DeleteCourse(req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) AttachGroupToCourse(ctx context.Context, req *pb.AttachGroupCourseRequest) (*emptypb.Empty, error) {
	err := s.store.Group().AttachGroupToCourse(req.CourseID, req.GroupID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) DetachGroupFromCourse(ctx context.Context, req *pb.DetachGroupCourseRequest) (*emptypb.Empty, error) {
	err := s.store.Group().DetachGroupToCourse(req.GroupID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) AttachTeacherToCourse(ctx context.Context, req *pb.AttachTeacherCourseRequest) (*emptypb.Empty, error) {
	err := s.store.Course().AssignTeacherToCourse(req.TeacherID, req.CourseID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) DetachTeacherToCourse(ctx context.Context, req *pb.DetachTeacherCourseRequest) (*emptypb.Empty, error) {
	err := s.store.Course().DetachTeacherToCourse(req.TeacherID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) GetGroupsData(ctx context.Context, req *emptypb.Empty) (*pb.GroupsDataResponse, error) {
	result := make([]*pb.GroupData, 0)

	groupsData, err := s.store.Group().GetGroups()

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	for _, group := range groupsData {
		g := &pb.GroupData{
			GroupID:   group.ID,
			GroupName: group.Name,
		}

		result = append(result, g)
	}

	return &pb.GroupsDataResponse{
		GroupsData: result,
	}, nil
}

func (s *InternalServer) GetCoursesData(ctx context.Context, req *emptypb.Empty) (*pb.CoursesDataResponse, error) {
	result := make([]*pb.CourseData, 0)

	courseData, err := s.store.Course().GetCourses()

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	for _, course := range courseData {
		g := &pb.CourseData{
			CourseID:   course.ID,
			CourseName: course.Name,
		}

		result = append(result, g)
	}

	return &pb.CoursesDataResponse{
		CoursesData: result,
	}, nil
}
