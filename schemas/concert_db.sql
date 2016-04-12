\connect dbconcerts

CREATE TABLE artists (
  artist_id                     SERIAL PRIMARY KEY,
  artist_name                   varchar(256) NOT NULL CHECK (artist_name <> ''),
  short_name                    varchar(20) NOT NULL CHECK (short_name <> ''),
  UNIQUE (short_name)
);

CREATE TABLE songs (
  song_id                       SERIAL PRIMARY KEY,
  artist_id                     integer NOT NULL references artists(artist_id),
  title                         varchar(256) NOT NULL,
  lyrics                        text
);

CREATE TABLE location (
  location_id                   SERIAL PRIMARY KEY,
  city                          varchar(100) NOT NULL CHECK (city <> ''),
  state                         varchar(100) NOT NULL CHECK (state <> ''),
  country                       varchar(100) NOT NULL CHECK (country <> ''),
  UNIQUE (city, state, country)
);

CREATE TABLE venues (
  venue_id                      SERIAL PRIMARY KEY,
  venue_name                    varchar(256) NOT NULL CHECK (venue_name <> ''),
  location_id                   integer NOT NULL references location(location_id)
);

CREATE TABLE concerts (
  concert_id                    SERIAL PRIMARY KEY,
  artist_id                     integer NOT NULL references artists(artist_id),
  venue_id                      integer NOT NULL references venues(venue_id),
  date                          date,
  notes                         text
);

CREATE TABLE concert_setlist_mapping (
  concert_id                    integer NOT NULL references concerts(concert_id),
  songs                         integer ARRAY
);

CREATE TABLE recording_types (
  recording_type_id             SERIAL PRIMARY KEY,
  recording_name                varchar(50) NOT NULL
);

CREATE TABLE source_types (
  source_type_id                SERIAL PRIMARY KEY,
  source_name                   varchar(50)
);

CREATE TABLE recording (
  recording_id                  SERIAL PRIMARY KEY,
  recording_type                integer NOT NULL references recording_types(recording_type_id),
  source_type                   integer NOT NULL references source_types(source_type_id),
  taper                         varchar(100),
  length                        integer,
  notes                         text
);

CREATE TABLE concert_recording_mapping (
  recording_id                  integer NOT NULL references recording(recording_id),
  concert_id                    integer NOT NULL references concerts(concert_id),
  UNIQUE (recording_id, concert_id)
);
