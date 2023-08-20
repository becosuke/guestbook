package injection

import (
	"context"
	"log"
	"sync"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	syncmap_repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/driver/grpcserver"
	syncmap_driver "github.com/becosuke/guestbook/api/internal/driver/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/pbgo"
)

type Injection interface {
	InjectConfig() *config.Config
	InjectLogger() *zap.Logger
	InjectAuthFunc() grpc_auth.AuthFunc
	InjectGrpcServer() *grpc.Server
	InjectController() pbgo.GuestbookServiceServer
	InjectControllerBoundary() controller.Boundary
	InjectUsecase() usecase.Usecase
	InjectSyncmapRepositoryGenerator() syncmap_repository.Generator
	InjectSyncmapRepositoryQuerier() syncmap_repository.Querier
	InjectSyncmapRepositoryCommander() syncmap_repository.Commander
	InjectSyncmapRepositoryBoundary() syncmap_repository.Boundary
	InjectSyncmapDriver() syncmap_driver.Syncmap
}

func NewInjection(serviceName, version string) Injection {
	return &injectionImpl{serviceName: serviceName, version: version}
}

type injectionImpl struct {
	container struct {
		Config                     *config.Config
		Logger                     *zap.Logger
		GrpcServer                 *grpc.Server
		GuestbookServiceServer     pbgo.GuestbookServiceServer
		ControllerBoundary         controller.Boundary
		Usecase                    usecase.Usecase
		SyncmapRepositoryGenerator syncmap_repository.Generator
		SyncmapRepositoryQuerier   syncmap_repository.Querier
		SyncmapRepositoryCommander syncmap_repository.Commander
		SyncmapRepositoryBoundary  syncmap_repository.Boundary
		SyncmapDriver              syncmap_driver.Syncmap
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
	actual, _ := i.store.LoadOrStore("grpc_server", &sync.Once{})
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
				i.InjectLogger(),
				i.InjectUsecase(),
				i.InjectControllerBoundary(),
			)
		})
	}
	return i.container.GuestbookServiceServer
}

func (i *injectionImpl) InjectControllerBoundary() controller.Boundary {
	actual, _ := i.store.LoadOrStore("controller_boundary", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.ControllerBoundary = controller.NewBoundary()
		})
	}
	return i.container.ControllerBoundary
}

func (i *injectionImpl) InjectUsecase() usecase.Usecase {
	actual, _ := i.store.LoadOrStore("usecase", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.Usecase = usecase.NewUsecase(i.InjectConfig(), i.InjectSyncmapRepositoryQuerier(), i.InjectSyncmapRepositoryCommander())
		})
	}
	return i.container.Usecase
}

func (i *injectionImpl) InjectSyncmapRepositoryGenerator() syncmap_repository.Generator {
	actual, _ := i.store.LoadOrStore("syncmap_repository_generator", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.SyncmapRepositoryGenerator = syncmap_repository.NewGenerator()
		})
	}
	return i.container.SyncmapRepositoryGenerator
}

func (i *injectionImpl) InjectSyncmapRepositoryQuerier() syncmap_repository.Querier {
	actual, _ := i.store.LoadOrStore("syncmap_repository_querier", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.SyncmapRepositoryQuerier = syncmap_repository.NewQuerier(
				i.InjectConfig(), i.InjectSyncmapDriver(), i.InjectSyncmapRepositoryBoundary())
		})
	}
	return i.container.SyncmapRepositoryQuerier
}

func (i *injectionImpl) InjectSyncmapRepositoryCommander() syncmap_repository.Commander {
	actual, _ := i.store.LoadOrStore("syncmap_repository_commander", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.SyncmapRepositoryCommander = syncmap_repository.NewCommander(
				i.InjectConfig(), i.InjectSyncmapDriver(), i.InjectSyncmapRepositoryBoundary(), i.InjectSyncmapRepositoryGenerator())
		})
	}
	return i.container.SyncmapRepositoryCommander
}

func (i *injectionImpl) InjectSyncmapRepositoryBoundary() syncmap_repository.Boundary {
	actual, _ := i.store.LoadOrStore("syncmap_repository_boundary", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.SyncmapRepositoryBoundary = syncmap_repository.NewBoundary()
		})
	}
	return i.container.SyncmapRepositoryBoundary
}

func (i *injectionImpl) InjectSyncmapDriver() syncmap_driver.Syncmap {
	actual, _ := i.store.LoadOrStore("syncmap_driver", &sync.Once{})
	once, ok := actual.(*sync.Once)
	if ok {
		once.Do(func() {
			i.container.SyncmapDriver = syncmap_driver.NewSyncmap()
		})
	}
	return i.container.SyncmapDriver
}
