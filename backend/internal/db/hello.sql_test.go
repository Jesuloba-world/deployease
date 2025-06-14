package db

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

)

func TestGetGreeting(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := New(mock)

	mock.ExpectQuery("SELECT 'Hello, world!' AS greeting").
		WillReturnRows(pgxmock.NewRows([]string{"greeting"}).AddRow("Hello, world!"))

	greeting, err := q.GetGreeting(context.Background())
	require.NoError(t, err)
	require.Equal(t, "Hello, world!", greeting)

	// we make sure that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
