package met

import (
	"log"
	"os"
	"testing"
)

func TestGetPayload(t *testing.T) {

	t.Run("Testing Get Payload", func(t *testing.T) {

		log := NewLogger("dummy.txt")

		payload := GetPayload(log, []float64{53.120607133644576, -3.9983421531498133})

		if len(payload) > 1 != true {
			t.Fatalf("payload empty")
		}
		log.Println(payload)
	})

}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}
