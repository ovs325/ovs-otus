package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	pb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/grpc/event_service"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
)

type GrpcServer struct {
	logic hd.AbstractLogic
	log   lg.Logger
	srv   *grpc.Server
}

func NewGrpcServer(l hd.AbstractLogic, logger lg.Logger) *GrpcServer {
	return &GrpcServer{logic: l, log: logger}
}

type eventServiceServer struct {
	pb.UnimplementedEventServiceServer
}

// Новый Event-сервис
func NewEventServiceServer() *eventServiceServer {
	return &eventServiceServer{}
}

func (s *GrpcServer) Start(cnf *cf.Config) error {
	s.log.Info("the GRPC-server starts")

	dsn := fmt.Sprintf("%s:%s", cnf.GrpcServer.Host, cnf.GrpcServer.Port)

	// Иннициализация gRPS-serverа.
	lis, err := net.Listen("tcp", dsn)
	if err != nil {
		s.log.Error("GRPC-Server failed to listen", "error", err.Error())
		return err
	}
	var opts []grpc.ServerOption
	s.srv = grpc.NewServer(opts...)
	// Регистрация, созданная автоматически
	pb.RegisterEventServiceServer(s.srv, NewEventServiceServer())
	// Конец Иннициализации gRPS-server-а

	s.log.Info("GRPC-Server started successfully!", "address", dsn)
	if err := s.srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func (s *GrpcServer) Stop() error {
	fmt.Println("GRPC-Server forced to shutdown")
	s.srv.Stop()
	fmt.Println("GRPC-Server Shutdown is successful!!")
	return nil
}
