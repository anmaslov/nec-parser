package main

import (
	"fmt"
	"github.com/anmaslov/nec-parser/config"
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

	session := initialiseMongo()
	mongoStore.session = session
	defer session.Close()

	//Создаем канал
	p := DataProducer{
		OutChan: make(chan CallInfo),
	}

	phones, err := getPhones()
	if err != nil {
		logger.Fatal("unable to get phones", zap.Error(err))
	}

	for _, phone := range phones {
		go stListener(phone, p, logger)
	}

	for data := range p.getOutChan() { //Ждем данных от канала
		err := insertCall(&data)
		if err != nil {
			logger.Fatal("unable to write db", zap.Error(err)) //Падаем, т.к. запись в базу - критична
		}

		logger.Debug("write to DB success, date end call", zap.String("date_end", data.Cvt.DateEnd.String()))
	}

}
