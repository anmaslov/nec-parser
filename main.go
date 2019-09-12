package main

import (
	"fmt"
	"github.com/anmaslov/nec-parser/config"
	"github.com/anmaslov/nec-parser/store"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var cfg config.Config

func init() {
	parser := flags.NewParser(&cfg, flags.Default)
	parser.SubcommandsOptional = true
	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error init: %s.\nFor help use -h\n", err)
		os.Exit(1)
	}
}

// initLogger создает и настривает новый экземпляр логера
func initLogger() (*zap.Logger, error) {
	lvl := zap.InfoLevel
	err := lvl.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal log-level: %s", err)
	}

	opts := zap.NewProductionConfig()
	opts.Level = zap.NewAtomicLevelAt(lvl)
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if !cfg.LogJSON {
		opts.Encoding = "console"
		opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return opts.Build()
}

func main() {
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("error on init logger: %s.\nFor help use -h\n", err.Error())
		os.Exit(1)
	}

	mongoConfig := &config.Db{
		Host:     cfg.DbAddress,
		Dbname:   cfg.DbName,
		Username: cfg.DbUser,
		Password: cfg.DbPassword,
	}

	mongo, err := store.NewMongo(mongoConfig)
	if err != nil {
		logger.Fatal("error creating connection to MongoDb", zap.Error(err))
	}
	defer mongo.Session.Close()

	//Создаем канал
	p := DataProducer{
		OutChan: make(chan store.CallInfo),
	}

	phones, err := mongo.GetPhones()
	if err != nil {
		logger.Fatal("unable to get phones", zap.Error(err))
	}

	for _, phone := range phones {
		go stListener(phone, p, logger)
	}

	for data := range p.getOutChan() { //Ждем данных от канала
		err := mongo.InsertCall(&data)
		if err != nil {
			logger.Fatal("unable to write db", zap.Error(err)) //Падаем, т.к. запись в базу - критична
		}

		logger.Debug("write to DB success, date end call", zap.String("date_end", data.Cvt.DateEnd.String()))
	}

}
