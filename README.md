# marvel-forwarder

## Technical design & considerations

### Caching strategy

A DB is used to persist the characters, and a cronjob (called `bifrost`, the bridge in the Nine Realms) is used to periodically fetch the data from Marvel API.

### Alternatives & comparisons

- Instead of DB, use cache like redis: easier to set up, but if bifrost is down/has problems, the app won't be able to return any data
- Instead of cronjob, keep the last sync time and let user queries later than that time (say, by more than 1 hour) trigger the sync: the first queries will definitely be delayed (whereas cronjob is controlled), and it will clutter the logic
- Instead of cronjob, let an admin trigger the sync: feasible, can be an addition to the cronjob

### Caveats

- Since data is never deleted, all characters live on here even if Marvel deletes them :)
- If Marvel updates the data right when bifrost is running, be it add/delete/update, if it's at the "page" that bifrost has finished reading, it won't be able to catch it. (it's also worth noting that it won't break, because the result is pruned of duplicates before inserting to DB)
- I intentionally tried to avoid dependencies to see how far I can go with Go itself. I didn't have time to add the swagger docs, I hope it's ok, but if not, let me know

### ORM

- I intentionally avoided the use of ORM to try out sqlx, the code might look messier

## Running the app

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker compose](https://docs.docker.com/compose/install/)
- `cp .env.dev .env` then replace `<insert>` with the actual private key

### Running the API

- `make run`

### Running unit testing

- `make test`

### View test coverage

- `make coverage`
