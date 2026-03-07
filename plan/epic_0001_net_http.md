# PRD Title
Refactor for removing `echo` dependency and use `net/http` etc  in controller layer instead.

## 1. Purpose
- use standard library and patterns.

## 2. Scope
- controller layer to replace `echo` functions with standard library functions

## 5. Technical Requirements
### 5.1 Code Changes
- List modules, packages, or services that will be touched.

### 5.3 Tests and Implementation Approach
- Unit / integration / regression tests required, iteration by controller struct.
- Expected coverage targets.
- Critical edge cases to ensure functionality is preserved.

## 6. Risks & Mitigations
- Risk: breaking existing functionality.
    - Mitigation: maintain full test coverage, incremental commits.

## 7. Success Criteria
- Code readability improved (examples if possible).
- Test coverage maintained or improved.

## Example of current code and proposed change