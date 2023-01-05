package main

import (
	"github.com/calvarado2004/go-concurrency/data"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var pageTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string
}{
	{
		name:               "Home Page",
		url:                "/",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.HomePage,
	},
	{
		name:               "Login Page",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:               "Register Page",
		url:                "/register",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.RegisterPage,
		expectedHTML:       `<h1 class="mt-5">Register</h1>`,
	},
	{
		name:               "Logout Page",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.LogoutPage,
		sessionData: map[string]any{
			"userID": 1,
			"user": &data.User{
				ID:       1,
				Email:    "admin@test.com",
				Password: "password",
			},
		},
	},
	{
		name:               "Activate Page",
		url:                "/activate",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.ActivateAccount,
	},
}

func Test_Pages(t *testing.T) {
	pathToTemplates = "./templates"

	for _, tt := range pageTests {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", tt.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := getContext(req)

		req = req.WithContext(ctx)

		if len(tt.sessionData) > 0 {
			for key, value := range tt.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		tt.handler.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("handler %v returned wrong status code: got %v want %v", tt.name, rr.Code, tt.expectedStatusCode)
		}

		if len(tt.expectedHTML) > 0 {
			html := rr.Body.String()

			if !strings.Contains(html, tt.expectedHTML) {
				t.Errorf("handler returned unexpected body: got %v want %v", html, tt.expectedHTML)
			}
		}

	}

}

func TestConfig_PostLoginPage(t *testing.T) {
	pathToTemplates = "./templates"

	postedData := url.Values{
		"email":    {"admin@test.com"},
		"password": {"password"},
	}

	rr := httptest.ResponseRecorder{}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	ctx := getContext(req)

	req = req.WithContext(ctx)

	handler := http.HandlerFunc(testApp.PostLoginPage)

	handler.ServeHTTP(&rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusSeeOther)
	}

	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("Expected session to exist")
	}

}
