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
	file, err := os.Create("../db/heist.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Database created")

	heistDatabase, _ := sql.Open("sqlite3", "../db/heist.db")
	defer heistDatabase.Close()
	createMembersTable(heistDatabase)
	createMemberSkillsTable(heistDatabase)
	createSkillsTable(heistDatabase)
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
	createSkillsTableSQL := `CREATE TABLE memberSkills (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL,
		"level" TEXT,
		"memberId" TEXT NOT NULL
	);`

	log.Println("Creating member skills table...")
	statement, err := db.Prepare(createSkillsTableSQL)
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

