package unittest

import (
	"data-connector/controller"
	"testing"
)

func TestGetTokenMisa_pass(t *testing.T) {
	getTokenMisa := controller.GetTokenMisa
	partner := `{"namePartner": "Misa", "clientId": "demoamisapp", "clientSecret": "L8rKc7hYtlKq+QOgD4RRHF9VM4Gzq6ix8HnEUTGYBIM="}`
	result := getTokenMisa(partner)
	if result == "" {
		t.Errorf("fail")
	} else {
		t.Logf("cilentid or clientSecret not correct: %s", result)
	}
}
func TestGetTokenMisa_fail(t *testing.T) {
	getTokenMisa := controller.GetTokenMisa
	partner := `{"namePartner": "Misa", "clientId": "demo", "clientSecret": "L8rKc7hYtlKq+QOgD4RRHF9VM4Gzq6ix8HnEUTGYBIM="}`
	result := getTokenMisa(partner)
	if result != "" {
		t.Errorf("fail")
	} else {
		t.Logf("pass; got: %s", result)
	}
}
