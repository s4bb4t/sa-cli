package scaffold

import "fmt"

func (p *Project) mainTemplate() string {
	if p.Mode.HasOpenAPI() && p.Mode.HasGRPC() {
		return p.mainTemplateBoth()
	}
	if p.Mode.HasOpenAPI() {
		return p.mainTemplateOpenAPI()
	}
	return p.mainTemplateGRPC()
}

func (p *Project) mainTemplateGRPC() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"os/signal"
	"syscall"

	"%s/internal/config"
	"github.com/go-faster/errors"
	"github.com/s4bb4t/srvmon"
	"github.com/s4bb4t/zapang"
	"go.uber.org/zap"
)

const serviceName = "%s"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		print("load config:", err)
		return
	}

	l := zapang.New(ctx, serviceName, cfg.Logger, nil)

	l.Info("starting service", zapang.Environment(cfg.Logger.Environment))

	if err := start(ctx, *cfg, l); err != nil {
		l.Fatal("start", zap.Error(err))
	}
}

func start(ctx context.Context, cfg config.Config, log *zap.Logger) error {
	mon := srvmon.New(cfg.SrvMon, log)
	go mon.Run(ctx)

	err := initRepositories(ctx, cfg, log)
	if err != nil {
		return errors.Wrap(err, "init repositories")
	}

	exampleConn, err := initServices(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "init services")
	}

	mon.AddDependencies(
		srvmon.NewConnChecker(exampleConn, "exampleConn", true),
	)

	run(ctx, cfg, log)
	mon.SetReady()

	<-ctx.Done()
	log.Info("shutting down gracefully")

	return nil
}

func run(ctx context.Context, cfg config.Config, log *zap.Logger, reg ...Registrar) {
	go func() {
		if err := serveGRPC(ctx, cfg.Server.GRPCAddress, log, reg...); err != nil {
			log.Error(errors.Wrap(err, "serve grpc").Error())
		}
	}()

	go func() {
		if err := startMetricsServer(cfg.Server.MetricsAddr, log); err != nil {
			log.Error(errors.Wrap(err, "start metrics server").Error())
		}
	}()
}
`, p.Module, p.Name)
}

func (p *Project) mainTemplateOpenAPI() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"embed"
	"os/signal"
	"syscall"

	"%s/internal/config"
	//v1 "%s/internal/presentation/rest/v1"
	//api "%s/pkg/openapi/v1"
	"github.com/go-faster/errors"
	"github.com/s4bb4t/srvmon"
	"github.com/s4bb4t/zapang"
	"go.uber.org/zap"
)

const serviceName = "%s"

//go:embed docs
var docsFS embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		print("load config:", err)
		return
	}

	l := zapang.New(ctx, serviceName, cfg.Logger, nil)

	l.Info("starting service", zapang.Environment(cfg.Logger.Environment))

	if err := start(ctx, *cfg, l); err != nil {
		l.Fatal("start", zap.Error(err))
	}
}

func start(ctx context.Context, cfg config.Config, log *zap.Logger) error {
	mon := srvmon.New(cfg.SrvMon, log)
	go mon.Run(ctx)

	err := initRepositories(ctx, cfg, log)
	if err != nil {
		return errors.Wrap(err, "init repositories")
	}

	// todo: uncomment after running make generate-ogen
	//handler := v1.New(log)
	//
	//srv, err := api.NewServer(
	//	handler,
	//	nil,
	//	api.WithPathPrefix("/api/v1"),
	//)
	//if err != nil {
	//	return errors.Wrap(err, "create ogen server")
	//}
	//
	//run(cfg, log, WithVersion(srv, "/api/v1", "v1"))

	run(cfg, log)
	mon.SetReady()

	<-ctx.Done()
	log.Info("shutting down gracefully")

	return nil
}

func run(cfg config.Config, log *zap.Logger, servers ...VersionedServer) {
	go func() {
		if err := serveHTTP(cfg, log, servers...); err != nil {
			log.Error(errors.Wrap(err, "serve http").Error())
		}
	}()

	go func() {
		if err := startMetricsServer(cfg.Server.MetricsAddr, log); err != nil {
			log.Error(errors.Wrap(err, "start metrics server").Error())
		}
	}()
}
`, p.Module, p.Module, p.Module, p.Name)
}

func (p *Project) mainTemplateBoth() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"embed"
	"os/signal"
	"syscall"

	"%s/internal/config"
	//v1 "%s/internal/presentation/rest/v1"
	//api "%s/pkg/openapi/v1"
	"github.com/go-faster/errors"
	"github.com/s4bb4t/srvmon"
	"github.com/s4bb4t/zapang"
	"go.uber.org/zap"
)

const serviceName = "%s"

