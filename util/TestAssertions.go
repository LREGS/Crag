package util

import "testing"

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
