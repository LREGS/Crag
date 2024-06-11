package met

import (
	"testing"
	"time"
)

func TestingSchedule(t *testing.T) {

	mockLastUpdateTime := time.Date(2024, time.June, 9, 2,0,0,0,time.Local)

	scheduler := &scheduler{}
	scheduler.start(mockLastUpdateTime)

}
S