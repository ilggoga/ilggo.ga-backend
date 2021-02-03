create user ilggoga@localhost;
create schema ilggoga;

grant all privileges on ilggoga.* to ilggoga@localhost;
use ilggoga;

create table users (
  id varchar(12) not null, -- identifier
  passwd varchar(128) not null, -- password
  display varchar(50) default null, -- display name
  is_admin boolean default false not null, -- is admin permited?
  created_at timestamp default current_timestamp not null -- account creation timestamp
);

create table novels (
  flags text default "", -- novel flags
  content text not null, -- content of novel
  id varchar(12) not null, -- identifier
  author varchar(12) not null, -- author identifier
  created_at timestamp default current_timestamp not null -- novel creation timestamp
);
