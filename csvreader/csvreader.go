package main

import (
	"os"
	"syscall"

	"fmt"

	"github.com/7joe7/csvstreamtest/common/catcher"
	"github.com/7joe7/csvstreamtest/common/logger"
	"github.com/7joe7/csvstreamtest/common/rpc"
	"github.com/7joe7/csvstreamtest/csvreader/controller"
	"github.com/7joe7/csvstreamtest/csvreader/service"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const (
	cfgLogLevel = "csvreader/log/level"
	cfgGRPCPort = "csvreader/grpc/imports/port"
	cfgGRPCHost = "csvreader/grpc/imports/host"
)

func main() {
	sig := catcher.NewSignal(syscall.SIGINT, syscall.SIGTERM)

	zlog := logger.NewZeroLog(os.Stderr, logger.Pretty)

	consulCfg := api.DefaultConfig()
	consulCfg.Address = "consul:8500"
	consul, err := api.NewClient(consulCfg)
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not initialize consul client")
	}
	// we could watch for changes ideally
	kv := consul.KV()
	pairs, _, err := kv.List("csvreader", nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not find configuration")
	}
	var grpcPort, grpcHost []byte
	for _, pair := range pairs {
		zlog.Info().Msgf("configuration key '%s', value '%s'", pair.Key, pair.Value)
		switch pair.Key {
		case cfgLogLevel:
			logger.SetGlobalLevel(string(pair.Value))
		case cfgGRPCPort:
			grpcPort = pair.Value
		case cfgGRPCHost:
			grpcHost = pair.Value
		}
	}

	zlog.Info().Msg("csvreader starting up...")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", string(grpcHost), string(grpcPort)), grpc.WithInsecure()) // TODO security
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not dial grpc importer")
	}

	dbmanagerRepo := rpc.NewImporterClient(conn)
	importSrc := service.NewImportsService(zlog, dbmanagerRepo)
	importCtr := controller.NewImportsController(zlog, importSrc)
	err = importCtr.Start()
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not initialize import controller")
	}

	zlog.Info().Msg("csvreader startup complete")

	sig.Wait()

	zlog.Info().Msg("csvreader shutting down...")

	importCtr.Stop()
	err = conn.Close()
	if err != nil {
		zlog.Error().Err(err).Msg("error closing grpc connection")
	}

	zlog.Info().Msg("csvreader shutdown complete")
}
