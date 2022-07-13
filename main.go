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
	defer db.Close()

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
	clients, err := client.Read(&repository.InstrumentSearchCriteria{
		Client_ID: 1111,
		Method_ID: "method",
	})
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
	err = client.Update(
		&repository.ClientInstrument{
			Client_ID:          1112,
			Instrument_Details: []byte("{\"test\": 123}"),
			Method_ID:          "method",
			Name:               "name",
			Is_Default:         false,
		}, &repository.InstrumentSearchCriteria{
			Client_ID: 1111,
		})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UPDATE SUCCESSFULLY")
}
