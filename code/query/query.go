package query

import (
	"database/sql"
	"log"
	"strings"
)

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

type Score struct {
	Critical int
	High     int
	Medium   int
	Low      int
	NA       int
}

type ScoreEntry struct {
	V3score string
	V2score string
}

type ScoreResult struct {
	V3list Score
	V2list Score
}

type Weakness struct {
	CWEID string
	Count int
	Name  string
}

// Performs a query on the database and returns a list of products.
func QueryVendors(dat *sql.DB) []string {
	rows, err := dat.Query("SELECT DISTINCT vendorProject FROM vulnerabilities ORDER BY vendorProject ASC;")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	vendors := []string{}
	for rows.Next() {
		var vendor string

		if err := rows.Scan(&vendor); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		vendors = append(vendors, vendor)
	}

	return vendors
}

// Performs a query on the database and returns a list of products.
func QueryProducts(dat *sql.DB, vendor string) []string {
	rows, err := dat.Query("SELECT DISTINCT product FROM vulnerabilities WHERE vendorProject = '" + vendor + "' ORDER BY product ASC;")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	products := []string{}
	for rows.Next() {
		var product string

		if err := rows.Scan(&product); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		products = append(products, product)
	}
	return products
}

// Performs a query on the database and returns a list of records given vendor, project
func QueryResults(dat *sql.DB, vendor string, product string) []ProdW {
	rows, err := dat.Query("SELECT * FROM vulnerabilities WHERE vendorProject = '" + vendor + "' AND product = '" + product + "';")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	result := []ProdW{}
	for rows.Next() {
		res := ProdW{}

		if err := rows.Scan(&res.Cvekey, &res.Cveid, &res.VendorProject, &res.Product, &res.VulnerabilityName, &res.DateAdded, &res.ShortDescription, &res.RequiredAction, &res.DueDate, &res.Notes); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		result = append(result, res)
	}

	return result
}

func QueryScores(dat *sql.DB, prodList []ProdW) ScoreResult {
	var cveList []string
	for _, product := range prodList {
		cveList = append(cveList, product.Cveid)
	}
	query := "SELECT v3score, v2score FROM base WHERE "

	for index, cve := range cveList {
		if index == 0 {
			queryStr := "\"cveid\" = '" + cve + "'"
			query = query + queryStr
		} else {
			queryStr := " OR \"cveid\" = '" + cve + "'"
			query = query + queryStr
		}
	}
	query = query + ";"

	//println(query)

	rows, err := dat.Query(query)
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	// This row will return the v3score & the v2score
	result := ScoreResult{}
	V3result := Score{Critical: 0, High: 0, Medium: 0, Low: 0, NA: 0}
	V2result := Score{Critical: 0, High: 0, Medium: 0, Low: 0, NA: 0}

	for rows.Next() {
		res := ScoreEntry{}
		// scan here
		if err := rows.Scan(&res.V3score, &res.V2score); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}

		// filter to add count to appropriate list
		if strings.Contains(res.V3score, "CRITICAL") {
			V3result.Critical += 1
		} else if strings.Contains(res.V3score, "HIGH") {
			V3result.High += 1
		} else if strings.Contains(res.V3score, "MEDIUM") {
			V3result.Medium += 1
		} else if strings.Contains(res.V3score, "LOW") {
			V3result.Low += 1
		} else {
			V3result.NA += 1
		}

		if strings.Contains(res.V2score, "CRITICAL") {
			V2result.Critical += 1
		} else if strings.Contains(res.V2score, "HIGH") {
			V2result.High += 1
		} else if strings.Contains(res.V2score, "MEDIUM") {
			V2result.Medium += 1
		} else if strings.Contains(res.V2score, "LOW") {
			V2result.Low += 1
		} else {
			V2result.NA += 1
		}
	}
	result.V3list = V3result
	result.V2list = V2result
	return result
}

func QueryWeaknesses(dat *sql.DB, prodList []ProdW) []Weakness {
	var cveList []string
	for _, product := range prodList {
		cveList = append(cveList, product.Cveid)
	}
	query := "SELECT source, cwename, COUNT(*) as count FROM cwe WHERE cveid IN ("

	for index, cve := range cveList {
		if index == 0 {
			query += ("'" + cve + "'")
		} else {
			query += (",'" + cve + "'")
		}
	}

	query = query + ") GROUP BY source, cwename;"

	//println(query)

	rows, err := dat.Query(query)
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	result := []Weakness{}

	for rows.Next() {
		res := Weakness{}
		// scan here
		if err := rows.Scan(&res.Name, &res.CWEID, &res.Count); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		result = append(result, res)
		//println("name: ", res.CWEID, "\tCount: ", res.Count, "\n")
	}
	return result
}
