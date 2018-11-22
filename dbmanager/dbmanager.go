package main

import (
	"net"
	"os"
	"strconv"
	"syscall"

	"fmt"

	"github.com/7joe7/csvstreamtest/common/catcher"
	"github.com/7joe7/csvstreamtest/common/dbconn"
	"github.com/7joe7/csvstreamtest/common/logger"
	"github.com/7joe7/csvstreamtest/common/rpc"
	"github.com/7joe7/csvstreamtest/dbmanager/repository"
	"github.com/7joe7/csvstreamtest/dbmanager/service"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const (
	cfgLogLevel      = "dbmanager/log/level"
	cfgMySQLUsername = "dbmanager/mysql/username"
	/* #nosec */
	cfgMySQLPassword = "dbmanager/mysql/password"
	cfgMySQLHost     = "dbmanager/mysql/host"
	cfgMySQLPort     = "dbmanager/mysql/port"
	cfgMySQLDatabase = "dbmanager/mysql/database"
	cfgGRPCPort      = "dbmanager/grpc/port"
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
	pairs, _, err := kv.List("dbmanager", nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not find configuration")
	}
	var dbUsername, dbPwd, dbHost, dbPort, database, grpcPort []byte
	for _, pair := range pairs {
		if pair.Key != cfgMySQLPassword {
			zlog.Info().Msgf("configuration key '%s', value '%s'", pair.Key, pair.Value)
		}
		switch pair.Key {
		case cfgLogLevel:
			logger.SetGlobalLevel(string(pair.Value))
		case cfgMySQLUsername:
			dbUsername = pair.Value
		case cfgMySQLPassword:
			zlog.Info().Msgf("configuration key '%s', value '***********'", pair.Key)
			dbPwd = pair.Value
		case cfgMySQLHost:
			dbHost = pair.Value
		case cfgMySQLPort:
			dbPort = pair.Value
		case cfgMySQLDatabase:
			database = pair.Value
		case cfgGRPCPort:
			grpcPort = pair.Value
		}
	}

	zlog.Info().Msg("dbmanager starting up...")

	dbPortInt, err := strconv.Atoi(string(dbPort))
	if err != nil {
		zlog.Fatal().Err(err).Msg("port is not a number")
	}
	crmDB, err := dbconn.GORMOpen(string(dbUsername), string(dbPwd), string(dbHost), string(database), dbPortInt)
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not initialize DB connection")
	}

	clients := repository.NewORMClient(crmDB)
	// db connection test
	_, err = clients.Find(nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("could not find clients")
	}
	clientsSrc := service.NewClientsService(zlog, clients)

	tcpSocket, err := net.Listen("tcp", fmt.Sprintf(":%s", string(grpcPort)))
	if err != nil {
		zlog.Fatal().Err(err).Msgf("could not listen on port %s", string(grpcPort))
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterImporterServer(grpcServer, clientsSrc)
	go func() {
		err = grpcServer.Serve(tcpSocket)
		if err != nil {
			zlog.Fatal().Err(err).Msg("grpc failure")
		}
	}()

	zlog.Info().Msg("dbmanager startup complete")

	sig.Wait()

	zlog.Info().Msg("dbmanager shutting down...")

	grpcServer.GracefulStop()

	zlog.Info().Msg("dbmanager shutdown complete")
}
