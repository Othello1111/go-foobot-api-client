/*
*
*/
--------------------------------------------------------------------------------
-------------------------- table definition ------------------------------------
--------------------------------------------------------------------------------
CREATE TABLE sensor_data (
    device_uuid VARCHAR(50),
    seconds_since_epoch INTEGER,
    datetime TIMESTAMPTZ,
    particles_ugm3 NUMERIC,
    temp_c NUMERIC,
    humidity_perct NUMERIC,
    co2_ppm NUMERIC,
    voc_ppb NUMERIC,
    pollution_score NUMERIC,
    PRIMARY KEY (device_uuid, seconds_since_epoch)
);

----- import with psql -f filename CONNECTSTRING -------------------------------
/*
INSERT into sensor_data (device_uuid, seconds_since_epoch, datetime, particles_ugm3, temp_c, humidity_perct, co2_ppm, voc_ppb, pollution_score) 
VALUES ('2C044665391009F1', 1504641600, to_timestamp(1504641600), 12.60001, 24.695, 54.011623, 454.375, 126.125, 12.760724);

select datetime::timestamp at time zone 'UTC' at time zone 'America/New_York' from sensor_data;

NOTE: First FOOBOT database at 2017-09-05 16:00 EDT, seconds since epoch: 1504642002
*/