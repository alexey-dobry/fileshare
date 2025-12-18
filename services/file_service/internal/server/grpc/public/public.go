package public

import (
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

func (s *PublicServer) UploadFile(
	stream pb.FileService_UploadFileServer,
) error {

	ctx := stream.Context()

	userID, err := utils.UserIDFromContext(ctx)
	if err != nil {
		return err
	}

	fileID := uuid.New()
	storageKey := fmt.Sprintf("files/%s", fileID)

	var (
		meta  *pb.FileMetadata
		pipeR *io.PipeReader
		pipeW *io.PipeWriter
		size  int64
	)

	pipeR, pipeW = io.Pipe()

	// Upload to MinIO in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.store.File().Put(
			storageKey,
			pipeR,
			-1,
			meta.GetMimeType(),
		)
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			pipeW.Close()
			break
		}
		if err != nil {
			return err
		}

		switch data := req.Data.(type) {
		case *pb.UploadFileRequest_Metadata:
			meta = data.Metadata

		case *pb.UploadFileRequest_Chunk:
			n, err := pipeW.Write(data.Chunk)
			if err != nil {
				return err
			}
			size += int64(n)
		}
	}

	if err := <-errCh; err != nil {
		return status.Errorf(codes.Internal, "upload failed: %v", err)
	}

	file := &model.File{
		UUID:       fileID.String(),
		Name:       meta.Filename,
		MimeType:   meta.MimeType,
		Size:       size,
		UploaderID: userID.String(),
		CourseID:   meta.CourseId,
		GroupID:    meta.GroupId,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
	}

	if err := s.store.Meta().Create(file); err != nil {
		_ = s.store.File().Delete(storageKey) // rollback
		return status.Errorf(codes.Internal, "db error: %v", err)
	}

	return stream.SendAndClose(&pb.UploadFileResponse{
		FileId: fileID.String(),
	})
}

func (s *PublicServer) DownloadFile(
	req *pb.DownloadFileRequest,
	stream pb.FileService_DownloadFileServer,
) error {
	fileID, err := uuid.Parse(req.FileId)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid file id")
	}

	file, err := s.store.Meta().GetByID(fileID.String())
	if err != nil {
		return status.Error(codes.NotFound, "file not found")
	}

	reader, err := s.store.File().Get(file.StorageKey)
	if err != nil {
		return status.Error(codes.Internal, "storage error")
	}
	defer reader.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if err := stream.Send(&pb.DownloadFileResponse{
				Chunk: buf[:n],
			}); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
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
