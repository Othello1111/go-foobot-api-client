# go-foobot-api-client
Go(lang) API client for <a href="https://foobot.io/">FooBot</a> air quality monitor, logs to database. 

This is currently a work in progress.... more to come.

Deployed as a WIP to <a href="http://airquality.clayshekleton.com/">http://airquality.clayshekleton.com/</a>

## Requirements

- A Heroku account

## Usage

- 
-


## To-Do

 - [x] Figure out web worker process crash on deploy 
 - [x] Create Procfile & change dyno type?
 - [ ] Add scheduler config to terraform deployment / app.json?
 - [ ] Further separate into additional functions besides main()
 - [ ] Automate table creation after database is instantiated (app.json post-deployment script? 'release' phase of Procfile?). Update SQL to create if not exists.
 - [ ] Tighten up logging
 - [ ] Implement OAuth