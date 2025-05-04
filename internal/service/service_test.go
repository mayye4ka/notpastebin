package service

import (
	"context"
	"strings"
	"testing"

	"github.com/mayye4ka/notpastebin/internal/errs"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

const (
	note1     = "note 1"
	note2     = "note 2"
	adminHash = "admin_hash"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	dbMock  *MockRepository
	service *Service
}

func (s *ServiceTestSuite) TestCreateNote_ValidationFailed() {
	adminHash, readerHash, err := s.service.CreateNote(s.ctx, "")
	s.Empty(adminHash)
	s.Empty(readerHash)
	s.ErrorContains(err, "invalid input: empty text")
}

func (s *ServiceTestSuite) TestCreateNote_Collisition() {
	s.dbMock.EXPECT().CreateNote(s.ctx, note1, gomock.Any(), gomock.Any()).Return(errs.ErrCollision)
	s.dbMock.EXPECT().CreateNote(s.ctx, note1, gomock.Any(), gomock.Any()).Return(nil)
	adminHash, readerHash, err := s.service.CreateNote(s.ctx, note1)
	s.NoError(err)
	s.NotEmpty(adminHash)
	s.NotEmpty(readerHash)
	s.NotEqual(adminHash, readerHash)
}

func (s *ServiceTestSuite) TestCreateNote() {
	s.dbMock.EXPECT().CreateNote(s.ctx, note1, gomock.Any(), gomock.Any()).Return(nil)
	adminHash, readerHash, err := s.service.CreateNote(s.ctx, note1)
	s.NoError(err)
	s.NotEmpty(adminHash)
	s.NotEmpty(readerHash)
	s.NotEqual(adminHash, readerHash)
}

func (s *ServiceTestSuite) TestUpdateNote_ValidationFailed() {
	err := s.service.UpdateNote(s.ctx, "", strings.Repeat("a", 70000))
	s.ErrorContains(err, "invalid input: text longer than 65535 symbols")
}

func (s *ServiceTestSuite) TestUpdateNote() {
	s.dbMock.EXPECT().UpdateNote(s.ctx, adminHash, note2).Return(nil)
	err := s.service.UpdateNote(s.ctx, adminHash, note2)
	s.NoError(err)
}

func (s *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.dbMock = NewMockRepository(ctrl)
	s.service = New(s.dbMock)
	s.ctx = context.Background()
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
