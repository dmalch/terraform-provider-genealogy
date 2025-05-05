## 0.13.2

IMPROVEMENTS:

* Updated batch processing functions to eliminate duplicate IDs using a hashset. This ensures optimized and accurate request handling for unions, profiles, and documents, simplifying the logic for single and multiple ID processing scenarios.

## 0.13.1

BUG FIXES:

* Fixed an issue where the `Profile` resource was not correctly handling computed subfields in the `current_residence` field.

## 0.13.0

FEATURES:

* Added support for `current_residence` field in the `Profile` resource.