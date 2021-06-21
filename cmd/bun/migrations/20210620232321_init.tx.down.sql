--
-- CLEAN UP/RESET
--
drop view if exists participant_point_total;
drop view if exists team_point_total;

drop table if exists season;
drop table if exists ledger_entry;
drop table if exists task_participant_restriction;
drop table if exists task;
drop table if exists task_collection;
drop table if exists task_type;
drop table if exists team;
drop table if exists league;
drop table if exists allowed_mod;
drop table if exists allowed_admin;

drop function if exists func_stamp_modified;

drop extension if exists btree_gist;
