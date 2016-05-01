# Live Recordings in Go

This repo is documenting going from 0 to 100 (miles, not miles per hour) in Go by using the live recordings site that I wrote in Python, which is need of a serious overhaul anyway, as a base.

## Requirements

- PostgreSQL database >= 9.5 (needed for upserts, aka `on conflict`)
- [pg driver](https://github.com/lib/pq) to interact with the database
- [httprouter](https://github.com/julienschmidt/httprouter) for handling routes

## Javascript Dependencies

Note: the following dependencies are bundled with the project.
Feel free to host them from a CDN!

- [Bootstrap 2.3.2](http://getbootstrap.com/2.3.2/)
- [jQuery 2.2.3](https://code.jquery.com/jquery-2.2.3.min.js)
- Sortable library from [RubaXa](https://github.com/RubaXa/Sortable)
