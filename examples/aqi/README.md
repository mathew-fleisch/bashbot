# Bashbot Example - Get Air Quality Index By Zip Code

In this example, a curl is executed via bash script from the [Air Now API](https://docs.airnowapi.org/) and the response json is parsed via jq to send a formatted message back to slack.

## Bashbot configuration

This command is triggered by sending `bashbot aqi [zip]` in a slack channel where Bashbot is also a member. The script is expected to exist before execution at the relative path `./examples/aqi` and requires the following environment variables to be set: `AIRQUALITY_API_KEY` The `zip` parameter is validated by building a list of all five digit integers with a for loop.

```json
{
  "name": "Air Quality Index",
  "description": "Get air quality index by zip code",
  "help": "bashbot aqi [zip-code]",
  "trigger": "aqi",
  "location": "./examples/aqi",
  "command": [
    "source /bashbot/.env",
    "&& ./aqi.sh ${zip}"
  ],
  "parameters": [
    {
      "name": "zip",
      "allowed": [],
      "description": "any integer between 10000 through 99999",
      "source": ["for i in {10000..99999}; do echo $i; done"]
    }
  ],
  "log": true,
  "ephemeral": false,
  "response": "text",
  "permissions": [
    "all"
  ]
}
```

## Bashbot scripts

There is one script ([aqi.sh](aqi.sh)) associated with this example and takes one argument/parameter to retrieve the air quality index value for specific zip codes.

```bash
if [ -z "$AIRQUALITY_API_KEY" ]; then
  echo "Missing Air Now API Key..."
  echo "<https://docs.airnowapi.org/|Air Now API>"
  exit 0
fi

zip=$1
if [ -z "$zip" ]; then
  echo "Usage: $0 [zip]"
  exit 0
fi

response=$(curl -s "http://www.airnowapi.org/aq/observation/zipCode/current/?zipCode=${zip}&distance=5&format=application/json&API_KEY=${AIRQUALITY_API_KEY}")

if [[ "$response" == "[]" ]]; then
  echo "There is no <https://docs.airnowapi.org/|aqi value> for this zip: $zip"
  exit 0
fi

aqi=$(echo "$response" | jq '.[0]')
reporting_area=$(echo "$aqi" | jq -r '.ReportingArea')
aqi_value=$(echo "$aqi" | jq -r '.aqi_value')
time_stamp="$(echo "$aqi" | jq -r '.DateObserved')$(echo "$aqi" | jq -r '.HourObserved'):00"
category=$(echo "$aqi" | jq -r '.Category.Name')
case $category in
  "Good") emoji=":large_green_circle:";;
  "Moderate") emoji=":large_yellow_circle:";;
  "Unhealthy for Sensitive Groups") emoji=":large_orange_circle:";;
  "Unhealthy") emoji=":red_circle:";;
  "Very Unhealthy") emoji=":large_purple_circle:";;
  "Hazardous") emoji=":black_circle:";;
esac

echo "$emoji The <https://docs.airnowapi.org/aq101|Air Quality Index> in $reporting_area is $aqi_value ($category) as of $time_stamp";

```