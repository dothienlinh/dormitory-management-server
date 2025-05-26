package common

import "time"

func GetExpireTime(minute int) string {
	return time.Now().Add(time.Duration(minute) * time.Minute).Format(time.RFC3339)
}

func ParsedTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
