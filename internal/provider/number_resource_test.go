// Copyright (c) Arthur Cesaré-Herriau
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sequence": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccNumberResource_Defaults(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "sequence_number" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "1"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "001"),
					resource.TestCheckResourceAttr("sequence_number.test", "id", "001"),
					resource.TestCheckResourceAttr("sequence_number.test", "width", "3"),
					resource.TestCheckResourceAttr("sequence_number.test", "start", "1"),
				),
			},
		},
	})
}

func TestAccNumberResource_CustomWidthAndStart(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sequence_number" "test" {
  start = 42
  width = 5
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "42"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "00042"),
				),
			},
		},
	})
}

func TestAccNumberResource_PrefixSuffix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sequence_number" "test" {
  prefix = "vm-"
  suffix = "-prod"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "vm-001-prod"),
					resource.TestCheckResourceAttr("sequence_number.test", "id", "vm-001-prod"),
				),
			},
		},
	})
}

func TestAccNumberResource_KeepersIncrement(t *testing.T) {
	cfg := func(region string) string {
		return fmt.Sprintf(`
resource "sequence_number" "test" {
  prefix = "vm-"
  keepers = {
    region = %q
  }
}`, region)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg("eu-west-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "1"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "vm-001"),
				),
			},
			{
				Config: cfg("us-east-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "2"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "vm-002"),
				),
			},
			{
				Config: cfg("ap-south-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "3"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "vm-003"),
				),
			},
		},
	})
}

func TestAccNumberResource_NegativeWidth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sequence_number" "test" {
  width = -1
}`,
				ExpectError: regexp.MustCompile(`must be >= 0`),
			},
		},
	})
}

func TestAccNumberResource_WidthChangeReformats(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "sequence_number" "test" { width = 3 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "001"),
				),
			},
			{
				Config: `resource "sequence_number" "test" { width = 6 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sequence_number.test", "number", "1"),
					resource.TestCheckResourceAttr("sequence_number.test", "formatted", "000001"),
				),
			},
		},
	})
}
