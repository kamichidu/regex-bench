## Pre-release validation: coregex v0.12.0-rc

Testing commit `61fe34e` from [coregex PR #112](https://github.com/coregx/coregex/pull/112) before release.

### Changes in coregex v0.12.0-rc

- Anti-quadratic guard for reverse searches
- DFA 4x loop unrolling
- Prefilter IsFast() gate
- DFA cache clear & continue
- OnePass capture limit fix (17 → 16)

### Purpose

CI-only validation. Do NOT merge until coregex v0.12.0 is tagged.
After tag: update go.mod to `v0.12.0`, then merge.
