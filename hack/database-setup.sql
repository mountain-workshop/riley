DROP VIEW IF EXISTS team_total;
DROP TABLE IF EXISTS points_log;
DROP TABLE IF EXISTS team;

CREATE TABLE team (
  discord_role_id bigint PRIMARY KEY not null,
  team_name text not null,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP not null,
  modified_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP not null
);

CREATE TABLE points_log (
  points_log_id serial PRIMARY KEY not null,
  user_id bigint not null,
  team_discord_role_id bigint REFERENCES team (discord_role_id) not null,
  points integer not null,
  effective_datetime timestamp with time zone DEFAULT CURRENT_TIMESTAMP not null,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP not null,
  modified_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP not null
);

CREATE OR REPLACE FUNCTION log_modification() RETURNS trigger AS $log_modification$
    BEGIN
        NEW.modified_at := current_timestamp;
        RETURN NEW;
    END;
$log_modification$ LANGUAGE plpgsql;

CREATE TRIGGER log_modification BEFORE UPDATE ON team
	FOR EACH ROW EXECUTE FUNCTION log_modification();

CREATE TRIGGER log_modification BEFORE UPDATE ON points_log
	FOR EACH ROW EXECUTE FUNCTION log_modification();
    
CREATE VIEW team_total AS
select
	t.discord_role_id
    , t.team_name
    , sum(pl.points) as TOTAL_POINTS
from points_log pl
	join team t on (pl.team_discord_role_id = t.discord_role_id)
group by t.discord_role_id
;
