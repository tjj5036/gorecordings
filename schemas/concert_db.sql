\connect dbconcerts

CREATE TABLE artists (
  artist_id                     SERIAL PRIMARY KEY,
  artist_name                   varchar(256) NOT NULL CHECK (artist_name <> '')
);

CREATE TABLE songs (
  song_id                       SERIAL PRIMARY KEY,
  artist_id                     integer NOT NULL references artists(artist_id),
  title                         varchar(256) NOT NULL,
  lyrics                        text
);

CREATE TABLE venues (
  venue_id                      SERIAL PRIMARY KEY,
  venue_details                 json
);

CREATE TABLE concerts (
  concert_id                    SERIAL PRIMARY KEY,
  artist_id                     integer NOT NULL references artists(artist_id),
  venue_id                      integer NOT NULL references venues(venue_id),
  notes                         text
);

CREATE TABLE concert_setlist_mapping (
  concert_id                    integer NOT NULL references concerts(concert_id),
  songs                         integer ARRAY
);
