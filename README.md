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

dpid原版
[nds.png](https://github.com/Mishini/dpid/blob/master/nds.png)
python版
[python.png](https://github.com/Mishini/dpid/blob/master/python.png)
golang版
[golang.png](https://github.com/Mishini/dpid/blob/master/golang.png)
放大后可看出区别
为了方便制作小尺寸头像顺便编译了其他平台的版本，论坛自带的头像处理缩小后太糊了

