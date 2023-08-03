package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var client = resty.New().SetHostURL("http://localhost:5984").SetBasicAuth("admin", "YOURPASSWORD").SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json")

func removeExpiredNonces() {
	totalRows := float64(1) // start value to enter the loop
	for totalRows > 0 {

		response, _ := client.R().Get(fmt.Sprintf("%s/%s", "mydb", "_design/nonce/_view/older_than?limit=500"))
		var expiredNonces map[string]interface{}
		json.Unmarshal(response.Body(), &expiredNonces)

		totalRows = expiredNonces["total_rows"].(float64)

		if totalRows == float64(0) {
			break
		}

		fmt.Printf("total rows: %d\n", int(totalRows))

		bulkDelete := []map[string]interface{}{}
		for _, nonceDoc := range expiredNonces["rows"].([]interface{}) {
			delteDoc := map[string]interface{}{
				"_id":      nonceDoc.(map[string]interface{})["id"],
				"_rev":     nonceDoc.(map[string]interface{})["value"].(string),
				"_deleted": true,
			}
			bulkDelete = append(bulkDelete, delteDoc)
		}

		deleteBody := map[string]interface{}{
			"docs": bulkDelete,
		}
		response, _ = client.R().SetBody(deleteBody).Post(fmt.Sprintf("%s/%s", "mydb", "_bulk_docs"))
		var deletedList []map[string]interface{}
		json.Unmarshal(response.Body(), &deletedList)

		fmt.Printf("Deleted %d nonces\n", len(deletedList))
	}

	return
}

func main() {
	for {
		removeExpiredNonces()
		fmt.Printf("No nonces found, waiting 5 seconds before checking again\n")
		time.Sleep(time.Second * 5)
	}
}
