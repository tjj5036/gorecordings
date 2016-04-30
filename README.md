# Live Recordings in Go

This repo is documenting going from 0 to 100 (miles, not miles per hour) in Go by using the live recordings site that I wrote in Python, which is need of a serious overhaul anyway, as a base.

## Requirements

- PostgreSQL database >= 9.5 (needed for upserts, aka `on conflict`)
- [pg driver](https://github.com/lib/pq) to interact with the database
- [httprouter](https://github.com/julienschmidt/httprouter) for handling routes
