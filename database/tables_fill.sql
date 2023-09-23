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