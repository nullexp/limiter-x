package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nullexp/limiter-x/internal/adapter/driven/cache"
	drivenDb "github.com/nullexp/limiter-x/internal/adapter/driven/db"
	repository "github.com/nullexp/limiter-x/internal/adapter/driven/db/repository"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	grpcDriver "github.com/nullexp/limiter-x/internal/adapter/driver/grpc"
	ratev1 "github.com/nullexp/limiter-x/internal/adapter/driver/grpc/proto/rate/v1"
	driver "github.com/nullexp/limiter-x/internal/adapter/driver/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting the server")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	m, err := migrate.New("file://internal/adapter/driven/db/migration",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal("Error loading migration files ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error calling up function ", err)
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	log.Println("DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	port := os.Getenv("APP_PORT")
	ip := os.Getenv("APP_IP")

	addr := fmt.Sprintf("%s:%v", ip, port)
	// Create a TCP listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	_ = lis

	log.Println("successfully initialized")

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the Greeter service

	txFactory := drivenDb.NewPostgresDbTransactionFactory(db)
	rateRepoFactory := repository.NewUserRateLimitRepositoryFactory()
	cache := cache.NewRedisClient("", "", "redis", os.Getenv("REDIS_URL"))
	cache.Connect()

	windowMilSecond, err := strconv.Atoi(os.Getenv("WINDOW_MILI_SEC"))
	if err != nil {
		log.Fatal(err)
	}
	rateService := driver.NewRateLimitService(rateRepoFactory, cache, txFactory, time.Duration(windowMilSecond)*time.Millisecond)
	grpcService := grpcDriver.NewRateLimiterService(rateService)
	ratev1.RegisterRateLimiterServiceServer(s, grpcService)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	// Log and start the server
	log.Printf("gRPC server listening on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
