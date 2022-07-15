package token

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

type EachState struct {
	RefreshToken string
}

type tokenList struct {
	tokens map[string]EachState
}

func (t tokenList) syncDB() {
	DBDir := getDatabaseDirectory()
	DBFilePath := path.Join(DBDir, databaseName)

	jsonData, err := json.Marshal(t.tokens)
	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile(DBFilePath, jsonData, 0755)
	if err != nil {
		log.Println(err)
	}
}

func (t tokenList) Add(state string, data *EachState) {
	t.tokens[state] = *data
	t.syncDB()
}

func getCurrentDatabase() map[string]EachState {
	DBDir := getDatabaseDirectory()
	DBFilePath := path.Join(DBDir, databaseName)

	file, err := os.ReadFile(DBFilePath)
	if err != nil {
		log.Println(err)
	}

	data := make(map[string]EachState)

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

var InitToken = tokenList {
	tokens: getCurrentDatabase(),
}
