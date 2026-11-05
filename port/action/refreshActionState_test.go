package action

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func TestWriteTriggerToResource_TimerPropertyExpiredEvent_NilFields(t *testing.T) {
	tests := []struct {
		name                string
		blueprintIdentifier *string
		propertyIdentifier  *string
		wantPanic           bool
	}{
		{
			name:                "nil BlueprintIdentifier panics before fix",
			blueprintIdentifier: nil,
			propertyIdentifier:  strPtr("eol_timer"),
			wantPanic:           false,
		},
		{
			name:                "nil PropertyIdentifier panics before fix",
			blueprintIdentifier: strPtr("myBlueprint"),
			propertyIdentifier:  nil,
			wantPanic:           false,
		},
		{
			name:                "both nil panics before fix",
			blueprintIdentifier: nil,
			propertyIdentifier:  nil,
			wantPanic:           false,
		},
		{
			name:                "both populated works",
			blueprintIdentifier: strPtr("myBlueprint"),
			propertyIdentifier:  strPtr("eol_timer"),
			wantPanic:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActionResource{}
			state := &ActionModel{}
			action := &cli.Action{
				Identifier: "test-action",
				Trigger: &cli.Trigger{
					Type: consts.Automation,
					Event: &cli.TriggerEvent{
						Type:                consts.TimerPropertyExpired,
						BlueprintIdentifier: tt.blueprintIdentifier,
						PropertyIdentifier:  tt.propertyIdentifier,
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("writeTriggerToResource() panicked unexpectedly: %v", r)
					}
				} else if tt.wantPanic {
					t.Error("writeTriggerToResource() did not panic as expected")
				}
			}()

			err := r.writeTriggerToResource(context.Background(), action, state)
			if err != nil {
				t.Errorf("writeTriggerToResource() returned error: %v", err)
			}

			if state.AutomationTrigger == nil {
				t.Fatal("AutomationTrigger should not be nil")
			}

			if state.AutomationTrigger.TimerPropertyExpiredEvent == nil {
				t.Fatal("TimerPropertyExpiredEvent should not be nil")
			}

			if tt.blueprintIdentifier == nil {
				if !state.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.IsNull() {
					t.Error("BlueprintIdentifier should be null when input is nil")
				}
			} else {
				if state.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.ValueString() != *tt.blueprintIdentifier {
					t.Errorf("BlueprintIdentifier = %v, want %v",
						state.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.ValueString(),
						*tt.blueprintIdentifier)
				}
			}

			if tt.propertyIdentifier == nil {
				if !state.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.IsNull() {
					t.Error("PropertyIdentifier should be null when input is nil")
				}
			} else {
				if state.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.ValueString() != *tt.propertyIdentifier {
					t.Errorf("PropertyIdentifier = %v, want %v",
						state.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.ValueString(),
						*tt.propertyIdentifier)
				}
			}
		})
	}
}

func TestWriteTriggerToResource_EntityEvents_NilBlueprintIdentifier(t *testing.T) {
	eventTypes := []string{
		consts.EntityCreated,
		consts.EntityUpdated,
		consts.EntityDeleted,
		consts.AnyEntityChange,
	}

	for _, eventType := range eventTypes {
		t.Run(eventType+"_nil_BlueprintIdentifier", func(t *testing.T) {
			r := &ActionResource{}
			state := &ActionModel{}
			action := &cli.Action{
				Identifier: "test-action",
				Trigger: &cli.Trigger{
					Type: consts.Automation,
					Event: &cli.TriggerEvent{
						Type:                eventType,
						BlueprintIdentifier: nil,
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("writeTriggerToResource() panicked for %s with nil BlueprintIdentifier: %v", eventType, r)
				}
			}()

			err := r.writeTriggerToResource(context.Background(), action, state)
			if err != nil {
				t.Errorf("writeTriggerToResource() returned error: %v", err)
			}
		})
	}
}

func TestWriteTriggerToResource_RunEvents_NilActionIdentifier(t *testing.T) {
	eventTypes := []string{
		consts.RunCreated,
		consts.RunUpdated,
		consts.AnyRunChange,
	}

	for _, eventType := range eventTypes {
		t.Run(eventType+"_nil_ActionIdentifier", func(t *testing.T) {
			r := &ActionResource{}
			state := &ActionModel{}
			action := &cli.Action{
				Identifier: "test-action",
				Trigger: &cli.Trigger{
					Type: consts.Automation,
					Event: &cli.TriggerEvent{
						Type:             eventType,
						ActionIdentifier: nil,
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("writeTriggerToResource() panicked for %s with nil ActionIdentifier: %v", eventType, r)
				}
			}()

			err := r.writeTriggerToResource(context.Background(), action, state)
			if err != nil {
				t.Errorf("writeTriggerToResource() returned error: %v", err)
			}
		})
	}
}

func TestWriteTriggerToResource_CustomerScenario(t *testing.T) {
	r := &ActionResource{}
	state := &ActionModel{}

	blueprintId := "myBlueprint"
	propertyId := "eol_timer"
	action := &cli.Action{
		Identifier: "test-automation",
		Trigger: &cli.Trigger{
			Type: consts.Automation,
			Event: &cli.TriggerEvent{
				Type:                consts.TimerPropertyExpired,
				BlueprintIdentifier: &blueprintId,
				PropertyIdentifier:  &propertyId,
			},
		},
	}

	err := r.writeTriggerToResource(context.Background(), action, state)
	if err != nil {
		t.Fatalf("writeTriggerToResource() returned error: %v", err)
	}

	if state.AutomationTrigger == nil {
		t.Fatal("AutomationTrigger should not be nil")
	}

	if state.AutomationTrigger.TimerPropertyExpiredEvent == nil {
		t.Fatal("TimerPropertyExpiredEvent should not be nil")
	}

	if state.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.ValueString() != blueprintId {
		t.Errorf("BlueprintIdentifier = %v, want %v",
			state.AutomationTrigger.TimerPropertyExpiredEvent.BlueprintIdentifier.ValueString(),
			blueprintId)
	}

	if state.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.ValueString() != propertyId {
		t.Errorf("PropertyIdentifier = %v, want %v",
			state.AutomationTrigger.TimerPropertyExpiredEvent.PropertyIdentifier.ValueString(),
			propertyId)
	}
}

func TestWriteTriggerToResource_SelfServiceTrigger(t *testing.T) {
	r := &ActionResource{}
	state := &ActionModel{}

	blueprintId := "myBlueprint"
	operation := "CREATE"
	action := &cli.Action{
		Identifier: "test-action",
		Trigger: &cli.Trigger{
			Type:                consts.SelfService,
			BlueprintIdentifier: &blueprintId,
			Operation:           &operation,
			UserInputs: &cli.ActionUserInputs{
				Properties: map[string]cli.ActionProperty{},
			},
		},
	}

	err := r.writeTriggerToResource(context.Background(), action, state)
	if err != nil {
		t.Fatalf("writeTriggerToResource() returned error: %v", err)
	}

	if state.SelfServiceTrigger == nil {
		t.Fatal("SelfServiceTrigger should not be nil")
	}

	if state.SelfServiceTrigger.BlueprintIdentifier.ValueString() != blueprintId {
		t.Errorf("BlueprintIdentifier = %v, want %v",
			state.SelfServiceTrigger.BlueprintIdentifier.ValueString(),
			blueprintId)
	}

	if state.SelfServiceTrigger.Operation != types.StringValue(operation) {
		t.Errorf("Operation = %v, want %v",
			state.SelfServiceTrigger.Operation,
			types.StringValue(operation))
	}
}

func strPtr(s string) *string {
	return &s
}
