# Citybike
Citybike is a web app showing data from journeys made with city bikes in the Helsinki Capital area.\
I preferred to use Golang for its performance and ease of use and SQLite as a database since it is lightweight and easy to use with Go.

Journey data imported from:
- <a href="https://dev.hsl.fi/citybikes/od-trips-2021/2021-05.csv">https://dev.hsl.fi/citybikes/od-trips-2021/2021-05.csv</a>
- <a href="https://dev.hsl.fi/citybikes/od-trips-2021/2021-06.csv">https://dev.hsl.fi/citybikes/od-trips-2021/2021-06.csv</a>
- <a href="https://dev.hsl.fi/citybikes/od-trips-2021/2021-07.csv">https://dev.hsl.fi/citybikes/od-trips-2021/2021-07.csv</a>

Bicycle station data imported from:
- <a href="https://opendata.arcgis.com/datasets/726277c507ef4914b0aec3cbcfcbfafc_0.csv">https://opendata.arcgis.com/datasets/726277c507ef4914b0aec3cbcfcbfafc_0.csv</a>
  - <a href="https://www.avoindata.fi/data/en/dataset/hsl-n-kaupunkipyoraasemat/resource/a23eef3a-cc40-4608-8aa2-c730d17e8902">License and information</a>

## Prerequisites
Golang version 1.19 or higher

## Usage
```
make all
```
Or you can run a docker image:
```
make docker
```
The application will start running on <a href="http://localhost:8000">localhost:8000</a>
## Features
- <a href="http://localhost:8000">Index page</a> provides links for journey and station lists.
- <a href="http://localhost:8000/journeys">The journey list view</a>  shows departure and return stations (names and IDs), covered distance in kilometres and duration in minutes.
- <a href="http://localhost:8000/stations">The station list view</a> shows the station ID, name and address information.
- Single station view, when clicked on the station ID either from the journey list or station list, shows the station name and address and the total number of journeys starting from and ending at the station.

## Testing
### Unit tests
```
make test
```
### cURL tests
First, run the service:
```
make all
```
Then, open a new terminal page and run the script:
```
./curl.sh
```
