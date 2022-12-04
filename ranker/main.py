import re
from flask import Flask, redirect, url_for, request
import find_trips
import json

app = Flask(__name__)


# @app.route('/success/<name>',methods=['POST'])
# def find_trips(name):
#     return find_trips.func()
#     # return 'welcome %s' % name


@app.route('/trips', methods=['POST'])
def trips():
    if request.method == 'POST':
        # sourceLat = request.args.get('sourcelat')
        # sourceLng = request.args.get('sourcelng')
        # endLat = request.args.get('endlat')
        # endLng = request.args.get('endlng')
        startLocation = request.json.get("startLocation")
        endLocation = request.json.get("endLocation")
        print(request.form)
        print(request.json)
        # sourceLat = 40.6943
        # sourceLng = -73.9249
        # endLat = 41.8373
        # endLng = -87.6862
        # print("sourceLat: ", sourceLat)
        # print("sourceLng: ", sourceLng)
        # print("endLat: ", endLat)
        # print("endLng: ", endLng)
        print(find_trips.ret_trips(startLocation, endLocation))
        return find_trips.ret_trips(startLocation, endLocation)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=8080)
