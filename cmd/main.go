package main

import (
	"flag"
	"fmt"
	"github.com/iamkirillnb/finStar/internal"
	"github.com/iamkirillnb/finStar/internal/handlers"
	"github.com/iamkirillnb/finStar/internal/repos"
	"log"
)
var cfgPath string

func init()  {
	flag.StringVar(&cfgPath, "config", "dev.yaml", "config path")

}


func main()  {
	flag.Parse()
	fmt.Println("finStaar Start's")

	config := internal.GetConfig(cfgPath)
	log.Println("config gotten successfully")

	// Postgres
	db := repos.NewPostgres(&config.Db)

	// Db Repo
	repo := repos.NewDbRepo(db)

	//usr := entities.Wallet{
	//	Id: "a143a66b-5e91-45ca-85b2-4ad738806787",
	//	Amount: 6,
	//}
	//err := repo.CreateUser(&usr)
	//if err != nil {
	//	log.Println("create user failed")
	//	log.Fatal(err)
	//}

	//err := repo.AddMoney(&usr)
	//if err != nil {
	//	log.Println("add money failed")
	//	log.Fatal(err)
	//}


	//from :=&entities.Wallet{
	//	Id:     "ab0df19a-eed4-4bcc-8b21-f51ee68af479",
	//	Amount: 1000,
	//}
	//to := &entities.Wallet{
	//	Id:     "a143a66b-5e91-45ca-85b2-4ad738806787",
	//}

	//err := repo.SendMoneyToUser(from, to)
	//if err != nil {
	//	log.Println("send money failed")
	//	log.Fatal(err)
	//
	//}

	// handler
	handler := handlers.NewHandler(&config.Server, repo)
	handler.Start()
}

