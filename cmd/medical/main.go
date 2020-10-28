package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/pkg/database"
	"test/pkg/medical"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP Listen Address")
	)

	flag.Parse()

	db := database.DbConn()
	defer db.Close()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s medical.Service
	{
		repo, err := database.NewRepo(db)
		if err != nil {
			logger.Log("exit", err)
		}
		s = medical.NewMedicalService(repo)
		s = medical.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = medical.MakeHTTpHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
