#!/usr/bin/python
# -*- coding: utf-8 -*-
import sys
import cv2
import math
import numpy as np
import os.path
a=len(sys.argv)
#img=cv2.imread('Reference.png')
filename=sys.argv[1]
img=cv2.imread(filename)
filename=os.path.basename(filename)
[iHeight, iWidth, channels]=img.shape
oWidth=128
oHeight=0
_lambda=1.0
if a>2:
    oWidth = int(sys.argv[2])
if a>3:
    oHeight = int(sys.argv[3])
if a>4:
    _lambda = float(sys.argv[4])



if oWidth  == 0:
    oWidth  = round(iWidth  * oHeight/iHeight)
if oHeight == 0:
    oHeight = round(iHeight * oWidth /iWidth)
outputFilename = str(filename+'_'+str(oWidth)+'x'+str(oHeight)+'_'+str(_lambda)+'.png')
print(filename,oWidth,oHeight,_lambda)
avgImage = np.zeros([oHeight, oWidth, channels])
oImage   = np.zeros([oHeight, oWidth, channels])

pWidth  = iWidth  / oWidth
pHeight = iHeight / oHeight
#calc average image
for py in range(oHeight):
    for px in range(oWidth):
        sx = max(px * pWidth, 0)
        ex = min((px+1) * pWidth, iWidth)
        sy = max(py * pHeight, 0)
        ey = min((py+1) * pHeight, iHeight)

        sxr = math.floor(sx)
        syr = math.floor(sy)
        exr = math.ceil(ex)
        eyr = math.ceil(ey)

        avgF = 0
        
        for iy in range(syr,eyr):
            for ix in range(sxr,exr):
                f=1
                if(ix < sx):
                    f = f * (1.0 - (sx - ix))
                if((ix+1) > ex):
                    f = f * (1.0 - ((ix+1) - ex))
                if(iy < sy):
                    f = f * (1.0 - (sy - iy))
                if((iy+1) > ey):
                    f = f * (1.0 - ((iy+1) - ey))
                avgImage[py, px, :] = avgImage[py, px, :] + (img[iy, ix, :] * f)
                avgF = avgF + f
        avgImage[py, px, :] = avgImage[py, px, :] / avgF
#cv2.imwrite("avg.png", avgImage)
#calc output image
for py in range(oHeight):
    for px in range(oWidth):
        avg=np.zeros([1, channels + 1])
        if(py > 0):
            if(px > 0):
                avg = avg + np.append(np.reshape(avgImage[py-1, px-1,   :], [1,channels]) * 1,1)
            avg = avg + np.append(np.reshape(avgImage[py-1, px+0, :], [1,channels]) * 2,2)
            if((px+1) < oWidth):
                avg = avg + np.append(np.reshape(avgImage[py-1, px+1, :], [1,channels]) * 1,1)
        if(px > 0):
            avg = avg + np.append(np.reshape(avgImage[py+0, px-1,   :], [1,channels]) * 2,2)
        avg = avg + np.append(np.reshape(avgImage[py+0, px+0, :], [1,channels]) * 4,4)
        if((px+1) < oWidth):
            avg = avg + np.append(np.reshape(avgImage[py+0, px+1, :], [1,channels]) * 2,2)

        if((py+1) < oHeight):
            if(px > 0):
                avg = avg + np.append(np.reshape(avgImage[py+1, px-1,   :], [1,channels]) * 1,1)
            avg = avg + np.append(np.reshape(avgImage[py+1, px+0, :], [1,channels]) * 2,2)
            if((px+1) < oWidth):
                avg = avg + np.append(np.reshape(avgImage[py+1, px+1, :], [1,channels]) * 1,1)           
        if avg[0][3]==4:
            print(avg[0][3])
        avg = avg / avg[0][3]
        avg = avg[0][0:channels]
        sx = max(px * pWidth, 0)
        ex = min((px+1) * pWidth, iWidth)
        sy = max(py * pHeight, 0)
        ey = min((py+1) * pHeight, iHeight)

        sxr = math.floor(sx)
        syr = math.floor(sy)
        exr = math.ceil(ex)
        eyr = math.ceil(ey)

        oF = 0

        for iy in range(syr,(eyr)):
            for ix in range(sxr,(exr)):
                if _lambda == 0:
                    f = 1
                else:
                    f=np.linalg.norm(avg - np.reshape(img[iy + 0, ix + 0, :], [1,channels]),2)
                    #f=f/441.6729559
                    f = f **_lambda
                
                if(ix < sx):
                    f = f * (1.0 - (sx - ix))
                if((ix+1) > ex):
                    f = f * (1.0 - ((ix+1) - ex))
                if(iy < sy):
                    f = f * (1.0 - (sy - iy))
                if((iy+1) > ey):
                    f = f * (1.0 - ((iy+1) - ey))
                
                oImage[py + 0, px + 0, :] = oImage[py + 0, px + 0, :] + (img[iy + 0, ix + 0, :] * f)
                oF = oF + f

        if (oF == 0):
            oImage[py + 0, px + 0, :] = avg
        else:
            oImage[py + 0, px + 0, :] = oImage[py + 0, px + 0, :] / oF
cv2.imwrite(outputFilename, oImage)
        
