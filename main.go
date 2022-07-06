package main

import (
	"fmt"
	"log"

	"github.com/callmehorhe/test/pkg/repository"
	_ "github.com/lib/pq"
)

func main() {
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	client := repository.NewClientInstrumentDB(db)

	// CREATE
	err = client.Create(&repository.ClientInstrument{
		Client_ID:          1111,
		Instrument_Details: []byte("{\"details\":\"id\"}"),
		Instrument_ID:      "id",
		Method_ID:          "method",
		Name:               "name",
		Is_Default:         false,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CREATE SUCCESSFULLY")

	// READ
	clients, err := client.Read("client_id=1111")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range *clients {
		fmt.Printf("%d, %s, %s, %s, %s, %t\n",
			c.Client_ID,
			string(c.Instrument_Details),
			c.Instrument_ID,
			c.Method_ID,
			c.Name,
			c.Is_Default,
		)
	}
	fmt.Println("READ SUCCESSFULLY")

	//UPDATE
	err = client.Update(&repository.ClientInstrument{
		Client_ID:  1112,
		Method_ID:  "method",
		Name:       "name",
		Is_Default: false,
	}, "client_id=1112")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UPDATE SUCCESSFULLY")

	// DELETE
	err = client.Delete("client_id=1112 OR client_id=1111")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("DELETE SUCCESSFULLY")
}
