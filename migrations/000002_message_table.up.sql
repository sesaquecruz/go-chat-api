create table if not exists messages (
	id varchar(36) primary key, 
	room_id varchar(36) not null references rooms(id), 
	sender_id varchar(36) not null, 
	sender_name varchar(50) not null, 
	text varchar(100) not null, 
	created_at timestamp with time zone not null
);
