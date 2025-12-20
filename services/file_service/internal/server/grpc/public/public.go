package public

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/file/pubfile"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/domain/model"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/domain/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PublicServer) UploadFileUnary(
	ctx context.Context,
	req *pb.UploadFileUnaryRequest,
) (*pb.UploadFileUnaryResponse, error) {

	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fileID := uuid.New()
	storageKey := fmt.Sprintf("files/%s", fileID)

	// Загружаем в MinIO
	err = s.store.File().Put(
		storageKey,
		bytes.NewReader(req.Content),
		int64(len(req.Content)),
		req.MimeType,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "minio upload failed: %v", err)
	}

	// Сохраняем метаданные в Postgres
	file := &model.File{
		UUID:       fileID.String(),
		Name:       req.Filename,
		MimeType:   req.MimeType,
		Size:       int64(len(req.Content)),
		UploaderID: userID.String(),
		CourseID:   req.CourseId,
		GroupID:    req.GroupId,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
	}

	if err := s.store.Meta().Create(file); err != nil {
		// rollback MinIO
		_ = s.store.File().Delete(storageKey)
		return nil, status.Errorf(codes.Internal, "db insert failed: %v", err)
	}

	return &pb.UploadFileUnaryResponse{
		FileId: fileID.String(),
	}, nil
}

func (s *PublicServer) DownloadFileUnary(
	ctx context.Context,
	req *pb.DownloadFileUnaryRequest,
) (*pb.DownloadFileUnaryResponse, error) {

	fileID, err := uuid.Parse(req.FileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid file id")
	}

	file, err := s.store.Meta().GetByID(fileID.String())
	if err != nil {
		return nil, status.Error(codes.NotFound, "file not found")
	}

	reader, err := s.store.File().Get(file.StorageKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "storage error: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "read failed: %v", err)
	}

	return &pb.DownloadFileUnaryResponse{
		Content:  content,
		Filename: file.Name,
		MimeType: file.MimeType,
	}, nil
}

func (s *PublicServer) GetFile(
	ctx context.Context,
	req *pb.GetFileRequest,
) (*pb.File, error) {

	fileID, err := uuid.Parse(req.FileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	file, err := s.store.Meta().GetByID(fileID.String())
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &pb.File{
		Id:         file.UUID,
		Name:       file.Name,
		MimeType:   file.MimeType,
		Size:       file.Size,
		UploaderId: file.UploaderID,
		CourseId:   file.CourseID,
		GroupId:    file.GroupID,
		CreatedAt:  timestamppb.New(file.CreatedAt),
	}, nil
}

func (s *PublicServer) DeleteFile(
	ctx context.Context,
	req *pb.DeleteFileRequest,
) (*pb.DeleteFileResponse, error) {
	fileID, err := uuid.Parse(req.FileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	file, err := s.store.Meta().GetByID(fileID.String())
	if err != nil {
		return nil, status.Error(codes.NotFound, "file not found")
	}

	if err := s.store.File().Delete(file.StorageKey); err != nil {
		return nil, status.Error(codes.Internal, "storage error")
	}

	return &pb.DeleteFileResponse{Success: true}, nil
}

func (s *PublicServer) ListFilesByCourse(
	ctx context.Context,
	req *pb.ListFilesByCourseRequest,
) (*pb.ListFilesResponse, error) {

	courseID, err := uuid.Parse(req.CourseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}

	files, err := s.store.Meta().ListByCourse(courseID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "db error")
	}

	f := make([]*pb.File, 0)
	for _, file := range files {
		pf := &pb.File{
			Id:         file.UUID,
			Name:       file.Name,
			MimeType:   file.MimeType,
			Size:       file.Size,
			UploaderId: file.UploaderID,
			CourseId:   file.CourseID,
			GroupId:    file.GroupID,
			CreatedAt:  timestamppb.New(file.CreatedAt),
		}
		f = append(f, pf)
	}

	return &pb.ListFilesResponse{
		Files: f,
		Total: int32(len(f)),
	}, nil
}

func (s *PublicServer) ListFilesByGroup(
	ctx context.Context,
	req *pb.ListFilesByGroupRequest,
) (*pb.ListFilesResponse, error) {

	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid group id")
	}

	files, err := s.store.Meta().ListByGroup(groupID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "db error")
	}

	f := make([]*pb.File, 0)
	for _, file := range files {
		pf := &pb.File{
			Id:         file.UUID,
			Name:       file.Name,
			MimeType:   file.MimeType,
			Size:       file.Size,
			UploaderId: file.UploaderID,
			CourseId:   file.CourseID,
			GroupId:    file.GroupID,
			CreatedAt:  timestamppb.New(file.CreatedAt),
		}
		f = append(f, pf)
	}

	return &pb.ListFilesResponse{
		Files: f,
		Total: int32(len(f)),
	}, nil
}
