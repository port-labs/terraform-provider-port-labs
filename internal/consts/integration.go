package consts

const (
	InstallationTypeOnPrem     = "OnPrem"
	InstallationTypeSaas       = "Saas"
	InstallationTypeSaasOAuth2 = "SaasOAuth2"
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
	ValidateSaasSpec       = "saasSpecValidation"
	ValidateSaasOAuth2Spec = "saasOAuth2SpecValidation"
)

const (
	CreatePortResourcesOriginEmpty = "Empty"
	CreatePortResourcesOriginPort  = "Port"
)

func IsSaas(installationType string) bool {
	return installationType == InstallationTypeSaas
}

func IsSaasOAuth2(installationType string) bool {
	return installationType == InstallationTypeSaasOAuth2
}

func IsHosted(installationType string) bool {
	return IsSaas(installationType) || IsSaasOAuth2(installationType)
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
