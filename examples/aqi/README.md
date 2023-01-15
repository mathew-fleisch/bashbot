# Bashbot Example - Get Air Quality Index By Zip Code

In this example, a curl is executed via bash script from the [Air Now API](https://docs.airnowapi.org/) and the response json is parsed via jq to send a formatted message back to slack.

<img src="https://i.imgur.com/YV13qKC.gif" />

## Bashbot configuration

This command is triggered by sending `!bashbot aqi [zip]` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/aqi` and requires the following environment variables to be set: `AIRQUALITY_API_KEY` The `zip` parameter is validated by building a list of all five digit integers with a for loop. This command requires [jq](https://stedolan.github.io/jq/) and [curl](https://curl.se/) to be installed on the host machine.

```yaml
name: Air Quality Index
description: Get air quality index by zip code
envvars:
  - AIRQUALITY_API_KEY
dependencies:
  - curl
  - jq
help: "!bashbot aqi [zip-code]"
trigger: aqi
location: /bashbot/vendor/bashbot/examples/aqi
command:
  - "./aqi.sh ${zip}"
parameters:
  - name: zip
    allowed: []
    description: any zip code
    match: (^\d{5}$)|(^\d{9}$)|(^\d{5}-\d{4}$)
log: true
ephemeral: false
response: text
permissions:
  - all
```

---

There is one script ([aqi.sh](aqi.sh)) associated with this example and takes one argument/parameter to retrieve the air quality index value for specific zip codes. The raw json response from the [Air Now API](https://docs.airnowapi.org/) before [aqi.sh](aqi.sh) parses the values that are displayed in slack:

```json
[
  {
    "DateObserved": "2021-08-25 ",
    "HourObserved": 3,
    "LocalTimeZone": "PST",
    "ReportingArea": "San Francisco",
    "StateCode": "CA",
    "Latitude": 37.75,
    "Longitude": -122.43,
    "ParameterName": "O3",
    "AQI": 32,
    "Category": {
      "Number": 1,
      "Name": "Good"
    }
  },
  {
    "DateObserved": "2021-08-25 ",
    "HourObserved": 3,
    "LocalTimeZone": "PST",
    "ReportingArea": "San Francisco",
    "StateCode": "CA",
    "Latitude": 37.75,
    "Longitude": -122.43,
    "ParameterName": "PM2.5",
    "AQI": 30,
    "Category": {
      "Number": 1,
      "Name": "Good"
    }
  }
]
```
