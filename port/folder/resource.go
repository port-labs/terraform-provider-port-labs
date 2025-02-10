package folder

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &FolderResource{}
var _ resource.ResourceWithImportState = &FolderResource{}

func NewFolderResource() resource.Resource {
	return &FolderResource{}
}

type FolderResource struct {
	portClient *cli.PortClient
}

func (r *FolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *FolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *FolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: FolderResourceMarkdownDescription,
		Attributes:          FolderSchema(),
	}
}

func (r *FolderResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state FolderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	betaFeaturesEnabledEnv := os.Getenv("PORT_BETA_FEATURES_ENABLED")
	if !(betaFeaturesEnabledEnv == "true") {
		resp.Diagnostics.AddError("Beta features are not enabled", "Folder resource is currently in beta and is subject to change in future versions. Use it by setting the Environment Variable PORT_BETA_FEATURES_ENABLED=true.")
		return
	}
}

func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *FolderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder := FolderModelToCLI(state)
	createdFolder, err := r.portClient.CreateFolder(ctx, folder)
	if err != nil {
		resp.Diagnostics.AddError("failed to create folder", err.Error())
		return
	}

	writeFolderComputedFieldsToState(state, createdFolder)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *FolderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f, statusCode, err := r.portClient.GetFolder(ctx, state.Identifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read folder", err.Error())
		return
	}

	err = refreshFolderToState(state, f)
	if err != nil {
		resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *FolderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder := FolderModelToCLI(state)
	updatedFolder, err := r.portClient.UpdateFolder(ctx, folder)
	if err != nil {
		resp.Diagnostics.AddError("failed to update folder", err.Error())
		return
	}

	writeFolderComputedFieldsToState(state, updatedFolder)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *FolderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Identifier.IsNull() {
		resp.Diagnostics.AddError("failed to extract folder identifier", "identifier is required")
		return
	}

	statusCode, err := r.portClient.DeleteFolder(ctx, state.ID.ValueString())
	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete folder", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func FolderModelToCLI(state *FolderModel) *cli.Folder {
	return &cli.Folder{
		Identifier: state.Identifier.ValueString(),
		Title:      state.Title.ValueString(),
		After:      state.After.ValueString(),
		Parent:     state.Parent.ValueString(),
	}
}

func writeFolderComputedFieldsToState(state *FolderModel, fr *cli.Folder) {
	state.ID = types.StringValue(fr.Identifier)
	state.Identifier = types.StringValue(fr.Identifier)

	// if fr.Parent == nil {
	// 	state.Parent = types.StringNull()
	// } else {
	// 	state.Parent = types.StringValue(fr.Parent)
	// }

	if fr.Parent != "" {
		state.Parent = types.StringValue(fr.Parent)
	}
	// else {
	// 	state.Parent = types.StringNull()
	// }

	if fr.After != "" {
		state.After = types.StringValue(fr.After)
	}
	// else {
	// 	state.After = types.StringNull()
	// }
	if fr.Title != "" {
		state.Title = types.StringValue(fr.Title)
	}
}

// package folder

// import (
// 	"context"

// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
// )

// var _ resource.Resource = &FolderResource{}
// var _ resource.ResourceWithImportState = &FolderResource{}

// func NewFolderResource() resource.Resource {
// 	return &FolderResource{}
// }

// type FolderResource struct {
// 	portClient *cli.PortClient
// }

// func (r *FolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_folder"
// }

// func (r *FolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}
// 	r.portClient = req.ProviderData.(*cli.PortClient)
// }

// func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var state *FolderModel

// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	f, statusCode, err := r.portClient.GetFolder(ctx, state.Identifier.ValueString())
// 	if err != nil {
// 		if statusCode == 404 {
// 			resp.State.RemoveResource(ctx)
// 			return
// 		}
// 		resp.Diagnostics.AddError("failed to read folder", err.Error())
// 		return
// 	}

// 	err = refreshFolderToState(state, f)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
// 		return
// 	}

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	var state *FolderModel
// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	folder := FolderToPortRequestCli(state)
// 	createdFolder, err := r.portClient.CreateFolder(ctx, folder)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to create folder", err.Error())
// 		return
// 	}

// 	writeFolderComputedFieldsToState(state, createdFolder)
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var state *FolderModel
// 	var previousState *FolderModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
// 	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	f := FolderToPortRequestCli(state)

// 	var bp *cli.Folder
// 	var err error
// 	if previousState.Identifier.IsNull() {
// 		bp, err = r.portClient.CreateFolder(ctx, f)
// 		if err != nil {
// 			resp.Diagnostics.AddError("failed to create folder", err.Error())
// 			return
// 		}
// 	} else {
// 		existingFolder, statusCode, err := r.portClient.GetFolder(ctx, previousState.Identifier.ValueString())
// 		if err != nil {
// 			if statusCode == 404 {
// 				resp.Diagnostics.AddError("Folder doesn't exists, it is required to update the folder", err.Error())
// 				return
// 			}
// 			resp.Diagnostics.AddError("failed getting folder", err.Error())
// 			return
// 		}

// 		f.Identifier = existingFolder.Identifier
// 		bp, err = r.portClient.UpdateFolder(ctx, f)
// 		if err != nil {
// 			resp.Diagnostics.AddError("failed to update folder", err.Error())
// 			return
// 		}
// 	}

// 	writeFolderComputedFieldsToState(state, bp)
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	var state *FolderModel
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	if state.Identifier.IsNull() {
// 		resp.Diagnostics.AddError("failed to extract folder identifier", "identifier is required")
// 		return
// 	}

// 	statusCode, err := r.portClient.DeleteFolder(ctx, state.ID.ValueString())
// 	if err != nil {
// 		if statusCode == 404 {
// 			resp.State.RemoveResource(ctx)
// 			return
// 		}
// 		resp.Diagnostics.AddError("failed to delete folder", err.Error())
// 		return
// 	}

// 	resp.State.RemoveResource(ctx)
// }

// func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resp.Diagnostics.Append(resp.State.SetAttribute(
// 		ctx, path.Root("identifier"), req.ID,
// 	)...)
// 	resp.Diagnostics.Append(resp.State.SetAttribute(
// 		ctx, path.Root("id"), req.ID,
// 	)...)
// }

// func FolderToPortRequestCli(state *FolderModel) *cli.Folder {
// 	return &cli.Folder{
// 		Identifier: state.Identifier.ValueString(),
// 		Title:      state.Title.ValueString(),
// 		After:      state.After.ValueString(),
// 		Parent:     state.Parent.ValueString(),
// 	}
// }

// func writeFolderComputedFieldsToState(state *FolderModel, fr *cli.Folder) {
// 	state.ID = types.StringValue(fr.Identifier)
// 	state.Identifier = types.StringValue(fr.Identifier)

// 	if fr.Parent != "" {
// 		state.Parent = types.StringValue(fr.Parent)
// 	} else {
// 		state.Parent = types.StringNull()
// 	}

// 	if fr.After != "" {
// 		state.After = types.StringValue(fr.After)
// 	} else {
// 		state.After = types.StringNull()
// 	}

// 	state.Title = types.StringValue(fr.Title)
// }
