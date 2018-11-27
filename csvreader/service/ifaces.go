package service

import "io"

type ImportSrc interface {
	ImportCSV(reader io.Reader) error
	ImportTar(reader io.Reader) error
}
