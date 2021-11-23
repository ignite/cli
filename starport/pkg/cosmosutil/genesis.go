package cosmosutil

import (
	"encoding/json"
	"os"
	"time"
)

const genesisTimeField = "genesis_time"

// SetGenesisTime sets the genesis time inside a genesis file
func SetGenesisTime(genesisPath string, genesisTime int64) error {
	// fetch and parse genesis
	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	genesis := make(map[string]interface{}, 0)
	if err := json.Unmarshal(genesisBytes, &genesis); err != nil {
		return err
	}

	// check the genesis time with the RFC3339 standard format
	formattedTime := time.Unix(genesisTime, 0).Format(time.RFC3339)

	// modify and save the new genesis
	genesis[genesisTimeField] = &formattedTime
	genesisBytes, err = json.Marshal(genesis)
	if err != nil {
		return err
	}
	return os.WriteFile(genesisPath, genesisBytes, 0644)
}
