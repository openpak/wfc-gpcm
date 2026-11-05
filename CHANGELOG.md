# Changelog

## v0.4.0 - 2026-11-05

- Link-first (prds#139), behind `WFC_LINK_FIRST=1`: an unlinked console is refused at login before a game profile is made, with its pairing code. Mario Kart Wii shows text with the code and link number (22100); other Wii titles show the 6-digit link number as the error code; the DS gets GP err `WFC_LINK_DS_ERR` (0 = the generic login error). A banned console is told so and never shown a code. A Wii on the DNS-exploit path logs in flagged unlinked (no address kept, no friends/status/profile), gets the payload, and is refused at `exlogin`.
- With `WFC_CORE_ADDRESS`/`WFC_CORE_KEY` the gate heals a half-done approval from the core, and Mario Kart Wii's message of the day ends with `OpenPak: @handle` for a linked console. The old pairing tail on the message of the day is gone.
- Internal op `kick-user` closes every session of a console (22101 unlinked, 22002 banned). The OpenPak ban lookup uses the console's link subject.
- Vendors wfc-common v0.8.0. Deploy after wfc-api v0.5.0.

## v0.3.0 - 2026-11-05

- Vendors wfc-common v0.7.0 (prds#131): users.email and users.unique_nick are derived and dropped at start-up, last addresses are blanked after 90 days. Deploy after wfc-api.

## v0.2.0 - 2026-11-04

- Mario Kart Wii: an unlinked console is shown a pairing code in the Wi-Fi menu message of the day (`wl:motd`), to enter at openpak.org/link. Linked consoles see the message unchanged. The message is capped to what the game's buffer takes, the code never cut. A database error means no code, never a failed login (wfc-common v0.5.0).

## v0.1.3 - 2026-11-03

- Internal `/internal/` polls are counted and summarised once a minute instead of logged per call (wfc-common v0.3.1). Parse-error and peer-message log lines now record byte counts instead of the raw message; gamestats no longer logs set-data values.

## v0.1.2 - 2026-11-03

- Every gpcm message now logs one line to stderr, `gpcm cmd <name> -> <result> <ms>ms` (verb and result only, never payload values) (results ok, error, denied, unknown). Internal HTTP logs as above through wfc-common v0.3.0 `httplog`.

## v0.1.1

- The session projection published to wfc-qr2 now carries the console user id, so qr2 can report
  online sessions to rewards.
