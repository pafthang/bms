# BMS Web SDK

TypeScript client SDK generated from `bms/openapi/openapi.json`.

## Install

```bash
pnpm install
```

## Generate client

```bash
pnpm run generate
```

Generated files are written to `src/generated`.

## Build

```bash
pnpm run build
```

## Services split

Client is generated with service separation by OpenAPI tags, e.g.:

- `AuthService`
- `UsersService`
- `UserSettingsService`
- `WorkspacesService`
- `WorkspaceUsersService`
- `BookmarksService`
- `TagsService`
- `AdminService`

Re-export entry point:

- `src/index.ts`
