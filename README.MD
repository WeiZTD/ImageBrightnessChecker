# Image Brightness Checker
A simple CLI tool to get average brightness of image

## Usage
Show usage message
```
imgBrightChk.exe -?
```
Get brightness from original image file. Range within 0 ~ 255 .
```
imgBrightChk.exe -p pathtoimage.jpg
```
Resize image width & height before calcuate brighness.
If only given one border setting, another will be same as original image.
```
imgBrightChk.exe -p pathtoimage.jpg -w 500 -h 500
```
```
imgBrightChk.exe -p pathtoimage.jpg -w 500 (-h = original image height)
```
## Built With
 [Imaging](httpsgithub.comdisintegrationimaging) - image processing package for Go
