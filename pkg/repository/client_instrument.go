package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	Read(criteria *InstrumentSearchCriteria) (*[]ClientInstrument, error)
	Update(client *ClientInstrument, criteria *InstrumentSearchCriteria) error
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
	var js map[string]interface{}
	if err := json.Unmarshal(client.Instrument_Details, &js); err != nil {
		return err
	}
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

type InstrumentSearchCriteria struct {
	Client_ID     int64
	Instrument_ID string
	Method_ID     string
	Name          string
}

// makeCondition ...
func makeCondition(criteria *InstrumentSearchCriteria) (string, error) {
	conditions := []string{}
	if criteria.Client_ID != 0 {
		conditions = append(conditions, fmt.Sprintf("client_id = %d", criteria.Client_ID))
	}
	if criteria.Instrument_ID != "" {
		conditions = append(conditions, fmt.Sprintf("instrument_id = '%s'", criteria.Instrument_ID))
	}
	if criteria.Method_ID != "" {
		conditions = append(conditions, fmt.Sprintf("method_id = '%s'", criteria.Method_ID))
	}
	if criteria.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name = '%s'", criteria.Name))
	}

	if len(conditions) == 0 {
		return "", errors.New("not enough criteria")
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND ")), nil
}

// Find client_instrument by criteria.
func (r *ClientInstrumentDB) Read(criteria *InstrumentSearchCriteria) (*[]ClientInstrument, error) {
	condition, err := makeCondition(criteria)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(`
	SELECT client_id, instrument_details, instrument_id, method_id, name, is_default
	FROM client_instruments
	%s
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

// Update client instruent by critetia.
func (r *ClientInstrumentDB) Update(client *ClientInstrument, criteria *InstrumentSearchCriteria) error {
	condition, err := makeCondition(criteria)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf(`
	UPDATE client_instruments
	SET client_id=$1, instrument_details=$2, instrument_id=$3, method_id=$4, name=$5, is_default=$6
	%s
	`, condition)

	var js map[string]interface{}
	if err := json.Unmarshal(client.Instrument_Details, &js); err != nil {
		return err
	}

	_, err = r.db.Exec(sql,
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
