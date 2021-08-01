package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var originalDnsName string = "terraform"

// var updatedDnsName string = "terraform.updated"
var originalAddress string = "10.255.255.1"
var updatedAddress string = "10.0.0.1"

func TestAccMikrotikDnsRecord_create(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_createAndPlanWithNonExistantRecord(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	removeDnsRecord := func() {

		c := client.NewClient(client.GetConfigFromEnv())
		dns, err := c.FindDnsRecord(originalDnsName)

		if err != nil {
			t.Fatalf("Error finding the DNS record: %s", err)
		}
		err = c.DeleteDnsRecord(dns.Id)
		if err != nil {
			t.Fatalf("Error removing the DNS record: %s", err)
		}

	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				PreConfig:          removeDnsRecord,
				Config:             testAccDnsRecord(),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMikrotikDnsRecord_updateAddress(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", originalAddress),
				),
			},
			{
				Config: testAccDnsRecordUpdatedAddress(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", updatedAddress)),
			},
		},
	})
}

func TestAccMikrotikDnsRecord_import(t *testing.T) {
	resourceName := "mikrotik_dns_record.bar"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMikrotikDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDnsRecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsRecord() string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
}
`, originalDnsName, originalAddress)
}

func testAccDnsRecordUpdatedAddress() string {
	return fmt.Sprintf(`
resource "mikrotik_dns_record" "bar" {
    name = "%s"
    address = "%s"
    ttl = "300"
}
`, originalDnsName, updatedAddress)
}

func testAccDnsRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mikrotik_dns_record does not exist in the statefile")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Unable to get the dns record with error: %v", err)
		}

		if dnsRecord == nil {
			return fmt.Errorf("Unable to get the dns record with name: %s", dnsRecord.Name)
		}

		if dnsRecord.Name == rs.Primary.ID {
			return nil
		}
		return nil
	}
}

func testAccCheckMikrotikDnsRecordDestroyNow(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dns record Id is set")
		}

		c := client.NewClient(client.GetConfigFromEnv())

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		_, ok = err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}
		err = c.DeleteDnsRecord(dnsRecord.Id)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckMikrotikDnsRecordDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_dns_record" {
			continue
		}

		dnsRecord, err := c.FindDnsRecord(rs.Primary.ID)

		_, ok := err.(*client.NotFound)
		if !ok && err != nil {
			return err
		}

		if dnsRecord != nil {
			return fmt.Errorf("dns recrod (%s) still exists", dnsRecord.Id)
		}
	}
	return nil
}
