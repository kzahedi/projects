import functools
import glob
import json
import operator
import random
from igraph import *
from enum import unique

handles = []
edges = []
vNames = []


def collectData(node):
    global edges
    global vNames
    if node['Children'] == None:
        return

    vNames.append(node['TwitterHandle'])
    for x in node['Children']:
        edges.append((x['ParentID'], x['ID']))
        vNames.append(x['TwitterHandle'])

    for x in node['Children']:
        collectData(x)


def processFile(f,n):
    global edges
    global vNames
    random.seed()
    out = f.replace(".json",".pdf")
    out = out.replace("data","pdfs")
    fd = open(f,"r")
    obj = json.load(fd)
    fd.close()

    edges = []
    vNames = []
    collectData(obj)

    uniqueNames = list(set(vNames))

    colour_dict = {}
    for n in uniqueNames:
        r = random.randint(0,255)
        g = random.randint(0,255)
        b = random.randint(0,255)
        colour_dict[n] = "rgb(" + str(r) + ", " + str(g) + ", " + str(b) + ")"

    if len(edges) < 5:
        return

    max = 0
    for e in edges:
        if e[0] > max:
            max = e[0]
        if e[1] > max:
            max = e[1]

    g = Graph()

    g.add_vertices(max+1)
    g.add_edges(edges)
    layout = g.layout_fruchterman_reingold(maxiter=10000)
    colours = [colour_dict[name] for name in vNames]
    visual_style = dict()
    visual_style["vertex_size"] = 10
    visual_style["vertex_label_size"] = 5
    visual_style["vertex_label_dist"] = 3
    visual_style["vertex_color"] = colours
    visual_style["vertex_label_color"] = colours
    visual_style["vertex_label"] = vNames
    visual_style["edge_width"] = 2
    visual_style["layout"] = layout
    # visual_style["bbox"] = (1200, 1000)
    # visual_style["margin"] = 100

    plot(g, out, **visual_style)

    print(out)

# files = glob.glob("sfi/data/*.json")
# for f in files:
#     processFile(f, 1)

processFile("sfi/data/1039261599012450305.json", 1)