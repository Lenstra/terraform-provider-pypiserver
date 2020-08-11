package pypiserver

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"pypiserver": testAccProvider,
	}
}

func TestAccUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pypiserver_user.test", "username", "test"),
				),
			},
			{
				Config: testResourceUser2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pypiserver_user.test", "username", "test"),
					resource.TestCheckResourceAttr("pypiserver_user.test", "email", "foo@bar.com"),
				),
			},
		},
	})
}

func TestAccUser_multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testMultipleUsers,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pypiserver_user.test.0", "username", "user0"),
					resource.TestCheckResourceAttr("pypiserver_user.test.0", "email", "user0@mail.fr"),
					resource.TestCheckResourceAttr("pypiserver_user.test.1", "username", "user1"),
					resource.TestCheckResourceAttr("pypiserver_user.test.1", "email", "user1@mail.fr"),
				),
			},
		},
	})
}

const testResourceUser = `
resource "pypiserver_user" "test" {
	username = "test"
}
`

const testResourceUser2 = `
resource "pypiserver_user" "test" {
	username = "test"
	email = "foo@bar.com"
}
`

const testMultipleUsers = `
resource "pypiserver_user" "test" {
	count = 2

	username = "user${count.index}"
	email = "user${count.index}@mail.fr"
}
`