//go:embed docs
var docsFS embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		print("load config:", err)
		return
	}

	l := zapang.New(ctx, serviceName, cfg.Logger, nil)

	l.Info("starting service", zapang.Environment(cfg.Logger.Environment))

	if err := start(ctx, *cfg, l); err != nil {
		l.Fatal("start", zap.Error(err))
	}
}

func start(ctx context.Context, cfg config.Config, log *zap.Logger) error {
	mon := srvmon.New(cfg.SrvMon, log)
	go mon.Run(ctx)

	err := initRepositories(ctx, cfg, log)
	if err != nil {
		return errors.Wrap(err, "init repositories")
	}

	exampleConn, err := initServices(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "init services")
	}

	mon.AddDependencies(
		srvmon.NewConnChecker(exampleConn, "exampleConn", true),
	)

	// todo: uncomment after running make generate-ogen
	//handler := v1.New(log)
	//
	//srv, err := api.NewServer(
	//	handler,
	//	nil,
	//	api.WithPathPrefix("/api/v1"),
	//)
	//if err != nil {
	//	return errors.Wrap(err, "create ogen server")
	//}
	//
	//run(ctx, cfg, log, WithVersion(srv, "/api/v1", "v1"))

	run(ctx, cfg, log)
	mon.SetReady()

	<-ctx.Done()
	log.Info("shutting down gracefully")

	return nil
}

func run(ctx context.Context, cfg config.Config, log *zap.Logger, servers []VersionedServer, registrar []Registrar,) {
	go func() {
		if err := serveHTTP(cfg, log, servers...); err != nil {
			log.Error(errors.Wrap(err, "serve http").Error())
		}
	}()

	go func() {
		if err := serveGRPC(ctx, cfg.Server.GRPCAddress, log, registrar...); err != nil {
			log.Error(errors.Wrap(err, "serve grpc").Error())
		}
	}()

	go func() {
		if err := startMetricsServer(cfg.Server.MetricsAddr, log); err != nil {
			log.Error(errors.Wrap(err, "start metrics server").Error())
		}
	}()
}
`, p.Module, p.Module, p.Module, p.Name)
}

func (p *Project) serverTemplate() string {
	if p.Mode.HasOpenAPI() && p.Mode.HasGRPC() {
		return p.serverTemplateBoth()
	}
	if p.Mode.HasOpenAPI() {
		return p.serverTemplateOpenAPI()
	}
	return p.serverTemplateGRPC()
}

func (p *Project) serverTemplateGRPC() string {
	return `package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConcurrentStreams = 1000
	maxMsgSize           = 4 << 20 // 4 MB

	keepaliveIdle    = 15 * time.Minute
	keepaliveMaxAge  = 30 * time.Minute
	keepaliveGrace   = 5 * time.Second
	keepaliveTime    = 5 * time.Second
	keepaliveTimeout = 1 * time.Second
	keepaliveMinTime = 5 * time.Second

	metricsReadTimeout  = 5 * time.Second
	metricsWriteTimeout = 5 * time.Second
	metricsIdleTimeout  = 60 * time.Second
)

type Registrar interface {
	Register(*grpc.Server)
}

func serveGRPC(ctx context.Context, address string, l *zap.Logger, servers ...Registrar) error {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     keepaliveIdle,
			MaxConnectionAge:      keepaliveMaxAge,
			MaxConnectionAgeGrace: keepaliveGrace,
			Time:                  keepaliveTime,
			Timeout:               keepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepaliveMinTime,
			PermitWithoutStream: true,
		}),
		grpc.MaxConcurrentStreams(maxConcurrentStreams),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	for _, reg := range servers {
		reg.Register(srv)
	}

	l.Info("starting gRPC server", zap.String("address", address))

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	if err := srv.Serve(lis); err != nil {
		return errors.Wrap(err, "grpc serve")
	}

	return nil
}

