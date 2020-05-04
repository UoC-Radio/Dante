-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.1
-- PostgreSQL version: 10.0
-- Project Site: pgmodeler.io
-- Model Author: ---

SET check_function_bodies = false;
-- ddl-end --


-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.
-- -- object: rastapank_db | type: DATABASE --
-- -- DROP DATABASE IF EXISTS rastapank_db;
-- CREATE DATABASE rastapank_db;
-- -- ddl-end --
-- 

-- object: radio | type: SCHEMA --
-- DROP SCHEMA IF EXISTS radio CASCADE;
CREATE SCHEMA radio;
-- ddl-end --
ALTER SCHEMA radio OWNER TO postgres;
-- ddl-end --

SET search_path TO pg_catalog,public,radio;
-- ddl-end --

-- object: radio.members | type: TABLE --
-- DROP TABLE IF EXISTS radio.members CASCADE;
CREATE TABLE radio.members(
	user_id integer NOT NULL,
	username varchar NOT NULL,
	real_name varchar NOT NULL,
	CONSTRAINT members_pk PRIMARY KEY (user_id),
	CONSTRAINT members_username_key UNIQUE (username)

);
-- ddl-end --
COMMENT ON TABLE radio.members IS 'The radio station''s members (people)';
-- ddl-end --
COMMENT ON COLUMN radio.members.user_id IS 'We leave this as an integer instead of serial so that it doesn''t get set automaticaly by the DB. The idea is to match this id with the uid on the station''s user registration/authentication backend. In our case that''s the uid used on LDAP.';
-- ddl-end --
COMMENT ON COLUMN radio.members.username IS 'The username used on the registration/authentication system';
-- ddl-end --
COMMENT ON COLUMN radio.members.real_name IS 'User''s real name';
-- ddl-end --
ALTER TABLE radio.members OWNER TO postgres;
-- ddl-end --

