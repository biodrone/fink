package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS services
(
	id SERIAL,
	name TEXT NOT NULL,
	endpoint TEXT NOT NULL,
	command TEXT NOT NULL,
	CONSTRAINT services_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	a = App{}
	a.Initialise(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM services")
	a.DB.Exec("ALTER SEQUENCE services_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/services", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentService(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/services/11", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Service not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Service not found'. Got '%s'", m["error"])
	}
}

func TestCreateService(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test_service","endpoint":"some_webhook","command":"ls -lah"}`)

	req, _ := http.NewRequest("POST", "/service", bytes.NewBuffer(payload))
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test_service" {
		t.Errorf("Expected service name to be 'test_service'. Got '%v'", m["name"])
	}

	if m["endpoint"] != "some_webhook" {
		t.Errorf("Expected service endpoint to be 'some_webhook'. Got '%v'", m["endpoint"])
	}

	if m["command"] != "ls -lah" {
		t.Errorf("Expected service command to be 'ls -lah'. Got '%v'", m["command"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected service ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetService(t *testing.T) {
	clearTable()
	addServices(1)

	req, _ := http.NewRequest("GET", "/service/1", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addServices(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO services(name, price) VALUES($1, $2)", "Service "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateService(t *testing.T) {
	clearTable()
	addServices(1)

	req, _ := http.NewRequest("GET", "/service/1", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)
	var originalService map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalService)

	payload := []byte(`{"name":"test_service_UPDATE","endpoint":"some_webhook","command":"ls -lah"}`)

	req, _ = http.NewRequest("PUT", "/service/1", bytes.NewBuffer(payload))
	req.Header.Set("APIKEY", "APIKEY")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalService["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalService["id"], m["id"])
	}

	if m["name"] == originalService["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalService["name"], m["name"], m["name"])
	}

	if m["endpoint"] == originalService["endpoint"] {
		t.Errorf("Expected the endpoint to change from '%v' to '%v'. Got '%v'", originalService["endpoint"], m["endpoint"], m["endpoint"])
	}

	if m["command"] == originalService["command"] {
		t.Errorf("Expected the endpoint to change from '%v' to '%v'. Got '%v'", originalService["command"], m["command"], m["command"])
	}
}

func TestDeleteService(t *testing.T) {
	clearTable()
	addServices(1)

	req, _ := http.NewRequest("GET", "/service/1", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/service/1", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/service/1", nil)
	req.Header.Set("APIKEY", "APIKEY")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
