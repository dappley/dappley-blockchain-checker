package helper

import (
	"io/ioutil"
	"strings"
	"strconv"
	"time"
	"log"
)

//Create a time frame relative to the current timestamp.
func TimeFrame(currTime time.Time) (before, now, after time.Time) {
	//Create layout for time.Time
	layout := "15:04:05"

	//Create lower bound
	lowerBound := "08:55:00"
	before, err := time.Parse(layout, lowerBound)
	if err != nil { log.Fatal(err) }

	//Create upper bound
	upperBound := "09:55:00"
	after, err = time.Parse(layout, upperBound)
	if err != nil { log.Fatal(err) }

	//Create current time
	curr := currTime.Format("15:04:05")
	now, err = time.Parse(layout, curr)
	if err != nil { log.Fatal(err) }
	
	return
}

//Checks if the timestamp recorded in the lastError.txt file is 24 hours before the current timestamp.
func ItsBeen24HrForAll(fileNames []string, currTime time.Time) bool {
	itsBeen24Hr   := true
	currTimeStamp := currTime.Unix()
	for _, fileName := range fileNames {
		serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
		data, err := ioutil.ReadFile("../lastError/lastError_" + serverType + ".txt")
		if err != nil {
			UpdateLastError(serverType, currTime)
			itsBeen24Hr = itsBeen24Hr && false
			continue
		}
		lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil { log.Fatal(err) }
		if (currTimeStamp - lastTimeStamp >= 86400) {
			itsBeen24Hr = itsBeen24Hr && true
			continue
		}
		itsBeen24Hr = itsBeen24Hr && false
	}
	return itsBeen24Hr
}