
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
    CragID INTEGER REFERENCES crag(Id)

);


