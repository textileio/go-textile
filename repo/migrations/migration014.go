package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor014 struct{}

func (Minor014) Up(repoPath string, pinCode string, testnet bool) error {
	var dbPath string
	if testnet {
		dbPath = path.Join(repoPath, "datastore", "testnet.db")
	} else {
		dbPath = path.Join(repoPath, "datastore", "mainnet.db")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	if pinCode != "" {
		_, err = db.Exec("pragma key='" + pinCode + "';")
		if err != nil {
			return err
		}
	}

	query := `
    alter table blocks add column data text not null default '';
    alter table blocks add column status integer not null default 0;
    alter table blocks add column attempts integer not null default 0;
    create index block_data on blocks (data);
    create index block_status on blocks (status);
    `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// target -> data
	_, err = db.Exec(`
    update blocks set data=target where type=7;
    update blocks set target='' where type=7;
    `)
	if err != nil {
		return err
	}

	// update version
	f15, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f15.Close()
	if _, err = f15.Write([]byte("15")); err != nil {
		return err
	}
	return nil
}

func (Minor014) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor014) Major() bool {
	return false
}
