resource "port_folder" "example_folder" {
  identifier = "example_folder"
  title      = "Example Folder"
}

resource "port_folder" "child_folder" {
  identifier = "child_folder"
  parent     = port_folder.example_folder.identifier
  title      = "Child Folder"
}

resource "port_folder" "another_folder" {
  identifier = "another_folder"
  after      = port_folder.example_folder.identifier
  title      = "Another Folder"
}