# Wii/DS GPCM Server

OpenPak's Wii/DS GPCM Server.

| | |
| --- | --- |
| Network | Wii, DS |
| Role | GameSpy login, player sessions, friend state, presence and leased QR2 login projections |
| Depends on | [`wfc-qr2`](https://github.com/openpak/wfc-qr2) (`WFC_QR2_URL`), the existing WFC database |
| Image | `ghcr.io/openpak/wfc-gpcm` |
| Based on | WiiLink WFC (AGPL-3.0) |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| Internal HTTP (`WFC_HTTP_ADDR`, default `:20200`) | `GET /healthz` (unauthenticated); protocol and internal paths need `WFC_INTERNAL_KEY` |
| GameSpy login | Login, player sessions, friend state and presence |

## Build

```sh
docker build -t wfc-gpcm .
```

## Run

```sh
cp example.env .env   # fill in the existing WFC values
```

Build, test and environment detail is in [Deploy](docs/deploy.md).

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
