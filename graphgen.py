#!/usr/bin/env python3.8
import random, sys, math

if sys.version_info[0] < 3:
    raise Exception("Python 3 or a more recent version is required.")
    exit()

if len(sys.argv) < 3 or len(sys.argv) > 6:
	raise ValueError('Please run the program as ./graphgen <num_nodes> <max_degree> [-suppress] [-sparse] [description]')

num_nodes = int(sys.argv[1])
max_degree = int(sys.argv[2])
name = f"Graph_N{num_nodes}_D{max_degree}"
description = f"A graph with {num_nodes} nodes and a max degree of {max_degree}"

def getRando(sparse, MD):
    if sparse:
        return round(random.uniform(0, MD-2))
    else:
        return 0
        

suppress = False
sparse = False
for k in range(3, len(sys.argv)):
    if sys.argv[k][0] == "-":
        if sys.argv[k] == "-suppress":
            suppress = True
        elif sys.argv[k] == "-sparse":
            sparse = True
        else:
            raise ValueError('Please run the program as ./graphgen <num_nodes> <max_degree> [-suppress] [-sparse] [description]')
    else: 
        description = sys.argv[k]
        break

if sparse:
    sparseName = "_sparse"
else:
    sparseName = ""

file = open(name+sparseName+".txt", "w+")
file.write(name + "\n")
file.write(str(description) + "\n")
file.write(str(max_degree) + "\n")

array = [[0 for _ in range(num_nodes)] for _ in range(num_nodes)]

#links first and last
array[0][num_nodes-1] = 1
array[num_nodes-1][0] = 1

#If sparse, starts at somewhat already populated
countsTracker = [getRando(sparse, max_degree) for _ in range(num_nodes)]
countsTracker[0] = 1
countsTracker[num_nodes-1] = 1

#links first to random sample
for num in random.sample(range(2, num_nodes-1), max_degree - 2):
    array[0][num] = 1
    array[num][0] = 1
    countsTracker[num] += 1
    countsTracker[0] += 1

for x in range(1, num_nodes):
    count = countsTracker[x]
    '''
    count = 0
    for y1 in range(0,num_nodes): #counts number of already matched pairings with x
		if array[x][y1] == 1:
			count += 1
	'''
    if count >= max_degree:
        continue
	#print(num, x+1, num_nodes)
    possvals = []
    for y2 in range(num_nodes):
        if y2 == x:
            continue
        count2 = countsTracker[y2]
        '''count2 = 0
		for i in array[y2]: #counts number of already matched pairings with y2
			count2 += i
        '''
        if count2 < max_degree:
            possvals.append(y2) #creates the allowable matches list
    if not suppress:
        print(possvals, max_degree-count)
    num = min(max_degree-count, len(possvals))
    for y3 in random.sample(possvals, num):
        array[x][y3] = 1
        array[y3][x] = 1
        countsTracker[x] += 1
        countsTracker[y3] += 1

bits = math.ceil(math.log(num_nodes, 26))

alphabets = list(map(chr, range(ord('A'), ord('Z')+1)))
ids = []

for y in range(num_nodes):
	id = ""
	for x in range(bits):
		id = alphabets[math.floor(y/math.pow(26, x)) % 26] + id
	
	ids.append(id)
	
for x in range(num_nodes):
	file.write(ids[x] + ":")
	done = 0
	for y in range(num_nodes):
		if array[x][y] == 1:
			if done == 1:
				file.write(",")
			done = 1
			file.write(ids[y])
	file.write("\n")
if not suppress:
    print(array)
else:
    print("Done.", name+sparseName)
			
