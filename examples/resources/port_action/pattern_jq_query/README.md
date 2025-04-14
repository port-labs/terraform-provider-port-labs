# Pattern JQ Query Examples

This example demonstrates how to use the `pattern_jq_query` attribute in different ways:

## Ways to Use `pattern_jq_query`

The `pattern_jq_query` field accepts a JQ expression that gets evaluated by the Port API at runtime. This JQ expression can produce:

### 1. A Regex Pattern

```hcl
pattern_jq_query = "if .environment == \"production\" then \"^[a-z][a-z0-9-]{3,20}$\" else \"^[a-z][a-z0-9-]{2,10}$\" end"
```

This approach generates a regex pattern dynamically based on context. It works like the static `pattern` field but allows the regex to change based on other properties.

### 2. A List of Allowed Values

```hcl
pattern_jq_query = "if .team == \"platform\" then [\"dev\", \"staging\", \"production\"] else [\"dev\", \"staging\"] end"
```

This approach generates a list of allowed values dynamically. It's similar to an enum but with values that can change based on context.

### 3. Direct JSON Array of Allowed Values

```hcl
pattern_jq_query = "[\"value1\", \"value2\", \"value3\"]"
```

This simpler format can be used to specify a fixed list of allowed values directly.

## Important Notes

- You cannot use both `pattern` and `pattern_jq_query` at the same time on the same property
- The JQ query is evaluated at runtime by the Port API
- The context available to the JQ expression depends on where the pattern is being used (entity context, action context, etc.)
- For regex patterns, make sure to escape special characters properly

## How to Use

1. Choose which approach is most appropriate for your use case
2. Write a JQ expression or JSON array as needed
3. Apply the pattern validation to the appropriate string property in your action definition 