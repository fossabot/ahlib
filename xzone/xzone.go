package xzone

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Parse Timezone string to time.Location. Format: `^[+-][0-9]{1,2}([0-9]{1,2})?$`
func ParseTimeZone(zone string) (*time.Location, error) {
	regex, err := regexp.Compile(`^([+-])([0-9]{1,2})(?::([0-9]{1,2}))?$`)
	if err != nil {
		return nil, err
	}

	wrongFmtErr := fmt.Errorf("timezone string has a wrong format")
	ok := regex.Match([]byte(zone))
	if !ok {
		return nil, wrongFmtErr
	}

	matches := regex.FindAllStringSubmatch(zone, 1)
	if len(matches) == 0 || len(matches[0][1:]) < 3 {
		return nil, wrongFmtErr
	}
	group := matches[0][1:]

	signStr := group[0]
	hourStr := group[1]
	minuteStr := group[2]
	if signStr != "+" && signStr != "-" {
		return nil, wrongFmtErr
	}
	if minuteStr == "" {
		minuteStr = "0"
	}

	sign := +1
	if signStr == "-" {
		sign = -1
	}
	hour, err1 := strconv.Atoi(hourStr)
	minute, err2 := strconv.Atoi(minuteStr)
	if err1 != nil || err2 != nil {
		return nil, wrongFmtErr
	}

	name := fmt.Sprintf("UTC%s%02d:%02d", signStr, hour, minute)
	offset := sign * (hour*3600 + minute*60)
	return time.FixedZone(name, offset), nil
}

// Parse timezone string and move time to specific timezone.
func MoveToZone(t time.Time, zone string) (time.Time, error) {
	loc, err := ParseTimeZone(zone)
	if err != nil {
		return t, err
	}

	return t.In(loc), nil
}
