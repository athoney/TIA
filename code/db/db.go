package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// we have to import the driver, but don't use it in our code
	// so we use the `_` symbol
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Product struct {
	VendorProject string
	Product       string
}

type ProdW struct {
	Cvekey            string
	Cveid             string
	VendorProject     string
	Product           string
	VulnerabilityName string
	DateAdded         string
	ShortDescription  string
	RequiredAction    string
	DueDate           string
	Notes             string
}

// Initiate db connection and load .csv file.
// Returns db context for queries in API
func Main() *sql.DB {
	curdir, err := os.Getwd()

	fmt.Println("current directory: ", curdir)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("DBPORT"))
	user := os.Getenv("DBUSER")
	pw := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	// Initiate server connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pw, dbname)
	conn, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Read and process entries from the exploited_vulnerabilities .csv file
	entries := readCsv("known_exploited_vulnerabilities.csv")
	processEntries(entries, conn)
	weaknesses := readCsv("weaknesses.csv")
	processWeak(weaknesses, conn)
	pars := readCsv("parent.csv")
	processPar(pars, conn)
	cwe := readCsv("cwe.csv")
	processCWE(cwe, conn)
	base := readCsv("newbasescore.csv")
	processBaseScore(base, conn)

	return conn
}

// Reads in from .csv file and returns list of formatted entries
func readCsv(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(file, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	entries, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(file, err)
	}
	return entries
}

// Creates the vulnerabilities table in db and inserts the .csv file entries
func processEntries(ent [][]string, dat *sql.DB) {
	len := len(ent)
	//Drop the table so we can then update it
	dat.Exec("DROP TABLE vulnerabilities;")
	//Create the table
	dat.Exec("CREATE TABLE vulnerabilities (cveKey INT8 PRIMARY KEY, cveID VARCHAR(1024), vendorProject VARCHAR(1024), product VARCHAR(1024), vulnerabilityName VARCHAR(1024), dateAdded DATE NOT NULL, shortDescription VARCHAR(2048), requiredAction VARCHAR(2048), dueDate DATE NOT NULL, notes VARCHAR(1024)) ;")
	totalRows := 0
	for i := 1; i < len; i++ {
		k := ent[i][0]
		key := removeCve(k)
		result, err := dat.Exec("INSERT INTO vulnerabilities (cveKey, cveID, vendorProject, product, vulnerabilityName, dateAdded, shortDescription, requiredAction, dueDate, notes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);", key, ent[i][0], ent[i][1], ent[i][2], ent[i][3], ent[i][4], ent[i][5], ent[i][6], ent[i][7], ent[i][8])
		if err != nil {
			log.Fatalf("row not affected: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		totalRows += int(rowsAffected)
	}
	fmt.Println("Total rows inserted for vulnerabilities table: ", totalRows)

}

// Creates CWE categories table in db and inserts the .csv file entries
func processWeak(ent [][]string, dat *sql.DB) {
	len := len(ent)
	fmt.Printf(strconv.Itoa(len))
	//Drop the table so we can then update it
	dat.Exec("DROP TABLE weak;")
	//Create the table
	dat.Exec("CREATE TABLE weak (ParentID INTEGER, CWEname VARCHAR(1024), CWEID INTEGER) ;")
	totalRows := 0
	for i := 1; i < len; i++ {
		result, err := dat.Exec("INSERT INTO weak (ParentID, CWEname, CWEID) VALUES ($1,$2,$3);", ent[i][0], ent[i][1], ent[i][2])
		if err != nil {
			log.Fatalf("row not affected: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		totalRows += int(rowsAffected)
	}
	fmt.Println("Total rows inserted for weak table: ", totalRows)
}

// Creates the CWE source table in db and inserts the .csv file entries
func processCWE(ent [][]string, dat *sql.DB) {
	len := len(ent)
	fmt.Printf(strconv.Itoa(len))
	// Drop the table so we can then update it
	dat.Exec("DROP TABLE IF EXISTS cwe;")
	// Create the table
	dat.Exec("CREATE TABLE cwe (CVEId VARCHAR(1024), CWEName VARCHAR(1024), Source VARCHAR(1024)) ;")
	totalRows := 0
	for i := 1; i < len; i++ {
		result, err := dat.Exec("INSERT INTO cwe (CVEId, CWEName, Source) VALUES ($1,$2,$3);", ent[i][0], ent[i][1], ent[i][2])
		if err != nil {
			log.Fatalf("row not affected: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		totalRows += int(rowsAffected)
	}
	fmt.Println("Total rows inserted for cwe table: ", totalRows)
}

// Creates the CWE categories table in db and inserts the .csv file entries
func processBaseScore(ent [][]string, dat *sql.DB) {
	len := len(ent)
	fmt.Printf(strconv.Itoa(len))
	//Drop the table so we can then update it
	dat.Exec("DROP TABLE IF EXISTS base;")
	//Create the table
	dat.Exec("CREATE TABLE base (CVEId VARCHAR(1024),  V3Score VARCHAR(1024), V3Vector VARCHAR(1024), V2Score VARCHAR(1024), V2Vector VARCHAR(1024)) ;")
	totalRows := 0
	for i := 1; i < len; i++ {
		result, err := dat.Exec("INSERT INTO base (CVEId, V3Score, V3Vector, V2Score, V2Vector) VALUES ($1,$2,$3,$4,$5);", ent[i][0], ent[i][1], ent[i][2], ent[i][3], ent[i][4])
		if err != nil {
			log.Fatalf("row not affected: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		totalRows += int(rowsAffected)
	}
	fmt.Println("Total rows inserted for base table: ", totalRows)
}

// Creates the CWE categories table in db and inserts the .csv file entries
func processPar(ent [][]string, dat *sql.DB) {
	len := len(ent)
	fmt.Printf(strconv.Itoa(len))
	//Drop the table so we can then update it
	dat.Exec("DROP TABLE par;")
	//Create the table
	dat.Exec("CREATE TABLE par (CWEname VARCHAR(1024), CWEID INTEGER) ;")
	totalRows := 0
	for i := 1; i < len; i++ {
		result, err := dat.Exec("INSERT INTO par (CWEname, CWEID) VALUES ($1,$2);", ent[i][0], ent[i][1])
		if err != nil {
			log.Fatalf("row not affected: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}
		totalRows += int(rowsAffected)
	}
	fmt.Println("Total rows inserted for par table: ", totalRows)
}

// Removes the cve from table entries
func removeCve(cve string) int {
	cut := strings.Trim(cve, "C")
	cut1 := strings.Trim(cut, "V")
	cut2 := strings.Trim(cut1, "E")
	cut3 := strings.ReplaceAll(cut2, "-", "")
	result, err := strconv.Atoi(cut3)
	if err != nil {
		log.Fatalf("couldn't convert %v", err)
	}
	return result
}
