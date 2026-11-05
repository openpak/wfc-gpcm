# Deploy: wfc-gpcm

Use example.env with the existing WFC database and shared WFC_TOKEN_KEY_FILE.
Never rotate token material or change profile IDs during ordinary service releases.
Dependencies: WFC_QR2_URL.

Link-first (prds#139): `WFC_LINK_FIRST=1` turns the gate on (default off: logins as before).
`WFC_CORE_ADDRESS`/`WFC_CORE_KEY` (optional) let the gate heal a half-done approval and add
`OpenPak: @handle` to Mario Kart Wii's message of the day. `WFC_LINK_DS_ERR` is the GP err value
sent to a refused DS (0 = generic login error). Deploy wfc-api first.

```sh
GOWORK=off go test -mod=vendor ./...
GOWORK=off go vet -mod=vendor ./...
docker build -t wfc-gpcm .
```

Shared wire types, errors, API clients and legacy database helpers are vendored
from wfc-common. The image and tests build from this repository alone. Refresh
vendor with the sibling library available after shared contract changes.
See [the full split guide](https://github.com/openpak/wfc-common/blob/main/docs/service-split.md) for full topology, failure/recovery behavior and
migration.
