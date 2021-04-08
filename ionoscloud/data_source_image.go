package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	images, _, err := client.ImageApi.ImagesGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud images %s", err)
	}

	name := d.Get("name").(string)
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")

	results := []ionoscloud.Image{}

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk {
		name_ver := fmt.Sprintf("%s-%s", name, version.(string))
		if images.Items != nil {
			for _, img := range *images.Items {
				if strings.Contains(strings.ToLower(*img.Properties.Name), strings.ToLower(name_ver)) {
					results = append(results, img)
				}
			}
		}
	} else {
		if images.Items != nil {
			for _, img := range *images.Items {
				if strings.Contains(strings.ToLower(*img.Properties.Name), strings.ToLower(name)) {
					results = append(results, img)
				}
			}
		}
	}

	if imageTypeOk {
		imageTypeResults := []ionoscloud.Image{}
		for _, img := range results {
			if img.Properties.ImageType != nil && *img.Properties.ImageType == imageType.(string) {
				imageTypeResults = append(imageTypeResults, img)
			}

		}
		results = imageTypeResults
	}

	if locationOk {
		locationResults := []ionoscloud.Image{}
		for _, img := range results {
			if img.Properties.Location != nil && *img.Properties.Location == location.(string) {
				locationResults = append(locationResults, img)
			}
		}
		results = locationResults
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one image that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no images that match the search criteria")
	}

	if results[0].Properties.Name != nil {
		err := d.Set("name", *results[0].Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for image %s: %s", d.Id(), err)
		}
	}

	if results[0].Properties.Description != nil {
		if err := d.Set("description", *results[0].Properties.Description); err != nil {
			return err
		}
	}

	if results[0].Properties.Size != nil {
		if err := d.Set("size", *results[0].Properties.Size); err != nil {
			return err
		}
	}

	if results[0].Properties.CpuHotPlug != nil {
		if err := d.Set("cpu_hot_plug", *results[0].Properties.CpuHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.CpuHotUnplug != nil {
		if err := d.Set("cpu_hot_unplug", *results[0].Properties.CpuHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.RamHotPlug != nil {
		if err := d.Set("ram_hot_plug", *results[0].Properties.RamHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.RamHotUnplug != nil {
		if err := d.Set("ram_hot_unplug", *results[0].Properties.RamHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.NicHotPlug != nil {
		if err := d.Set("nic_hot_plug", *results[0].Properties.NicHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.NicHotUnplug != nil {
		if err := d.Set("nic_hot_unplug", *results[0].Properties.NicHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscVirtioHotPlug != nil {
		if err := d.Set("disc_virtio_hot_plug", *results[0].Properties.DiscVirtioHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscVirtioHotUnplug != nil {
		if err := d.Set("disc_virtio_hot_unplug", *results[0].Properties.DiscVirtioHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscScsiHotPlug != nil {
		if err := d.Set("disc_scsi_hot_plug", *results[0].Properties.DiscScsiHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscScsiHotUnplug != nil {
		if err := d.Set("disc_scsi_hot_unplug", *results[0].Properties.DiscScsiHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.LicenceType != nil {
		if err := d.Set("license_type", *results[0].Properties.LicenceType); err != nil {
			return err
		}
	}

	if results[0].Properties.Public != nil {
		if err := d.Set("public", *results[0].Properties.Public); err != nil {
			return err
		}
	}

	d.SetId(*results[0].Id)

	return nil
}
