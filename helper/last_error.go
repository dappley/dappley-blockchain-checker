package helper

import(
	"strconv"
	"strings"
	"time"
	"fmt"
	"log"
	"os"
)

//If the lastError folder does not exists in the repository, assume that the lastError occured
//24 hours prior to the current testing time then create lastError directory and lastError files.
func IsLastErrorExist(fileNames []string, currTime time.Time) {
	lastErrorExists := isExist("../lastError")
	if !lastErrorExists {
		fmt.Println("Creating the lastError directory...")
		err := os.Mkdir("../lastError", os.ModePerm)
		if err != nil {
			fmt.Println("Could not create the lastError directory!")
			panic(err)
		}
		for _, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
			yesterday  := currTime.AddDate(0, 0, -1)
			UpdateLastError(serverType, yesterday)
		}
	}
}

//Update/Create a lastError file with the current timestamp
func UpdateLastError(serverType string, currTime time.Time) {
	//Create/OverRide the file
	f, err := os.Create("../lastError/lastError_" + serverType + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	//Write current timestamp to the lastError file
	_, err = f.WriteString(strconv.FormatInt(currTime.Unix(), 10))
	if err != nil {
		log.Fatal(err)
	}
}

func isExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}