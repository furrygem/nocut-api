package apiserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/furrygem/nocut-api/internal/links"
	"github.com/furrygem/nocut-api/internal/links/db"
	"github.com/furrygem/nocut-api/pkg/client/mongodb"
	"github.com/furrygem/nocut-api/pkg/logging"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Apiserver struct {
	Logger        logging.Logger
	Router        *mux.Router
	MongoDatabase *mongo.Database
}

// New creates and configures new apiserver
func New(config *Config) *Apiserver {

	mongoDBClient, err := mongodb.NewClient(
		context.Background(),
		config.MongoDB.Host,
		config.MongoDB.Port,
		config.MongoDB.Username,
		config.MongoDB.Password,
		config.MongoDB.Database,
		config.MongoDB.AuthDB,
	)

	if err != nil {
		log.Fatal(err)
	}

	s := &Apiserver{
		Router:        mux.NewRouter(),
		Logger:        logging.GetLogger(),
		MongoDatabase: mongoDBClient,
	}

	s.configureRouterMongo(config)
	s.configureLogger(config)

	return s

}

func (as *Apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	as.Logger.Infof("%s -> %s %s", r.RemoteAddr, r.Method, r.URL)
	as.Router.ServeHTTP(w, r)
}

// Start creates net listener and starts the apiserver
func (as *Apiserver) Start(config *Config) error {
	httpserver := http.Server{
		Handler:      as,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	listenAddr := fmt.Sprintf("%s:%d", config.BindAddr, config.BindPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	return httpserver.Serve(listener)
}

func (as *Apiserver) configureLogger(c *Config) error {
	err := as.Logger.SetLevel(c.LogLevel)
	if err != nil {
		return err
	}
	return nil
}

func (as *Apiserver) configureRouterMongo(c *Config) error {
	storage := db.NewStorage(as.MongoDatabase, c.MongoDB.Collection, &as.Logger)
	han := links.NewHandler(as.Logger, storage, c.MongoDB.LinkTTL)
	han.Register(as.Router, c.APIPrefix)
	return nil
}
