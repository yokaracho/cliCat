package root

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/yokaracho/cliCat/cmd/logger"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

const (
	filePermission = 0644
)

var getLogger = logger.GetLogger()

func loadEnv() {
	// load env variables
	if err := godotenv.Load("/home/dev/cliCat/.env"); err != nil {
		getLogger.Fatalf("error loading env variables: %s", err.Error())
	}
	getLogger.Infof("Env variables loading successfully")
}

func getCatBreeds() ([]byte, error) {
	resp, err := http.Get(os.Getenv("API"))
	if err != nil {
		getLogger.Errorf("Error fetching cat breeds: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		getLogger.Errorf("Error fetching cat breeds. Status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		getLogger.Errorf("Error reading response body: %v", err)
		return nil, err
	}

	return body, nil
}

// decode json
func decodeJSON(data []byte) (ApiResponse, error) {
	var apiResponse ApiResponse
	err := json.Unmarshal(data, &apiResponse)
	if err != nil {
		return apiResponse, err
	}

	return apiResponse, nil
}

// group cat by country
func groupBreedsByCountry(apiResponse ApiResponse) map[string][]string {
	breedsByCountry := make(map[string][]string)

	for _, breedInfo := range apiResponse.Data {
		breedsByCountry[breedInfo.Country] = append(breedsByCountry[breedInfo.Country], breedInfo.Breed)
	}

	return breedsByCountry
}

// sort name cat by length
func sortBreedsByLength(breedsByCountry map[string][]string) {
	for _, breeds := range breedsByCountry {
		sort.Slice(breeds, func(i, j int) bool {
			return len(breeds[i]) < len(breeds[j])
		})
	}
	getLogger.Infof("Breed names sorted by length")

}

func runCmd(cmd *cobra.Command, args []string) {
	loadEnv()

	breeds, err := getCatBreeds()
	if err != nil {
		getLogger.Fatalf("Error getting cat breeds: %v", err)
		return
	}

	apiResponse, err := decodeJSON(breeds)
	if err != nil {
		getLogger.Fatalf("Error decoding JSON: %v", err)
		return
	}

	breedsByCountry := groupBreedsByCountry(apiResponse)
	sortBreedsByLength(breedsByCountry)

	outputJSON, err := json.MarshalIndent(breedsByCountry, "", "    ")
	if err != nil {
		getLogger.Fatalf("Error encoding JSON: %v", err)
		return
	}
	//write json to file
	err = ioutil.WriteFile(os.Getenv("FILE"), outputJSON, filePermission)
	if err != nil {
		getLogger.Fatalf("Error decoding JSON: %v", err)
		return
	}

	getLogger.Infof("Data written to out.json successfully.")
}
