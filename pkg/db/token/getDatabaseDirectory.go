package token

import (
	"log"
	"os"
	"path"

	"github.com/neutrixs/spotifinfo-server/pkg/env"
)

const databaseName = "token.json"

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	DBDir := getDatabaseDirectory()
	DBFilePath := path.Join(DBDir, databaseName)

	_, err := os.ReadFile(DBFilePath)
	if err != nil {
		err := os.MkdirAll(DBDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(DBFilePath, []byte("{}"), 0755)
	}
}

func getDatabaseDirectory() string {
	databaseDirectory, err := env.Get("DB_DIR")
	if err != nil {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		databaseDirectory = path.Join(homedir, "spotifinfo-db")
	}

	return databaseDirectory
}