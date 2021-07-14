package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main (){
	// Note: due to the usage of relative paths, this script has to be run from this directory (go run main.go).
	// Running from Goland directly may cause incorrect behaviour.

	log.Println("Creating the database...")
	file, err := os.Create("../db/moneyHeist.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Database created")

	heistDatabase, _ := sql.Open("sqlite3", "../db/moneyHeist.db")
	defer heistDatabase.Close()
	createMembersTable(heistDatabase)
	createMemberSkillsTable(heistDatabase)
	createSkillsTable(heistDatabase)
	createHeistMembersTable(heistDatabase)
	createHeistsTable(heistDatabase)
	createHeistSkillsTable(heistDatabase)
}

func createMembersTable (db *sql.DB) {
	createMembersTableSQL := `CREATE TABLE members (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL,
		"sex" TEXT NOT NULL,
		"email" TEXT NOT NULL,
		"mainSkill" TEXT,
		"status" TEXT NOT NULL
	);`

	log.Println("Creating members table...")
	statement, err := db.Prepare(createMembersTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Members table created")
}

func createMemberSkillsTable (db *sql.DB) {
	createMemberSkillsTableSQL := `CREATE TABLE memberSkills (
    	"memberId" TEXT NOT NULL FOREIGN KEY,
		"skillId" TEXT NOT NULL FOREIGN KEY,
		"level" TEXT,
		"experience" TEXT
	);`

	log.Println("Creating member skills table...")
	statement, err := db.Prepare(createMemberSkillsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Members skills table created")
}

func createSkillsTable (db *sql.DB) {
	createSkillsTableSQL := `CREATE TABLE skills (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL
	);`

	log.Println("Creating skills table...")
	statement, err := db.Prepare(createSkillsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Skills table created")
}

func createHeistMembersTable (db *sql.DB) {
	createHeistMembersTableSQL := `CREATE TABLE heistMembers (
		"memberId" TEXT NOT NULL FOREIGN KEY,
		"heistId" TEXT NOT NULL FOREIGN KEY
	);`

	log.Println("Creating heist members table...")
	statement, err := db.Prepare(createHeistMembersTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Heist members table created")
}

func createHeistsTable (db *sql.DB) {
	createHeistsTableSQL := `CREATE TABLE heists (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL,
		"location" TEXT NOT NULL,
		"startTime" TIME NOT NULL,
		"endTime" TIME N0T NULL,
		"status" TEXT NOT NULL,
		"outcome" TEXT
	);`

	log.Println("Creating heist table...")
	statement, err := db.Prepare(createHeistsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Heists table created")
}

func createHeistSkillsTable (db *sql.DB) {
	createHeistSkillsTableSQL := `CREATE TABLE heistSkills (
    	"heistId" TEXT NOT NULL FOREIGN KEY,
		"skillId" TEXT NOT NULL FOREIGN KEY,
		"members" INT NOT NULL,
		"level" TEXT,
	);`

	log.Println("Creating heist skills table...")
	statement, err := db.Prepare(createHeistSkillsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Heist skills table created")
}

