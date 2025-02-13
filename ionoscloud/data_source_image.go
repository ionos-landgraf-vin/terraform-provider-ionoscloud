package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
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
				Computed: true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image_aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_init": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	images, apiResponse, err := client.ImagesApi.ImagesGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud images %s", err))
	}

	name, nameOk := d.GetOk("name")
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")
	cloudInit, cloudInitOk := d.GetOk("cloud_init")

	var results []ionoscloud.Image

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk && nameOk {
		nameVer := fmt.Sprintf("%s-%s", name.(string), version.(string))
		if images.Items != nil {
			for _, img := range *images.Items {
				if img.Properties != nil && img.Properties.Name != nil && strings.ToLower(*img.Properties.Name) == strings.ToLower(nameVer) {
					results = append(results, img)
				}
			}
		}
		if results == nil {
			return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name %s and version %s (%s)", name.(string), version.(string), nameVer))
		}
	} else if nameOk {
		if images.Items != nil {
			for _, img := range *images.Items {
				if img.Properties != nil && img.Properties.Name != nil && strings.ToLower(*img.Properties.Name) == strings.ToLower(name.(string)) {
					results = append(results, img)
					break
				}
			}
		}
		if results == nil {
			return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name %s", name.(string)))
		}
	}

	if imageTypeOk {
		var imageTypeResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.ImageType != nil && *img.Properties.ImageType == imageType.(string) {
				imageTypeResults = append(imageTypeResults, img)
			}

		}
		results = imageTypeResults
	}

	if locationOk {
		var locationResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.Location != nil && *img.Properties.Location == location.(string) {
				locationResults = append(locationResults, img)
			}
		}
		results = locationResults
	}

	if cloudInitOk {
		var cloudInitResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.CloudInit != nil && *img.Properties.CloudInit == cloudInit.(string) {
				cloudInitResults = append(cloudInitResults, img)
			}
		}
		results = cloudInitResults
	}

	var image ionoscloud.Image

	if results == nil || len(results) == 0 {
		return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name = %s, type = %s, location = %s, version = %s, cloudInit = %s", name.(string), imageType.(string), location.(string), version.(string), cloudInit.(string)))
	} else if len(results) > 1 {
		return diag.FromErr(fmt.Errorf("more than one cluster found with the specified criteria: name = %s, type = %s, location = %s, version = %s, cloudInit = %s", name.(string), imageType.(string), location.(string), version.(string), cloudInit.(string)))
	} else {
		image = results[0]
	}

	if err := ImageSetData(d, &image); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ImageSetData(d *schema.ResourceData, image *ionoscloud.Image) error {
	if image.Id != nil {
		d.SetId(*image.Id)
	}

	if image.Properties != nil {

		if image.Properties.Name != nil {
			err := d.Set("name", *image.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for image %s: %s", d.Id(), err)
			}
		}

		if image.Properties.Description != nil {
			if err := d.Set("description", *image.Properties.Description); err != nil {
				return err
			}
		}

		if image.Properties.Size != nil {
			if err := d.Set("size", *image.Properties.Size); err != nil {
				return err
			}
		}

		if image.Properties.CpuHotPlug != nil {
			if err := d.Set("cpu_hot_plug", *image.Properties.CpuHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.CpuHotUnplug != nil {
			if err := d.Set("cpu_hot_unplug", *image.Properties.CpuHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.RamHotPlug != nil {
			if err := d.Set("ram_hot_plug", *image.Properties.RamHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.RamHotUnplug != nil {
			if err := d.Set("ram_hot_unplug", *image.Properties.RamHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.NicHotPlug != nil {
			if err := d.Set("nic_hot_plug", *image.Properties.NicHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.NicHotUnplug != nil {
			if err := d.Set("nic_hot_unplug", *image.Properties.NicHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.DiscVirtioHotPlug != nil {
			if err := d.Set("disc_virtio_hot_plug", *image.Properties.DiscVirtioHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.DiscVirtioHotUnplug != nil {
			if err := d.Set("disc_virtio_hot_unplug", *image.Properties.DiscVirtioHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.DiscScsiHotPlug != nil {
			if err := d.Set("disc_scsi_hot_plug", *image.Properties.DiscScsiHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.DiscScsiHotUnplug != nil {
			if err := d.Set("disc_scsi_hot_unplug", *image.Properties.DiscScsiHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.LicenceType != nil {
			if err := d.Set("license_type", *image.Properties.LicenceType); err != nil {
				return err
			}
		}

		if image.Properties.Public != nil {
			if err := d.Set("public", *image.Properties.Public); err != nil {
				return err
			}
		}

		if image.Properties.ImageAliases != nil && len(*image.Properties.ImageAliases) > 0 {
			if err := d.Set("image_aliases", *image.Properties.ImageAliases); err != nil {
				return err
			}
		}

		if image.Properties.CloudInit != nil {
			if err := d.Set("cloud_init", *image.Properties.CloudInit); err != nil {
				return err
			}
		}

		if image.Properties.ImageType != nil {
			err := d.Set("type", *image.Properties.ImageType)
			if err != nil {
				return err
			}
		}

		if image.Properties.Location != nil {
			err := d.Set("location", *image.Properties.Location)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
