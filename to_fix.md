# Issues to Fix

Found via CodeRabbit code review on 2026-02-07.

---

## Code Issues

### 1. Resource leak in `backend-go/pkg/massive/client.go:303-308`
`defer resp.Body.Close()` inside a loop keeps all response bodies open until the function returns. Replace with explicit `resp.Body.Close()` calls after each read/error path.

### 2. `NEXT_PUBLIC_*` env var ineffective at runtime (`docker-compose.yml:42-44`)
`NEXT_PUBLIC_API_URL` is inlined at build time by Next.js — setting it at runtime in docker-compose has no effect on client-side code. Pass it as a build argument instead.

### 3. `setState` called during render (`frontend/app/page.tsx:33-36`)
`setSelectedExpiration` is called directly during render. Wrap in `useEffect` with `[data, expirations, selectedExpiration]` as dependencies.

### 4. Orphaned process in `Makefile:155-162`
`make run_go_front` captures the `make` PID but not the actual Go server child — killing it may leave the server running on port 8080. Use `trap 'kill 0' EXIT` or run the Go binary directly.

---

## Accessibility

### 5. Missing ARIA attributes on dropdown (`frontend/components/ExpirationDropdown.tsx:39-58`)
Button needs `type="button"`, `aria-haspopup="listbox"`, and `aria-expanded={isOpen}`.

---

## Documentation

### 6. Timezone limitation buried in Notes (`current_plan.md:302`)
Move the "America-only" date assumption to the Known Issues section.

### 7. Data validation deprioritized (`current_plan.md:104-108`)
Validation is scheduled for Week 1-2 but the app is already displaying financial data — move to immediate priorities.
