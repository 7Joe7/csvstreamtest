package service

import (
	"archive/tar"
	"context"
	"encoding/csv"
	"io"
	"path/filepath"
	"reflect"

	"regexp"
	"strings"

	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/common/rpc"
	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var userImportHeaders = []string{"id", "name", "email", "mobile_number"}

// ImportsService is a service for managing imports.
type ImportsService struct {
	log       *zerolog.Logger
	importer  rpc.ImporterClient
	importSrc ImportSrc
}

// NewImportsService creates a new import service.
func NewImportsService(log *zerolog.Logger, importer rpc.ImporterClient) *ImportsService {
	src := &ImportsService{
		log:      log,
		importer: importer,
	}
	src.importSrc = src
	return src
}

// ImportTar imports from .tar file.
func (src *ImportsService) ImportTar(reader io.Reader) error {
	src.log.Info().Msg("unpacking tar")
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "could not stream the file")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// not interested in directories
		case tar.TypeReg:
			src.log.Info().Msgf("parsing file: %s", header.Name)
			// we depend on the suffix of the file, ideal would be to check encoding as well whether it is a text file
			switch filepath.Ext(header.Name) { // potentially we can support other import file formats
			case ".csv":
				err = src.importSrc.ImportCSV(tarReader)
				if err != nil {
					src.log.Error().Err(err).Msg("error importing CSV file")
				}
			}
		default:
			src.log.Error().Msgf("unable to figure out type: %c of file %s", header.Typeflag, header.Name)
		}
	}
	return nil
}

// ImportCSV imports contents of a CSV file
func (src *ImportsService) ImportCSV(reader io.Reader) (err error) {
	src.log.Info().Msg("importing CSV file")
	r := csv.NewReader(reader)

	dec, err := csvutil.NewDecoder(r)
	if err != nil {
		return errors.Wrap(err, "could not create initialize CSV decoder")
	}
	header := dec.Header()
	// based on header we import the structure we found in the csv file
	switch {
	// comparing fields through a map would be more efficient but this is not the bottleneck so not worth it
	case reflect.DeepEqual(header, userImportHeaders):
		stream, err := src.importer.ImportClients(context.Background())
		if err != nil {
			return errors.Wrap(err, "could not connect to importer")
		}
		defer func() {
			var report *model.ImportReport
			report, err = stream.CloseAndRecv()
			if err != nil {
				err = errors.Wrap(err, "could not close and receive from the server")
			}
			src.log.Info().Msgf("CSV import: %v", report)
		}()
		for {
			client := &model.Client{}
			if err = dec.Decode(&client); err == io.EOF {
				break
			} else if err != nil {
				return errors.Wrap(err, "could not read entry")
			}
			client.MobileNumber = src.unifyClientMobileNumber(client.MobileNumber)
			src.log.Debug().Msgf("sending client: %v", client)
			err = stream.Send(client)
			if err != nil {
				return errors.Wrap(err, "could not send another client")
			}
		}
	default:
		// we could support other headers
	}
	return nil
}

// unifyClientMobileNumber converts client mobile number to the same format
// removing any brackets or spaces and prefixing with +44 as country code
func (src *ImportsService) unifyClientMobileNumber(mobileNumber string) string {
	re := regexp.MustCompile("[() ]")
	mobileNumber = re.ReplaceAllString(mobileNumber, "")
	if !strings.HasPrefix(mobileNumber, "+") {
		mobileNumber = "+44" + mobileNumber
	}
	return mobileNumber
}
