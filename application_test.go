package main

import (
	"errors"
	"github.com/quii/mockingjay-server/mockingjay"
	"log"
	"net/http"
	"testing"
)

const someMonkeyConfigString = "Hello, world"

func TestItFailsWhenTheConfigFileCantBeLoaded(t *testing.T) {
	app := testApplication()
	app.configLoader = failingIOUtil

	err := app.Run("mockingjay config path", 1234, "", "")

	if err == nil {
		t.Error("Didnt get an error when the config file failed to load")
	}
}

func TestItFailsWhenTheConfigIsInvalid(t *testing.T) {
	app := testApplication()
	app.mockingjayLoader = failingMockingjayLoader

	err := app.Run("mockingjay config path", 1234, "", "")

	if err == nil {
		t.Error("Didnt get an error when the mockingjay config failed to load")
	}
}

func TestItFailsWhenTheMonkeyConfigIsInvalid(t *testing.T) {
	app := testApplication()

	err := app.Run("mockingjay config path", 1234, "", "monkey config path")

	if err == nil {
		t.Error("Didnt get an error when the monkey config failed to load")
	}
}

func testApplication() *application {
	app := new(application)
	app.configLoader = passingIOUtil
	app.mockingjayLoader = passingMockingjayLoader
	app.mockingjayServerMaker = mockingjay.NewServer
	app.monkeyServerMaker = failingMonkeyServerMaker
	return app
}

func testMockingJayConfig() []mockingjay.FakeEndpoint {

	m, err := mockingjay.NewFakeEndpoints([]byte(testYAML("hello, world")))

	if err != nil {
		log.Fatal(err)
	}

	return m
}

func passingIOUtil(path string) ([]byte, error) {
	return []byte(someMonkeyConfigString), nil
}

func failingIOUtil(path string) ([]byte, error) {
	return nil, errors.New("Couldnt load file")
}

func failingMockingjayLoader([]byte) ([]mockingjay.FakeEndpoint, error) {
	return nil, errors.New("Couldn't load file")
}

func passingMockingjayLoader([]byte) ([]mockingjay.FakeEndpoint, error) {
	return testMockingJayConfig(), nil
}

func failingMonkeyServerMaker(http.Handler, string) (http.Handler, error) {
	log.Println("what")
	return nil, errors.New("Couldn't load monkey config")
}
