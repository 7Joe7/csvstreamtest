package controller

import "io"

type importSrc interface {
	ImportCSV(reader io.Reader) error
	ImportTar(reader io.Reader) error
}
