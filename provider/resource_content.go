package provider

import (
	"context"
	"fmt"
	pc "polycode-provider/client"
	"polycode-provider/client/models/content"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentCreate,
		ReadContext:   resourceContentRead,
		UpdateContext: resourceContentUpdate,
		DeleteContext: resourceContentDelete,
		Schema: map[string]*schema.Schema{
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update of the resource",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content name",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if len(i.(string)) < 3 {
						return nil, []error{fmt.Errorf("name must be at least 3 characters long")}
					}
					return nil, nil
				},
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content description",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if len(i.(string)) < 3 {
						return nil, []error{fmt.Errorf("description must be at least 3 characters long")}
					}
					return nil, nil
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content type",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(string) != "exercise" {
						return nil, []error{fmt.Errorf("type must be exercise (more types will come in the future)")}
					}
					return nil, nil
				},
			},
			"reward": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The content reward",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(int) < 0 {
						return nil, []error{fmt.Errorf("reward must be a positive integer")}
					}
					return nil, nil
				},
			},
			"container": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The content component",
				Elem:        resourceContentContainer(1),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// `resourceContentContainer` is a recursive function that will create a schema for the content container to a max level of 3 nested containers
func resourceContentContainer(i int) *schema.Resource {
	if i > 3 {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the component",
			},
			"position": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The position where the component will be rendered",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(int) < 0 {
						return nil, []error{fmt.Errorf("position must be a positive integer")}
					}
					return nil, nil
				},
			},
			"orientation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The orientation of the container",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(string) != "horizontal" && i.(string) != "vertical" {
						return nil, []error{fmt.Errorf("orientation must be horizontal or vertical")}
					}
					return nil, nil
				},
			},
			"markdown": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The markdown component",
				Elem:        resourceContentDataMarkdown(),
			},
			"editor": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The editor component",
				Elem:        resourceContentDataEditor(),
			},
			"container": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The container component",
				Elem:        resourceContentContainer(i + 1),
			},
		},
	}
}

func resourceContentDataMarkdown() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the component",
			},
			"position": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The position where the component will be rendered",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(int) < 0 {
						return nil, []error{fmt.Errorf("position must be a positive integer")}
					}
					return nil, nil
				},
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content of the markdown",
			},
		},
	}
}

func resourceContentDataEditor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentCreate,
		ReadContext:   resourceContentRead,
		UpdateContext: resourceContentUpdate,
		DeleteContext: resourceContentDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the component",
			},
			"position": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The position where the component will be rendered",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(int) < 0 {
						return nil, []error{fmt.Errorf("position must be a positive integer")}
					}
					return nil, nil
				},
			},
			"language_settings": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of languages for the editor",
				Elem:        resourceContentLanguage(),
			},
			"hint": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of hints id for the editor",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"validator": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of validators for the editor",
				Elem:        resourceContentValidator(),
			},
		},
	}
}

func resourceContentValidator() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the validator",
			},
			"inputs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of inputs for the validator",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"outputs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of outputs for the validator",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_hidden": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the validator is hidden",
			},
		},
	}
}

func resourceContentLanguage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default code of the language",
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The language",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					if i.(string) != "PYTHON" && i.(string) != "NODE" && i.(string) != "JAVA" && i.(string) != "RUST" {
						return nil, []error{fmt.Errorf("language must be one of python, javascript, java or c")}
					}
					return nil, nil
				},
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The version of the language",
			},
		},
	}
}

func resourceContentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*pc.Client)

	var diags diag.Diagnostics

	rootComponent, err := serializeRootComponent(d.Get("container.0").(map[string]interface{}), ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to serialize child components",
			Detail:   fmt.Sprintf("Error when serializing child components: %s", err.Error()),
		})
		return diags
	}

	co := content.Content{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Type:          d.Get("type").(string),
		Reward:        int64(d.Get("reward").(int)),
		RootComponent: *rootComponent,
		Data:          content.ContentData{},
	}

	createdContent, err := c.CreateContent(co)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create content",
			Detail:   fmt.Sprintf("Error when creating content: %s", err.Error()),
		})
		return diags
	}

	d.SetId(createdContent.ID)

	tflog.Info(ctx, fmt.Sprintf("Created Content %s", d.Id()))

	return resourceContentRead(ctx, d, m)
}

func resourceContentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*pc.Client)

	var diags diag.Diagnostics

	tflog.Debug(ctx, fmt.Sprintf("Reading Content %s", d.Id()))

	content, err := c.GetContent(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Content",
			Detail:   fmt.Sprintf("Error when getting Content: %s", err.Error()),
		})
		return diags
	}

	err = d.Set("name", content.Name)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set name",
			Detail:   fmt.Sprintf("Error when setting name: %s", err.Error()),
		})
		return diags
	}
	err = d.Set("description", content.Description)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set description",
			Detail:   fmt.Sprintf("Error when setting description: %s", err.Error()),
		})
		return diags
	}
	err = d.Set("type", content.Type)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set type",
			Detail:   fmt.Sprintf("Error when setting type: %s", err.Error()),
		})
		return diags
	}
	err = d.Set("reward", content.Reward)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set reward",
			Detail:   fmt.Sprintf("Error when setting reward: %s", err.Error()),
		})
		return diags
	}
	err = d.Set("container", deserializeRootComponent(content.RootComponent, 0, ctx))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set container",
			Detail:   fmt.Sprintf("Error when setting container: %s", err.Error()),
		})
	}

	return diags
}

func resourceContentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*pc.Client)

	rootComponent, err := serializeRootComponent(d.Get("container.0").(map[string]interface{}), ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to serialize child components",
			Detail:   fmt.Sprintf("Error when serializing child components: %s", err.Error()),
		})
		return diags
	}

	co := content.Content{
		ID:            d.Id(),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Type:          d.Get("type").(string),
		Reward:        int64(d.Get("reward").(int)),
		RootComponent: *rootComponent,
		Data:          content.ContentData{},
	}

	_, err = c.UpdateContent(co)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update content",
			Detail:   fmt.Sprintf("Error when updating content: %s", err.Error()),
		})
		return diags
	}

	err = d.Set("last_update", time.Now().Format(time.RFC850))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set last_update",
			Detail:   fmt.Sprintf("Error when setting last_update: %s", err.Error()),
		})
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("Updated Content %s", d.Id()))

	return resourceContentRead(ctx, d, m)
}

func resourceContentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*pc.Client)

	err := c.DeleteContent(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Content",
			Detail:   fmt.Sprintf("Error when deleting Content: %s", err.Error()),
		})
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted Content %s", d.Id()))

	return diags
}

