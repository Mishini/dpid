# Description
"Rapid, Detail-Preserving Image Downscaling" Implementations in Python and Goland
# Examples
## python:
    Python dpid.py "myImage.jpg" 

## Golang:
    dpid "myImage.jpg"              // downscales using default values
    dpid "myImage.jpg" 256          // downscales to 256px width, keeping aspect ratio
    dpid "myImage.jpg" 0 256        // downscales to 256px height, keeping aspect ratio
    dpid "myImage.jpg" 128 0 0.5    // downscales to 128px width, keeping aspect ratio, using lamdba=0.5
    dpid "myImage.jpg" 128 128      // downscales to 128x128px, ignoring aspect ratio
|  org 200x300   | dpid lamdba=0.5  |dpid lamdba=0.8|
|  ----  | ----  |----|
| ![raw200](https://github.com/Mishini/dpid/blob/master/AuroraTree_Wallace_200.png)  | ![dpid 0.5](https://github.com/Mishini/dpid/blob/master/AuroraTree_Wallace_2048.jpg_200x300_0.5.png) |![dpid 0.8](https://github.com/Mishini/dpid/blob/master/AuroraTree_Wallace_2048.jpg_200x300_0.8.png)|


