#!/usr/bin/env python3

import glob

files = glob.glob("data/*.json")

allstr = "all: "
jobs = ""
index = 0
for f in files:
    job = "j" + str(index)
    jobs = jobs + job + ":\n\t" + "./convert_json_file.py -f " + f + " -d pdfs\n\n"
    allstr = allstr + " " + job
    index = index + 1


fd = open("Makefile", "w")
fd.write(allstr + "\n\n"+jobs)
fd.close()
