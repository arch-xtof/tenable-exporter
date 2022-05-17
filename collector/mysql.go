package collector

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/arch-xtof/tenable-exporter/models"
	_ "github.com/go-sql-driver/mysql"
)

func SqlConnect(user string, pass string) (db *sql.DB) {
	db, err := sql.Open("mysql", "tenable:tenable@tcp(mysql:3306)/tenable")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		_, err = db.Query("SELECT * FROM Asset")
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	log.Println("Connected to MySQL")
	return
}

func SqlDelete(db *sql.DB) {
	_, err := db.Query("DELETE FROM Asset")
	if err != nil {
		fmt.Println(err)
	}
}

func SqlInsert(db *sql.DB, vulns []models.Vuln, assets models.AssetMap) {
	sqlStrAsset := "INSERT IGNORE INTO Asset (uuid, hostname, asset_group) VALUES "
	valsAsset := []interface{}{}
	sqlStrVuln := "INSERT INTO Vuln (name, asset_uuid, severity) VALUES "
	valsVuln := []interface{}{}
	for _, vuln := range vulns {
		if vuln.SeverityModificationType == "ACCEPTED" {
			continue
		}
		asset := assets[vuln.Asset.Uuid]
		if asset["Type"] != "Servers" {
			continue
		}
		sqlStrAsset += "(?, ?, ?),"
		sqlStrVuln += "(?, ?, ?),"
		valsAsset = append(valsAsset, vuln.Asset.Uuid, vuln.Asset.Hostname, asset["Servers"])
		valsVuln = append(valsVuln, vuln.Plugin.Name, vuln.Asset.Uuid, vuln.Severity)
	}
	if len(valsAsset) == 0 || len(valsVuln) == 0 {
		return
	}
	sqlStrAsset = sqlStrAsset[0 : len(sqlStrAsset)-1]
	sqlStrVuln = sqlStrVuln[0 : len(sqlStrVuln)-1]

	stmtAsset, _ := db.Prepare(sqlStrAsset)
	stmtVuln, _ := db.Prepare(sqlStrVuln)

	_, err := stmtAsset.Exec(valsAsset...)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmtVuln.Exec(valsVuln...)
	if err != nil {
		fmt.Println(err)
	}
}
