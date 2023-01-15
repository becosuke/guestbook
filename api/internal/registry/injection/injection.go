package injection

import (
	"github.com/becosuke/guestbook/api/internal/adapters/boundary"
	"github.com/becosuke/guestbook/api/internal/adapters/controller"
	"github.com/becosuke/guestbook/api/internal/adapters/gateway"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/drivers/grpcserver"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/pkg/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/api/pb"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type Injection interface {
	InjectConfig() *config.Config
	InjectLogger() logger.Logger
	InjectGrpcServer() *grpc.Server
	InjectController() pb.GuestbookServiceServer
	InjectUsecase() post.Usecase
	InjectBoundary() boundary.Boundary
	InjectRepository() post.Repository
	InjectGenerator() post.Generator
	InjectSyncmap() syncmap.Syncmap
}

func NewInjection(serviceName, version string) Injection {
	return &injectionImpl{serviceName: serviceName, version: version}
}

type injectionImpl struct {
	container struct {
		Config                 *config.Config
		Logger                 logger.Logger
		GrpcServer             *grpc.Server
		GuestbookServiceServer pb.GuestbookServiceServer
		Usecase                post.Usecase
		Boundary               boundary.Boundary
		Repository             post.Repository
		Generator              post.Generator
		Syncmap                syncmap.Syncmap
	}
	serviceName string
	version     string
	store       sync.Map
}

func (i *injectionImpl) InjectConfig() *config.Config {
	actual, _ := i.store.LoadOrStore("config", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Config = config.NewConfig()
		})
	}
	return i.container.Config
}

func (i *injectionImpl) InjectLogger() logger.Logger {
	actual, _ := i.store.LoadOrStore("logger", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			l, err := logger.NewLogger(
				i.InjectConfig().LogLevel,
				i.serviceName,
				i.version,
				i.InjectConfig().Environment.String(),
			)
			if err != nil {
				log.Fatal(err)
			}
			i.container.Logger = l
		})
	}
	return i.container.Logger
}

func (i *injectionImpl) InjectGrpcServer() *grpc.Server {
	actual, _ := i.store.LoadOrStore("grpcServer", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.GrpcServer = grpcserver.NewGrpcServer()
		})
	}
	return i.container.GrpcServer
}

func (i *injectionImpl) InjectController() pb.GuestbookServiceServer {
	actual, _ := i.store.LoadOrStore("controller", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.GuestbookServiceServer = controller.NewGuestbookServiceServer(
				i.InjectConfig(),
				i.InjectUsecase(),
				i.InjectBoundary(),
			)
		})
	}
	return i.container.GuestbookServiceServer
}

func (i *injectionImpl) InjectUsecase() post.Usecase {
	actual, _ := i.store.LoadOrStore("usecase", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Usecase = usecase.NewUsecase(i.InjectConfig(), i.InjectRepository())
		})
	}
	return i.container.Usecase
}

func (i *injectionImpl) InjectBoundary() boundary.Boundary {
	actual, _ := i.store.LoadOrStore("boundary", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Boundary = boundary.NewBoundary()
		})
	}
	return i.container.Boundary
}

func (i *injectionImpl) InjectRepository() post.Repository {
	actual, _ := i.store.LoadOrStore("repository", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Repository = gateway.NewRepository(i.InjectConfig(), i.InjectSyncmap(), i.InjectGenerator())
		})
	}
	return i.container.Repository
}

func (i *injectionImpl) InjectGenerator() post.Generator {
	actual, _ := i.store.LoadOrStore("generator", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Generator = gateway.NewGenerator()
		})
	}
	return i.container.Generator
}

func (i *injectionImpl) InjectSyncmap() syncmap.Syncmap {
	actual, _ := i.store.LoadOrStore("syncmap", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Syncmap = syncmap.NewSyncmap()
		})
	}
	return i.container.Syncmap
}
