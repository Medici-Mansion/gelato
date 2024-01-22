package time

import (
	"fmt"
	"time"
)

func PrintCurrentTime() string {
	t := time.Now()
	hour, minute, second := t.Hour(), t.Minute(), t.Second()
	ampm := "오전"

	if hour >= 12 {
		hour -= 12
		ampm = "오후"
	}

	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("현재 시각은 %d년 %d월 %d일 %s %d시 %d분 %d초입니다.", t.Year(), t.Month(), t.Day(), ampm, hour, minute, second)
}
