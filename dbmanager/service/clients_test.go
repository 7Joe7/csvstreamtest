package service

import (
	"fmt"
	"io/ioutil"
	"testing"

	"io"

	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/dbmanager/service/mocks"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// TestNewClientsService tests NewClientsService
func TestNewClientsService(t *testing.T) {
	log := zerolog.New(ioutil.Discard)
	clientsRepoMock := &mocks.ClientRepo{}
	src := NewClientsService(&log, clientsRepoMock)
	assert.Equal(t, &log, src.log)
	assert.Equal(t, clientsRepoMock, src.clients)
}

// TestClientsService_ImportClients tests ImportClients
func TestClientsService_ImportClients(t *testing.T) {
	testingClient := &model.Client{Name: "some", Email: "test", MobileNumber: "number"}
	tests := []struct {
		name         string
		expectations func(*mocks.ClientRepo, *mocks.Importer_ImportClientsServer)
		expectError  bool
	}{
		{
			name: "no clients received return success",
			expectations: func(repo *mocks.ClientRepo, stream *mocks.Importer_ImportClientsServer) {
				repo.On("NewWithTransaction").Return(repo, nil)
				stream.On("SendAndClose", &model.ImportReport{Success: true}).Return(nil)
				stream.On("Recv").Return(nil, io.EOF)
				repo.On("Commit").Return(nil)
			},
		},
		{
			name: "clients received should store",
			expectations: func(repo *mocks.ClientRepo, stream *mocks.Importer_ImportClientsServer) {
				repo.On("NewWithTransaction").Return(repo, nil)
				stream.On("SendAndClose", &model.ImportReport{Success: true}).Return(nil)
				stream.On("Recv").Return(testingClient, nil).Once()
				repo.On("Store", testingClient).Return(nil)
				stream.On("Recv").Return(nil, io.EOF).Once()
				repo.On("Commit").Return(nil)
			},
		},
		{
			name: "clients received with error should rollback",
			expectations: func(repo *mocks.ClientRepo, stream *mocks.Importer_ImportClientsServer) {
				repo.On("NewWithTransaction").Return(repo, nil)
				stream.On("SendAndClose", &model.ImportReport{Error: "could not create client: test: some error"}).Return(nil)
				stream.On("Recv").Return(testingClient, nil)
				repo.On("Store", testingClient).Return(errors.New("some error"))
				repo.On("Rollback").Return(nil)
			},
		},
	}

	for idx, test := range tests {
		idx, test := idx, test
		t.Run(fmt.Sprintf("%d-%v", idx, test.name), func(t *testing.T) {
			log := zerolog.New(ioutil.Discard)
			clientsRepoMock := &mocks.ClientRepo{}
			streamMock := &mocks.Importer_ImportClientsServer{}
			src := &ClientsService{
				log:     &log,
				clients: clientsRepoMock,
			}
			// we provide all mocks with expected function calls and return values declared in the test definition
			if test.expectations != nil {
				test.expectations(clientsRepoMock, streamMock)
			}
			err := src.ImportClients(streamMock)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			clientsRepoMock.AssertExpectations(t)
			streamMock.AssertExpectations(t)
		})
	}
}
