package consts

const (
	InstallationTypeOnPrem = "OnPrem"
	InstallationTypeSaas   = "Saas"
)

const (
	IntegrationStatusCreating  = "Creating"
	IntegrationStatusRunning   = "Running"
	IntegrationStatusError     = "Error"
	IntegrationStatusUnHealthy = "UnHealthy"
	IntegrationStatusUpdating  = "Updating"
	IntegrationStatusDeleting  = "Deleting"
)

const (
	ValidateIntegrationSpecModeFull = "full"
)

const (
	CreatePortResourcesOriginEmpty = "Empty"
	CreatePortResourcesOriginPort  = "Port"
)

func IsSaas(installationType string) bool {
	return installationType == InstallationTypeSaas
}
