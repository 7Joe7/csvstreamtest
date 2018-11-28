package service

import (
	"io"

	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/common/rpc"
	"github.com/7joe7/csvstreamtest/common/types"
	"github.com/7joe7/csvstreamtest/dbmanager/repository"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ClientsService is a service for managing clients.
type ClientsService struct {
	log     *zerolog.Logger
	clients repository.ClientRepo
}

// NewClientsService creates a new client service.
func NewClientsService(log *zerolog.Logger, clients repository.ClientRepo) *ClientsService {
	return &ClientsService{
		log:     log,
		clients: clients,
	}
}

// ImportClients imports clients.
func (src *ClientsService) ImportClients(srv rpc.Importer_ImportClientsServer) (err error) {
	src.log.Info().Msg("importing clients")
	defer func() {
		if err != nil {
			err = srv.SendAndClose(&model.ImportReport{
				Error: err.Error(),
			})
			if err != nil {
				src.log.Error().Err(err).Msg("could not send import report")
			}
			return
		}
		err = srv.SendAndClose(&model.ImportReport{Success: true})
		if err != nil {
			src.log.Error().Err(err).Msg("could not send import report")
		}
	}()
	var trClients repository.ClientRepo
	trClients, err = src.clients.NewWithTransaction()
	if err != nil {
		src.log.Error().Err(err).Msg("could not start transaction")
		return errors.Wrap(err, "could not start transaction")
	}
	defer func() {
		switch err {
		case nil, io.EOF:
			err = trClients.Commit()
		default:
			rollbackErr := trClients.Rollback()
			if rollbackErr != nil {
				err = errors.Wrapf(err, "could not rollback transaction: %v", rollbackErr)
			}
		}
		if err != nil {
			src.log.Error().Err(err).Msg("could not start transaction")
		}
	}()
	for {
		var client *model.Client
		client, err = srv.Recv()
		src.log.Debug().Msgf("received client: %v", client)
		switch err {
		case nil:
			if client.Id > 0 {
				_, err = trClients.Retrieve(client.Id)
				switch err {
				case nil:
					src.log.Info().Int32("Id", client.Id).Str("email", client.Email).Msg("updating client")
					err = trClients.Update(client)
					if err != nil {
						err = errors.Wrapf(err, "could not update client: %d", client.Id)
					}
				case types.ErrNotFound:
					src.log.Info().Int32("Id", client.Id).Str("email", client.Email).Msg("creating client")
					err = trClients.Store(client)
					if err != nil {
						err = errors.Wrapf(err, "could not create client: %d", client.Id)
					}
				default:
					err = errors.Wrapf(err, "could not retrieve client: %d", client.Id)
				}
				if err != nil {
					return
				}
				continue
			}
			src.log.Info().Str("email", client.Email).Msg("creating client")
			err = trClients.Store(client)
			if err != nil {
				err = errors.Wrapf(err, "could not create client: %s", client.Email)
				return
			}
		default:
			return
		}
	}
}
