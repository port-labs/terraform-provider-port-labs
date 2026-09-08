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

func IsSaas(installationType string) bool {
	return installationType == InstallationTypeSaas
}
