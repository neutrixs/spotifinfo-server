package state

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

type statesScope struct {
	states map[string]string
}

func (s statesScope) Add(state, scope string) {
	s.states[state] = scope

	DBDIR := getDatabaseDirectory()
	DBFilePath := path.Join(DBDIR, databaseName)

	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	os.WriteFile(DBFilePath, jsonData, 0755)
}

var InitStates = statesScope {
	states: make(map[string]string),
}