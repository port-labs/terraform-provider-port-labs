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

// ServerManagedAppSpecKeys are appSpec fields Port assigns during SaaS
// provisioning. Users cannot set them in Terraform; they must not appear in
// state or cause plan drift.
var ServerManagedAppSpecKeys = map[string]struct{}{
	"liveEventsUuid":           {},
	"liveEventsIngestHostname": {},
}

func IsServerManagedAppSpecKey(key string) bool {
	_, ok := ServerManagedAppSpecKeys[key]
	return ok
}
