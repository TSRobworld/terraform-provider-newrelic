//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(testAccExpectedApplicationName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedApplicationName),
				),
			},
		},
	})
}

// This test case checks if an entity with a single quote "'" in its name
// is created, and is subsequently, fetched by the data source.
func TestAccNewRelicSingleQuotedEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccSingleQuotedPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(testAccExpectedSingleQuotedApplicationName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedSingleQuotedApplicationName),
				),
			},
		},
	})
}

func TestAccNewRelicEntityData_Missing(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicEntityDataConfig(strings.ToUpper(testAccExpectedApplicationName), testAccountID),
				ExpectError: regexp.MustCompile(`no entities found`),
			},
		},
	})
}

func TestAccNewRelicEntityData_IgnoreCase(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig_IgnoreCase(strings.ToUpper(testAccExpectedApplicationName), testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedApplicationName),
				),
			},
		},
	})
}

func testAccCheckNewRelicEntityDataExists(t *testing.T, n string, appName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["guid"] == "" {
			return fmt.Errorf("expected to get an entity GUID")
		}

		if a["application_id"] == "" {
			return fmt.Errorf("expected to get an application ID")
		}

		if a["name"] != appName {
			return fmt.Errorf("expected the entity name to be: %s, but got: %s", appName, a["name"])
		}

		return nil
	}
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "application"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_IgnoreCase(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	ignore_case = true
	type = "application"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_InvalidType(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "app"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_InvalidDomain(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "application"
	domain = "VIZ"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}
