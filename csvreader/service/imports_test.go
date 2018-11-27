package service

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/csvreader/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// TestNewImportsService tests NewImportsService
func TestNewImportsService(t *testing.T) {
	log := zerolog.New(ioutil.Discard)
	importerMock := &mocks.ImporterClient{}

	src := NewImportsService(&log, importerMock)
	assert.Equal(t, &log, src.log)
	assert.Equal(t, importerMock, src.importer)
}

// TestImportsService_ImportTar tests ImportTar
func TestImportsService_ImportTar(t *testing.T) {
	emptyBuffer := bytes.NewBuffer([]byte{})
	tarWriter := tar.NewWriter(emptyBuffer)
	err := tarWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	fullBuffer := bytes.NewBuffer([]byte{})
	fullTarWriter := tar.NewWriter(fullBuffer)
	fullTarWriter.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "test.csv",
	})
	err = fullTarWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}
	err = fullTarWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		tarReader    io.Reader
		expectations []func(*mocks.ImporterClient, *mocks.ImportSrc)
		expectError  bool
	}{
		{
			name:        "no tar headers should fail",
			tarReader:   bytes.NewBufferString("hello world"),
			expectError: true,
		},
		{
			name:      "no files in tar should do nothing",
			tarReader: emptyBuffer,
		},
		{
			name:      "tar with a CSV file should call ImportCSV",
			tarReader: fullBuffer,
			expectations: []func(*mocks.ImporterClient, *mocks.ImportSrc){
				func(importer *mocks.ImporterClient, src *mocks.ImportSrc) {
					// prepare expected param
					b := bytes.NewBuffer([]byte{})
					w := tar.NewWriter(b)
					w.WriteHeader(&tar.Header{
						Typeflag: tar.TypeReg,
						Name:     "test.csv",
					})
					w.Flush()
					w.Close()
					r := tar.NewReader(b)
					r.Next()
					src.On("ImportCSV", r).Return(nil)
				},
			},
		},
	}

	for idx, test := range tests {
		idx, test := idx, test
		t.Run(fmt.Sprintf("%d-%v", idx, test.name), func(t *testing.T) {
			log := zerolog.New(ioutil.Discard)
			importerMock := &mocks.ImporterClient{}
			srcMock := &mocks.ImportSrc{}
			src := &ImportsService{
				log:       &log,
				importer:  importerMock,
				importSrc: srcMock,
			}
			// we provide all mocks with expected function calls and return values declared in the test definition
			for _, expectation := range test.expectations {
				expectation(importerMock, srcMock)
			}
			err := src.ImportTar(test.tarReader)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			importerMock.AssertExpectations(t)
			srcMock.AssertExpectations(t)
		})
	}
}

// TestImportsService_ImportCSV tests ImportCSV
func TestImportsService_ImportCSV(t *testing.T) {
	tests := []struct {
		name         string
		csvReader    io.Reader
		expectations []func(*mocks.ImporterClient, *mocks.ImportSrc)
		expectError  bool
	}{
		{
			name:      "unknown headers should do nothing",
			csvReader: bytes.NewBufferString("hello,world"),
		},
		{
			name:      "correct headers no data should call import clients",
			csvReader: bytes.NewBufferString(strings.Join(userImportHeaders, ",")),
			expectations: []func(*mocks.ImporterClient, *mocks.ImportSrc){
				func(importer *mocks.ImporterClient, src *mocks.ImportSrc) {
					streamMock := &mocks.Importer_ImportClientsClient{}
					importer.On("ImportClients", context.Background()).Return(streamMock, nil)
					streamMock.On("CloseAndRecv").Return(&model.ImportReport{Success: true}, nil)
				},
			},
		},
		{
			name:      "correct headers with a client should call import clients and provide client with correct mobile number",
			csvReader: bytes.NewBufferString(fmt.Sprintf("%s\n1,test,test@test.com,(123) 456789", strings.Join(userImportHeaders, ","))),
			expectations: []func(*mocks.ImporterClient, *mocks.ImportSrc){
				func(importer *mocks.ImporterClient, src *mocks.ImportSrc) {
					streamMock := &mocks.Importer_ImportClientsClient{}
					importer.On("ImportClients", context.Background()).Return(streamMock, nil)
					streamMock.On("Send", &model.Client{Id: 1, Name: "test", Email: "test@test.com", MobileNumber: "+44123456789"}).Return(nil)
					streamMock.On("CloseAndRecv").Return(&model.ImportReport{Success: true}, nil)
				},
			},
		},
	}

	for idx, test := range tests {
		idx, test := idx, test
		t.Run(fmt.Sprintf("%d-%v", idx, test.name), func(t *testing.T) {
			log := zerolog.New(ioutil.Discard)
			importerMock := &mocks.ImporterClient{}
			srcMock := &mocks.ImportSrc{}
			src := &ImportsService{
				log:       &log,
				importer:  importerMock,
				importSrc: srcMock,
			}
			// we provide all mocks with expected function calls and return values declared in the test definition
			for _, expectation := range test.expectations {
				expectation(importerMock, srcMock)
			}
			err := src.ImportCSV(test.csvReader)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			importerMock.AssertExpectations(t)
			srcMock.AssertExpectations(t)
		})
	}
}
