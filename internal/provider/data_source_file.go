package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"path"
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
			"patterns": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Glob patterns to filter files to extract. Defaults to `[\"**\"]` (all files included).",
			},
			"excludes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Glob patterns to exclude files to extract. Defaults to `[]` (no file excluded).",
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
				Description: "Extracted files.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the file.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of the file.",
						},
					},
				},
			},
		},
	}
}

func dataSourceUnarchiveFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	type_ := d.Get("type").(string)
	sourceFile := d.Get("source_file").(string)
	patterns := toStringSlice(d.Get("patterns").([]interface{}))
	excludes := toStringSlice(d.Get("excludes").([]interface{}))
	outputDir := d.Get("output_dir").(string)

	if type_ != "zip" {
		diag.Errorf("type not supported")
	}

	if len(patterns) == 0 {
		patterns = []string{"**"}
	}

	fileNames, err := UnzipSource(sourceFile, patterns, excludes, outputDir)
	if err != nil {
		return diag.FromErr(err)
	}
	outputFiles := []map[string]string{}
	for _, f := range fileNames {
		p := path.Join(outputDir, f)
		if err != nil {
			return diag.FromErr(err)
		}
		outputFiles = append(outputFiles, map[string]string{
			"name": f,
			"path": p,
		})
	}
	err = d.Set("output_files", outputFiles)

	// Calculate hashes
	sha1, err := GenerateHash(sourceFile)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(sha1)

	return nil
}

func toStringSlice(values []interface{}) []string {
	ret := []string{}
	for _, v := range values {
		ret = append(ret, v.(string))
	}
	return ret
}
