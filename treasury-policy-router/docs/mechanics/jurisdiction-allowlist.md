# Jurisdiction Allowlist

Jurisdiction is represented as a 2-letter ISO country code in the allowlist. If a country is not present, it is treated as allowed by default. To block a country, explicitly set it to false.

## Example

```rust
set_jurisdiction("US", true, admin)
set_jurisdiction("IR", false, admin)
```
