package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	clientID := 1

	client, err := selectClient(db, clientID)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, client.ID, clientID)
	assert.NotEmpty(t, client.FIO)
	assert.NotEmpty(t, client.Login)
	assert.NotEmpty(t, client.Birthday)
	assert.NotEmpty(t, client.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	clientID := -1

	client, err := selectClient(db, clientID)
	if err != nil {
		require.Equal(t, err, sql.ErrNoRows)
	}
	assert.Empty(t, client.ID)
	assert.Empty(t, client.FIO)
	assert.Empty(t, client.Login)
	assert.Empty(t, client.Birthday)
	assert.Empty(t, client.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	insertId, err := insertClient(db, cl)
	if err != nil {
		require.NoError(t, err)
	}
	require.NotEmpty(t, insertId)
	cl.ID = insertId

	client, err := selectClient(db, cl.ID)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, cl.ID, client.ID)
	assert.Equal(t, cl.FIO, client.FIO)
	assert.Equal(t, cl.Login, client.Login)
	assert.Equal(t, cl.Birthday, client.Birthday)
	assert.Equal(t, cl.Email, client.Email)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	insertId, err := insertClient(db, cl)
	if err != nil {
		require.NoError(t, err)
	}
	require.NotEmpty(t, insertId)

	_, err = selectClient(db, insertId)
	if err != nil {
		require.NoError(t, err)
	}

	err = deleteClient(db, insertId)
	if err != nil {
		require.NoError(t, err)
	}
	_, err = selectClient(db, insertId)
	require.Equal(t, sql.ErrNoRows, err)
}
