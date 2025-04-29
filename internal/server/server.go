package server

import (
	"context"

	"github.com/mayye4ka/notpastebin/internal/errs"
	api "github.com/mayye4ka/notpastebin/pkg/api/go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	service Service
	api.NotPasteBinServer
}

type Service interface {
	CreateNote(ctx context.Context, text string) (string, string, error)
	GetNote(ctx context.Context, hash string) (string, error)
	UpdateNote(ctx context.Context, hash, text string) error
	DeleteNote(ctx context.Context, hash string) error
}

func New(service Service) *server {
	return &server{
		service: service,
	}
}

func (s *server) CreateNote(ctx context.Context, req *api.CreateNoteRequest) (*api.CreateNoteResponse, error) {
	adminHash, readerHash, err := s.service.CreateNote(ctx, req.Text)
	if err != nil {
		return &api.CreateNoteResponse{}, errs.ToStatusError(err)
	}
	return &api.CreateNoteResponse{
		AdminHash:  adminHash,
		ReaderHash: readerHash,
	}, nil
}

func (s *server) GetNote(ctx context.Context, req *api.GetNoteRequest) (*api.GetNoteResponse, error) {
	text, err := s.service.GetNote(ctx, req.Hash)
	if err != nil {
		return &api.GetNoteResponse{}, errs.ToStatusError(err)
	}
	return &api.GetNoteResponse{
		Text: text,
	}, nil
}

func (s *server) UpdateNote(ctx context.Context, req *api.UpdateNoteRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateNote(ctx, req.AdminHash, req.Text)
	if err != nil {
		return &emptypb.Empty{}, errs.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) DeleteNote(ctx context.Context, req *api.DeleteNoteRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteNote(ctx, req.AdminHash)
	if err != nil {
		return &emptypb.Empty{}, errs.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}
