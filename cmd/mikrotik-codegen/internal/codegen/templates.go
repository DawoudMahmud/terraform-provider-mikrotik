package codegen

const (
	generatedNotice = "// This code was generated. Review it carefully."

	resourceDefinitionTemplate = `
package {{ .Package }}

import (
	{{ range $import := .Imports -}}
		"{{ $import }}"
	{{ end }}
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"

)
{{ $resourceStructName := .ResourceName | firstLower}}
type {{$resourceStructName}} struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &{{$resourceStructName}}{}
	_ resource.ResourceWithConfigure   = &{{$resourceStructName}}{}
	_ resource.ResourceWithImportState = &{{$resourceStructName}}{}
)

// New{{.ResourceName}}Resource is a helper function to simplify the provider implementation.
func New{{.ResourceName}}Resource() resource.Resource {
	return &{{$resourceStructName}}{}
}

func (r *{{$resourceStructName}}) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *{{$resourceStructName}}) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{.ResourceName | lowercase}}"
}

// Schema defines the schema for the resource.
func (s *{{$resourceStructName}}) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik {{.ResourceName}}.",
		Attributes: map[string]schema.Attribute{
			{{range .Fields -}}
			"{{.AttributeName}}": schema.{{.Type}}Attribute{
				Required: {{.Required}},
				Optional: {{.Optional}},
				Computed: {{.Computed}},
				{{if .Computed -}}
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				{{- end}}
				Description: "Identifier of this resource assigned by RouterOS",
			},
			{{end}}
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *{{$resourceStructName}}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan {{$resourceStructName}}Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := r.client.Add{{.ResourceName}}(modelTo{{.ResourceName}}(&plan))
	if err != nil {
		resp.Diagnostics.AddError("creation failed", err.Error())
		return
	}

	resp.Diagnostics.Append({{.ResourceName | firstLower }}ToModel(created, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *{{$resourceStructName}}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state {{$resourceStructName}}Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.Find{{.ResourceName}}(state.{{.TerraformIDField.Name}}.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remote resource",
			fmt.Sprintf("Could not read {{.ResourceName}} with {{.TerraformIDField.AttributeName}} %q", state.{{.TerraformIDField.Name}}.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append({{$resourceStructName}}ToModel(resource, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *{{$resourceStructName}}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan {{$resourceStructName}}Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.Update{{.ResourceName}}(modelTo{{.ResourceName}}(&plan))
	if err != nil {
		resp.Diagnostics.AddError("update failed", err.Error())
		return
	}

	resp.Diagnostics.Append({{.ResourceName | firstLower}}ToModel(updated, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *{{$resourceStructName}}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state {{$resourceStructName}}Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Delete{{.ResourceName}}(state.{{.DeleteField.Name}}.ValueString()); err != nil {
		resp.Diagnostics.AddError("Could not delete {{.ResourceName}}", err.Error())
		return
	}
}

func (r *{{$resourceStructName}}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("{{.TerraformIDField.AttributeName}}"), req, resp)
}

type {{$resourceStructName}}Model struct {
	{{range .Fields -}}
	{{.Name}}        tftypes.{{.Type}} ` + "`" + `tfsdk:"{{.AttributeName}}"` + "`" + `
	{{end}}
}

func {{.ResourceName | firstLower}}ToModel(r *client.{{.ResourceName}}, m *{{$resourceStructName}}Model) diag.Diagnostics {
	var diags diag.Diagnostics
	if r == nil {
		diags.AddError("{{.ResourceName}} cannot be nil", "Cannot build model from nil object")
		return diags
	}

	{{range .Fields -}}
	m.{{.Name}} = tftypes.{{.Type}}Value(r.{{.MikrotikField.OriginalName}})
	{{end}}

	return diags
}

func modelTo{{.ResourceName}}(m *{{$resourceStructName}}Model) *client.{{.ResourceName}} {
	return &client.{{.ResourceName}}{
	{{range .Fields -}}
		{{.MikrotikField.OriginalName}}: m.{{.Name}}.Value{{.Type}}(),
	{{end}}
	}
}
`
)
