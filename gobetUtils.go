package gobetUtils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"net/http"
	"strings"

	"encoding/json"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	RequestHeader RequestHeaderStruct `yaml:"Request Header"`
	EventTypes    []string            `yaml:"Event Types"`
	TimeRange     []string            `yaml:"Time Range"`
	Cycle         []int               `yaml:"Cycle"`
	Test          []string            `yaml:"Test"`
}

type RequestHeaderStruct struct {
	SessionToken      string `yaml:"session_token"`
	HeaderContentType string `yaml:"header_Content-Type"`
	HeaderAccept      string `yaml:"header_Accept"`
}

var Config ConfigStruct
var HeaderToken string

//func ReadBasicYamlConfig(f string) map[string]string {
func ReadBasicYamlConfig(f string) {

	filePath := f
	//fmt.Printf("Reading config file %s\n", filePath)

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}

	//sb := string(yamlFile)
	//log.Println(sb)

	//var c ConfigStruct

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	//fmt.Printf("SessionToken:%v - HeaderContentType:%v - HeaderAccept:%v\n", Config.RequestHeader.SessionToken, Config.RequestHeader.HeaderContentType, Config.RequestHeader.HeaderAccept)
	//fmt.Printf("arrayslice 0 for event types is %v", Config.EventTypes[0])

	//var configMap = make(map[string]string)
	//configMap["SessionToken"] = Config.RequestHeader.SessionToken
	//configMap["HeaderContentType"] = Config.RequestHeader.HeaderContentType
	//configMap["HeaderAccept"] = Config.RequestHeader.HeaderAccept

	if os.Getenv("BAPI_AUTH_TOKEN") != "" {
		//fmt.Print("BAPI Auth token set, adding to config map")
		//configMap["HeaderToken"] = os.Getenv("BAPI_AUTH_TOKEN")
		HeaderToken = os.Getenv("BAPI_AUTH_TOKEN")
	} else {
		log.Fatal("No BAPI auth token set, quitting")
		os.Exit(1)
	}

	//return configMap

}

//func SubmitRequest(requestType string, headerMap map[string]string, endpoint string, filter string) []byte {

func SubmitAPIRequest(configFile string, APIType string, requestType string, endpoint string, filter string) []byte {

	/*fmt.Println("Map received")
	fmt.Println(headerMap)

	fmt.Println("requestType is" + requestType)
	fmt.Println("endpoitn is " + endpoint)
	fmt.Println("filter is " + filter)

	*/

	// Let's read in our basic config first, without this, we're not going anywhere
	//configMap := ReadBasicYamlConfig(configFile)
	ReadBasicYamlConfig(configFile)

	// Now, let's work out what endpoint we want to hit to build the url
	var url string
	switch APIType {

	case "EXCHANGE":
		//fmt.Println("Exchange API request received")
		url = "https://api.betfair.com/exchange/betting/rest/v1.0/"
	case "ACCOUNT":
		//fmt.Println("Account API request received")
		url = "https://api.betfair.com/exchange/account/rest/v1.0/"
	default:
		fmt.Println("Unknown API Type requested, exiting")
		os.Exit(1)
	}

	responseBody := strings.NewReader(filter)
	client := http.Client{}
	req, err := http.NewRequest(requestType, url+endpoint, responseBody)

	if err != nil {
		log.Fatalf("An error occurred %v", err)
	}

	req.Header = http.Header{
		//"X-Application":    {configMap["HeaderToken"]},
		"X-Application": {HeaderToken},
		//"X-Authentication": {configMap["SessionToken"]},
		"X-Authentication": {Config.RequestHeader.SessionToken},
		//"Content-Type":     {configMap["HeaderContentType"]},
		"Content-Type": {Config.RequestHeader.HeaderContentType},
		//"Accept":           {configMap["HeaderAccept"]},
		"Accept": {Config.RequestHeader.HeaderAccept},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("default client error")
	}

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return body

}

