package swaggerui

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerHandler(t *testing.T) {
	h := SwaggerHandler()

	// Test if index template gets outputted
	{
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("status not ok")
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), `<title>Swagger UI</title>`) {
			t.Fatal("expected Swagger UI index page")
		}
	}

	// Shallow test if filesystem does it's job
	{
		r := httptest.NewRequest("GET", "/test123file", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("status not 404")
		}
	}

	// Shallow test if filesystem does it's job
	{
		files := []string{"favicon-16x16.png", "favicon-32x32.png",
			"oauth2-redirect.html", "swagger-ui.css",
			"swagger-ui-bundle.js", "swagger-ui-standalone-preset.js"}

		for _, f := range files {
			r := httptest.NewRequest("GET", "/"+f, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			res := w.Result()
			if res.StatusCode != http.StatusOK {
				t.Fatal("status not ok")
			}
		}
	}
}
func TestSwaggerUrlHandler(t *testing.T) {
	h := SwaggerUrlHandler("test_url_is_test")

	// Test if the index has set up the url
	{
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("status not ok")
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), `test_url_is_test`) {
			t.Fatal("expected url to be test_url_is_test")
		}
	}

	// Ensure swagger.json shows 404 when url is used
	{
		r := httptest.NewRequest("GET", "/swagger.json", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("status not 404")
		}
	}
}

func TestSwaggerFileUIHandler(t *testing.T) {

	// Test if it converts a yaml file to a swagger json file
	{
		h := SwaggerFileHandler("./fixtures/swagger_test.yml")
		r := httptest.NewRequest("GET", "/swagger.json", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("status not ok")
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !isJSON(string(body)) {
			t.Fatal("response was not a json format")
		}
	}

	// Check if a existing swagger json file gets outputted
	{
		h := SwaggerFileHandler("./fixtures/swagger.json")
		r := httptest.NewRequest("GET", "/swagger.json", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("status not ok")
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !isJSON(string(body)) {
			t.Fatal("response was not a json format")
		}
	}

	// Check if it errors when given file does not exist
	{
		h := SwaggerFileHandler("./not_existing")
		r := httptest.NewRequest("GET", "/swagger.json", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("status not 404")
		}
	}

	// Check if it displays 500 if the yaml file is not correct format
	{
		h := SwaggerFileHandler("./fixtures/swagger_not.yml")
		r := httptest.NewRequest("GET", "/swagger.json", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		res := w.Result()
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("status not 500")
		}
	}
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
