package stix2

import (
	"encoding/json"
	"fmt"

	"github.com/TcM1911/stix2"
	"gitlab.ssc.dev/Research/tia/code/query"
)

//runs on the api call get Vulns and returns the json to populate the STIX graph
func Stixify(vulns []query.ProdW, ven string, pro string) string {
	//creates the collection that we will add to our bundle so like an array to store our STIX objects
	collection := stix2.New()

	//creates vendor object and names it the vendor
	vendor, err := stix2.NewIdentity(ven, stix2.OptionClass("organization"))
	if err != nil {
		fmt.Println(err)
	}
	collection.Add(vendor)

	//creates infrastructure object and names it the vendor's product
	infra, err := stix2.NewInfrastructure(pro)
	if err != nil {
		fmt.Println(err)
	}
	collection.Add(infra)

	//creates belongs-to relationship for the vendor and product
	ref, err := stix2.NewRelationship("owns", vendor.ID, infra.ID)
	if err != nil {
		fmt.Println(err)
	}
	collection.Add(ref)

	//for loop to run through all the vulnerabilites make them into objects and then link back to the vendor's product
	listlen := len(vulns)
	for i := 0; i < listlen; i++ {
		vuln := vulns[i].VulnerabilityName
		vu, err := stix2.NewVulnerability(vuln, stix2.OptionDescription(vulns[i].ShortDescription))
		if err != nil {
			fmt.Println(err)
		}
		collection.Add(vu)

		ref, err := infra.AddHas(vu.ID)
		if err != nil {
			fmt.Println(err)
		}
		collection.Add(ref)
	}

	//adds collection to the Bundle and then turns it into JSON data
	b, err := collection.ToBundle()
	if err != nil {
		fmt.Println(err)
	}
	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("\nvulnerabilities: ")
	// fmt.Println(vulns)
	// fmt.Println("vendor: " + ven)
	// fmt.Println("product: " + pro)
	// fmt.Println(string(data))

	return string(data)
}