func TrackMarket(market string, marketName string, runners map[int]string) {

	runnerBackOddsMap := make(map[string][]float64)
	runnerLayOddsMap := make(map[string][]float64)
	fmt.Printf("recieved runnermap for market %v and market name %v\n", market, marketName)
	fmt.Printf("runners map %v\n", runners)

	/*type RunnerStruct struct {
		SelectionId int `json:"selectionId"`
		Status	string	`json:"status"`
		LastPriceTraded	float32	`json:"lastPriceTraded"`
		TotalMatched	float32	`json:"totalMatched"`
		Prices
	}

	type ListMarketBookStruct struct {
		MarketId       string         `json:"marketId"`
		Status         string         `json:"status"`
		LastMatchTime  string         `json:"lastMatchTime"`
		TotalMatched   float32        `json:"totalMatched"`
		TotalAvailable float32        `json:"totalAvailable"`
		Runners        []RunnerStruct `json:"runners"`
	}

	*/

	type ListMarketBookStruct struct {
		MarketID              string    `json:"marketId"`
		IsMarketDataDelayed   bool      `json:"isMarketDataDelayed"`
		Status                string    `json:"status"`
		BetDelay              int       `json:"betDelay"`
		BspReconciled         bool      `json:"bspReconciled"`
		Complete              bool      `json:"complete"`
		Inplay                bool      `json:"inplay"`
		NumberOfWinners       int       `json:"numberOfWinners"`
		NumberOfRunners       int       `json:"numberOfRunners"`
		NumberOfActiveRunners int       `json:"numberOfActiveRunners"`
		LastMatchTime         time.Time `json:"lastMatchTime"`
		TotalMatched          float64   `json:"totalMatched"`
		TotalAvailable        float64   `json:"totalAvailable"`
		CrossMatching         bool      `json:"crossMatching"`
		RunnersVoidable       bool      `json:"runnersVoidable"`
		Version               int64     `json:"version"`
		Runners               []struct {
			SelectionID     int     `json:"selectionId"`
			Handicap        float64 `json:"handicap"`
			Status          string  `json:"status"`
			LastPriceTraded float64 `json:"lastPriceTraded"`
			TotalMatched    float64 `json:"totalMatched"`
			Ex              struct {
				AvailableToBack []struct {
					Price float64 `json:"price"`
					Size  float64 `json:"size"`
				} `json:"availableToBack"`
				AvailableToLay []struct {
					Price float64 `json:"price"`
					Size  float64 `json:"size"`
				} `json:"availableToLay"`
				TradedVolume []interface{} `json:"tradedVolume"`
			} `json:"ex"`
		} `json:"runners"`
	}

	fmt.Printf("Received market ID: %v to track\n", market)

	// This function takes a market and continually polls it looking for interesting changes in odds (interesting as defined by user in config).
	// If the change has been found, a back or lay will be made

	filter := `{"marketIds":["` + market + `"],"priceProjection":{"priceData":["EX_BEST_OFFERS"]}, "id": 1}`

	for range time.Tick(time.Second * 10) {

		var listMarketOutput []ListMarketBookStruct
		//backOddsSlice := make([]float64, 0)

		fmt.Printf("Market data for %v...\n", market)
		body := SubmitAPIRequest("config.yaml", "EXCHANGE", "POST", "listMarketBook/", filter)

		//sb := string(body)
		//log.Println(sb)

		err := json.Unmarshal(body, &listMarketOutput)
		if err != nil {
			print(err)
		}

		for _, value := range listMarketOutput {
			fmt.Printf("Market ID:%v - Status:%v - Matched:%f\n", value.MarketID, value.Status, value.TotalMatched)
			fmt.Printf("Runner data...\n")
			for _, runner := range value.Runners {
				backCount := 1
				fmt.Printf("MARKET ID: %v --- Runner ID: %v - Last Price Traded: %f\n", value.MarketID, runners[runner.SelectionID], runner.LastPriceTraded)
				for _, backOdds := range runner.Ex.AvailableToBack {
					fmt.Printf("MARKET ID: %v --- Runner ID: %v - BackPrice %f - BackSize %f\n", value.MarketID, runners[runner.SelectionID], backOdds.Price, backOdds.Size)
					// let's create an a map of slices to keep our last minutes worth of odds movement
					if backCount == 1 {
						// This wil be the current available back price
						_, found := runnerBackOddsMap[runners[runner.SelectionID]]
						if !found {
							fmt.Printf("No back odds map found for %v\n", runners[runner.SelectionID])
							//s := strconv.FormatFloat(backOdds.Price, 'f', -1, 64)
							backOddsSlice := []float64{backOdds.Price}
							runnerBackOddsMap[runners[runner.SelectionID]] = backOddsSlice
							fmt.Printf("back oddds map for %v is %v\n", runners[runner.SelectionID], runnerBackOddsMap[runners[runner.SelectionID]])
						} else {
							fmt.Printf("back odds map found we need a routine that will build the string properly here\n")
							fmt.Printf("back odds map for %v is %v\n", runners[runner.SelectionID], runnerBackOddsMap[runners[runner.SelectionID]])
						}
					}
					backCount++
				}

				layCount := 1
				for _, layOdds := range runner.Ex.AvailableToLay {
					fmt.Printf("MARKET ID: %v --- Runner ID: %v - LayPrice %f - LaySize %f\n", value.MarketID, runners[runner.SelectionID], layOdds.Price, layOdds.Size)
					if layCount == 1 {
						// This wil be the current available back price
						_, found := runnerLayOddsMap[runners[runner.SelectionID]]
						if !found {
							fmt.Printf("No lay odds map found for %v\n", runners[runner.SelectionID])
							//s := strconv.FormatFloat(backOdds.Price, 'f', -1, 64)
							layOddsSlice := []float64{layOdds.Price}
							runnerLayOddsMap[runners[runner.SelectionID]] = layOddsSlice
							fmt.Printf("lay oddds map for %v is %v\n", runners[runner.SelectionID], runnerLayOddsMap[runners[runner.SelectionID]])
						} else {
							fmt.Printf("lay odds map found we need a routine that will build the string properly here\n")
							fmt.Printf("lay odds map for %v is %v\n", runners[runner.SelectionID], runnerLayOddsMap[runners[runner.SelectionID]])
						}
					}
					layCount++
				}
			}
		}

	}

}

/*
type MarketTrackingUpdate struct {
	Job          string
	MarketStatus string
}

var marketTrackingJobs = make(chan string, 10)
var marketTrackingReturn = make(chan MarketTrackingUpdate, 10)

func TrackMarket(marketID string) MarketTrackingUpdate {

	fmt.Printf("Track Market called on %v", marketID)
	marketUpdate := MarketTrackingUpdate{marketID, "OPEN"}
	return marketUpdate

}

func MarketWorker(wg *sync.WaitGroup) {
	for job := range marketTrackingJobs {
		output := TrackMarket(job)
		marketTrackingReturn <- output
	}
	wg.Done()
}

func CreateMarketWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go MarketWorker(&wg)
	}
	wg.Wait()
	//close(marketTrackingReturn)
}

func AllocateTrackMarketJob(market string) {
	//for i := 0; i < noOfJobs; i++ {
	//randomno := rand.Intn(999)
	//job := Job{i, randomno}
	//jobs <- job

	job := market
	marketTrackingJobs <- job

	//}

	//close(marketTrackingJobs)
}

func MarketReturn(done chan bool) {
	for result := range marketTrackingReturn {
		fmt.Printf("Market id %v, Status %v\n", result.Job, result.MarketStatus)
	}
	done <- true
}

*/
