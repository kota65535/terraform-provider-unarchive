package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUnarchiveFile() *schema.Resource {
	return &schema.Resource{
		Description: "Extract an archive and then enumerate the files.",
		ReadContext: dataSourceUnarchiveFileRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the archive file. NOTE: `zip` is supported",
			},
			"source_file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path of the archive file.",
			},
			"pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Glob pattern to filter files to extract.",
			},
			"output_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the directory where files are extracted.",
				Default:     ".",
			},
			"output_files": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Paths of the extracted files.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceUnarchiveFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	type_ := d.Get("type").(string)
	sourceFile := d.Get("source_file").(string)
	pattern := d.Get("pattern").(string)
	outputDir := d.Get("output_dir").(string)

	if type_ != "zip" {
		diag.Errorf("type not supported")
	}

	extracted, err := UnzipSource(sourceFile, pattern, outputDir)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("output_files", extracted)
	if err != nil {
		return diag.FromErr(err)
	}

	// Calculate hashes
	sha1, _, _, err := GenerateHashes(sourceFile)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(sha1)

	return nil
}
