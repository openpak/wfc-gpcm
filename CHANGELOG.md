# Changelog

## v0.2.0 - 2026-11-04

- Mario Kart Wii: an unlinked console is shown a pairing code in the Wi-Fi menu message of the day (`wl:motd`), to enter at openpak.org/link. Linked consoles see the message unchanged. The message is capped to what the game's buffer takes, the code never cut. A database error means no code, never a failed login (wfc-common v0.5.0).

## v0.1.3 - 2026-11-03

- Internal `/internal/` polls are counted and summarised once a minute instead of logged per call (wfc-common v0.3.1). Parse-error and peer-message log lines now record byte counts instead of the raw message; gamestats no longer logs set-data values.

## v0.1.2 - 2026-11-03

- Every gpcm message now logs one line to stderr, `gpcm cmd <name> -> <result> <ms>ms` (verb and result only, never payload values) (results ok, error, denied, unknown). Internal HTTP logs as above through wfc-common v0.3.0 `httplog`.

## v0.1.1

- The session projection published to wfc-qr2 now carries the console user id, so qr2 can report
  online sessions to rewards.
