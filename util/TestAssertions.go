package util

import (
	"net/http/httptest"
	"testing"
)

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct status, got %d, wanted %d", got, want)
	}
}
func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response is wrong, got %s, wanted %s", got, want)
	}
}

func CheckErrs(t testing.TB, r *httptest.ResponseRecorder, e string) string {
	var res Response

	_, err := DecodeResponse(r.Body, &res)
	if err != nil {
		t.Fatalf("could not decode because of err :%s", err)
	}

	if res.Error != e {
		t.Fatalf("expected error %s but got error %s", e, res.Error)
	}
	return ""

}
