package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/javadmohebbi/goAuth/server"
	"github.com/javadmohebbi/goAuth/server/debugger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dbPath := flag.String("path", "/tmp/go-auth.db", "Path to sqlite database")
	host := flag.String("host", "0.0.0.0", "API server listen host")
	port := flag.Int("port", 7161, "API server listen port")
	debug := flag.Bool("debug", true, "Enable/Disable debugging mode")

	flag.Parse()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful: false, // Disable color
		},
	)

	// open or create sqlite database
	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{
		Logger: newLogger,
	})

	// check if no successful
	if err != nil {
		fmt.Printf("Can not open %v due to error: %v", *dbPath, err)
	}

	dbg := debugger.New(*debug, logrus.New(), "log")
	s := server.New(*host, *port, dbg, db)

	s.ServeHTTP()

}
