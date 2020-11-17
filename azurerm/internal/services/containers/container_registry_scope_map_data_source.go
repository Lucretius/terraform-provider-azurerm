package containers

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmContainerRegistryScopeMap() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmContainerRegistryRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"container_registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateAzureRMContainerRegistryName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"description": {
				Type: schema.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}


func dataSourceArmContainerRegistrySourceMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Container Registry token %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on token %q (Azure Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("container_registry_name", containerRegistryName)
	d.Set("description", resp.Description)
	d.Set("actions", utils.FlattenStringSlice(resp.Actions))

	return nil
}
