# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

## [0.5.0] 2024-12-19

- Security patch golang.net/x/net

## [0.4.0] 2023-10-12

- Security patch golang.net/x/net

## [0.3.0] 2022-11-10

- Rename field Sensor.Pause to Interval
- Add option --write-example-script
- Add option --interval
- cmd/sense executes ./.onchange.sh if found
- Add options to cmd/sense

## [0.2.0] 2022-07-17

- Add constructor NewSensor and hide fields Root and Visit
- Add MIT license

## [0.1.1] 2021-08-19

- Handle nil FileInfo

## [0.1.0] 2020-02-18

- Initial sensor of file modifications using periodic scanning
