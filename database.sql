create user ilggoga@localhost;
create schema ilggoga;

grant all privileges on ilggoga.* to ilggoga@localhost;
use ilggoga;

create table users (
  id varchar(12) not null primary key, -- identifier
  passwd varchar(128) not null, -- password
  display varchar(50) default null, -- display name
  is_admin boolean default false not null, -- is admin permited?
  created_at timestamp default current_timestamp not null -- account creation timestamp
);

create table novels (
  id int not null primary key, -- identifier
  likes text default "", -- liked users
  flags text default "", -- novel flags
  content text not null, -- content of novel
  author varchar(12) not null, -- author identifier
  created_at timestamp default current_timestamp not null -- novel creation timestamp
);

create table comments (
  id int not null primary key, -- identifier
  novel int not null, -- target novel
  author varchar(12) not null, -- comment author
  content text not null, -- content of comment
  created_at timestamp default current_timestamp not null -- comment creation timestamp 
);
