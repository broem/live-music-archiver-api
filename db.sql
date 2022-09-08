-- DROP SCHEMA event;

CREATE SCHEMA event AUTHORIZATION postgres;

-- builders definition

-- Drop table

-- DROP TABLE builders;

CREATE TABLE builders (
	builder_map jsonb NULL,
	user_id varchar NOT NULL
);

-- captured definition

-- Drop table

-- DROP TABLE captured;

CREATE TABLE captured (
	user_id varchar NULL,
	map_id uuid NULL,
	title varchar NULL,
	description varchar NULL,
	description_url varchar NULL,
	"date" varchar NULL,
	"time" varchar NULL,
	venue varchar NULL,
	venue_address varchar NULL,
	venue_contact_info varchar NULL,
	ticket_cost varchar NULL,
	ticket_url varchar NULL,
	other_performers varchar NULL,
	age_required varchar NULL,
	facebook_url varchar NULL,
	twitter_url varchar NULL,
	instagram_url varchar NULL,
	misc varchar NULL,
	images _varchar NULL,
	approved bool NULL,
	captured_id uuid NOT NULL,
	url varchar NULL,
	capture_date timestamptz NULL,
	cbsa varchar NULL,
	state_fips varchar NULL,
	county_fips varchar NULL,
	CONSTRAINT captured_pk PRIMARY KEY (captured_id)
);
-- mappers definition

-- Drop table

-- DROP TABLE mappers;

CREATE TABLE mappers (
	map_id uuid NOT NULL,
	title_selector varchar NULL,
	description_selector varchar NULL,
	description_url_selector varchar NULL,
	date_selector varchar NULL,
	time_selector varchar NULL,
	venue_name_selector varchar NULL,
	venue_address_selector varchar NULL,
	venue_contact_info_selector varchar NULL,
	ticket_cost_selector varchar NULL,
	ticket_url_selector varchar NULL,
	other_performers_selector varchar NULL,
	age_required_selector varchar NULL,
	facebook_url_selector varchar NULL,
	twitter_url_selector varchar NULL,
	instagram_url varchar NULL,
	misc_selector varchar NULL,
	images_selector _varchar NULL,
	user_id varchar NULL,
	venue_base_url varchar NULL,
	full_event_selector varchar NULL,
	approved bool NOT NULL DEFAULT false,
	cbsa varchar NULL,
	state_fips varchar NULL,
	county_fips varchar NULL,
	CONSTRAINT mapper_pk PRIMARY KEY (map_id)
);
-- runners definition

-- Drop table

-- DROP TABLE runners;

CREATE TABLE runners (
	map_id uuid NOT NULL,
	user_id varchar NOT NULL,
	chron numeric NULL,
	last_run timestamptz NULL,
	run_type varchar NULL,
	enabled bool NOT NULL DEFAULT false,
	CONSTRAINT runner_pk PRIMARY KEY (map_id, user_id)
);

-- users definition

-- Drop table

-- DROP TABLE users;

CREATE TABLE users (
	user_id varchar NOT NULL,
	email varchar NULL,
	install_date timestamptz NULL,
	last_run timestamptz NULL,
	user_name varchar NULL,
	"password" varchar NULL,
	CONSTRAINT users_pk PRIMARY KEY (user_id)
);


-- ig_mappers definition

-- Drop table

-- DROP TABLE ig_mappers;

CREATE TABLE ig_mappers (
	map_id uuid NOT NULL,
	user_id varchar NULL,
	user_email varchar NULL,
	ig_user_name varchar NULL
);


-- ig_runners definition

-- Drop table

-- DROP TABLE ig_runners;

CREATE TABLE ig_runners (
	map_id uuid NOT NULL,
	user_id varchar NOT NULL,
	chron numeric NULL,
	last_run timestamptz NULL,
	enabled bool NOT NULL DEFAULT false,
	CONSTRAINT ig_runners_pk PRIMARY KEY (map_id, user_id)
);


-- ig_captured definition

-- Drop table

-- DROP TABLE ig_captured;

CREATE TABLE ig_captured (
	map_id uuid NOT NULL,
	user_id varchar NULL,
	capture_date timestamptz NULL,
	ig_username varchar NULL,
	raw_scraped_payload jsonb NULL
);