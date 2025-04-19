package entity

import "time"

// TimeNow returns the current time
// This function is used to make testing easier by allowing time to be mocked
func TimeNow() time.Time {
	return time.Now()
}
