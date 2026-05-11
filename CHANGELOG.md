## 0.19.1 (Unreleased)

## 0.19.0

FEATURES:

* New data source `geni_profile`: look up an existing Geni profile by `id` or `guid` (exactly one) without managing it. Honors the provider's `auto_update_merged_profiles` flag — when set, the data source walks the `merged_into` chain (up to 10 hops) to return the surviving profile.

IMPROVEMENTS:

* Profile resource: exposes `guid` as a computed attribute alongside the existing `id`, surfacing the globally unique identifier Geni assigns to each profile.

## 0.18.2

BUG FIXES:

* Profile: events (`birth`, `death`, `baptism`, `burial`) that the Geni API auto-creates from a sibling input — most visibly `death` when `cause_of_death` is set — no longer flap the refresh plan. The Read path now treats an event that carries only a server-generated name (no `date` and no `location`) as no event, so `state.death` stays null when the user did not author a `death` block in HCL. User-set events still round-trip because the schema validator already requires them to carry `date` or `location`. (#91)

## 0.18.1

BUG FIXES:

* Profile: `terraform import`, batched managed-resource Read, and the v0.18.0 `geni_profile` list resource now request `project_ids` from the Geni API and copy it into the `projects` attribute. Previously the field was omitted from the shared `fields=` query parameter, so state landed with `projects = null` and every subsequent `terraform plan` against an HCL `projects = [...]` showed a spurious in-place update — in large workspaces this masked real drift behind thousands of no-op changes. (#89)

IMPROVEMENTS:

* Testing: added acceptance tests that verify `terraform plan -generate-config-out` produces a no-diff HCL config for `geni_profile`, `geni_document`, and `geni_union` via the `terraform-plugin-testing` 1.16 `GenerateConfig` ImportState mode. The framework's auto-implemented `GenerateResourceConfiguration` RPC (shipped on `terraform-plugin-framework` v1.19.0 since v0.17.0) round-trips cleanly for every attribute the API returns — the seed configs intentionally exclude set-only attributes (`projects` on profile/document; `text` / `file` / `file_name` on document) because the generated HCL omits them and the framework requires a no-op plan. (#85)

## 0.18.0

FEATURES:

* List resources: added Terraform 1.14 `list "geni_profile" "..." {}` and `list "geni_document" "..." {}` blocks so `terraform query` discovers existing managed profiles and uploaded documents on the Geni account without needing numeric IDs up-front. Each result carries an `identity = { id = "..." }` that pastes directly into a `import {}` block — closing the discover-then-import workflow. Backed by `/api/user/managed-profiles` and `/api/user/uploaded-documents`; results stream page-by-page through the existing rate-limited `*geni.Client`. Union list is deferred — the Geni API has no enumeration endpoint. (#82, #87)

IMPROVEMENTS:

* Profile and document: exposed `NewEmptyResourceModel()` constructors that return a model whose collection fields carry typed-null defaults matching the schema. Used by the list resources to seed a from-scratch model before `ValueFrom` runs; the managed-resource Read paths continue to seed via `req.State.Get`. (#87)

## 0.17.1

BUG FIXES:

* Import: plannable `import {}` blocks that pass the resource via the typed `identity` attribute (Terraform 1.12+ `import { identity = { id = "..." } to = ... }`) now work. The import-time validation introduced in v0.17.0 only inspected the legacy string ID, which is empty in plannable mode, so every plannable import failed with a spurious `<Resource> not found` diagnostic. The handler now reads from `req.Identity` when its `Raw` is non-null. (#84, #86)
* Profile: importing or refreshing a profile that the Geni API returns with flat top-level name fields (`first_name`, `last_name`, etc.) and an empty localized `names` map now hydrates an en-US entry from those flat fields. Previously the `names` attribute came back null and the next plan showed a spurious in-place update recreating the locale entry; this also unblocks plannable imports for the common single-locale case. (#86)

IMPROVEMENTS:

* Testing: modernized identity coverage on `geni_document`, `geni_profile`, and `geni_union` — every happy-path apply now asserts `statecheck.ExpectIdentity` + `statecheck.ExpectIdentityValueMatchesState`, and one acceptance step per resource exercises the Terraform 1.12+ plannable-import flow via `ImportStateKind: resource.ImportBlockWithResourceIdentity`. (#84, #86)
* Testing: added unit-test coverage for the new flat-name fallback (both `namesWithFlatFallback` directly and end-to-end through `ValueFrom`), in the repo's gomega/`t.Run` convention. (#86)
* Testing: pointed the project acceptance tests at writable sandbox projects (`project-8`, `project-9`); the previously hardcoded `project-6` is not accessible to the test account and was failing with `Access Denied`. Also switched the multi-project assertion from positional `AtSliceIndex` to `SetExact` since `projects` is a set. (#86)

## 0.17.0

BUG FIXES:

* Batch read: a `geni_document`, `geni_profile`, or `geni_union` that no longer exists on Geni is now removed from state cleanly when refreshed through the concurrent batch path. The bulk endpoints silently omit missing IDs from their response; the batch client now translates that omission into `ErrResourceNotFound` so the Read handlers' existing not-found branch runs `resp.State.RemoveResource`. Previously the synthesized error was generic, the resource stayed in state with all fields zeroed out, and every subsequent `terraform plan` failed refresh on the same row until it was manually `terraform state rm`'d. (#80)
* Import: `terraform import` for `geni_document`, `geni_profile`, and `geni_union` now validates the supplied ID against the Geni API before writing state. A non-existent ID surfaces a clear `<Resource> not found` diagnostic instead of succeeding silently and producing a zombie state row that fails every subsequent refresh. (#80)

BEHAVIORAL CHANGES:

* Import: each `terraform import` now performs one additional GET against the Geni API to verify the ID exists before falling through to the framework's standard import-then-refresh flow. When `use_document_cache` / `use_profile_cache` is enabled, the framework's follow-up Read is served from cache so only the validation GET hits the API; unions have no cache and always cost two GETs per import. The trade-off is intentional: it eliminates the silent-success path that previously left unrecoverable zombie state rows.

IMPROVEMENTS:

* Testing: added unit-test coverage for the batch client's response-to-request dispatch logic and for the new import validation paths, both in the repo's gomega/`t.Run` convention. Added acceptance tests that exercise the import-non-existent-id path against the live sandbox API for all three resources.
* Documentation: documented the unit-test conventions (gomega dot-import, `RegisterTestingT(t)` per sub-test, same-package tests) in `CLAUDE.md` so contributors find the pattern without reverse-engineering it.
* Maintenance: dependency updates (`labstack/echo/v4`, grouped `go_modules` bumps) and CI action updates (`goreleaser/goreleaser-action`, `actions/setup-go`).

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
