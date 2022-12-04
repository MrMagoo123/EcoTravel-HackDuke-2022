import pandas as pd
import numpy as np
import networkx as nx
import pickle
import sqlite3 as sql3
import distance

CITIES = 100
df = pd.read_csv(r'data/cities.csv')[0:CITIES]
CitiesG = nx.Graph()

for i in range(CITIES):
    CitiesG.add_node(i,city=df['city_ascii'][i],lat=df['lat'][i],lng=df['lng'][i])
for i in range(CITIES):
    for j in range(i+1,CITIES):
        CitiesG.add_edge(i,j,dist=distance.distance(CitiesG,i,j))

pickle.dump(CitiesG,open('graphs/city-graph.txt','wb'))


# G = nx.Graph()
# G.add_node(0,color='yellow')
# G.add_nodes_from([(1,{'color':'yellow'}),
#                   (2,{'tint':'green'})])
# print(G.nodes[0]['color'])
#
# pickle.dump(G,open('graphs/large-city-graph.txt','wb'))
# G=pickle.load(open('graphs/large-city-graph.txt','rb'))
# print(G)
