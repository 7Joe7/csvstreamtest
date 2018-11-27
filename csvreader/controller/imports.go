package controller

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog"
)

const (
	watchedFolderPath = "/opt/data"
)

// ImportsService is a service for managing imports.
type ImportsController struct {
	log     *zerolog.Logger
	src     importSrc
	watcher *watcher.Watcher
}

// NewImportsController creates a new import service.
func NewImportsController(log *zerolog.Logger, src importSrc) *ImportsController {
	return &ImportsController{
		log:     log,
		src:     src,
		watcher: watcher.New(),
	}
}

// Start starts the controller watching for a new file in data folder to be imported
// This feature was not required, just seemed like an ok thing to do, we need to import the file from somewhere
// http request seemed not the way as the file could be really huge potentially but could be of course as well
func (ctr *ImportsController) Start() error {
	ctr.watcher.SetMaxEvents(1)
	ctr.watcher.FilterOps(watcher.Move)
	// Watch this folder for changes.
	if err := ctr.watcher.Add(watchedFolderPath); err != nil {
		return errors.Wrapf(err, "could not watch %s", watchedFolderPath)
	}

	go func() {
		for {
			select {
			case event := <-ctr.watcher.Event:
				ctr.log.Info().Msg(event.String())
			case err := <-ctr.watcher.Error:
				ctr.log.Error().Err(err).Msg("error while watching data folder")
			case <-ctr.watcher.Closed:
				return
			}
		}
	}()

	// automatically import all files currently in the folder
	for path, f := range ctr.watcher.WatchedFiles() {
		if !f.IsDir() {
			err := ctr.Import(path)
			if err != nil {
				ctr.log.Error().Err(err).Msgf("could not import file %s", path)
			}
		}
	}
	return nil
}

// Import imports the file specified by path
// careful about potential file inclusion via variable - if ever the path is provided by user
func (ctr *ImportsController) Import(path string) error {
	ctr.log.Info().Msgf("found file %s", path)
	/* #nosec */
	f, err := os.Open(path)
	if err != nil {
		ctr.log.Error().Err(err).Msg("could not open the file")
	}
	defer f.Close()

	switch filepath.Ext(path) {
	case ".tar":
		err = ctr.src.ImportTar(f)
	case ".csv":
		err = ctr.src.ImportCSV(f)
	default:
		ctr.log.Info().Msgf("file %s skipped, unsupported extension %s", path, filepath.Ext(path))
	}
	if err != nil {
		return err
	}
	return nil
}

// Stop stops watching the data folder
func (ctr *ImportsController) Stop() {
	ctr.watcher.Close()
}
