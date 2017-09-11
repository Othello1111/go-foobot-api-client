/*
*
*/
--------------------------------------------------------------------------------
-------------------------- table definition ------------------------------------
--------------------------------------------------------------------------------
CREATE TABLE "testing" (
 "device_uuid" INT NOT NULL PRIMARY KEY,
 "device_name" varchar NOT NULL,
 "device_mac" varchar
);

----- import with psql -f filename CONNECTSTRING -------------------------------