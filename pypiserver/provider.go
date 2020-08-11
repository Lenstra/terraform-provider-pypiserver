package pypiserver

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PYPISERVER_ADDRESS", nil),
			},

			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PYPISERVER_TOKEN", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"pypiserver_user": resourceUser(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := map[string]string{
		"address": d.Get("address").(string),
		"token":   d.Get("token").(string),
	}

	return config, nil
}

func getClient(meta interface{}) *Client {
	config := meta.(map[string]string)
	return NewClient(config["address"], config["token"])
}
