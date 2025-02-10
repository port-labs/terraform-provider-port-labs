package folder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	// client, ok := req.ProviderData.(*cli.PortClient)
	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Resource Configure Type",
	// 		fmt.Sprintf("Expected *cli.PortClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	// 	)
	// 	return
	// }

	// r.portClient = client
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

// func (r *FolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		MarkdownDescription: "Folder resource",
// 		Attributes:          FolderSchema(),
// 	}
// }

func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *FolderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	f, err := FolderToPortRequest(state)

	// folderId := state.ID.ValueString()

	if err != nil {
		resp.Diagnostics.AddError("failed to create blueprint", err.Error())
		return
	}

	folder, err := r.portClient.CreateFolder(ctx, f)
	if err != nil {
		resp.Diagnostics.AddError("failed to create folder", err.Error())
		return
	}

	writeFolderComputedFieldsToState(state, folder)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	// var state *FolderModel

	// resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// if state.Identifier.IsNull() {
	// 	state.Identifier = types.StringValue(utils.GenID())
	// }

	// folder, err := FolderToPortBody(state)
	// if err != nil {
	// 	resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
	// 	return
	// }

	// f, err := r.portClient.CreateFolder(ctx, folder)
	// if err != nil {
	// 	resp.Diagnostics.AddError("failed to create folder", err.Error())
	// 	return
	// }

	// // state.ID = types.StringValue(f.Identifier)
	// // state.Identifier = types.StringValue(f.Identifier)
	// err = refreshFolderToState(state, f)
	// if err != nil {
	// 	resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
	// 	return
	// }

	// resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func writeFolderComputedFieldsToState(state *FolderModel, fr *cli.Folder) {
	state.ID = types.StringValue(fr.Identifier)
	state.Identifier = types.StringValue(fr.Identifier)

	if fr.Parent != "" {
		state.Parent = types.StringValue(fr.Parent)
	} else {
		state.Parent = types.StringNull()
	}

	if fr.After != "" {
		state.After = types.StringValue(fr.After)
	} else {
		state.After = types.StringNull()
	}

	state.Title = types.StringValue(fr.Title)
}

func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *FolderModel
	var previousState *FolderModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	f, err := FolderToPortRequest(state)

	if err != nil {
		resp.Diagnostics.AddError("failed to transform folder", err.Error())
		return
	}

	var fr *cli.Folder
	if previousState.Identifier.IsNull() {
		fr, err = r.portClient.CreateFolder(ctx, f)
		if err != nil {
			resp.Diagnostics.AddError("failed to create folder", err.Error())
			return
		}
	} else {
		existingFolder, statusCode, err := r.portClient.GetFolder(ctx, previousState.Identifier.ValueString())
		if err != nil {
			if statusCode == 404 {
				resp.Diagnostics.AddError("Folder doesn't exists, it is required to update the folder", err.Error())
				return
			}
			resp.Diagnostics.AddError("failed reading folder", err.Error())
			return
		}

		// f.Title = existingFolder.Title
		f.Identifier = existingFolder.Identifier
		fr, err = r.portClient.UpdateFolder(ctx, f)
		if err != nil {
			resp.Diagnostics.AddError("failed to update folder", err.Error())
			return
		}
	}

	/*
		resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		folder, err := FolderToPortBody(state)
		if err != nil {
			resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
			return
		}

		f, err := r.portClient.UpdateFolder(ctx, folder)
		if err != nil {
			resp.Diagnostics.AddError("failed to update folder", err.Error())
			return
		}

		// state.ID = types.StringValue(f.Identifier)
		// err = refreshFolderToState(state, f)
		// if err != nil {
		// 	resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
		// 	return
		// }
	*/
	writeFolderComputedFieldsToState(state, fr)

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

	//Matan 10/2
	// resp.State.RemoveResource(ctx)
}

func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("identifier"), req.ID,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
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

// func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resp.Diagnostics.Append(resp.State.SetAttribute(
// 		ctx, path.Root("identifier"), req.ID,
// 	)...)

// 	resp.Diagnostics.Append(resp.State.SetAttribute(
// 		ctx, path.Root("id"), req.ID,
// 	)...)
// }

// func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var state *FolderModel

// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	f, statusCode, err := r.portClient.GetFolder(ctx, state.ID.ValueString())

