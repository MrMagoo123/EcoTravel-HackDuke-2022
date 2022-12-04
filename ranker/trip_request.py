import requests
import json

def fetch_destinations(startLat,startLng, endLat,endLng):
    """Fetch list of destinations from API."""

    print("start lat: ", startLat)
    print("start lng: ", startLng)
    print("end lat: ", endLat)
    print("end lng: ", endLng)

    # payload = {

    # }

    return requests.post("https://api.ecotravel.tech/api/createQuery", json={
        "startLocation": {
            "formatted": "Columbia, SC, USA",
            "country": "US",
            "location": {
    "latitude": 34.0007104,
    "longitude": -81.0348144
            },
            "streetAddress": {
                "number": "",
                "name": "",
                "apt": ""
            },
            "postalCode": "undefined",
            "subdivision": "SC",
            "city": "Columbia"
        },
        "endLocation": {
            "formatted": "Atlanta, GA, USA",
            "country": "US",
            "location": {
    "latitude": 33.7489954,
    "longitude": -84.3879824
            },
            "streetAddress": {
                "number": "",
                "name": "",
                "apt": ""
            },
            "postalCode": "undefined",
            "subdivision": "GA",
            "city": "Atlanta"
        }
    }, headers={
        "content-type": "application/json",
        "user-agent": "wen jia hu's evil kingdom"
    },verify=False).json()

if __name__ == "__main__":
    print(fetch_destinations("New York","Los Angeles"))