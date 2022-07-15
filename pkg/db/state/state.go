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

func (s statesScope) syncDB() {
	DBDIR := getDatabaseDirectory()
	DBFilePath := path.Join(DBDIR, databaseName)

	jsonData, err := json.Marshal(s.states)
	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile(DBFilePath, jsonData, 0755)
	if err != nil {
		log.Println(err)
	}
}

func (s statesScope) Add(state, scope string) {
	s.states[state] = scope
	s.syncDB()
}

func (s statesScope) Remove(state string) {
	delete(s.states, state)
	s.syncDB()
}

func (s statesScope) Get(state string) string {
	return s.states[state]
}

func getCurrentDatabase() map[string]string {
	DBDir := getDatabaseDirectory()
	DBFilePath := path.Join(DBDir, databaseName)

	file, err := os.ReadFile(DBFilePath)
	if err != nil {
		log.Println(err)
	}

	var data = make(map[string]string)

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

var InitStates = statesScope {
	states: getCurrentDatabase(),
}