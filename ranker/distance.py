from math import cos, asin, sqrt, pi


def distance(G, c1, c2):
    lat1 = G.nodes[c1]['lat']
    lng1 = G.nodes[c1]['lng']
    lat2 = G.nodes[c2]['lat']
    lng2 = G.nodes[c2]['lng']
    p = pi/180
    a = 0.5 - cos((lat2-lat1)*p)/2 + cos(lat1*p) * cos(lat2*p) * (1-cos((lng2-lng1)*p))/2
    return 12742 * asin(sqrt(a))

def distance(G, c1, lat2, lng2):
    lat1 = G.nodes[c1]['lat']
    lng1 = G.nodes[c1]['lng']
    p = pi/180
    a = 0.5 - cos((lat2-lat1)*p)/2 + cos(lat1*p) * cos(lat2*p) * (1-cos((lng2-lng1)*p))/2
    return 12742 * asin(sqrt(a))