-- object: radio.week_days | type: TABLE --
-- DROP TABLE IF EXISTS radio.week_days CASCADE;
CREATE TABLE radio.week_days(
	id serial NOT NULL,
	label varchar,
	name varchar,
	CONSTRAINT week_days_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.week_days IS 'The 7 days of the week';
-- ddl-end --
ALTER TABLE radio.week_days OWNER TO postgres;
-- ddl-end --

INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Mon', E'Monday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Tue', E'Tuesday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Wed', E'Wednesday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Thu', E'Thursday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Fri', E'Friday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Sat', E'Saturday');
-- ddl-end --
INSERT INTO radio.week_days (id, label, name) VALUES (DEFAULT, E'Sun', E'Sunday');
-- ddl-end --

-- object: radio.shows | type: TABLE --
-- DROP TABLE IF EXISTS radio.shows CASCADE;
CREATE TABLE radio.shows(
	id serial NOT NULL,
	title varchar NOT NULL,
	description text,
	producer_nickname varchar NOT NULL,
	logo_filename varchar,
	active boolean NOT NULL DEFAULT true,
	last_aired timestamp with time zone,
	times_aired integer,
	CONSTRAINT shows_pk PRIMARY KEY (id),
	CONSTRAINT name_unique UNIQUE (title)

);
-- ddl-end --
COMMENT ON TABLE radio.shows IS 'Radio shows';
-- ddl-end --
COMMENT ON COLUMN radio.shows.title IS 'Show title';
-- ddl-end --
COMMENT ON COLUMN radio.shows.description IS 'Show''s description';
-- ddl-end --
COMMENT ON COLUMN radio.shows.producer_nickname IS 'How producers are referenced for their show - shown to the audience';
-- ddl-end --
COMMENT ON COLUMN radio.shows.logo_filename IS 'Filepath of the show''s logo image (optional)';
-- ddl-end --
COMMENT ON COLUMN radio.shows.active IS 'Active shows are shows that are still aired even out-of-schedule (e.g. on a per-case basis)';
-- ddl-end --
COMMENT ON COLUMN radio.shows.last_aired IS 'Last time the show aired';
-- ddl-end --
COMMENT ON COLUMN radio.shows.times_aired IS 'How many times the show aired';
-- ddl-end --
COMMENT ON CONSTRAINT name_unique ON radio.shows  IS 'Show name is unique';
-- ddl-end --
ALTER TABLE radio.shows OWNER TO postgres;
-- ddl-end --

-- object: radio.show_producers | type: TABLE --
-- DROP TABLE IF EXISTS radio.show_producers CASCADE;
CREATE TABLE radio.show_producers(
	user_id_members integer NOT NULL,
	id_shows integer NOT NULL,
	CONSTRAINT show_producers_pk PRIMARY KEY (user_id_members,id_shows)

);
-- ddl-end --
COMMENT ON TABLE radio.show_producers IS 'Each show has one or more producers, registered with their user ids

Constraints:
If a show is removed, remove all its producer entries
Don''t allow removing a member that''s still associated with a show';
-- ddl-end --

-- object: members_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_producers DROP CONSTRAINT IF EXISTS members_fk CASCADE;
ALTER TABLE radio.show_producers ADD CONSTRAINT members_fk FOREIGN KEY (user_id_members)
REFERENCES radio.members (user_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_producers DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.show_producers ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.show_urls | type: TABLE --
-- DROP TABLE IF EXISTS radio.show_urls CASCADE;
CREATE TABLE radio.show_urls(
	id serial NOT NULL,
	id_shows integer,
	url_uri varchar NOT NULL,
	url_text varchar,
	CONSTRAINT show_urls_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.show_urls IS 'Each show may have URLs associated with it (e.g. blog, facebook page etc)

Constraints:
If the show is removed, remove all URLs associated with it
';
-- ddl-end --
COMMENT ON COLUMN radio.show_urls.url_uri IS 'The url';
-- ddl-end --
COMMENT ON COLUMN radio.show_urls.url_text IS 'Text to be displayed on the link';
-- ddl-end --
ALTER TABLE radio.show_urls OWNER TO postgres;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_urls DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.show_urls ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.show_weekdays | type: TABLE --
-- DROP TABLE IF EXISTS radio.show_weekdays CASCADE;
CREATE TABLE radio.show_weekdays(
	id serial NOT NULL,
	id_shows integer,
	start_time time with time zone NOT NULL,
	duration interval MINUTE  NOT NULL,
	id_week_days integer,
	CONSTRAINT show_schedule_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.show_weekdays IS 'Shows aired on a weekly basis

Constraints:
If the show is removed, remove all shedule NOT NULL references to it
Don''t allow deleting a day (shouldn''t happen anyway) if it contains scheduled shows
';
-- ddl-end --
ALTER TABLE radio.show_weekdays OWNER TO postgres;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_weekdays DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.show_weekdays ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: week_days_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_weekdays DROP CONSTRAINT IF EXISTS week_days_fk CASCADE;
ALTER TABLE radio.show_weekdays ADD CONSTRAINT week_days_fk FOREIGN KEY (id_week_days)
REFERENCES radio.week_days (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.show_oneshot | type: TABLE --
-- DROP TABLE IF EXISTS radio.show_oneshot CASCADE;
CREATE TABLE radio.show_oneshot(
	id serial NOT NULL,
	id_shows integer,
	scheduled_time timestamp with time zone NOT NULL,
	duration interval MINUTE  NOT NULL,
	CONSTRAINT shows_oneshot_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.show_oneshot IS 'Shows aired once

Constrints:
If the show is removed, remove all one-shot schedule entries associated with it';
-- ddl-end --
COMMENT ON COLUMN radio.show_oneshot.scheduled_time IS 'Can be in the future';
-- ddl-end --
ALTER TABLE radio.show_oneshot OWNER TO postgres;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_oneshot DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.show_oneshot ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.show_messages | type: TABLE --
-- DROP TABLE IF EXISTS radio.show_messages CASCADE;
CREATE TABLE radio.show_messages(
	id serial NOT NULL,
	id_shows integer,
	received_datetime timestamp with time zone DEFAULT now(),
	user_agent varchar NOT NULL,
	ip_addr inet NOT NULL,
	nickname varchar NOT NULL,
	message text NOT NULL,
	CONSTRAINT show_messages_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.show_messages IS 'Listener messages
Notes: These are private messages sent from listeners to the show''s producers (all of them). Producers can access those from the dashboard application. Also this is the only table that uses data from outsiders.

Constraints:
If the show is removed, remove all messages associated with it';
-- ddl-end --
COMMENT ON COLUMN radio.show_messages.user_agent IS 'We want this to distinguish listeners using a browser or our mobile app';
-- ddl-end --
COMMENT ON COLUMN radio.show_messages.ip_addr IS 'Listener''s IPv4 address';
-- ddl-end --
COMMENT ON COLUMN radio.show_messages.nickname IS 'Listener''s nickname';
-- ddl-end --
ALTER TABLE radio.show_messages OWNER TO postgres;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.show_messages DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.show_messages ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.playlist_types | type: TABLE --
-- DROP TABLE IF EXISTS radio.playlist_types CASCADE;
CREATE TABLE radio.playlist_types(
	id integer NOT NULL,
	label varchar NOT NULL,
	intermediate bool DEFAULT false,
	remote bool DEFAULT true,
	CONSTRAINT playlist_types_pk PRIMARY KEY (id)

);
-- ddl-end --
ALTER TABLE radio.playlist_types OWNER TO postgres;
-- ddl-end --

INSERT INTO radio.playlist_types (id, label, intermediate, remote) VALUES (E'1', E'main', DEFAULT, DEFAULT);
-- ddl-end --
INSERT INTO radio.playlist_types (id, label, intermediate, remote) VALUES (E'2', E'intermediate', E'true', DEFAULT);
-- ddl-end --
INSERT INTO radio.playlist_types (id, label, intermediate, remote) VALUES (E'3', E'fallback', DEFAULT, E'false');
-- ddl-end --

-- object: radio.playlists | type: TABLE --
-- DROP TABLE IF EXISTS radio.playlists CASCADE;
CREATE TABLE radio.playlists(
	id serial NOT NULL,
	title varchar NOT NULL,
	file_path varchar NOT NULL,
	fade_in_secs integer DEFAULT 2,
	fade_out_secs integer DEFAULT 2,
	min_level numeric DEFAULT 0.0,
	max_level numeric DEFAULT 1.0,
	shuffle boolean DEFAULT true,
	id_playlist_types integer,
	description text,
	comments text,
	CONSTRAINT playlists_pk PRIMARY KEY (id),
	CONSTRAINT unique_title UNIQUE (title)

);
-- ddl-end --
COMMENT ON TABLE radio.playlists IS 'Playlists


Constraint: Don''t allow removing a playlist type (shouldn''t happen anyway) if there are playlists associated with it';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.title IS 'Short title e.g. "70s Funk"';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.file_path IS 'Path to a .pls or .m3u file';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.fade_in_secs IS 'Duration of fade in in secs (zero for no fade-in)';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.fade_out_secs IS 'Duration of fade out in secs (zero for no fade-out)';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.min_level IS 'Fader min level';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.max_level IS 'Fader max level';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.shuffle IS 'Shuffle songs or not';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.description IS 'Short description text (optional)';
-- ddl-end --
COMMENT ON COLUMN radio.playlists.comments IS 'Comments text (optional, internal)';
-- ddl-end --
ALTER TABLE radio.playlists OWNER TO postgres;
-- ddl-end --

-- object: playlist_types_fk | type: CONSTRAINT --
-- ALTER TABLE radio.playlists DROP CONSTRAINT IF EXISTS playlist_types_fk CASCADE;
ALTER TABLE radio.playlists ADD CONSTRAINT playlist_types_fk FOREIGN KEY (id_playlist_types)
REFERENCES radio.playlist_types (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.playlist_maintainers | type: TABLE --
-- DROP TABLE IF EXISTS radio.playlist_maintainers CASCADE;
CREATE TABLE radio.playlist_maintainers(
	user_id_members integer NOT NULL,
	id_playlists integer NOT NULL,
	CONSTRAINT playlist_maintainers_pk PRIMARY KEY (user_id_members,id_playlists)

);
-- ddl-end --
COMMENT ON TABLE radio.playlist_maintainers IS 'Playlist maintainers

Constraints:
Don''t allow deletion of a member that is referenced as a playlist maintainer
If a playlist is removed, all its maintainer entries';
-- ddl-end --

-- object: members_fk | type: CONSTRAINT --
-- ALTER TABLE radio.playlist_maintainers DROP CONSTRAINT IF EXISTS members_fk CASCADE;
ALTER TABLE radio.playlist_maintainers ADD CONSTRAINT members_fk FOREIGN KEY (user_id_members)
REFERENCES radio.members (user_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.playlist_maintainers DROP CONSTRAINT IF EXISTS playlists_fk CASCADE;
ALTER TABLE radio.playlist_maintainers ADD CONSTRAINT playlists_fk FOREIGN KEY (id_playlists)
REFERENCES radio.playlists (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.combo_playlists | type: TABLE --
-- DROP TABLE IF EXISTS radio.combo_playlists CASCADE;
CREATE TABLE radio.combo_playlists(
	id serial NOT NULL,
	date_created timestamp with time zone DEFAULT now(),
	date_modified timestamp with time zone,
	title character varying NOT NULL,
	description text,
	comments text,
	CONSTRAINT combo_playlists_pkey PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.combo_playlists IS 'We use the term zone to represent a musical context, think of it as a "mood".';
-- ddl-end --
COMMENT ON COLUMN radio.combo_playlists.title IS 'Short title, e.g. "Waking up"';
-- ddl-end --
COMMENT ON COLUMN radio.combo_playlists.description IS 'Short description text (optional)';
-- ddl-end --
COMMENT ON COLUMN radio.combo_playlists.comments IS 'Comments text (optional, internal)';
-- ddl-end --
ALTER TABLE radio.combo_playlists OWNER TO postgres;
-- ddl-end --

-- object: radio.day_zones | type: TABLE --
-- DROP TABLE IF EXISTS radio.day_zones CASCADE;
CREATE TABLE radio.day_zones(
	id serial NOT NULL,
	time_start time NOT NULL,
	duration interval MINUTE  NOT NULL,
	id_week_days integer NOT NULL,
	id_combo_playlists integer NOT NULL,
	CONSTRAINT day_zones_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.day_zones IS 'Each day is split into zones to represent different moods within the day.

Constraints:
Don''t allow removing a day (this shouldn''t happen anyway) if it still has zones scheduled
If a zone is removed, remove all schedule entries referring to it';
-- ddl-end --

-- object: combo_playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.day_zones DROP CONSTRAINT IF EXISTS combo_playlists_fk CASCADE;
ALTER TABLE radio.day_zones ADD CONSTRAINT combo_playlists_fk FOREIGN KEY (id_combo_playlists)
REFERENCES radio.combo_playlists (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: week_days_fk | type: CONSTRAINT --
-- ALTER TABLE radio.day_zones DROP CONSTRAINT IF EXISTS week_days_fk CASCADE;
ALTER TABLE radio.day_zones ADD CONSTRAINT week_days_fk FOREIGN KEY (id_week_days)
REFERENCES radio.week_days (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.zone_playlists | type: TABLE --
-- DROP TABLE IF EXISTS radio.zone_playlists CASCADE;
CREATE TABLE radio.zone_playlists(
	sched_weight numeric NOT NULL,
	id_playlists integer NOT NULL,
	id_combo_playlists integer NOT NULL,
	CONSTRAINT zone_playlists_pk PRIMARY KEY (id_combo_playlists,id_playlists),
	CONSTRAINT check_sched_weight CHECK (sched_weight > 0.0 AND sched_weight <= 1.0)

);
-- ddl-end --
COMMENT ON TABLE radio.zone_playlists IS 'Table for main playlists. Each zone may contain various playlists, it must include at least one main playlist and optionaly a fallback playlist and various intermediate ones.

Constraints:
If a zone is removed, remove all its playlist entries
Don''t allow removing a playlist still associated with a zone';
-- ddl-end --
COMMENT ON COLUMN radio.zone_playlists.sched_weight IS 'Scheduling weight for ''main'' playlists. Must be  0 < weight <= 1 and the sum of all ''main'' playlist weights on a zone must be 1';
-- ddl-end --

-- object: combo_playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.zone_playlists DROP CONSTRAINT IF EXISTS combo_playlists_fk CASCADE;
ALTER TABLE radio.zone_playlists ADD CONSTRAINT combo_playlists_fk FOREIGN KEY (id_combo_playlists)
REFERENCES radio.combo_playlists (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.zone_playlists DROP CONSTRAINT IF EXISTS playlists_fk CASCADE;
ALTER TABLE radio.zone_playlists ADD CONSTRAINT playlists_fk FOREIGN KEY (id_playlists)
REFERENCES radio.playlists (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.combo_playlist_intermediate | type: TABLE --
-- DROP TABLE IF EXISTS radio.combo_playlist_intermediate CASCADE;
CREATE TABLE radio.combo_playlist_intermediate(
	id_playlists integer NOT NULL,
	sched_interval_mins integer NOT NULL,
	sceduled_items_cardinality integer NOT NULL,
	id_combo_playlists integer NOT NULL,
	CONSTRAINT check_interval CHECK (sched_interval_mins > 0),
	CONSTRAINT check_cardinality CHECK (sceduled_items_cardinality >= 0),
	CONSTRAINT combo_playlist_intermediate_pk PRIMARY KEY (id_playlists,id_combo_playlists)

);
-- ddl-end --
COMMENT ON TABLE radio.combo_playlist_intermediate IS 'Table for intermediate playlists. Each zone may contain various playlists, it must include at least one main playlist and optionaly a fallback playlist and various intermediate ones.

Constraints:
If a zone is removed, remove all its intermediate playlist entries
Don''t allow removing an intermediate playlist still associated with a zone';
-- ddl-end --
COMMENT ON COLUMN radio.combo_playlist_intermediate.sched_interval_mins IS 'Scheduling interval for intermediate playlists in mins';
-- ddl-end --
COMMENT ON COLUMN radio.combo_playlist_intermediate.sceduled_items_cardinality IS 'Number of items to shedule each time zero is a special case in which we don''t schedule an intermediate playlist based on time but from a "hint" encoded in the main playlist';
-- ddl-end --

-- object: playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_intermediate DROP CONSTRAINT IF EXISTS playlists_fk CASCADE;
ALTER TABLE radio.combo_playlist_intermediate ADD CONSTRAINT playlists_fk FOREIGN KEY (id_playlists)
REFERENCES radio.playlists (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: combo_playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_intermediate DROP CONSTRAINT IF EXISTS combo_playlists_fk CASCADE;
ALTER TABLE radio.combo_playlist_intermediate ADD CONSTRAINT combo_playlists_fk FOREIGN KEY (id_combo_playlists)
REFERENCES radio.combo_playlists (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.combo_playlist_fallback | type: TABLE --
-- DROP TABLE IF EXISTS radio.combo_playlist_fallback CASCADE;
CREATE TABLE radio.combo_playlist_fallback(
	id_playlists integer NOT NULL,
	id_combo_playlists integer NOT NULL,
	CONSTRAINT combo_playlist_fallback_pk PRIMARY KEY (id_combo_playlists,id_playlists)

);
-- ddl-end --
COMMENT ON TABLE radio.combo_playlist_fallback IS 'Table for fallback playlists. Each zone may contain various playlists, it must include at least one main playlist and optionaly a fallback playlist and various intermediate ones.

Constraints:
If a zone is removed, remove all its falback playlist entries
Don''t allow removing a fallback playlist still associated with a zone';
-- ddl-end --

-- object: combo_playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_fallback DROP CONSTRAINT IF EXISTS combo_playlists_fk CASCADE;
ALTER TABLE radio.combo_playlist_fallback ADD CONSTRAINT combo_playlists_fk FOREIGN KEY (id_combo_playlists)
REFERENCES radio.combo_playlists (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_fallback DROP CONSTRAINT IF EXISTS playlists_fk CASCADE;
ALTER TABLE radio.combo_playlist_fallback ADD CONSTRAINT playlists_fk FOREIGN KEY (id_playlists)
REFERENCES radio.playlists (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.combo_playlist_maintainers | type: TABLE --
-- DROP TABLE IF EXISTS radio.combo_playlist_maintainers CASCADE;
CREATE TABLE radio.combo_playlist_maintainers(
	user_id_members integer NOT NULL,
	id_combo_playlists integer NOT NULL,
	CONSTRAINT combo_playlist_maintainers_pk PRIMARY KEY (user_id_members,id_combo_playlists)

);
-- ddl-end --
COMMENT ON TABLE radio.combo_playlist_maintainers IS 'Zone maintainers

Constraints:
Don''t allow removing a member that is referenced as a zone maintainer
If a zone is removed, remove all its maintainer entries';
-- ddl-end --

-- object: members_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_maintainers DROP CONSTRAINT IF EXISTS members_fk CASCADE;
ALTER TABLE radio.combo_playlist_maintainers ADD CONSTRAINT members_fk FOREIGN KEY (user_id_members)
REFERENCES radio.members (user_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: combo_playlists_fk | type: CONSTRAINT --
-- ALTER TABLE radio.combo_playlist_maintainers DROP CONSTRAINT IF EXISTS combo_playlists_fk CASCADE;
ALTER TABLE radio.combo_playlist_maintainers ADD CONSTRAINT combo_playlists_fk FOREIGN KEY (id_combo_playlists)
REFERENCES radio.combo_playlists (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: radio.create_show | type: FUNCTION --
-- DROP FUNCTION IF EXISTS radio.create_show(IN varchar,IN text,IN varchar,IN anyarray) CASCADE;
CREATE FUNCTION radio.create_show (IN name varchar, IN description text, IN nickname varchar, IN producers anyarray)
	RETURNS smallint
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
insert into radio.shows(title, producer_nickname) values (name, nickname) returning id into show_id;
	FOREACH x IN ARRAY producers
  	LOOP
    	raise notice 'Adding % to show %', x, name;
		insert into radio.show_producers (user_id_members, id_shows) values (x, show_id);
END LOOP;
$$;
-- ddl-end --
ALTER FUNCTION radio.create_show(IN varchar,IN text,IN varchar,IN anyarray) OWNER TO postgres;
-- ddl-end --

-- object: radio.shows_log | type: TABLE --
-- DROP TABLE IF EXISTS radio.shows_log CASCADE;
CREATE TABLE radio.shows_log(
	id serial NOT NULL,
	id_shows integer,
	start_time timestamp with time zone NOT NULL,
	end_time timestamp with time zone NOT NULL,
	recording_path varchar,
	playlist varchar[],
	commnents varchar,
	CONSTRAINT shows_log_pk PRIMARY KEY (id)

);
-- ddl-end --
COMMENT ON TABLE radio.shows_log IS 'A log for each show that was aired';
-- ddl-end --
ALTER TABLE radio.shows_log OWNER TO postgres;
-- ddl-end --

-- object: shows_fk | type: CONSTRAINT --
-- ALTER TABLE radio.shows_log DROP CONSTRAINT IF EXISTS shows_fk CASCADE;
ALTER TABLE radio.shows_log ADD CONSTRAINT shows_fk FOREIGN KEY (id_shows)
REFERENCES radio.shows (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


