package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil

}

func (app *Config) gRPCListen() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRpcPort))
	if err != nil {
		log.Println("failed to listen to grpc", err)
		return err
	}
	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.models})

	log.Println("started grpc on port 50001")

	err = s.Serve(listen)
	if err != nil {
		log.Println("failed to listen to grpc", err)
		return err
	}
	return nil
}