// `serializeRootComponent` takes the schema of a root component and returns itself as content.Component struct
func serializeRootComponent(rootComponent map[string]interface{}, ctx context.Context) (*content.Component, error) {
	length := 0
	for key, val := range rootComponent {
		if (key == "container" || key == "markdown" || key == "editor") && val != nil {
			length += len(val.([]interface{}))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Serializing root Component %s with %d child components", rootComponent["id"], length))

	childComponents := make([]content.Component, length)
	positions := make([]bool, length)

	for key, val := range rootComponent {
		switch key {
		case "markdown":
			for _, v := range val.([]interface{}) {
				markdown := v.(map[string]interface{})

				position := markdown["position"].(int)
				if position > length {
					return nil, fmt.Errorf("position %d is greater than the number of child components %d", position, length)
				}
				positions[position-1] = true

				childComponents[position-1] = content.Component{
					ID:   markdown["id"].(string),
					Type: "markdown",
					Data: content.ComponentData{
						Markdown: markdown["content"].(string),
					},
				}
			}
		case "editor":
			for _, v := range val.([]interface{}) {
				editor := v.(map[string]interface{})

				languages := make([]content.Language, 0)
				for _, language := range editor["language_settings"].([]interface{}) {
					languages = append(languages, content.Language{
						DefaultCode: language.(map[string]interface{})["default_code"].(string),
						Language:    language.(map[string]interface{})["language"].(string),
						Version:     language.(map[string]interface{})["version"].(string),
					})
				}

				validators := make([]content.Validator, 0)
				for _, validator := range editor["validator"].([]interface{}) {
					inputs := make([]string, len(validator.(map[string]interface{})["inputs"].([]interface{})))
					for key, val := range validator.(map[string]interface{})["inputs"].([]interface{}) {
						inputs[key] = val.(string)
					}
					outputs := make([]string, len(validator.(map[string]interface{})["outputs"].([]interface{})))
					for key, val := range validator.(map[string]interface{})["outputs"].([]interface{}) {
						outputs[key] = val.(string)
					}

					validators = append(validators, content.Validator{
						ID:       validator.(map[string]interface{})["id"].(string),
						IsHidden: validator.(map[string]interface{})["is_hidden"].(bool),
						Input: content.ValidatorInput{
							Stdin: inputs,
						},
						Output: content.ValidatorOutput{
							Stdout: outputs,
						},
					})
				}

				hints := make([]content.ItemIdentifier, 0)
				for _, item := range editor["hint"].([]interface{}) {
					hints = append(hints, content.ItemIdentifier{
						ID: item.(string),
					})
				}

				position := editor["position"].(int)
				if position > length {
					return nil, fmt.Errorf("position %d is greater than the number of child components %d", position, length)
				}
				positions[position-1] = true

				childComponents[position-1] = content.Component{
					ID:   editor["id"].(string),
					Type: "editor",
					Data: content.ComponentData{
						EditorSettings: content.EditorSettings{
							Languages: languages,
						},
						Validators: validators,
						Items:      hints,
					},
				}
			}
		case "container":
			for _, v := range val.([]interface{}) {
				container := v.(map[string]interface{})

				position := container["position"].(int)
				if position > length {
					return nil, fmt.Errorf("position %d is greater than the number of child components %d", position, length)
				}
				positions[position-1] = true

				containerComponent, err := serializeRootComponent(container, ctx)
				if err != nil {
					return nil, err
				}

				childComponents[position-1] = *containerComponent
			}
		}
	}

	for i, position := range positions {
		if !position {
			return nil, fmt.Errorf("child component at position %d is missing, this is probably due to duplicate position in the container, please check that your positions go from 1 to %d", i+1, length)
		}
	}

	return &content.Component{
		ID:          rootComponent["id"].(string),
		Orientation: rootComponent["orientation"].(string),
		Type:        "container",
		Data: content.ComponentData{
			Components: childComponents,
		},
	}, nil
}

// `deserializeRootComponent` takes a root component and convert it into a schema
func deserializeRootComponent(rootComponent content.Component, position int, ctx context.Context) []interface{} {
	markdown := make([]interface{}, 0)
	editor := make([]interface{}, 0)
	container := make([]interface{}, 0)

	tflog.Debug(ctx, fmt.Sprintf("Deserializing root Component %s", rootComponent.ID))

	for key, childComponent := range rootComponent.Data.Components {
		switch childComponent.Type {
		case "markdown":
			markdown = append(markdown, map[string]interface{}{
				"id":       childComponent.ID,
				"content":  childComponent.Data.Markdown,
				"position": key + 1,
			})
		case "editor":
			languages := make([]interface{}, 0)
			for _, language := range childComponent.Data.EditorSettings.Languages {
				languages = append(languages, map[string]interface{}{
					"default_code": language.DefaultCode,
					"language":     language.Language,
					"version":      language.Version,
				})
			}

			validators := make([]interface{}, 0)
			for _, validator := range childComponent.Data.Validators {
				inputs := make([]interface{}, len(validator.Input.Stdin))
				for key, val := range validator.Input.Stdin {
					inputs[key] = val
				}
				outputs := make([]interface{}, len(validator.Output.Stdout))
				for key, val := range validator.Output.Stdout {
					outputs[key] = val
				}

				validators = append(validators, map[string]interface{}{
					"id":        validator.ID,
					"is_hidden": validator.IsHidden,
					"inputs":    inputs,
					"outputs":   outputs,
				})
			}

			hints := make([]interface{}, 0)
			for _, item := range childComponent.Data.Items {
				hints = append(hints, item.ID)
			}

			editor = append(editor, map[string]interface{}{
				"id":                childComponent.ID,
				"language_settings": languages,
				"validator":         validators,
				"hint":              hints,
				"position":          key + 1,
			})
		case "container":
			container = deserializeRootComponent(childComponent, key+1, ctx)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"id":          rootComponent.ID,
			"orientation": rootComponent.Orientation,
			"position":    position,
			"markdown":    markdown,
			"editor":      editor,
			"container":   container,
		},
	}
}
