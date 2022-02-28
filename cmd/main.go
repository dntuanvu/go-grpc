package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/joho/godotenv"
	"github.com/oklog/oklog/pkg/group"
	"github.com/victordinh/gokit-grpc/pb"
	"github.com/victordinh/gokit-grpc/service"
	"github.com/victordinh/gokit-grpc/transports"
	"google.golang.org/grpc"

	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"github.com/victordinh/gokit-grpc/endpoints"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

func main() {
	godotenv.Load()

	var (
		logger   log.Logger
		server   = os.Getenv("SERVER_HOST")
		httpAddr = net.JoinHostPort(server, envString("SERVER_HTTP_PORT", defaultHTTPPort))
		grpcAddr = net.JoinHostPort(server, envString("SERVER_GRPC_PORT", defaultGRPCPort))
	)

	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	addservice := service.NewService(logger)
	addendpoint := endpoints.MakeEndpoints(addservice)
	grpcServer := transports.NewGRPCServer(addendpoint, logger)
	httpServer := transports.NewHTTPServer(addendpoint, logger)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpServer)
		}, func(error) {
			httpListener.Close()
		})
	}

	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by
			// the here demonstrated zipkin tracing middleware.
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterVictorServiceServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