func startMetricsServer(address string, log *zap.Logger) error {
	reg := prometheus.NewRegistry()
	if err := reg.Register(collectors.NewGoCollector()); err != nil {
		return errors.Wrap(err, "register go collector")
	}

	if err := reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return errors.Wrap(err, "register process collector")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  metricsReadTimeout,
		WriteTimeout: metricsWriteTimeout,
		IdleTimeout:  metricsIdleTimeout,
	}

	log.Info("starting metrics server", zap.String("address", address))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "metrics server")
	}

	return nil
}
`
}

func (p *Project) serverTemplateOpenAPI() string {
	return fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"%s/internal/config"
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	httpReadTimeout  = 15 * time.Second
	httpWriteTimeout = 15 * time.Second
	httpIdleTimeout  = 60 * time.Second

	metricsReadTimeout  = 5 * time.Second
	metricsWriteTimeout = 5 * time.Second
	metricsIdleTimeout  = 60 * time.Second
)

// VersionedServer — any ogen-generated server that can serve HTTP
// and report its API version for routing and docs.
type VersionedServer interface {
	http.Handler
	Prefix() string  // e.g. "/api/v1"
	Version() string // e.g. "v1"
}

func serveHTTP(cfg config.Config, log *zap.Logger, servers ...VersionedServer) error {
	r := mux.NewRouter()

	mountDocs(r, servers)
	mountAPI(r, servers)

	srv := &http.Server{
		Addr:         cfg.Server.HTTPAddress,
		Handler:      r,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
	}

	log.Info("starting HTTP server", zap.String("address", cfg.Server.HTTPAddress))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http serve")
	}

	return nil
}

func startMetricsServer(address string, log *zap.Logger) error {
	reg := prometheus.NewRegistry()
	if err := reg.Register(collectors.NewGoCollector()); err != nil {
		return errors.Wrap(err, "register go collector")
	}

	if err := reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return errors.Wrap(err, "register process collector")
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:         address,
		Handler:      httpMux,
		ReadTimeout:  metricsReadTimeout,
		WriteTimeout: metricsWriteTimeout,
		IdleTimeout:  metricsIdleTimeout,
	}

	log.Info("starting metrics server", zap.String("address", address))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "metrics server")
	}

	return nil
}

func mountDocs(r *mux.Router, servers []VersionedServer) {
	swaggerHTML, err := fs.ReadFile(docsFS, "docs/swagger.html")
	if err != nil {
		panic("embedded swagger.html not found")
	}

	for _, s := range servers {
		spec := loadSpec(s.Version())
		path := fmt.Sprintf("/docs/%%s/openapi.json", s.Version())
		r.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(spec)
		})
	}

	versions := make([]string, 0, len(servers))
	for _, s := range servers {
		versions = append(versions, s.Version())
	}
	versionsJSON, _ := json.Marshal(versions)
	r.HandleFunc("/docs/versions.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(versionsJSON)
	})

	r.HandleFunc("/docs/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(swaggerHTML)
	})

	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
}

func loadSpec(version string) []byte {
	data, err := fs.ReadFile(docsFS, fmt.Sprintf("docs/%%s/openapi.json", version))
	if err != nil {
		panic(fmt.Sprintf("embedded spec for %%s not found — run `+"`"+`make bundle-spec`+"`"+`", version))
	}
	return data
}

func mountAPI(r *mux.Router, servers []VersionedServer) {
	for _, s := range servers {
		r.PathPrefix(s.Prefix()).Handler(s)
	}
}

type versionedServer struct {
	http.Handler
	prefix  string
	version string
}

func (s *versionedServer) Prefix() string  { return s.prefix }
func (s *versionedServer) Version() string { return s.version }

// WithVersion wraps an http.Handler (e.g. *api.Server) into a VersionedServer.
func WithVersion(h http.Handler, prefix, version string) VersionedServer {
	return &versionedServer{Handler: h, prefix: prefix, version: version}
}
`, p.Module)
}

func (p *Project) serverTemplateBoth() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"%s/internal/config"
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConcurrentStreams = 1000
	maxMsgSize           = 4 << 20 // 4 MB

	keepaliveIdle    = 15 * time.Minute
	keepaliveMaxAge  = 30 * time.Minute
	keepaliveGrace   = 5 * time.Second
	keepaliveTime    = 5 * time.Second
	keepaliveTimeout = 1 * time.Second
	keepaliveMinTime = 5 * time.Second

	httpReadTimeout  = 15 * time.Second
	httpWriteTimeout = 15 * time.Second
	httpIdleTimeout  = 60 * time.Second

	metricsReadTimeout  = 5 * time.Second
	metricsWriteTimeout = 5 * time.Second
	metricsIdleTimeout  = 60 * time.Second
)

type (
	Registrar interface {
		Register(*grpc.Server)
	}

	// VersionedServer — any ogen-generated server that can serve HTTP
	// and report its API version for routing and docs.
	VersionedServer interface {
		http.Handler
		Prefix() string  // e.g. "/api/v1"
		Version() string // e.g. "v1"
	}
)

func serveGRPC(ctx context.Context, address string, l *zap.Logger, servers ...Registrar) error {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     keepaliveIdle,
			MaxConnectionAge:      keepaliveMaxAge,
			MaxConnectionAgeGrace: keepaliveGrace,
			Time:                  keepaliveTime,
			Timeout:               keepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepaliveMinTime,
			PermitWithoutStream: true,
		}),
		grpc.MaxConcurrentStreams(maxConcurrentStreams),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	for _, reg := range servers {
		reg.Register(srv)
	}

	l.Info("starting gRPC server", zap.String("address", address))

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "listen")
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	if err := srv.Serve(lis); err != nil {
		return errors.Wrap(err, "grpc serve")
	}

	return nil
}

