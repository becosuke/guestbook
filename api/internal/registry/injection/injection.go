package injection

import (
	"context"
	"log"
	"sync"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapters/controller"
	"github.com/becosuke/guestbook/api/internal/adapters/gateway"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/drivers/grpcserver"
	"github.com/becosuke/guestbook/api/internal/drivers/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/pbgo"
)

type Injection interface {
	InjectConfig() *config.Config
	InjectLogger() *zap.Logger
	InjectAuthFunc() grpc_auth.AuthFunc
	InjectGrpcServer() *grpc.Server
	InjectController() pbgo.GuestbookServiceServer
	InjectUsecase() post.Usecase
	InjectBoundary() controller.Boundary
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
		Logger                 *zap.Logger
		GrpcServer             *grpc.Server
		GuestbookServiceServer pbgo.GuestbookServiceServer
		Usecase                post.Usecase
		Boundary               controller.Boundary
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

func (i *injectionImpl) InjectLogger() *zap.Logger {
	actual, _ := i.store.LoadOrStore("logger", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			cfg := zap.NewProductionConfig()
			cfg.Level = zap.NewAtomicLevelAt(i.InjectConfig().LogLevel)
			cfg.DisableStacktrace = true
			cfg.Sampling = nil
			cfg.OutputPaths = []string{"stdout"}
			cfg.ErrorOutputPaths = []string{"stderr"}
			cfg.InitialFields = map[string]interface{}{
				"service": i.serviceName,
				"version": i.version,
				"env":     i.InjectConfig().Environment.String(),
			}
			l, err := cfg.Build()
			if err != nil {
				log.Fatal(err)
			}
			i.container.Logger = l
		})
	}
	return i.container.Logger
}

func (i *injectionImpl) InjectAuthFunc() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
}

func (i *injectionImpl) InjectGrpcServer() *grpc.Server {
	actual, _ := i.store.LoadOrStore("grpcServer", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.GrpcServer = grpcserver.NewGrpcServer(i.InjectLogger(), i.InjectAuthFunc())
		})
	}
	return i.container.GrpcServer
}

func (i *injectionImpl) InjectController() pbgo.GuestbookServiceServer {
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

func (i *injectionImpl) InjectBoundary() controller.Boundary {
	actual, _ := i.store.LoadOrStore("boundary", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Boundary = controller.NewBoundary()
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
