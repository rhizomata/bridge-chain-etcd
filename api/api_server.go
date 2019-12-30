package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rhizomata/bridge-chain-etcd/kernel"
	"github.com/rhizomata/bridge-chain-etcd/protocol"
)

// Server ..
type Server struct {
	router         *gin.Engine
	builtinService *BuiltinService
	err            chan error
}

func (server *Server) Error() <-chan error {
	return server.err
}

// StartServer ..
func StartServer(kernel *kernel.Kernel, listenAddress string) (server *Server) {
	server = new(Server)
	server.err = make(chan error)
	server.builtinService = &BuiltinService{kernel: kernel}
	server.router = gin.Default()

	v1 := server.router.Group(protocol.V1Path)
	{
		v1.HEAD(protocol.HealthPath, server.builtinService.health)
		v1.GET(protocol.HealthPath, server.builtinService.health)
		v1.POST(protocol.AddJobPath, server.builtinService.addJob)
	}

	go func() {
		err := server.router.Run(listenAddress)
		if err != nil {
			log.Fatal("Cannot Start API Server")
		}
	}()

	return server
}
