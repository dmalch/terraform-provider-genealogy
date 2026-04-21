## 0.16.5 (Unreleased)

## 0.16.4

FEATURES:

* Union: added `foster_children` and `adopted_children` attributes on `geni_union` that map to Geni's `relationship_modifier=foster|adopt` edges. Each child appears in exactly one of `children`, `foster_children`, or `adopted_children`; the three sets must be disjoint. The provider passes the correct modifier through to `AddChild`/`AddSibling` on create and update, and splits the API's subset arrays back out on read so drift surfaces naturally.
* Union: changing a child's relationship modifier between applies (e.g. moving an id from `foster_children` to `adopted_children`) now emits an attribute warning — Geni has no API to re-tag an existing edge, so the change must be made on Geni.com first.

IMPROVEMENTS:

* Tooling: upgrade `golangci-lint` to v2.11.4 and migrate `.golangci.yml` to the v2 schema. The pinned binary is installed locally via `make lint` (into `bin/`) and in CI via `golangci-lint-action`, so the tool version no longer drifts.
* CI: add a linting step to the `ci.yaml` workflow that was previously build+test only.
* Code quality: fixed pre-existing lint findings surfaced by the upgrade (forcetypeassert on a schema downcast, gofmt drift in the profile package, an unused blank in a range clause, and a `WriteString(fmt.Sprintf(...))` callsite in the geni client).

## 0.16.3

SECURITY:

* Provider: generate a cryptographically random OAuth2 state parameter and validate it in the callback to prevent CSRF attacks.
* CI: add explicit read-only permissions to the CI workflow.

IMPROVEMENTS:

* Provider: support `GENI_ACCESS_TOKEN` and `GENI_USE_SANDBOX` environment variables as fallbacks for provider configuration, following the standard Terraform provider pattern.
* Testing: extract shared test helpers, add `CheckDestroy` verification for profiles and documents, add import state tests for all resource types, and read access token from environment variable instead of requiring source code edits.
* Testing: add unit tests for OAuth2 callback state validation, event converters, provider utils, and union validator.
* Maintenance: dependency updates (`google.golang.org/grpc`, `terraform-plugin-framework`, `terraform-plugin-go`, `terraform-plugin-testing`, `setup-terraform`, `setup-go`, `ghaction-import-gpg`) and grouped Dependabot updates by ecosystem.

## 0.16.2

BUG FIXES:

* Provider: redact `access_token` from debug log output to prevent credential leakage.
* Provider: replace panic with a returned error when cache initialization fails.
* Profile: fix `FistName` typo in internal `NameModel` struct field (renamed to `FirstName`).

## 0.16.1

IMPROVEMENTS:

* Maintenance update to dependencies (`golang.org/x/oauth2`, `golang.org/x/crypto`, `terraform-plugin-framework`, `terraform-plugin-testing`, `terraform-plugin-log`, `labstack/echo`, `onsi/gomega`, `cloudflare/circl`) and CI actions (`actions/checkout`, `actions/setup-go`, `golangci-lint-action`, `goreleaser-action`) to ensure compatibility and security.

## 0.16.0

FEATURES:

* Resources: added identity handling for `geni_profile`, `geni_document`, and `geni_union`. Resources now implement the Terraform Plugin Framework identity APIs, expose an identity schema (id) required for import, and persist identity data returned by the API.

## 0.15.3

IMPROVEMENTS:

* Profile: added validation of the `names` and `about` map to ensure locale keys are formatted correctly as per BCP 47 standards.

## 0.15.2

FEATURES:

* Profile: added a state upgrade path to automatically migrate prior state where `about` was a single string into the new `about` map shape. When upgrading, the previous `about` string is moved into the `"en-US"` locale key. Empty values are converted into a null map.

## 0.15.1

BUG FIXES:

* Fixed the build issue in the previous release.

## 0.15.0

FEATURES:

* Profile: added support for localized "about" text. The `about` attribute for the `Profile` resource is now a map of locale -> string and is backed by the Geni API's `detail_strings` field. The provider converts between `detail_strings` and the Terraform `about` map and falls back to `en-US` when necessary.

IMPROVEMENTS:

* Provider: token cache path now respects `use_sandbox_env` and uses a separate cache file (`~/.genealogy/geni_sandbox_token.json`) when sandbox mode is enabled.

BREAKING CHANGES:

* The `about` attribute type changed from a single string to a map (locale -> string). Review your state or code that references `geni_profile.about` to ensure it handles the new map shape.

## 0.14.10

IMPROVEMENTS:

* Maintenance update to dependencies to ensure compatibility and security.

## 0.14.9

IMPROVEMENTS:

* Retry logic for Geni API requests now handles "connection reset by peer" errors as retryable, improving resilience to transient network issues.

## 0.14.8

BUG FIXES:

* Fixed an issue where the `Profile` resource was not correctly removing `cause_of_death` attribute.

## 0.14.7

BUG FIXES:

* Fixed an issue where the `Profile` resource was not correctly removing events.

## 0.14.6

IMPROVEMENTS:

* Retry logic for Geni API requests has been improved to handle network timeouts and other transient errors more effectively.

## 0.14.5

IMPROVEMENTS:

* Retry logic for Geni API requests has been improved to handle DNS resolution and broken pipe errors more effectively.

## 0.14.4

IMPROVEMENTS:

* Update the timeout handling to avoid rate limiting issues with Geni API.

## 0.14.3

BUG FIXES:

* Fixed status refresh of profiles that were already removed in Geni.

## 0.14.2

BUG FIXES:

* Fixed the deletion of profiles that were already removed in Geni.

## 0.14.1

IMPROVEMENTS:

* Optimized the automatic updating of unions for merged profiles during state refresh.

## 0.14.0

FEATURES:

* Implemented support for nicknames as a new optional field in profiles.

## 0.13.2

IMPROVEMENTS:

* Updated batch processing functions to eliminate duplicate IDs using a hashset. This ensures optimized and accurate
  request handling for unions, profiles, and documents, simplifying the logic for single and multiple ID processing
  scenarios.

## 0.13.1

BUG FIXES:

* Fixed an issue where the `Profile` resource was not correctly handling computed subfields in the `current_residence`
  field.

## 0.13.0

FEATURES:

* Added support for `current_residence` field in the `Profile` resource.
