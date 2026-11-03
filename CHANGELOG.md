# Changelog

## v0.1.2 - 2026-11-03

- Every gpcm message now logs one line to stderr, `gpcm cmd <name> -> <result> <ms>ms` (verb and result only, never payload values) (results ok, error, denied, unknown). Internal HTTP logs as above through wfc-common v0.3.0 `httplog`.

## v0.1.1

- The session projection published to wfc-qr2 now carries the console user id, so qr2 can report
  online sessions to rewards.
