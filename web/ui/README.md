# BMS Web UI

Svelte 5 + SvelteKit frontend for BMS.

## Stack

- Svelte 5
- SvelteKit
- Tailwind CSS (v4, Vite plugin)
- shadcn-svelte
- `@pafthang/bms-sdk` from workspace

## Commands

From repository root:

```bash
pnpm --filter @pafthang/bms-ui dev
pnpm --filter @pafthang/bms-ui check
pnpm --filter @pafthang/bms-ui build
```

## Environment

Use `PUBLIC_BMS_API_BASE` to point UI to backend base URL.

Example:

```bash
PUBLIC_BMS_API_BASE=http://localhost:8080
```
