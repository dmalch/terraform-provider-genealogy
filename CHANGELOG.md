## 0.14.5 (Unreleased)

IMPROVEMENTS:

* Retry logic for Geni API requests has been improved to handle DNS resolution errors more effectively.

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