func serveHTTP(cfg config.Config, log *zap.Logger, servers ...VersionedServer) error {
	r := mux.NewRouter()

	mountDocs(r, servers)
	mountAPI(r, servers)

	srv := &http.Server{
		Addr:         cfg.Server.HTTPAddress,
		Handler:      r,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
	}

	log.Info("starting HTTP server", zap.String("address", cfg.Server.HTTPAddress))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http serve")
	}

	return nil
}

func startMetricsServer(address string, log *zap.Logger) error {
	reg := prometheus.NewRegistry()
	if err := reg.Register(collectors.NewGoCollector()); err != nil {
		return errors.Wrap(err, "register go collector")
	}

	if err := reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return errors.Wrap(err, "register process collector")
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:         address,
		Handler:      httpMux,
		ReadTimeout:  metricsReadTimeout,
		WriteTimeout: metricsWriteTimeout,
		IdleTimeout:  metricsIdleTimeout,
	}

	log.Info("starting metrics server", zap.String("address", address))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "metrics server")
	}

	return nil
}

func mountDocs(r *mux.Router, servers []VersionedServer) {
	swaggerHTML, err := fs.ReadFile(docsFS, "docs/swagger.html")
	if err != nil {
		panic("embedded swagger.html not found")
	}

	for _, s := range servers {
		spec := loadSpec(s.Version())
		path := fmt.Sprintf("/docs/%%s/openapi.json", s.Version())
		r.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(spec)
		})
	}

	versions := make([]string, 0, len(servers))
	for _, s := range servers {
		versions = append(versions, s.Version())
	}
	versionsJSON, _ := json.Marshal(versions)
	r.HandleFunc("/docs/versions.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(versionsJSON)
	})

	r.HandleFunc("/docs/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(swaggerHTML)
	})

	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
}

func loadSpec(version string) []byte {
	data, err := fs.ReadFile(docsFS, fmt.Sprintf("docs/%%s/openapi.json", version))
	if err != nil {
		panic(fmt.Sprintf("embedded spec for %%s not found — run `+"`"+`make bundle-spec`+"`"+`", version))
	}
	return data
}

func mountAPI(r *mux.Router, servers []VersionedServer) {
	for _, s := range servers {
		r.PathPrefix(s.Prefix()).Handler(s)
	}
}

type versionedServer struct {
	http.Handler
	prefix  string
	version string
}

func (s *versionedServer) Prefix() string  { return s.prefix }
func (s *versionedServer) Version() string { return s.version }

// WithVersion wraps an http.Handler (e.g. *api.Server) into a VersionedServer.
func WithVersion(h http.Handler, prefix, version string) VersionedServer {
	return &versionedServer{Handler: h, prefix: prefix, version: version}
}
`, p.Module)
}

func (p *Project) servicesTemplate() string {
	if p.Mode.HasGRPC() {
		return p.servicesTemplateGRPC()
	}
	return p.servicesTemplateOpenAPI()
}

func (p *Project) servicesTemplateGRPC() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"time"

	"%s/internal/config"
	"github.com/go-faster/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type connEntry struct {
	name    string
	address string
}

func initRepositories(ctx context.Context, cfg config.Config, log *zap.Logger) error {
	// todo: create repos

	go func() {
		<-ctx.Done()

		// todo: close repos
	}()

	return nil
}

func initServices(ctx context.Context, cfg config.Config) (exampleConn *grpc.ClientConn, err error) {
	entries := []connEntry{
		//{name: "exampleConn", address: cfg.exampleConn.Address},
	}

	conns := make([]*grpc.ClientConn, 0, len(entries))

	for _, e := range entries {
		conn, dialErr := grpcConn(e.address)
		if dialErr != nil {
			return  nil, errors.Wrapf(dialErr, "connect to %%s", e.name)
		}

		conns = append(conns, conn)
	}

	for _, conn := range conns {
		go func() {
			<-ctx.Done()

			_ = conn.Close()
		}()
	}

	return conns[0], nil
}

func grpcConn(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.6,
					Jitter:     0.1,
					MaxDelay:   10 * time.Second,
				},
				MinConnectTimeout: 5 * time.Second,
			},
		),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc connection")
	}

	return conn, nil
}
`, p.Module)
}

func (p *Project) servicesTemplateOpenAPI() string {
	return fmt.Sprintf(`package main

import (
	"context"

	"%s/internal/config"
	"go.uber.org/zap"
)

func initRepositories(ctx context.Context, cfg config.Config, log *zap.Logger) error {
	// todo: create repos

	go func() {
		<-ctx.Done()

		// todo: close repos
	}()

	return nil
}
`, p.Module)
}
