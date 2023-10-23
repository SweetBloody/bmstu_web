-- drop database if exists formula1;
-- create database formula1;

drop table if exists GrandPrix cascade;
create table public.GrandPrix(
                                 gp_id serial not null primary key,
                                 gp_season int not null check (gp_season > 1949),
                                 gp_name text not null,
                                 gp_date_num int not null check (gp_date_num between 0 and 32),
                                 gp_month text not null,
                                 gp_place text not null,
                                 gp_track_id int not null
);

drop table if exists Tracks cascade;
create table public.Tracks(
                              track_id serial not null primary key,
                              track_name text not null,
                              track_country text not null,
                              track_town text not null
);

drop table if exists QualificationResults cascade;
create table public.QualificationResults(
                                            qual_id serial not null primary key,
                                            qual_driver_place int not null,
                                            driver_id int not null,
                                            team_id int not null,
                                            q1_time time,
                                            q2_time time,
                                            q3_time time,
                                            gp_id int not null
);

drop table if exists RaceResults cascade;
create table public.RaceResults(
                                   race_id serial not null primary key,
                                   race_driver_place int,
                                   driver_id int not null,
                                   team_id int not null,
                                   gp_id int not null
);

drop table if exists Drivers cascade;
create table public.Drivers(
                               driver_id serial not null primary key,
                               driver_name text not null,
                               driver_country text,
                               driver_birth_date date
);

drop table if exists Teams cascade;
create table public.Teams(
                             team_id serial not null primary key,
                             team_name text not null,
                             team_country text not null,
                             team_base text not null
);

drop table if exists TeamsDrivers cascade;
create table public.TeamsDrivers(
                                    td_id serial not null primary key,
                                    driver_id int not null,
                                    team_id int not null,
                                    team_driver_season int not null check (team_driver_season > 1949)
);

drop table if exists Users cascade;
create table public.Users(
                             user_id serial not null primary key,
                             login text not null,
                             password text not null,
                             role text not null
);

drop table if exists season_standings cascade;
create table public.season_standings(
                                        st_id serial not null primary key,
                                        season int not null,
                                        driver_id int not null,
                                        team_id int not null,
                                        score int not null
);

drop view if exists race_results_view cascade;
create view race_results_view as
select race_id, race_driver_place, driver_name, team_name, gp_name
from raceresults r
         join drivers d on r.driver_id = d.driver_id
         join grandprix g on r.gp_id = g.gp_id
         join teams t on r.team_id = t.team_id
where race_driver_place = 1;

drop view if exists drivers_of_season cascade;
create view drivers_of_season as
select d.driver_id, driver_name, driver_country, driver_birth_date
from drivers d
         join teamsdrivers t on d.driver_id = t.driver_id
where team_driver_season = 2022;

drop table if exists RaceResultsTmp cascade;
create table public.RaceResultsTmp(
                                      race_id serial not null primary key,
                                      race_driver_place int,
                                      driver_id int not null,
                                      team_id int not null,
                                      gp_id int not null
);

set datestyle to 'dmy';
alter table GrandPrix add foreign key (gp_track_id) references public.Tracks(track_id);
alter table QualificationResults add foreign key (gp_id) references public.GrandPrix(gp_id);
alter table RaceResults add foreign key (gp_id) references public.GrandPrix(gp_id);
alter table TeamsDrivers add foreign key (driver_id) references public.Drivers(driver_id);
alter table TeamsDrivers add foreign key (team_id) references public.Teams(team_id);
alter table season_standings add foreign key (driver_id) references public.Drivers(driver_id);
alter table season_standings add foreign key (team_id) references public.Teams(team_id);
-- alter table TeamsDrivers add primary key (driver_id, team_id);


copy public.tracks(track_name,
    track_country,
    track_town)
    from '/db_data/data/tracks.csv' delimiter ';' CSV;

copy public.grandprix(gp_season,
    gp_name,
    gp_date_num,
    gp_month,
    gp_place,
    gp_track_id)
    from '/db_data/data/gp.csv' delimiter ';' CSV;

copy public.drivers(driver_name,
    driver_country,
    driver_birth_date)
    from '/db_data/data/drivers.csv' delimiter ';' CSV;

copy public.qualificationresults(qual_driver_place,
    driver_id,
    team_id,
    q1_time,
    q2_time,
    q3_time,
    gp_id)
    from '/db_data/data/qualifications.csv' delimiter ';' CSV;

copy public.raceresultstmp(race_driver_place,
    driver_id,
    team_id,
    gp_id)
    from '/db_data/data/races.csv'  delimiter ';' CSV;

copy public.raceresults(race_driver_place,
    driver_id,
    team_id,
    gp_id)
    from '/db_data/data/races.csv'  delimiter ';' CSV;

copy public.teams(team_name,
    team_country,
    team_base)
    from '/db_data/data/teams.csv' delimiter ';' CSV;

copy public.teamsdrivers(driver_id,
    team_id,
    team_driver_season)
    from '/db_data/data/teams_drivers.csv' delimiter ';' CSV;


create user "default_guest";
create user "default_user";
create user "default_admin";

alter role "default_guest" password '11111111';
alter role "default_user" password '12344321';
alter role "default_admin" password '12345678';

grant select on table grandprix to "default_guest";
grant select on race_results_view to "default_guest";
grant select on drivers_of_season to "default_guest";

grant select on table drivers to "default_user";
grant select on table grandprix to "default_user";
grant select on table qualificationresults to "default_user";
grant select on table raceresults to "default_user";
grant select on table season_standings to "default_user";
grant select on table teams to "default_user";
grant select on table teamsdrivers to "default_user";
grant select on table tracks to "default_user";
grant select on race_results_view to "default_user";
grant select on drivers_of_season to "default_user";

alter role "default_admin" superuser;

-- drop owned by "default_user";
-- drop owned by "default_admin";
-- drop user "default_user";
-- drop user "default_admin";


-- Получение количества очков по гоночному результату гонщика

create or replace function GetScore(place int)
returns int
as $$
begin
    if place = 1 then
        return 25;
    elsif place = 2 then
        return 18;
    elsif place = 3 then
        return 15;
    elsif place = 4 then
        return 12;
    elsif place = 5 then
        return 10;
    elsif place = 6 then
        return 8;
    elsif place = 7 then
        return 6;
    elsif place = 8 then
        return 4;
    elsif place = 9 then
        return 2;
    elsif place = 10 then
        return 1;
else return 0;
end if;
end
$$ language plpgsql;

-- Получение сезона по id гонки
create or replace function GetSeason(id int)
returns int
as $$
declare res int;
begin
select gp_season
into res
from raceresults r
         join grandprix g on g.gp_id = r.gp_id
where race_id = id;
return res;
end
$$ language plpgsql;


-- Функция триггера
create or replace function UpdateTrigger()
returns trigger
as $$
begin
    raise notice 'New =  %, season = %, driver = %', new, GetSeason(new.race_id), new.driver_id;
    if GetSeason(new.race_id) = 2022 then
update season_standings
set score = score + GetScore(new.race_driver_place)
where driver_id = new.driver_id;
end if;
return new;
end
$$ language plpgsql;

-- Триггер
create trigger update_season_standing
    after insert on raceresults
    for each row
    execute procedure UpdateTrigger();


create or replace function DeleteTrigger()
returns trigger
as $$
begin
    if GetSeason(old.race_id) = 2022 then
update season_standings
set score = score - GetScore(old.race_driver_place)
where driver_id = old.driver_id;
end if;
return old;
end
$$ language plpgsql;

create trigger delete_season_standing
    before delete on raceresults
    for each row
    execute procedure DeleteTrigger();