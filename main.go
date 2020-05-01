package main

import (
	"./server"
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"time"
)

func NewLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	return lg
}

func main() {

	stopchan := make(chan os.Signal)
	//router := chi.NewRouter()
	lg := NewLogger()

	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("blog")

	serv := server.New(lg, *flagRootDir, db)

	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
	log.Print("gracefull shutdown")
}
