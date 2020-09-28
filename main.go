package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	a := len(os.Args)
	//imagePath := "Reference.png"
	imagePath := os.Args[1]
	file, err := os.Open(imagePath)
	_, fileName := filepath.Split(imagePath)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	img, imgType, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	iHeight := b.Max.Y
	iWidth := b.Max.X //?
	channels := 3
	fmt.Println(iWidth, iHeight, channels)
	fmt.Println(imgType)

	var oWidth int = 128
	oHeight := 0
	_lambda := 1.0
	if a > 2 {
		oWidth, _ = strconv.Atoi(os.Args[2])
	}
	if a > 3 {
		oHeight, _ = strconv.Atoi(os.Args[3])
	}
	if a > 4 {
		_lambda, _ = strconv.ParseFloat(os.Args[4], 64)
	}
	if oWidth == 0 {
		oWidth = int(math.Round(float64(iWidth * oHeight / iHeight)))
	}
	if oHeight == 0 {
		oHeight = int(math.Round(float64(iHeight * oWidth / iWidth)))
	}
	outputFilename := fileName + "_" + strconv.FormatInt(int64(oWidth), 10) + "x" + strconv.FormatInt(int64(oHeight), 10) + "_" + strconv.FormatInt(int64(_lambda), 10) + ".png"
	fmt.Println(outputFilename)
	avgImage := make([][][]int, oWidth)
	for i := range avgImage {
		avgImage[i] = make([][]int, oHeight)
		for j := range avgImage[i] {
			avgImage[i][j] = make([]int, channels)
		}
	}
	oImage := make([][][]int, oWidth)
	for i := range oImage {
		oImage[i] = make([][]int, oHeight)
		for j := range oImage[i] {
			oImage[i][j] = make([]int, channels)
		}
	}
	var pWidth float64 = float64(float64(iWidth) / float64(oWidth))
	var pHeight float64 = float64(float64(iHeight) / float64(oHeight))

	//calc average image
	for py := 0; py < oHeight; py++ {
		for px := 0; px < oWidth; px++ {

			sx := math.Max(float64(px)*pWidth, 0)
			ex := math.Min(float64(px+1)*pWidth, float64(iWidth))
			sy := math.Max(float64(py)*pHeight, 0)
			ey := math.Min(float64(py+1)*pHeight, float64(iHeight))

			sxr := math.Floor(sx)
			syr := math.Floor(sy)
			exr := math.Ceil(ex)
			eyr := math.Ceil(ey)

			avgF := 0.0

			for iy := syr; iy < eyr; iy++ {
				for ix := sxr; ix < exr; ix++ {

					f := 1.0
					if ix < sx {
						f = f * (1.0 - (sx - ix))
					}
					if (ix + 1) > ex {
						f = f * (1.0 - ((ix + 1) - ex))
					}
					if iy < sy {
						f = f * (1.0 - (sy - iy))
					}
					if (iy + 1) > ey {
						f = f * (1.0 - ((iy + 1) - ey))
					}
					var pix [4]uint32
					pix[0], pix[1], pix[2], pix[3] = img.At(int(ix), int(iy)).RGBA()

					for p := 0; p < channels; p++ {
						avgImage[px][py][p] = avgImage[px][py][p] + int(float64(pix[p])*f/0x101)
					}
					avgF = avgF + f
				}
			}

			for p := 0; p < channels; p++ {
				avgImage[px][py][p] = int(float64(avgImage[px][py][p]) / avgF)
			}

		}
	}

	avgOut := image.NewRGBA(image.Rect(0, 0, int(oWidth), int(oHeight)))
	for ax := 0; ax < oHeight; ax++ {
		for ay := 0; ay < oWidth; ay++ {
			avgOut.SetRGBA(ay, ax, color.RGBA{uint8(avgImage[ay][ax][0]), uint8(avgImage[ay][ax][1]), uint8(avgImage[ay][ax][2]), 255})

		}
	}

	outputFile, _ := os.Create("avg.png")

	png.Encode(outputFile, avgOut)
	outputFile.Close()
	//calc output image
	for py := 0; py < oHeight; py++ {
		for px := 0; px < oWidth; px++ {
			avg := make([]int, channels+1)
			if px > 0 {
				if py > 0 {
					avg[0] += avgImage[px-1][py-1][0] * 1
					avg[1] += avgImage[px-1][py-1][1] * 1
					avg[2] += avgImage[px-1][py-1][2] * 1
					avg[3] += 1
				}
				avg[0] += avgImage[px-1][py-0][0] * 2
				avg[1] += avgImage[px-1][py-0][1] * 2
				avg[2] += avgImage[px-1][py-0][2] * 2
				avg[3] += 2
				if (py + 1) < oHeight {
					avg[0] += avgImage[px-1][py+1][0] * 1
					avg[1] += avgImage[px-1][py+1][1] * 1
					avg[2] += avgImage[px-1][py+1][2] * 1
					avg[3] += 1
				}
			}
			if py > 0 {
				avg[0] += avgImage[px+0][py-1][0] * 2
				avg[1] += avgImage[px+0][py-1][1] * 2
				avg[2] += avgImage[px+0][py-1][2] * 2
				avg[3] += 2
			}
			avg[0] += avgImage[px-0][py-0][0] * 4
			avg[1] += avgImage[px-0][py-0][1] * 4
			avg[2] += avgImage[px-0][py-0][2] * 4
			avg[3] += 4

			if (py + 1) < oHeight {
				avg[0] += avgImage[px-0][py+1][0] * 2
				avg[1] += avgImage[px-0][py+1][1] * 2
				avg[2] += avgImage[px-0][py+1][2] * 2
				avg[3] += 2
			}
			if (px + 1) < oWidth {
				if py > 0 {
					avg[0] += avgImage[px+1][py-1][0] * 1
					avg[1] += avgImage[px+1][py-1][1] * 1
					avg[2] += avgImage[px+1][py-1][2] * 1
					avg[3] += 1
				}

				avg[0] += avgImage[px+1][py-0][0] * 2
				avg[1] += avgImage[px+1][py-0][1] * 2
				avg[2] += avgImage[px+1][py-0][2] * 2
				avg[3] += 2
				if (py + 1) < oHeight {
					avg[0] += avgImage[px+1][py+1][0] * 1
					avg[1] += avgImage[px+1][py+1][1] * 1
					avg[2] += avgImage[px+1][py+1][2] * 1
					avg[3] += 1
				}
			}
			if avg[3] == 4 {
				fmt.Println(avg[3])
			}
			avg[0] = int(float64(avg[0]) / float64(avg[3]))
			avg[1] = int(float64(avg[1]) / float64(avg[3]))
			avg[2] = int(float64(avg[2]) / float64(avg[3]))
			avg = avg[:3]

			sx := math.Max(float64(px)*pWidth, 0)
			ex := math.Min(float64(px+1)*pWidth, float64(iWidth))
			sy := math.Max(float64(py)*pHeight, 0)
			ey := math.Min(float64(py+1)*pHeight, float64(iHeight))
			sxr := math.Floor(sx)
			syr := math.Floor(sy)
			exr := math.Ceil(ex)
			eyr := math.Ceil(ey)
			oF := 0.0
			for iy := syr; iy < eyr; iy++ {
				for ix := sxr; ix < exr; ix++ {
					var f float64 = 1
					if _lambda == 0 {
						f = 1.0
					} else {
						// max(svd(X)) or 2-norm?
						/*
							var v [3]uint32 //plan b "gonum.org/v1/gonum/mat"
							v[0], v[1], v[2], _ = img.At(int(ix), int(iy)).RGBA()
							//var v64 [3]float64
							//fmt.Println(v)
							//fmt.Println(avg)
							v64 := make([]float64, 3)
							v64[0], v64[1], v64[2] = float64(avg[0])-float64(v[0])/0x101, float64(avg[1])-float64(v[1])/0x101, float64(avg[2])-float64(v[2])/0x101
							zero := mat.NewDense(1, 3, v64)
							f = mat.Norm(zero, 2)*/
						r, g, b, _ := img.At(int(ix), int(iy)).RGBA() //plan a
						f = math.Sqrt(float64(math.Pow(float64(avg[0])-float64(r)/0x101, 2) + math.Pow(float64(avg[1])-float64(g)/0x101, 2) + math.Pow(float64(avg[2])-float64(b)/0x101, 2)))

						//f = 1.0
						//fmt.Println(r, g, b)
						//f=f/441.6729559
						f = math.Pow(f, _lambda)
					}

					if ix < sx {
						f = f * (1.0 - (sx - ix))
					}
					if (ix + 1) > ex {
						f = f * (1.0 - ((ix + 1) - ex))
					}
					if iy < sy {
						f = f * (1.0 - (sy - iy))
					}
					if (iy + 1) > ey {
						f = f * (1.0 - ((iy + 1) - ey))
					}

					oF = oF + f
					var pix [4]uint32
					pix[0], pix[1], pix[2], pix[3] = img.At(int(ix+0), int(iy+0)).RGBA()

					for p := 0; p < channels; p++ {
						oImage[px][py][p] = oImage[px][py][p] + int(float64(pix[p])*f/0x101)
					}
				}
			}
			if oF == 0 {
				for p := 0; p < channels; p++ {
					oImage[px][py][p] = avgImage[px][py][p]
				}
			} else {

				for p := 0; p < channels; p++ {
					oImage[px][py][p] = int(float64(oImage[px][py][p]) / oF)
				}
			}
		}
	}
	oOut := image.NewRGBA(image.Rect(0, 0, int(oWidth), int(oHeight)))
	for i := 0; i < oHeight; i++ {
		for j := 0; j < oWidth; j++ {
			oOut.SetRGBA(j, i, color.RGBA{uint8(oImage[j][i][0]), uint8(oImage[j][i][1]), uint8(oImage[j][i][2]), 255})
		}
	}
	o2, _ := os.Create(outputFilename)
	png.Encode(o2, oOut)
	o2.Close()
}
