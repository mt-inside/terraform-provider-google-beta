// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("ComputeInstanceGroupNamedPort", &resource.Sweeper{
		Name: "ComputeInstanceGroupNamedPort",
		F:    testSweepComputeInstanceGroupNamedPort,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepComputeInstanceGroupNamedPort(region string) error {
	resourceName := "ComputeInstanceGroupNamedPort"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	// Setup variables to replace in list template
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":  config.Project,
			"region":   region,
			"location": region,
			"zone":     "-",
		},
	}

	listTemplate := strings.Split("https://www.googleapis.com/compute/beta/projects/{{project}}/aggregated/instanceGroups/{{group}}", "?")[0]
	listUrl, err := replaceVars(d, config, listTemplate)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error preparing sweeper list url: %s", err)
		return nil
	}

	res, err := sendRequest(config, "GET", config.Project, listUrl, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", listUrl, err)
		return nil
	}

	resourceList, ok := res["namedPorts"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}
	var rl []interface{}
	zones := resourceList.(map[string]interface{})
	// Loop through every zone in the list response
	for _, zonesValue := range zones {
		zone := zonesValue.(map[string]interface{})
		for k, v := range zone {
			// Zone map either has resources or a warning stating there were no resources found in the zone
			if k != "warning" {
				resourcesInZone := v.([]interface{})
				rl = append(rl, resourcesInZone...)
			}
		}
	}

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// items who don't match the tf-test prefix
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["name"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource name was nil", resourceName)
			return nil
		}

		name := GetResourceNameFromSelfLink(obj["name"].(string))
		// Only sweep resources with the test prefix
		if !strings.HasPrefix(name, "tf-test") {
			nonPrefixCount++
			continue
		}

		deleteTemplate := "https://www.googleapis.com/compute/beta/projects/{{project}}/zones/{{zone}}/instanceGroups/{{group}}/setNamedPorts"
		if obj["zone"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource zone was nil", resourceName)
			return nil
		}
		zone := GetResourceNameFromSelfLink(obj["zone"].(string))
		deleteTemplate = strings.Replace(deleteTemplate, "{{zone}}", zone, -1)

		deleteUrl, err := replaceVars(d, config, deleteTemplate)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error preparing delete url: %s", err)
			return nil
		}
		deleteUrl = deleteUrl + name

		// Don't wait on operations as we may have a lot to delete
		_, err = sendRequest(config, "DELETE", config.Project, deleteUrl, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, name)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf_test prefix remain.", nonPrefixCount)
	}

	return nil
}
