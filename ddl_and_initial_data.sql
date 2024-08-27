create database if not exists blog;

use blog;

create table blog_user(
	id bigint auto_increment,
	username varchar(255) not null,
	pass varchar(255) not null,
	primary key(id)
);


-- username y password son el mismo valor
insert into blog_user(username, pass) values ('nico','$2y$10$Ou6GpNeMB..Sd5Nv0dLvxOXrUxGTnMQG.6cCZ6j6C2pbRhIf0/Rre'); 
insert into blog_user(username, pass) values ('gerar','$2y$10$csZuIGL3piR1VGviTDN42Ovb3NWdA0h.d.k9CusHmecCPM5wvcxMO');

create table post(
	id bigint auto_increment,
	title varchar(200) not null,
	body varchar(2000),
	author bigint not null,
	primary key(id),
	foreign key(author) references blog_user(id)
);


create table tag(
	id bigint auto_increment,
	description varchar(30) not null,
	primary key(id)
);

insert into tag(description) values ('Java'), ('Golang'), ('Python'), ('Backend');

create table post_tag(
	id_post bigint,
	id_tag bigint,
	foreign key(id_post) references post(id),
	foreign key(id_tag) references tag(id)
);

create table comment(
	id bigint auto_increment,
	content varchar(500) not null,
	id_post bigint,
	id_user bigint,
	primary key(id),
	foreign key(id_post) references post(id),
	foreign key(id_user) references blog_user(id)
);

