// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVersionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVersionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.unraid_version.test", "unraid"),
					resource.TestCheckResourceAttrSet("data.unraid_version.test", "api"),
					resource.TestCheckResourceAttrSet("data.unraid_version.test", "kernel"),
				),
			},
		},
	})
}

const testAccVersionDataSourceConfig = `
data "unraid_version" "test" {}
`
