package xoa

import (
	"github.com/aslak11/terraform-provider-xenorchestra/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validActionOptions = []string{
	"admin",
	"operator",
	"viewer",
}

func resourceAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAclCreate,
		Read:   resourceAclRead,
		Delete: resourceAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subject": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validActionOptions, false),
			},
		},
	}
}

func resourceAclCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.XOClient)

	acl, err := c.CreateAcl(client.Acl{
		Subject: d.Get("subject").(string),
		Object:  d.Get("object").(string),
		Action:  d.Get("action").(string),
	})
	if err != nil {
		return err
	}
	return aclToData(acl, d)
}

func resourceAclRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.XOClient)

	acl, err := c.GetAcl(client.Acl{
		Id: d.Id(),
	})

	if _, ok := err.(client.NotFound); ok {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	return aclToData(acl, d)
}

func resourceAclDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.XOClient)

	err := c.DeleteAcl(client.Acl{
		Id: d.Id(),
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func aclToData(acl *client.Acl, d *schema.ResourceData) error {
	d.SetId(acl.Id)
	if err := d.Set("subject", acl.Subject); err != nil {
		return err
	}
	if err := d.Set("object", acl.Object); err != nil {
		return err
	}
	if err := d.Set("action", acl.Action); err != nil {
		return err
	}
	return nil
}
