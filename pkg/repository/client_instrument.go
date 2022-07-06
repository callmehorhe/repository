package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// ClientInstrument ...
type ClientInstrument struct {
	Client_ID          int64
	Instrument_Details json.RawMessage
	Instrument_ID      string
	Method_ID          string
	Name               string
	Is_Default         bool
}

// ClientInstrumentRepository ...
type ClientInstrumentRepository interface {
	Create(client *ClientInstrument) error
	Read(condition string) (*[]ClientInstrument, error)
	Update(client *ClientInstrument, instrumentId string) error
	Delete(instrumentId string) error
}

// ClientInstrumentDB ...
type ClientInstrumentDB struct {
	db *sql.DB
}

// NewClientInstrumentDB creates new ClientInstrumentDB.
func NewClientInstrumentDB(db *sql.DB) *ClientInstrumentDB {
	return &ClientInstrumentDB{
		db: db,
	}
}

// Create client_instrument.
func (r *ClientInstrumentDB) Create(client *ClientInstrument) error {
	sql := `
	INSERT INTO client_instruments
	(client_id, instrument_details, instrument_id, method_id, name, is_default)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(sql,
		client.Client_ID,
		client.Instrument_Details,
		client.Instrument_ID,
		client.Method_ID,
		client.Name,
		client.Is_Default,
	)
	if err != nil {
		return err
	}
	return nil
}

// Find client_instrument by your own condition.
// For example:
//	r.Read("client_id=1003")
func (r *ClientInstrumentDB) Read(condition string) (*[]ClientInstrument, error) {
	sql := fmt.Sprintf(`
	SELECT client_id, instrument_details, instrument_id, method_id, name, is_default
	FROM client_instruments
	WHERE %s
	`, condition)
	var clients []ClientInstrument
	var client ClientInstrument
	rows, err := r.db.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&client.Client_ID,
			&client.Instrument_Details,
			&client.Instrument_ID,
			&client.Method_ID,
			&client.Name,
			&client.Is_Default,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return &clients, nil
}

// Delete client instruent by condition. For example:
//	r.Update(client, "is_default=false")
func (r *ClientInstrumentDB) Update(client *ClientInstrument, condition string) error {
	sql := fmt.Sprintf(`
	UPDATE client_instruments
	SET client_id=$1, instrument_details=$2, instrument_id=$3, method_id=$4, name=$5, is_default=$6
	WHERE %s
	`, condition)

	instrumentDetails, err := json.Marshal(&client.Instrument_Details)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql,
		client.Client_ID,
		instrumentDetails,
		client.Instrument_ID,
		client.Method_ID,
		client.Name,
		client.Is_Default,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete client_instrument by condition.
func (r *ClientInstrumentDB) Delete(condition string) error {
	sql := fmt.Sprintf(`
	DELETE
	FROM client_instruments
	WHERE %s
	`, condition)
	_, err := r.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
