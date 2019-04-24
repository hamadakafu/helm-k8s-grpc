package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	api "sample-grpc/proto"
	"sample-grpc/service"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer() error {
	dbConfig := getDBConfig()
	serverConfig := getServerConfig()

	lis, err := net.Listen("tcp", ":"+serverConfig.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.host, dbConfig.user, dbConfig.password, dbConfig.dbname, dbConfig.sslmode))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	server := service.NewBookServiceServer(db)
	s := grpc.NewServer()
	api.RegisterBookServiceServer(s, server)
	reflection.Register(s)
	log.Println("starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}

type dbConfig struct {
	host     string
	user     string
	password string
	dbname   string
	sslmode  string
}

type serverConfig struct {
	port string
}

func getServerConfig() serverConfig {
	return serverConfig{
		port: os.Getenv("SERVER_PORT"),
	}
}

func getDBConfig() dbConfig {
	return dbConfig{
		host:     os.Getenv("POSTGRES_HOST"),
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		dbname:   os.Getenv("POSTGRES_DB"),
		sslmode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}
