## Scaffold

- [ ] Copy deployment template from `train-data` repo
- [ ] Restore tests


## Design

- [ ] Idempotency should be offloaded to db, check why `IdempotencyRepository` is needed.
- [ ] Remove `echo`, use plain net/http as much as possible
- [ ] Check why createdAt etc timestamps are set by domain and not db, which approach better for consistency


## Study

- [ ] `sqlc` generation
- [ ] Agentic Coding
- [ ] How to use test containers? Nested containers - if running tests within a container, what is the impact of running tests with test containers? 