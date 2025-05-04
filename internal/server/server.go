package server

import (
	"context"

	"github.com/mayye4ka/notpastebin/internal/errs"
	"github.com/mayye4ka/notpastebin/internal/service"
	api "github.com/mayye4ka/notpastebin/pkg/api/go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	service *service.Service
	api.NotPasteBinServer
}

func New(service *service.Service) *server {
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
	resp, err := s.service.GetNote(ctx, req.Hash)
	if err != nil {
		return &api.GetNoteResponse{}, errs.ToStatusError(err)
	}
	return &api.GetNoteResponse{
		Text:       resp.Note,
		IsAdmin:    resp.IsAdmin,
		ReaderHash: resp.ReaderHash,
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
