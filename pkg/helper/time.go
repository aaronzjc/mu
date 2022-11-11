package helper

import "time"

const (
	LocalTimezone = "Asia/Shanghai"
	LayoutISO     = "2006-01-02 15:04:05"
)

var localTimezone *time.Location

func init() {
	localTimezone, _ = time.LoadLocation(LocalTimezone)
}

func TimeToLocalStr(t time.Time) string {
	return t.In(localTimezone).Format(LayoutISO)
}

func CurrentTimeStr() string {
	return time.Now().In(localTimezone).Format(LayoutISO)
}