// 	if err != nil {
// 		if statusCode == 404 {
// 			resp.State.RemoveResource(ctx)
// 			return
// 		}
// 		resp.Diagnostics.AddError("failed to get folder", err.Error())
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
// 	var state FolderModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	folder, err := FolderToPortBody(&state)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
// 		return
// 	}

// 	f, err := r.portClient.CreateFolder(ctx, folder)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to create folder", err.Error())
// 		return
// 	}

// 	state.ID = types.StringValue(f.Identifier)
// 	err = refreshFolderToState(&state, f)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to write folder fields to resource", err.Error())
// 		return
// 	}

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	var state FolderModel
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
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

// // func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// // 	var state FolderModel

// // 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
// // 	if resp.Diagnostics.HasError() {
// // 		return
// // 	}

// // 	folder, err := FolderToPortBody(&state)
// // 	if err != nil {
// // 		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
// // 		return
// // 	}

// // 	f, err := r.portClient.CreateFolder(ctx, folder)
// // 	if err != nil {
// // 		resp.Diagnostics.AddError("failed to create folder", err.Error())
// // 		return
// // 	}
// // 	if f == nil {
// // 		f, _, err = r.portClient.GetFolder(ctx, state.ID.ValueString())
// // 		if err != nil {
// // 			resp.Diagnostics.AddError("failed to get folder", err.Error())
// // 			return
// // 		}
// // 	}

// // 	state.ID = types.StringValue(f.Identifier)
// // 	refreshFolderToState(&state, f)

// // 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// // }

// // func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// // 	var state FolderModel
// // 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// // 	if resp.Diagnostics.HasError() {
// // 		return
// // 	}

// // 	statusCode, err := r.portClient.DeleteFolder(ctx, state.ID.ValueString())
// // 	if err != nil {
// // 		if statusCode == 404 {
// // 			resp.State.RemoveResource(ctx)
// // 			return
// // 		}
// // 		resp.Diagnostics.AddError("failed to delete folder", err.Error())
// // 		return
// // 	}

// // 	resp.State.RemoveResource(ctx)
// // }

// // func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// // 	var state *FolderModel

// // 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

// // 	if resp.Diagnostics.HasError() {
// // 		return
// // 	}

// // 	statusCode, err := r.portClient.DeleteFolder(ctx, state.ID.ValueString())
// // 	if err != nil {
// // 		if statusCode == 404 {
// // 			resp.State.RemoveResource(ctx)
// // 			return
// // 		}
// // 		resp.Diagnostics.AddError("failed to delete folder", err.Error())
// // 		return
// // 	}

// // 	resp.State.RemoveResource(ctx)
// // }

// // func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// // 	var state *FolderModel

// // 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

// // 	if resp.Diagnostics.HasError() {
// // 		return
// // 	}

// // 	folder, err := FolderToPortBody(state)
// // 	if err != nil {
// // 		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
// // 		return
// // 	}

// // 	f, err := r.portClient.CreateFolder(ctx, folder)
// // 	if err != nil {
// // 		resp.Diagnostics.AddError("failed to create folder", err.Error())
// // 		return
// // 	}
// // 	if f == nil {
// // 		f, _, err = r.portClient.GetFolder(ctx, state.ID.ValueString())
// // 		if err != nil {
// // 			resp.Diagnostics.AddError("failed to get folder", err.Error())
// // 			return
// // 		}
// // 	}

// // 	refreshFolderToState(state, f)
// // 	// state.Identifier = types.StringValue(f.Identifier)
// // 	// state.Sidebar = types.StringValue(f.Sidebar)
// // 	// state.Parent = types.StringPointerValue(f.Parent)
// // 	// state.After = types.StringPointerValue(f.After)
// // 	// state.Title = types.StringPointerValue(f.Title)
// // 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// // }

// func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var state *FolderModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	f, _, err := r.portClient.GetFolder(ctx, state.ID.ValueString())
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to get folder", err.Error())
// 		return
// 	}

// 	folder, err := FolderToPortBody(state)
// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to convert folder resource to body", err.Error())
// 		return
// 	}

// 	_, err = r.portClient.UpdateFolder(ctx, folder)

// 	if err != nil {
// 		resp.Diagnostics.AddError("failed to update folder", err.Error())
// 		return
// 	}

// 	state.Sidebar = types.StringValue(f.Sidebar)
// 	refreshFolderToState(state, f)

// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }
