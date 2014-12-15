package golidera_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	g "github.com/wlMalk/golidera"
)

var translator g.TranslatorMap = g.TranslatorMap{
	"ar-SA": map[string]string{
		"username-cant-be":        "اسم المستخدم لا يمكن أن يكون {:par1}",
		"username-length-between": "ينبغي ان يكون طول اسم المستخدم بين {:par1} و {:par2} خانة",
		"password-cant-be":        "كلمة المرور لا يمكن أن تكون {:par1}",
		"password-length-between": "ينبغي ان يكون طول كلمة المرور بين {:par1} و {:par2} خانة",
	},
}

var validera = g.NewGolidera()

func init() {
	validera.SetLocalized(true)

	validera.Field("username", T(g.Not("username", "admin"), "username-cant-be")).
		When(
		g.FieldNotEmpty("id",
			T(g.LenBetween(4, 12), "username-length-between")))

	validera.Field("password",
		T(g.LenBetween(8, 16), "password-length-between"),
		T(g.Not("password"), "password-cant-be"))

	validera.Field("id")
}

func T(c g.CheckerFunc, s string) g.Checker {
	return g.T(c, s, translator)
}

type Tester func(method string, params url.Values) *httptest.ResponseRecorder

func GenerateTester(t testing.TB, handleFunc http.Handler) Tester {

	return func(method string, params url.Values) *httptest.ResponseRecorder {
		req, err := http.NewRequest(method, "", strings.NewReader(params.Encode()))
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

func FormPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passes, errors := validera.ValidateAll(validera.Request(r), r.FormValue("locale"))
		if !passes {
			w.WriteHeader(http.StatusBadRequest)
			for _, errs := range errors {
				for _, err := range errs {
					w.Write([]byte(err.Error()))
				}
			}
		}
	})
}

func TestGolideraHTTP(t *testing.T) {
	test := GenerateTester(t, FormPost())
	w := test("POST", url.Values{"username": []string{"admin"}, "password": []string{"password"}, "id": []string{"10"}, "locale": []string{"ar-SA"}})
	if w.Code != http.StatusBadRequest {
		t.Errorf("Home page didn't return %v", http.StatusBadRequest)
	}
}

func BenchmarkGolideraHTTP(t *testing.B) {
	test := GenerateTester(t, FormPost())
	for i := 0; i < t.N; i++ {
		w := test("POST", url.Values{"username": []string{"admin"}, "password": []string{""}, "id": []string{"10"}, "locale": []string{"ar-SA"}})
		if w.Code != http.StatusBadRequest {
			t.Errorf("Home page didn't return %v", http.StatusBadRequest)
		}
	}
}

func BenchmarkGolidera(t *testing.B) {
	for i := 0; i < t.N; i++ {
		f := g.NewForm()
		f.Field("username", "username")
		f.Field("password", "password")
		f.Field("id", "")
		_, _ = validera.ValidateAll(f, "ar-SA")
	}
}
