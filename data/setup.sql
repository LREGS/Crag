drop table crag cascade if extists; 
drop table climb; 
drop table report;
drop table forecast;

create table crag(
    Id          serial primary key, 
    Name        text, 
    Latitude double precision,
    Longitude double precision
);

create table climb(
    Id          serial primary key,
    Name        varchar(255),
    Grade       varchar(255),
    CragID      integer references crag(Id)
);

create table report(
    Id          serial primary key, 
    Content     varchar(255),
    Author      varchar(255),
    CragID      integer references crag(Id)
);

create table forecast(
    Id                  serial primary key, 
    Time                varchar(255),
    ScreenTemperature   double precision,
    FeelsLikeTemp       double precision, 
    WindSpeed           double precision,
    WindDirection       double precision,
    totalPrecipitation  double precision,
    ProbofPrecipitation int,
    Latitude            double precision,
    Longitude           double precision, 

);
