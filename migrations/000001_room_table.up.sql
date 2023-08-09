create type room_category_enum as ENUM (
	'General', 
	'Tech', 
	'Game',
	'Book',
	'Movie',
	'Music',
	'Language',
	'Science'
);

create table if not exists rooms (
	id varchar(36) primary key, 
	admin_id varchar(36) not null, 
	name varchar(50) not null, 
	category room_category_enum not null, 
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null
);
