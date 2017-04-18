# Simple Load Tester

## Goal
To create a configurable, deployable, and continuous performance testing tool to reveal issues with new versions of a system.

## The Why
- I'm tired of manually kicking off performance tests
- I'm tired of performance tests hiding the acutal limits of systems
- I want my system in non-prod environments to be under heavy load or comparable load to production

## The Plan
The long term plan for this tool is to have a deployable (docker?) artifact that will continuously stress a microservice. Using a stairstep function with configurable feedback, it will ramp up until the system starts to falter. Then it will back off/ramp up according to the configurable stair-step function. In theory, as the tool continues to run, we can figure out what the max capacity for a microservice is and how it responds to failure over longer periods of time.

## What does it do right now?
Very Little -- basically calls off to a service and either doubles or halves its request count

## Other thoughts 
I'm not quite convinced this is necessary, but only time will tell...


# TODO
* Remove the result channel and make it an interface *****
* Figure out a better interface than the 'Start' methods scattered everywhere -- missing a simplification somewhere
* look @ a feeder style setup from Gatling
* validate the engine is reaching the correct RPS (so that we don't keep adding reqs and then topple over)
