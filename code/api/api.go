package api

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"gitlab.ssc.dev/Research/tia/code/db"
	"gitlab.ssc.dev/Research/tia/code/query"
	"gitlab.ssc.dev/Research/tia/code/stix2"

	"log"
	"os/user"
	"strings"
)

var conn = db.Main()

// TEST
func TestViz(c *gin.Context) {
	stix := `{
		"type": "bundle",
		"id": "bundle--601cee35-6b16-4e68-a3e7-9ec7d755b4c3",
		"objects": [
			{
				"type": "threat-actor",
				"spec_version": "2.1",
				"id": "threat-actor--dfaa8d77-07e2-4e28-b2c8-92e9f7b04428",
				"created": "2014-11-19T23:39:03.893Z",
				"modified": "2014-11-19T23:39:03.893Z",
				"name": "Disco Team Threat Actor Group",
				"description": "This organized threat actor group operates to create profit from all types of crime.",
				"threat_actor_types": [
					"crime-syndicate"
				],
				"aliases": [
					"Equipo del Discoteca"
				],
				"roles": [
					"agent"
				],
				"goals": [
					"Steal Credit Card Information"
				],
				"sophistication": "expert",
				"resource_level": "organization",
				"primary_motivation": "personal-gain"
			},
			{
				"type": "identity",
				"spec_version": "2.1",
				"id": "identity--733c5838-34d9-4fbf-949c-62aba761184c",
				"created": "2016-08-23T18:05:49.307Z",
				"modified": "2016-08-23T18:05:49.307Z",
				"name": "Disco Team",
				"description": "Disco Team is the name of an organized threat actor crime-syndicate.",
				"identity_class": "organization",
				"contact_information": "disco-team@stealthemail.com"
			},
			{
				"type": "relationship",
				"spec_version": "2.1",
				"id": "relationship--a2e3efb5-351d-4d46-97a0-6897ee7c77a0",
				"created": "2020-02-29T18:01:28.577Z",
				"modified": "2020-02-29T18:01:28.577Z",
				"relationship_type": "attributed-to",
				"source_ref": "threat-actor--dfaa8d77-07e2-4e28-b2c8-92e9f7b04428",
				"target_ref": "identity--733c5838-34d9-4fbf-949c-62aba761184c"}]}`
	c.HTML(
		http.StatusOK,
		"homepage/viz.html",
		gin.H{
			"stix": stix,
		},
	)
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping:": "pong",
	})
}

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Home",
	})
}

func Home(c *gin.Context) {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	usernamesys := currentUser.Username
	username := usernamesys[strings.LastIndex(usernamesys, "\\")+1:]
	c.HTML(
		http.StatusOK,
		"homepage/index.html",
		gin.H{
			"Title": username,
		},
	)
}

func Graph(c *gin.Context) {
	vendors := query.QueryVendors(conn)
	c.HTML(
		http.StatusOK,
		"homepage/graph.html",
		gin.H{
			"vendorProject": vendors,
		},
	)
}

func GetProducts(c *gin.Context) {
	vendor := c.Param("vendor")
	products := query.QueryProducts(conn, vendor)
	c.HTML(
		http.StatusOK,
		"homepage/prodSelect.html",
		gin.H{
			"vendor":  vendor,
			"product": products,
		},
	)
}

func GetVulns(c *gin.Context) {
	vendor := c.Param("vendor")
	product := c.Param("product")
	vulnerabilities := query.QueryResults(conn, vendor, product)
	c.HTML(
		http.StatusOK,
		"homepage/vulnTable.html",
		gin.H{
			"vendorID":        vendor,
			"productID":       product,
			"vulnerabilities": vulnerabilities,
		},
	)
}

func GetStix(c *gin.Context) {
	vendor := c.Param("vendor")
	product := c.Param("product")
	vulnerabilities := query.QueryResults(conn, vendor, product)
	stix := stix2.Stixify(vulnerabilities, vendor, product)

	//Get current file in resources-tmp
	oldFiles, _ := ioutil.ReadDir("../website/assets/resources-tmp/")
	//Create new filename
	filename := createFilename(vendor, product)
	//Rename old file in resources-tmp to filename
	err := os.Rename("../website/assets/resources-tmp/"+oldFiles[0].Name(), filename)
	if err != nil {
		log.Fatal(err)
	}
	//Write STIX JSON string to filename
	err = ioutil.WriteFile("../website/assets/resources-tmp/"+filename, []byte(stix), 0644)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"stix":     stix,
			"filename": filename,
		},
	)
}

func GetScores(c *gin.Context) {
	vendor := c.Param("vendor")
	product := c.Param("product")
	vulnerabilities := query.QueryResults(conn, vendor, product)
	scores := query.QueryScores(conn, vulnerabilities)
	//fmt.Println(scores)
	c.JSON(
		http.StatusOK,
		gin.H{
			"cvss2": scores.V2list,
			"cvss3": scores.V3list,
		},
	)
}

func GetCwes(c *gin.Context) {
	vendor := c.Param("vendor")
	product := c.Param("product")
	vulnerabilities := query.QueryResults(conn, vendor, product)
	weakness := query.QueryWeaknesses(conn, vulnerabilities)
	//fmt.Println(weakness)
	c.JSON(
		http.StatusOK,
		gin.H{
			"cwes": weakness,
		},
	)
}

func About(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"homepage/about.html",
		gin.H{},
	)
}

func Help(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"homepage/help.html",
		gin.H{},
	)
}

func FOF(c *gin.Context) {
	c.HTML(
		http.StatusNotFound,
		"homepage/404.html",
		gin.H{},
	)
}

func createFilename(vendor, product string) string {
	char, _ := regexp.Compile("[a-zA-Z0-9\\s()-,]")
	space, _ := regexp.Compile("[\\s]")
	productArr := char.FindAllString(product, -1)
	vendorArr := char.FindAllString(vendor, -1)
	vendorString := space.ReplaceAllString(strings.Join(vendorArr, ""), "_")
	productString := space.ReplaceAllString(strings.Join(productArr, ""), "_")
	filename := (vendorString + "__" + productString + ".json")
	return filename
}
