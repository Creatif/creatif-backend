package server

import (
	"context"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartServer(srv *echo.Echo) {
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("Starting server on %s:%v...\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
		srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c

	fmt.Println(fmt.Sprintf("Received signal: %d", sig))

	appShutdown()

	ctx := context.Background()
	err := srv.Shutdown(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Server is terminated.")

	os.Exit(0)
}

func appShutdown() {
	sqlDB, err := storage.SQLDB()
	if err != nil {
		fmt.Println("Unable to disconnect from the database: ", err)
	}

	// Close
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Unable to disconnect from the database", err)
	}
}
