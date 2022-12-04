import pandas as pd
import numpy as np
import networkx as nx
import pickle
import sqlite3 as sql3
import distance
import trip_request

def ret_trips(startLocation,endLocation):
    sourceLat = startLocation["location"]["latitude"]
    sourceLng = startLocation["location"]["longitude"]
    endLat = endLocation["location"]["latitude"]
    endLng = endLocation["location"]["longitude"]
    CitiesG = pickle.load(open('graphs/city-graph.txt','rb'))
    CITIES = len(CitiesG.nodes)
    MAX_TRANSITS = 1
    MAX_SPLITS = 3

    # BFS through cities
    #     collect list of paths
    # For path
    #     fetch trips for every edge

    source_node = 0
    end_node = 0
    source_min_dist = np.Inf
    end_min_dist = np.Inf


    for i in range(CITIES):
        source_dist = distance.distance(CitiesG,i,sourceLat,sourceLng)
        end_dist = distance.distance(CitiesG,i,endLat,endLng)
        if source_dist < source_min_dist:
            source_node = i
            source_min_dist = source_dist
        if end_dist < end_min_dist:
            end_node = i
            end_min_dist = end_dist


    paths_list = []

    cities_queue = [source_node]
    levels_queue = [0]
    path_queue = [[source_node]]
    while len(cities_queue) > 0:
        curr_node = cities_queue.pop(0)
        curr_node_level = levels_queue.pop(0)
        curr_path = path_queue.pop(0)
        for edge in CitiesG.edges:
            next_node = edge[1]
            #if next_node exists in node chain: continue
            if next_node in curr_path:
                continue
            #if not exceeded maximum transit
            if curr_node_level <= MAX_TRANSITS:
                cities_queue.append(next_node)
                levels_queue.append(curr_node_level+1)
                curr_path.append(next_node)
                path_queue.append(curr_path)
            if next_node == end_node:
                paths_list.append(curr_path)

    return_trips = []
    for path in paths_list:
        trips_by_layer = []
        for i in range(len(path)-1):
            trips_by_layer.append([])
            fetch_s_city = path[i]
            fetch_e_city = path[i+1]
            edge_list = sorted(trip_request.fetch_destinations(CitiesG.nodes[fetch_s_city]["lat"],CitiesG.nodes[fetch_s_city]["lng"],CitiesG.nodes[fetch_e_city]["lat"],CitiesG.nodes[fetch_e_city]["lng"]), key=lambda x:x["endTime"])[0:MAX_SPLITS]
            print("edge list: ", edge_list)
            #edge_list = fetch(CitiesG.nodes[fetch_s_city]["city_ascii"],CitiesG.nodes[fetch_e_city]["city_ascii"])
            #keep trips with time > prev_edge arrival

            if i == 0:
                for trip in edge_list:
                    trips_by_layer[i].append([trip])
            else:
                for j in range(len(trips_by_layer[i-1])):
                    prev_trip = trips_by_layer[i-1][j][i-1]
                    for trip in edge_list:
                        if trip['startTime'] > prev_trip['endTime']:
                            new_trips = trips_by_layer[i-1][j]
                            new_trips.append(trip)
                            trips_by_layer[i].append(new_trips)
        for trip_list in trips_by_layer[len(path)-2]:
            return_trips.append([10,trip_list])
    return return_trips

