create table crag(
    Id          serial primary key, 
    Name        text, 
    location    point[],
);

create table climb(
    Id          serial primary key,
    Name        varchar(255),
    Grade       varchar(255),
    CragID      integer references crag(Id)
)

create table report(
    Id          serial primary key, 
    Content     varchar(255),
    Author      varchar(255),
    CragID      integer references crag(Id)
)
