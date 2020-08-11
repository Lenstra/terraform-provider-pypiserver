package pypiserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	c := getClient(m)
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	user := &User{
		Username: username,
		Email:    email,
	}
	err := c.CreateUser(user)
	if err != nil {
		return fmt.Errorf("Failed to create '%s': %v", username, err)
	}

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := getClient(m)
	username := d.Get("username").(string)
	user, err := c.User(username)
	if err != nil {
		if strings.Contains(err.Error(), "(404)") {
			log.Printf("[INFO] Removing '%s'", username)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to fetch '%s': %v", username, err)
	}

	log.Printf("[INFO] Saving email for '%s': '%s'", user.Username, user.Email)
	if err = d.Set("email", user.Email); err != nil {
		return fmt.Errorf("Failed to set 'email': %v", err)
	}

	d.SetId(user.Username)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	c := getClient(m)
	oldUsername := d.Id()
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	user := &User{
		Username: username,
		Email:    email,
	}
	err := c.UpdateUser(oldUsername, user)
	if err != nil {
		return fmt.Errorf("Failed to update user '%s': %v", oldUsername, err)
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	c := getClient(m)
	username := d.Get("username").(string)
	return c.DeleteUser(username)
}
