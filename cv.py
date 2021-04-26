# -*- coding: utf-8 -*-
"""
Created on Thu Apr 22 07:31:42 2021

@author: tyler
"""
import math

def calcColor(me, parent):
	print("calc color", me, parent)
	i = calcIndex(me, parent) #1
	shamt = int(math.ceil(math.log2(max(i,1))))
	return (((me >> i) & 1) << shamt) + i

def calcColorAlt(me, parent):
    #print("calc color alt", me, parent)
    temp1 = me ^ parent
    trail = trailing_zeros(temp1)
    orig_bit = 1 & (me >> trail)
    return orig_bit | trail << 1
    print(trailing_zeros(temp1))

def trailing_zeros(longint):
    manipulandum = str(longint)
    return len(manipulandum)-len(manipulandum.rstrip('0'))

def calcIndex(me, parent):
    print("\t",me, parent)
    if me % 2 != parent % 2:
    	return 0
    else:
    	return 1 + calcIndex(me >> 1, parent >> 1)

if __name__ == "__main__":
    print("CalcIndex: ", calcIndex(1, 7))
    print("calcColorAlt: ", calcColorAlt(1, 7), "\n\n")
    print("CalcIndex: ", calcIndex(10, 11))
    print("calcColorAlt: ", calcColorAlt(10, 11), "\n\n")
    print("CalcIndex: ", calcIndex(11, 10))
    print("calcColorAlt: ", calcColorAlt(11, 10), "\n\n")
    print("CalcIndex: ", calcIndex(8, 11))
    print("calcColorAlt: ", calcColorAlt(8, 11), "\n\n")
    
    for k in range(7):
        print("calcColorAlt: ", k, ": ", calcColorAlt(k, 10), "\n\n")