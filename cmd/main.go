package main

import (
	"bank/api"
	"bank/config"
	"bank/service"
	"bank/store"
	"log"
	"os"
)



func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	config, err := config.NewConfig()
	if err != nil {
		errLog.Fatal(err)
	}

	db, err := store.InitializeDB(config)
	if err != nil {
		// trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		// errLog.Output(2, trace)
		errLog.Fatal(err)
	}
	defer db.Close()

	store := store.NewStore(db)
	service := service.NewService(store)
	handler := api.NewHandler(infoLog, errLog, service)

	errLog.Fatal(handler.RunServer(config.Server))
	
}