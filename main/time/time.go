package time

import "time"

func PrintCurrentTime() string {
	t := time.Now()

	return t.Format("2006년 1월 2일 오후 3시 4분") + "입니다."
}
