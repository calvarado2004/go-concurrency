package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := getContext(req)

	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.AddDefaultData(&TemplateData{}, req)

	if td.Flash != "flash" {
		t.Error("flash not found in session")
	}

	if td.Warning != "warning" {
		t.Error("warning not found in session")
	}

	if td.Error != "error" {
		t.Error("error not found in session")
	}

}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := getContext(req)

	req = req.WithContext(ctx)

	if testApp.IsAuthenticated(req) {
		t.Error("got authenticated when not expecting it")
	}

	testApp.Session.Put(ctx, "userID", 1)

	if !testApp.IsAuthenticated(req) {
		t.Error("got not authenticated when expecting it")
	}

}

func TestConfig_Render(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := getContext(req)

	req = req.WithContext(ctx)

	pathToTemplates = "./templates"

	rr := httptest.NewRecorder()

	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	if rr.Code != http.StatusOK {
		t.Errorf("got %d want %d", rr.Code, http.StatusOK)
	}

}
