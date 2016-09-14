CREATE TABLE events (
  id        TEXT PRIMARY KEY,
  timestamp TEXT,
  source    json
);

CREATE INDEX events_timestamp_index
  ON events (timestamp);

CREATE VIRTUAL TABLE events_fts USING fts4(id, source);
