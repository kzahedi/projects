#!/usr/bin/env python3

import functools
import os
import glob
import json
import operator
import random
from igraph import *
from enum import unique
import argparse

handles = []
edges = []
vNames = []
eid = 0


def collectData(node):
    global edges
    global vNames
    global eid
    if node['Children'] == None:
        return

    mid = eid
    for x in node['Children']:
        eid = eid + 1
        edges.append((mid, eid))
        vNames.append(x['TwitterHandle'])
        collectData(x)


def processFile(file, dir):
    global edges
    global vNames
    global eid
    random.seed()
    out = file.replace(".json", ".pdf")
    out = os.path.basename(out)

    fd = open(file, "r")
    obj = json.load(fd)
    fd.close()

    outdir = dir + "/" + obj['TwitterHandle']
    out = outdir + "/" + out

    if os.path.isfile(out) == True:
        print("File %s already exists." % out)
        exit(0)

    eid = 0
    edges = []
    vNames = [obj['TwitterHandle']]
    collectData(obj)

    uniqueNames = list(set(vNames))

    colour_dict = {}
    for n in uniqueNames:
        r = random.randint(0, 255)
        g = random.randint(0, 255)
        b = random.randint(0, 255)
        colour_dict[n] = "rgb(" + str(r) + ", " + str(g) + ", " + str(b) + ")"

    if len(edges) < 5:
        return

    if os.path.isdir(outdir) == False:
        os.mkdir(outdir)

    max = 0
    for e in edges:
        if e[0] > max:
            max = e[0]
        if e[1] > max:
            max = e[1]

    g = Graph()

    # print(len(vNames))
    # print(edges)
    g.add_vertices(len(vNames))
    g.add_edges(edges)
    layout = g.layout_fruchterman_reingold(maxiter=10000)
    colours = [colour_dict[name] for name in vNames]
    visual_style = dict()
    visual_style["vertex_size"] = 10
    visual_style["vertex_label_size"] = 5
    visual_style["vertex_label_dist"] = 3
    colours[0] = "#FF0000"
    visual_style["vertex_label_color"] = colours
    colours[0] = "#FFFFFF"
    visual_style["vertex_color"] = colours
    visual_style["vertex_label"] = vNames
    visual_style["edge_width"] = 2
    visual_style["layout"] = layout
    # visual_style["bbox"] = (1200, 1000)
    # visual_style["margin"] = 100

    plot(g, out, **visual_style)

    print(out)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", type=str, default=None, help="Input File")
    parser.add_argument("-d", type=str, default=None, help="Input File")
    args = parser.parse_args()

    if args.f == None or args.d == None:
        print("Please check command line parameters.")
        os.exit(0)


    processFile(args.f, args.d)
