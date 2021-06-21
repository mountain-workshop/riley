--
-- EXTENSIONS
--
create extension if not exists btree_gist;

--
-- FUNCTIONS
--
create function func_stamp_modified() returns trigger as $$
	begin
    	NEW.modified_at := now();
        return NEW;
    end;
$$ language plpgsql;

--
-- TABLES AND TRIGGERS
--
---- allowed_admin
create table allowed_admin (
    discord_guild_id bigint not null
    , discord_role_id bigint not null
    , created_at timestamptz default now() not null
    , modified_at timestamptz default now() not null
    
    , primary key (discord_guild_id, discord_role_id)
);

create trigger trig_allowedadmin_stampmodified before update on allowed_admin
    for each row execute function func_stamp_modified();

--
---- allowed_mod
create table allowed_mod (
    discord_guild_id bigint not null
    , discord_role_id bigint not null
    , created_at timestamptz default now() not null
    , modified_at timestamptz default now() not null
    
    , primary key (discord_guild_id, discord_role_id)
);

create trigger trig_allowedmod_stampmodified before update on allowed_mod
    for each row execute function func_stamp_modified();

--
---- league
create table league (
    discord_guild_id bigint not null
    , league_id int generated always as identity
    , friendly_name text not null
    , created_at timestamptz default now() not null
    , modified_at timestamptz default now() not null
    
    , primary key (discord_guild_id, league_id)
);

create trigger trig_league_stampmodified before update on league
    for each row execute function func_stamp_modified();

--
---- team
create table team (
    discord_guild_id bigint not null
    , team_id int generated always as identity
	, discord_role_id bigint not null unique
	, league_id int
 	, created_at timestamptz default now() not null
 	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, team_id)
  	, foreign key (discord_guild_id, league_id) references league (discord_guild_id, league_id)
);

create trigger trig_team_stampmodified before update on team
	for each row execute function func_stamp_modified();

--
---- task_type
create table task_type (
  	discord_guild_id bigint not null
  	, task_type_id int generated always as identity
  	, friendly_name varchar(255) not null
	, parent_task_type_id int
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, task_type_id)
  	, foreign key (discord_guild_id, parent_task_type_id) references task_type (discord_guild_id, task_type_id)
);

create trigger trig_taskclassification_stampmodified before update on task_type
	for each row execute function func_stamp_modified();

--
---- task_collection
create table task_collection (
    discord_guild_id bigint not null
  	, task_collection_id int generated always as identity
  	, friendly_name varchar(255) not null
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, task_collection_id)
	, unique (discord_guild_id, friendly_name)
);

create trigger trig_taskset_stampmodified before update on task_collection
	for each row execute function func_stamp_modified();

--
---- task
create table task (
    discord_guild_id bigint not null
	, task_id bigint generated always as identity
  	, friendly_name varchar(255) not null
  	, point_base int not null check (point_base > 0)
  	, point_bonus int not null check (point_bonus >= 0)
  	, open_datetime timestamptz default now() not null
  	, close_datetime timestamptz
  	, per_participant_limit int check (per_participant_limit is null or per_participant_limit > 0)
  	, task_type_id int
  	, task_collection_id int
	, additional_info json
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, task_id)
  	, foreign key (discord_guild_id, task_type_id) references task_type (discord_guild_id, task_type_id)
  	, foreign key (discord_guild_id, task_collection_id) references task_collection (discord_guild_id ,task_collection_id)
	, unique (discord_guild_id, friendly_name)
);

create trigger trig_task_stampmodified before update on task
	for each row execute function func_stamp_modified();

--
---- task_participant_restriction
create table task_participant_restriction (
    discord_guild_id bigint not null
	, task_participant_restriction_id int generated always as identity
  	, task_id bigint not null
  	, discord_user_id bigint not null
  	, expiration_datetime timestamptz
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, task_participant_restriction_id)
  	, foreign key (discord_guild_id, task_id) references task (discord_guild_id, task_id)
	, unique (task_id, discord_user_id)
);

create trigger trig_taskparticipantrestriction_stampmodified before update on task_participant_restriction
	for each row execute function func_stamp_modified();

--
---- ledger_entry
create table ledger_entry (
    discord_guild_id bigint not null
  	, ledger_entry_id bigint generated always as identity
  	, effective_datetime timestamptz default now() not null
  	, discord_user_id bigint
	, team_id int
  	, task_id bigint not null
  	, diff_point_base int not null
  	, diff_point_bonus int not null
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, ledger_entry_id)
  	, foreign key (discord_guild_id, team_id) references team (discord_guild_id, team_id)
  	, foreign key (discord_guild_id, task_id) references task (discord_guild_id, task_id)

  	, constraint attributable check (discord_user_id is not null or team_id is not null)
);

-- TODO: add trigger that checks task.per_participant_limit?

create trigger trig_ledgerentry_stampmodified before update on ledger_entry
	for each row execute function func_stamp_modified();

--
---- season
create table season (
    discord_guild_id bigint not null
  	, season_id int generated always as identity
  	, friendly_name varchar(255) not null
	, datetime_range tstzrange not null
  	, created_at timestamptz default now() not null
  	, modified_at timestamptz default now() not null
  
  	, primary key (discord_guild_id, season_id)
	, exclude using gist (discord_guild_id with =, datetime_range with &&)
);

create trigger trig_season_stampmodified before update on season
	for each row execute function func_stamp_modified();

--
-- VIEWS
--
create view team_point_total as
select
	t.discord_guild_id
	, t.discord_role_id
    , sum(le.diff_point_base) + sum(le.diff_point_bonus) as POINT_TOTAL
from ledger_entry le
    join team t on (le.discord_guild_id = t.discord_guild_id and le.team_id = t.team_id)
group by
	t.discord_guild_id
	, t.discord_role_id
;

create view participant_point_total as
select
	le.discord_guild_id
	, le.discord_user_id
    , sum(le.diff_point_base) + sum(le.diff_point_bonus) as POINT_TOTAL
from ledger_entry le
group by
	le.discord_guild_id
	, le.discord_user_id
;
