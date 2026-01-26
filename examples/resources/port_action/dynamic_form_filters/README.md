# Dynamic Form Filters Example

This example demonstrates how to use dynamic JQ expressions in dataset rules to create dependent form fields where the options in one field are filtered based on the selection in another field.

## Use Case

When building self-service actions, you often need to filter entity selectors based on:
- User's form selections (`.form.field_name`)
- Current user's context (`.user.teams`, `.user.email`)
- Entity context (`.entity.properties.field`)

## How It Works

1. User selects a team from the "Select Team" dropdown (e.g., "engineering")
2. The "Target Service" entity selector automatically filters to show only services where:
   - `team` matches the selected team (dynamic)
   - `environment` is NOT "production" (literal string comparison)
3. The filtering is dynamic - changing the team selection updates the available services

## Key Configuration

```hcl
dataset = {
  combinator = "and"
  rules = [
    # Dynamic JQ expression - evaluates form input at runtime
    {
      property = "team"
      operator = "="
      value = {
        jq_query = ".form.selected_team"
      }
    },
    # Literal string comparison - note the escaped quotes
    {
      property = "environment"
      operator = "!="
      value = {
        jq_query = "\"production\""
      }
    }
  ]
}
```

## JQ Expression Types

| Type | Example | Description |
|------|---------|-------------|
| Dynamic | `.form.selected_team` | Evaluates to the form field value at runtime |
| Literal | `"\"production\""` | A literal string "production" (quotes inside) |

## Running the Example

```bash
# Set credentials
export TF_VAR_port_client_id="your-client-id"
export TF_VAR_port_client_secret="your-client-secret"

# Apply
terraform init
terraform apply

# Cleanup
terraform destroy
```

## Testing

1. Go to Port UI → Self-Service Actions
2. Find "Select Service by Team"
3. Select "engineering" → see only `dev-api-service` (prod is filtered out)
4. Select "platform" → see `staging-web-frontend` and `dev-web-frontend`

## Common JQ Expressions

| Expression | Description |
|------------|-------------|
| `.form.field_name` | Value from another form field |
| `.user.teams` | Current user's team memberships |
| `.user.email` | Current user's email |
| `.entity.identifier` | Current entity's identifier |
| `"literal"` | A literal string value |
