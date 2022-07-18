package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := NewClientInstrumentDB(db)

	tests := []struct {
		name    string
		mock    func()
		input   *ClientInstrument
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("INSERT INTO client_instruments").
					WithArgs(1111, []byte("{\"test\":123}"), "id", "method", "name", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &ClientInstrument{
				Client_ID:          1111,
				Instrument_Details: []byte("{\"test\":123}"),
				Instrument_ID:      "id",
				Method_ID:          "method",
				Name:               "name",
				Is_Default:         false,
			},
			wantErr: false,
		},
		{
			name: "json error",
			mock: func() {
				mock.ExpectExec("INSERT INTO client_instruments").
					WithArgs(1111, []byte("test"), "id", "method", "name", false).
					WillReturnError(errors.New("json error"))
			},
			input: &ClientInstrument{
				Client_ID:          1111,
				Instrument_Details: []byte("test"),
				Instrument_ID:      "id",
				Method_ID:          "method",
				Name:               "",
				Is_Default:         false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := r.Create(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, err)
			}
		})
	}
}

func TestRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := NewClientInstrumentDB(db)

	tests := []struct {
		name    string
		mock    func()
		input   *InstrumentSearchCriteria
		want    *[]ClientInstrument
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"client_id", "instrument_details", "instrument_id", "method_id", "name", "is_default",
				}).AddRow(1111, []byte("{\"test\":123}"), "test", "test", "test", false)
				condition, _ := makeCondition(&InstrumentSearchCriteria{
					Client_ID: 1111,
					Name:      "test",
				})
				mock.ExpectQuery("SELECT (.+) FROM client_instruments " + condition).
					WillReturnRows(rows)
			},
			input: &InstrumentSearchCriteria{
				Client_ID: 1111,
				Name:      "test",
			},
			want: &[]ClientInstrument{
				{
					1111, []byte("{\"test\":123}"), "test", "test", "test", false,
				},
			},
			wantErr: false,
		},
		{
			name: "empty criteria",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"client_id", "instrument_details", "instrument_id", "method_id", "name", "is_default",
				}).RowError(0, errors.New("err"))
				condition, _ := makeCondition(&InstrumentSearchCriteria{})
				mock.ExpectQuery("SELECT (.+) FROM client_instruments " + condition).
					WillReturnRows(rows)
			},
			input:   &InstrumentSearchCriteria{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			get, err := r.Read(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, get)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := NewClientInstrumentDB(db)

	type args struct {
		*ClientInstrument
		*InstrumentSearchCriteria
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				condition, _ := makeCondition(&InstrumentSearchCriteria{
					Client_ID: 1111,
				})
				mock.ExpectExec("UPDATE client_instruments SET (.+) "+condition).
					WithArgs(1111, []byte("{\"test\":123}"), "test", "test", "test", false).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				InstrumentSearchCriteria: &InstrumentSearchCriteria{Client_ID: 1111},
				ClientInstrument: &ClientInstrument{
					1111, []byte("{\"test\":123}"), "test", "test", "test", false,
				},
			},
			wantErr: false,
		},
		{
			name: "JSON error",
			mock: func() {
				condition, _ := makeCondition(&InstrumentSearchCriteria{
					Client_ID: 1111,
				})
				mock.ExpectExec("UPDATE client_instruments SET (.+) "+condition).
					WithArgs(1111, []byte("test"), "test", "test", "test", false).
					WillReturnError(errors.New("json error"))
			},
			input: args{
				InstrumentSearchCriteria: &InstrumentSearchCriteria{Client_ID: 1111},
				ClientInstrument: &ClientInstrument{
					1111, []byte("test"), "test", "test", "test", false,
				},
			},
			wantErr: true,
		},
		{
			name: "empty instrument_details",
			mock: func() {
				condition, _ := makeCondition(&InstrumentSearchCriteria{
					Client_ID: 1111,
				})
				mock.ExpectExec("UPDATE client_instruments SET (.+) "+condition).
					WithArgs(1111, "test", "test", "test", false).
					WillReturnError(errors.New("json error"))
			},
			input: args{
				InstrumentSearchCriteria: &InstrumentSearchCriteria{Client_ID: 1111},
				ClientInstrument: &ClientInstrument{
					1111, []byte("test"), "test", "test", "test", false,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.ClientInstrument, tt.input.InstrumentSearchCriteria)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
