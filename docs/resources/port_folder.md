# port_folder

A full list of the available folder types and their identifiers can be found [here](https://docs.getport.io/customize-pages-dashboards-and-plugins/folder/catalog-folder).

~> **WARNING**
The folder resource is currently in beta and is subject to change in future versions.
Use it by setting the Environment Variable `PORT_BETA_FEATURES_ENABLED=true`.
If this Environment Variable isn't specified, you won't be able to use the resource.

## Example Usage

### Basic Folder

```hcl
resource "port_folder" "example_folder" {
  sidebar_identifier    = "example_sidebar"
  folder_identifier     = "example_folder"
  title                 = "Example Folder"
  description           = "This is an example folder"
}
```

### Folder with Parent

Create a folder inside another folder.

```hcl
resource "port_folder" "child_folder" {
  sidebar_identifier    = "example_sidebar"
  folder_identifier     = "child_folder"
  parent                = port_folder.example_folder.folder_identifier
  title                 = "Child Folder"
  description           = "This is a child folder"
}
```

## Argument Reference

The following arguments are supported:

* `sidebar_identifier` - (Required) The Identifier of the sidebar.
* `folder_identifier` - (Required) The Identifier of the folder.
* `title` - (Optional) The title of the folder.
* `parent` - (Optional) The identifier of the parent folder.
* `after` - (Optional) The identifier of the folder after which the folder should be placed.
* `description` - (Optional) The folder description.

## Attribute Reference

The following attributes are exported:

* `sidebar_identifier` - The Identifier of the sidebar.
* `folder_identifier` - The ID of the folder.
* `title` - The title of the folder.
* `parent` - The identifier of the parent folder.
* `after` - The identifier of the folder after which the folder should be placed.
* `description` - The folder description.
* `created_at` - The creation date of the folder.
* `created_by` - The creator of the folder.
* `updated_at` - The last update date of the folder.
* `updated_by` - The last updater of the folder.
