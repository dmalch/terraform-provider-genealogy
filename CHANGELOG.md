## 0.15.3  (Unreleased)

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
