package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

type UserInputsMapper struct {
	resource *ActionResource
}

func NewUserInputsMapper(portClient *cli.PortClient) *UserInputsMapper {
	return &UserInputsMapper{resource: &ActionResource{portClient: portClient}}
}

func UserPropertiesToBody(ctx context.Context, userProperties *UserPropertiesModel, requiredJqQuery types.String) (map[string]cli.ActionProperty, any, error) {
	if userProperties == nil {
		userProperties = &UserPropertiesModel{}
	}

	trigger := &cli.Trigger{UserInputs: &cli.ActionUserInputs{}}
	data := &SelfServiceTriggerModel{
		UserProperties:  userProperties,
		RequiredJqQuery: requiredJqQuery,
	}

	if err := actionPropertiesToBody(ctx, trigger, data); err != nil {
		return nil, nil, err
	}

	return trigger.UserInputs.Properties, trigger.UserInputs.Required, nil
}

func UserInputTitlesToBody(ctx context.Context, titles map[string]ActionTitle) (map[string]cli.ActionTitle, error) {
	if titles == nil {
		return nil, nil
	}

	trigger := &cli.Trigger{UserInputs: &cli.ActionUserInputs{}}
	data := &SelfServiceTriggerModel{Titles: titles}

	if err := actionTitlesToBody(ctx, trigger, data); err != nil {
		return nil, err
	}

	return trigger.UserInputs.Titles, nil
}

// configured reports whether a user_properties block was declared, which decides
// whether an empty result becomes an empty model or nil.
func (m *UserInputsMapper) UserPropertiesToState(ctx context.Context, userInputs *cli.ActionUserInputs, configured bool) (*UserPropertiesModel, error) {
	if userInputs == nil {
		return nil, nil
	}

	action := &cli.Action{Trigger: &cli.Trigger{UserInputs: userInputs}}
	state := &ActionModel{}
	if configured {
		state.SelfServiceTrigger = &SelfServiceTriggerModel{UserProperties: &UserPropertiesModel{}}
	}

	return m.resource.buildUserProperties(ctx, action, state)
}

func (m *UserInputsMapper) UserInputTitlesToState(userInputs *cli.ActionUserInputs) (map[string]ActionTitle, error) {
	if userInputs == nil {
		return nil, nil
	}

	return m.resource.buildActionTitles(&cli.Action{Trigger: &cli.Trigger{UserInputs: userInputs}})
}

func UserInputsRequiredToState(userInputs *cli.ActionUserInputs) types.String {
	if userInputs == nil {
		return types.StringNull()
	}

	requiredJqQuery, _ := buildRequired(userInputs)
	return requiredJqQuery
}

func UserInputsOrderToState(ctx context.Context, order []string) types.List {
	if len(order) == 0 {
		return types.ListNull(types.StringType)
	}

	return flex.GoArrayStringToTerraformList(ctx, order)
}

func UserInputsOrderToBody(ctx context.Context, order types.List) ([]string, error) {
	if order.IsNull() {
		return nil, nil
	}

	values, err := utils.TerraformListToGoArray(ctx, order, "string")
	if err != nil {
		return nil, err
	}

	return utils.InterfaceToStringArray(values), nil
}
