\connect dbconcerts

/*
 * Populates a database with some useful constants / types / etc
 */
INSERT INTO recording_types (recording_name) VALUES
  ('AVI'),
  ('DVD'),
  ('FLAC'),
  ('FLV'),
  ('HD'),
  ('MKV'),
  ('MP3'),
  ('MP4'),
  ('MPG'),
  ('OGG'),
  ('SHN'),
  ('VIDEO - OTHER'),
  ('WAV')
;

INSERT INTO source_types (source_name) VALUES
  ('Professional'),
  ('Amateur')
;
