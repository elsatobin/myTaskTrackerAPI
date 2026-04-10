# Task Tracker API — Recurring Tasks

A task management REST API built in Go with clean architecture. This extension adds recurrence support to tasks.

## Running

```bash
docker-compose up
```

The API is available at `http://localhost:8081`. Swagger UI at `http://localhost:8081/swagger/`.

Migrations must be applied in order:

```
migrations/0001_create_tasks.up.sql
migrations/0002_add_recurrence_to_tasks.up.sql
```

---

## Recurrence Feature

### Design decisions and assumptions

**Recurrence is stored on the task itself (JSONB column), not as a separate table.**

The task requirements describe recurrence as a *setting* on a task, not a separate entity. Keeping it embedded avoids a join on every read, keeps the model simple, and is easy to extend. A separate `recurrence_rules` table would make sense if rules were shared across many tasks — that's not the case here.

**No background scheduler — occurrences are computed on demand.**

The task says "tasks in the tracker should be created taking these settings into account." I interpreted this as: the recurrence config is stored with the task, and the system can tell you *when* that task is scheduled to occur. A background process that auto-creates duplicate task rows every day would add operational complexity (idempotency keys, deduplication, cron management) without clear benefit for the stated use cases. The `/occurrences` endpoint gives the client full visibility into the schedule.

**`start_date` is required; `end_date` is optional.** Without a start date, daily/even-odd recurrences have no anchor point. Without an end date, the recurrence is open-ended (valid for many real-world cases like daily patient calls).

**Day-of-month values are capped at 30 (not 31).** The requirement explicitly states "from 1 to 30". Days that don't exist in a given month (e.g. day 30 in February) are silently skipped.

**Dates use `YYYY-MM-DD` format throughout the API.** Time-of-day is irrelevant for scheduling purposes.

---

## Recurrence types

| Type | Required fields | Description |
|---|---|---|
| `daily` | `interval` (>= 1), `start_date` | Every N days starting from start_date |
| `monthly` | `days_of_month` (1–30), `start_date` | On specific days of each month |
| `specific_dates` | `specific_dates` (list), `start_date` | Only on the listed dates |
| `even_odd` | `even_odd` ("even"/"odd"), `start_date` | Even or odd days of the month |

---

## API

### Create task with recurrence

```
POST /api/v1/tasks
```

```json
{
  "title": "Daily patient calls",
  "description": "Call all patients on the ward",
  "recurrence": {
    "type": "daily",
    "interval": 1,
    "start_date": "2026-04-01",
    "end_date": "2026-12-31"
  }
}
```

Other recurrence examples:

```json
{ "type": "monthly", "days_of_month": [1, 15], "start_date": "2026-04-01" }
{ "type": "specific_dates", "specific_dates": ["2026-04-10", "2026-04-20"], "start_date": "2026-04-01" }
{ "type": "even_odd", "even_odd": "even", "start_date": "2026-04-01" }
```

### Get scheduled occurrences

```
GET /api/v1/tasks/{id}/occurrences?from=2026-04-01&to=2026-04-30
```

Response:

```json
{
  "dates": ["2026-04-01", "2026-04-03", "2026-04-05", "..."]
}
```

Returns an empty array for tasks without recurrence.

### Other endpoints

```
GET    /api/v1/tasks
GET    /api/v1/tasks/{id}
PUT    /api/v1/tasks/{id}
DELETE /api/v1/tasks/{id}
```

Task statuses: `new`, `in_progress`, `done`.

---

## Edge cases handled

- Months shorter than the requested day (e.g. day 30 in February) — skipped
- `end_date` before `start_date` — rejected with 400
- `interval < 1` for daily — rejected with 400
- `even_odd` value other than "even"/"odd" — rejected with 400
- `specific_dates` recurrence with empty list — rejected with 400
- `to` before `from` in occurrences query — rejected with 400
- Tasks without recurrence return `"recurrence": null` and empty occurrences list
