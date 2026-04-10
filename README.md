# Task Tracker API – Recurring Tasks Extension

## 📌 Overview

This service is a Task Tracker API built with Go and Clean Architecture.
It is used in a medical information system where doctors manage daily operational tasks such as patient calls, rounds, and reporting.

This extension adds support for **recurring tasks**, allowing tasks to be automatically generated based on predefined recurrence rules.

---

## 🎯 Goal of the Feature

The system must allow users to define recurrence rules for tasks so that tasks are automatically created in the system based on those rules.

Supported scenarios:

- Daily recurring tasks (every N days)
- Monthly recurring tasks (specific days of month)
- Tasks on specific dates
- Tasks on even/odd days of month

---

## 🧠 Key Design Decision

### Separation of Concerns

We separate two concepts:

- **Task** → a single executable unit of work
- **RecurrenceRule** → defines how tasks are generated over time

### Why this approach?

- Keeps Task entity simple and stable
- Makes recurrence logic reusable and extensible
- Avoids duplication of scheduling logic inside Task

---

## 🧩 Domain Model

### Task

- ID
- Title
- Description
- Status
- DueDate
- RecurrenceRuleID (optional)

---

### RecurrenceRule

Defines how tasks should be generated.

Fields:

- ID
- Type:
  - daily
  - monthly
  - specific_dates
  - even_odd
- Interval (for daily recurrence)
- DaysOfMonth (for monthly recurrence)
- SpecificDates (list of exact dates)
- EvenOdd (even / odd rule)
- StartDate
- EndDate (optional)

---

## 🔄 Recurrence Types

### 1. Daily

Task is generated every N days based on interval.

Example:
- Every 1 day (daily)
- Every 2 days

---

### 2. Monthly

Task is generated on specific days of month.

Example:
- 1st, 15th, 30th of each month

---

### 3. Specific Dates

Task is generated only on predefined dates.

Example:
- 2026-04-10
- 2026-04-20

---

### 4. Even / Odd Days

Task is generated only on even or odd day numbers.

Example:
- Even: 2, 4, 6...
- Odd: 1, 3, 5...

---

## ⚙️ Task Generation Strategy

A **scheduler-based approach** is used.

A background process runs once per day and evaluates all recurrence rules.

### Flow:

1. Load all active recurrence rules
2. Check if each rule matches the current date
3. If matched → create a task
4. Ensure no duplicate tasks are created

---

## 🧠 Why Scheduler-Based Approach

We chose a scheduler instead of on-demand generation because:

- Predictable and consistent behavior
- Avoids side effects during API reads
- Easier to scale and debug
- Clean separation between API and background processing

---

## ⚠️ Assumptions

Since requirements are not fully specified, the following assumptions were made:

- Tasks are generated once per day
- System timezone is UTC
- Invalid calendar dates are ignored (e.g. Feb 30)
- Recurrence requires a valid start date
- End date limits task generation
- Each task is generated only once (idempotent behavior)

---

## 🚧 Edge Cases Considered

- Leap years (Feb 29 handling)
- Months with fewer than 31 days
- Duplicate task prevention
- Invalid recurrence configuration
- Timezone consistency

---

## 🔁 Idempotency

Task generation is idempotent.

This ensures that:

- Running the scheduler multiple times does not create duplicates
- Each task is uniquely tied to (rule + date)

---

## 🧪 Testing Strategy

- Unit tests for recurrence logic
- Tests per recurrence type
- Integration tests for task generation flow
- Edge case validation tests

---

## 📈 Future Improvements

- Add recurrence preview endpoint
- Add audit logging for task generation
- Support user-specific recurrence ownership
- Support cron-expression-based recurrence
- Replace scheduler with queue-based processing (e.g. worker system)

---

## 📌 Summary

This solution focuses on:

- Clean architecture compliance
- Maintainability and scalability
- Predictable background processing
- Extensible recurrence model

The design ensures that new recurrence types can be added without modifying core task logic.