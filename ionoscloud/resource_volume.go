package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeImporter,
		},
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" {
						return true
					}
					return false
				},
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"image_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssh_key_path": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"sshkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bus": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
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
			"backup_unit_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pci_slot": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"boot_server": {
				Type:        schema.TypeString,
				Description: "The UUID of the attached server.",
				Computed:    true,
			},
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	volume := ionoscloud.Volume{
		Properties: &ionoscloud.VolumeProperties{},
	}

	var sshKeyPath []interface{}
	var publicKeys []string
	var image, imageAlias string
	var isSnapshot bool
	var diags diag.Diagnostics

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	imagePassword := d.Get("image_password").(string)
	sshKeyPath = d.Get("ssh_key_path").([]interface{})
	imageInput := d.Get("image_name").(string)
	licenceType := d.Get("licence_type").(string)

	if len(sshKeyPath) != 0 {
		for _, path := range sshKeyPath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error fetching sshkey from file (%s) (%s)", path, err.Error()))
				return diags
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}

	if imageInput != "" {
		image, imageAlias, isSnapshot, diags = checkImage(ctx, client, imageInput, imagePassword, licenceType, dcId, sshKeyPath)
		if diags != nil {
			return diags
		}
	}

	if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
		diags := diag.FromErr(fmt.Errorf("you can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot"))
		return diags
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		volume.Properties.Name = &name
	}

	if v, ok := d.GetOk("size"); ok {
		size := float32(v.(int))
		volume.Properties.Size = &size
	}

	if v, ok := d.GetOk("disk_type"); ok {
		diskType := v.(string)
		volume.Properties.Type = &diskType
	}

	if v, ok := d.GetOk("bus"); ok {
		bus := v.(string)
		volume.Properties.Bus = &bus
	}

	if _, ok := d.GetOk("availability_zone"); ok {
		raw := d.Get("availability_zone").(string)
		volume.Properties.AvailabilityZone = &raw
	}

	if imagePassword != "" {
		volume.Properties.ImagePassword = &imagePassword
	} else {
		volume.Properties.ImagePassword = nil
	}

	if licenceType != "" {
		volume.Properties.LicenceType = &licenceType
	} else {
		volume.Properties.LicenceType = nil
	}

	if image != "" {
		volume.Properties.Image = &image
	} else {
		volume.Properties.Image = nil
	}

	if imageAlias != "" {
		volume.Properties.ImageAlias = &imageAlias
	} else {
		volume.Properties.ImageAlias = nil
	}

	if len(publicKeys) != 0 {
		volume.Properties.SshKeys = &publicKeys

	} else {
		volume.Properties.SshKeys = nil
	}

	if userData, ok := d.GetOk("user_data"); ok {
		if image == "" && imageAlias == "" {
			diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image that has cloud-init compatibility in conjunction with user_data property "))
			return diags
		} else {
			userData := userData.(string)
			volume.Properties.UserData = &userData
		}
	}

	if backupUnitId, ok := d.GetOk("backup_unit_id"); ok {
		if IsValidUUID(backupUnitId.(string)) {
			if image == "" && imageAlias == "" {
				diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image that has cloud-init compatibility in conjunction with backup_unit_id property "))
				return diags
			} else {
				backupUnitId := backupUnitId.(string)
				volume.Properties.BackupunitId = &backupUnitId
			}
		} else {
			diags := diag.FromErr(fmt.Errorf("the backup_unit_id that you specified is not a valid UUID"))
			return diags
		}
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPost(ctx, dcId).Volume(volume).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a volume: %s", err))
		return diags
	}

	d.SetId(*volume.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	volumeToAttach := ionoscloud.Volume{Id: volume.Id}
	volume, apiResponse, err = client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverId).Volume(volumeToAttach).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverId, *volumeToAttach.Id, err))
		return diags
	}

	sErr := d.Set("server_id", serverId)

	if sErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting serverId %s: %s", serverId, sErr))
		return diags
	}

	// Wait, catching any errors
	_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			sErr := d.Set("server_id", "")
			if sErr != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting serverId: %s", sErr))
				return diags
			}
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcId, volumeID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching volume with ID %s: %s", d.Id(), err))
		return diags
	}

	_, apiResponse, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, serverID, volumeID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if err2 := d.Set("server_id", ""); err2 != nil {
			diags := diag.FromErr(err2)
			return diags
		}
	}

	if err := setVolumeData(d, &volume); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	properties := ionoscloud.VolumeProperties{}
	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		properties.Name = &newValueStr
	}
	if d.HasChange("disk_type") {
		diags := diag.FromErr(fmt.Errorf("disk_type is immutable"))
		return diags
	}
	if d.HasChange("size") {
		_, newValue := d.GetChange("size")
		newValueFloat32 := float32(newValue.(int))
		properties.Size = &newValueFloat32
	}
	if d.HasChange("bus") {
		_, newValue := d.GetChange("bus")
		newValueStr := newValue.(string)
		properties.Bus = &newValueStr
	}
	if d.HasChange("availability_zone") {
		diags := diag.FromErr(fmt.Errorf("availability_zone is immutable"))
		return diags
	}

	if d.HasChange("user_data") {
		diags := diag.FromErr(fmt.Errorf("user_data property of resource volume is immutable "))
		return diags
	}

	if d.HasChange("backup_unit_id") {
		diags := diag.FromErr(fmt.Errorf("backup_unit_id property of resource volume is immutable "))
		return diags
	}

	if d.HasChange("image_name") {
		diags := diag.FromErr(fmt.Errorf("backup_unit_id property of resource volume is immutable "))
		return diags
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, dcId, d.Id()).Volume(properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating volume with ID %s: %s", d.Id(), err))
		return diags

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode > 299 {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a volume ID %s %s", d.Id(), err))
		return diags
	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		volumeToAttach := ionoscloud.Volume{Id: volume.Id}
		_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverID).Volume(volumeToAttach).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s",
				dcId, serverID, *volume.Id, err))
			return diags
		}

		// Wait, catching any errors
		_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a volume ID %s %s", d.Id(), err))
		return diags

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil
}

func resourceVolumeImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{volume}", d.Id())
	}

	client := meta.(SdkBundle).CloudApiClient

	dcId := parts[0]
	srvId := parts[1]
	volumeId := parts[2]

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcId, volumeId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("volume does not exist %q", volumeId)
		}
		return nil, fmt.Errorf("an error occured while trying to find the volume %q", volumeId)
	}

	log.Printf("[INFO] volume found: %+v", volume)

	d.SetId(*volume.Id)
	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if err := d.Set("server_id", srvId); err != nil {
		return nil, err
	}

	if err := setVolumeData(d, &volume); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setVolumeData(d *schema.ResourceData, volume *ionoscloud.Volume) error {

	if volume.Id != nil {
		d.SetId(*volume.Id)
	}

	if volume.Properties.Name != nil {
		err := d.Set("name", *volume.Properties.Name)
		if err != nil {
			return fmt.Errorf("error while setting name property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Type != nil {
		err := d.Set("disk_type", *volume.Properties.Type)
		if err != nil {
			return fmt.Errorf("error while setting type property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Size != nil {
		err := d.Set("size", *volume.Properties.Size)
		if err != nil {
			return fmt.Errorf("error while setting size property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Bus != nil {
		err := d.Set("bus", *volume.Properties.Bus)
		if err != nil {
			return fmt.Errorf("error while setting bus property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Image != nil {
		err := d.Set("image", *volume.Properties.Image)
		if err != nil {
			return fmt.Errorf("error while setting image property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.ImageAlias != nil {
		err := d.Set("image_alias", *volume.Properties.ImageAlias)
		if err != nil {
			return fmt.Errorf("error while setting image_alias property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.AvailabilityZone != nil {
		err := d.Set("availability_zone", *volume.Properties.AvailabilityZone)
		if err != nil {
			return fmt.Errorf("error while setting availability_zone property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.CpuHotPlug != nil {
		err := d.Set("cpu_hot_plug", *volume.Properties.CpuHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting cpu_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.RamHotPlug != nil {
		err := d.Set("ram_hot_plug", *volume.Properties.RamHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting ram_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.NicHotPlug != nil {
		err := d.Set("nic_hot_plug", *volume.Properties.NicHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting nic_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.NicHotUnplug != nil {
		err := d.Set("nic_hot_unplug", *volume.Properties.NicHotUnplug)
		if err != nil {
			return fmt.Errorf("error while setting nic_hot_unplug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotPlug != nil {
		err := d.Set("disc_virtio_hot_plug", *volume.Properties.DiscVirtioHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting disc_virtio_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotUnplug != nil {
		err := d.Set("disc_virtio_hot_unplug", *volume.Properties.DiscVirtioHotUnplug)
		if err != nil {
			return fmt.Errorf("error while setting disc_virtio_hot_unplug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.BackupunitId != nil {
		err := d.Set("backup_unit_id", *volume.Properties.BackupunitId)
		if err != nil {
			return fmt.Errorf("error while setting backup_unit_id property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.UserData != nil {
		err := d.Set("user_data", *volume.Properties.UserData)
		if err != nil {
			return fmt.Errorf("error while setting user_data property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.DeviceNumber != nil {
		err := d.Set("device_number", *volume.Properties.DeviceNumber)
		if err != nil {
			return fmt.Errorf("error while setting device_number property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.BootServer != nil {
		err := d.Set("boot_server", *volume.Properties.BootServer)
		if err != nil {
			return fmt.Errorf("error while setting boot_server property for volume %s: %s", d.Id(), err)
		}
	}
	return nil
}

func resolveImageName(ctx context.Context, client *ionoscloud.APIClient, imageName string, location string) (*ionoscloud.Image, error) {

	if imageName == "" {
		return nil, fmt.Errorf("imageName not suplied")
	}

	images, apiResponse, err := client.ImagesApi.ImagesGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of images %s", err))
		return nil, err
	}

	if len(*images.Items) > 0 {
		for _, i := range *images.Items {
			imgName := ""
			if i.Properties != nil && i.Properties.Name != nil && *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && *i.Properties.ImageType == "HDD" && *i.Properties.Location == location {
				return &i, err
			}

			if imgName != "" && strings.ToLower(imageName) == strings.ToLower(*i.Id) && *i.Properties.ImageType == "HDD" && *i.Properties.Location == location {
				return &i, err
			}

		}
	}
	return nil, err
}

func getSnapshotId(ctx context.Context, client *ionoscloud.APIClient, snapshotName string) string {

	if snapshotName == "" {
		return ""
	}

	snapshots, apiResponse, err := client.SnapshotsApi.SnapshotsGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of snapshots %s", err))
	}

	if len(*snapshots.Items) > 0 {
		for _, i := range *snapshots.Items {
			imgName := ""
			if i.Properties != nil && i.Properties.Name != nil && *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return *i.Id
			}
		}
	}
	return ""
}

func getImageAlias(ctx context.Context, client *ionoscloud.APIClient, imageAlias string, location string) string {

	if imageAlias == "" {
		return ""
	}
	parts := strings.SplitN(location, "/", 2)
	if len(parts) != 2 {
		log.Print(fmt.Errorf("invalid location id %s", location))
	}

	locations, apiResponse, err := client.LocationsApi.LocationsFindByRegionIdAndId(ctx, parts[0], parts[1]).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of locations %s", err))
	}

	if len(*locations.Properties.ImageAliases) > 0 {
		for _, i := range *locations.Properties.ImageAliases {
			alias := ""
			if i != "" {
				alias = i
			}

			if alias != "" && strings.ToLower(alias) == strings.ToLower(imageAlias) {
				return i
			}
		}
	}
	return ""
}

func checkImage(ctx context.Context, client *ionoscloud.APIClient, imageInput, imagePassword, licenceType, dcId string, sshKeyPath []interface{}) (image, imageAlias string, isSnapshot bool, diags diag.Diagnostics) {
	isSnapshot = false

	if imageInput != "" {
		if !IsValidUUID(imageInput) {
			dc, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error fetching datacenter %s: (%s)", dcId, err))
				return image, imageAlias, isSnapshot, diags
			}
			img, err := resolveImageName(ctx, client, imageInput, *dc.Properties.Location)
			if err != nil {
				diags := diag.FromErr(err)
				return image, imageAlias, isSnapshot, diags
			}
			if img != nil {
				image = *img.Id
			}
			// if no image id was found with that name we look for a matching snapshot
			if image == "" {
				image = getSnapshotId(ctx, client, imageInput)
				if image != "" {
					isSnapshot = true
				} else {
					imageAlias = getImageAlias(ctx, client, imageInput, *dc.Properties.Location)
				}
			}

			if image == "" && imageAlias == "" {
				diags := diag.FromErr(fmt.Errorf("could not find an image/imagealias/snapshot that matches %s ", imageInput))
				return image, imageAlias, isSnapshot, diags
			}

			if imagePassword == "" && len(sshKeyPath) == 0 && isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public {
				diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided"))
				return image, imageAlias, isSnapshot, diags
			}

		} else {
			img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageInput).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
					snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, imageInput).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
							diags := diag.FromErr(fmt.Errorf("image/snapshot %s not found: %s", imageInput, err))
							return image, imageAlias, isSnapshot, diags
						} else {
							diags := diag.FromErr(fmt.Errorf("an error occured while fetching snapshot %s: %s", imageInput, err))
							return image, imageAlias, isSnapshot, diags
						}
					}
					image = *snapshot.Id
					isSnapshot = true
				} else {
					diags := diag.FromErr(fmt.Errorf("error fetching image %s: %s", imageInput, err))
					return image, imageAlias, isSnapshot, diags
				}
			} else {
				if isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public == true {
					if imagePassword == "" && len(sshKeyPath) == 0 {
						diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'sshkey' must be provided"))
						return image, imageAlias, isSnapshot, diags
					}
				}
				image = *img.Id
			}
		}
	}

	if imageInput == "" && licenceType == "" && isSnapshot == false {
		diags := diag.FromErr(fmt.Errorf("either 'image_name', or 'licence_type' must be set"))
		return image, imageAlias, isSnapshot, diags
	}

	return image, imageAlias, isSnapshot, diags
}
