package main

import (
  "math"
  "math/rand"
  "image"
  "fmt"
  "image/color"
  _  "image/jpeg"
  _  "image/png"
  "github.com/disintegration/gift"
)

// return the coordinate of the min location as well
func getMin( img image.Image ) (float64, int, int) {
  bounds := img.Bounds()
  minVal := 65536.0
  minx := bounds.Min.X
  miny := bounds.Min.Y
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      r,g,b,_ := img.At(x,y).RGBA()
      avg := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
      // gray := color.Gray{uint8(math.Ceil(avg))}
      if (avg < minVal) {
        minVal = avg
        minx = x
        miny = y
      }
    }
  }

  return minVal,minx,miny
}

func getMax( img image.Image ) (float64,int,int) {
  bounds := img.Bounds()
  maxVal := 0.0
  maxx := bounds.Min.X
  maxy := bounds.Min.Y
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      r,g,b,_ := img.At(x,y).RGBA()
      avg := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
      //gray := color.Gray{uint8(math.Ceil(avg))}
      if (avg > maxVal) {
        maxVal = avg
        maxx = x
        maxy = y
      }
    }
  }

  return maxVal, maxx, maxy
}

func blur( img image.Image, howmuch float32 ) image.Image {
  g := gift.New( gift.Grayscale() );
  g.Add(gift.GaussianBlur(howmuch))
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)
  return(dst)
}

func mean( img image.Image, disk int ) image.Image {
  g := gift.New( gift.Grayscale() )
  g.Add(gift.Mean(disk, false)) // use square neighborhood
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)
  return(dst)
}

// a type of image that is a matrix of floats
type imageF [][]float32

func meanF( img image.Image, disk int ) imageF {

  g := gift.New( gift.Grayscale() )
  g.Add(gift.Mean(disk, false)) // use square neighborhood
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)

  // now convert this to float array of arrays
  floatData := make([][]float32, dst.Bounds().Max.Y-dst.Bounds().Min.Y)
  for i := range floatData {
    floatData[i] = make([]float32, dst.Bounds().Max.X-dst.Bounds().Min.X)
  }
  bounds := dst.Bounds()
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      pr,_,_,_ := dst.At(x,y).RGBA()
      floatData[x][y] = float32(pr) // 0.2125*float32(pr) + 0.7154*float32(pg) + 0.0721*float32(pb)
    }
  }

  return floatData
}


func varianceF( img image.Image, disk int ) (imageF, imageF) {
  m := meanF(img, disk) // gets a grayscale copy of local mean

  // create a grayscale version of the original
  //g := gift.New( gift.Grayscale() )
  //v := image.NewRGBA(g.Bounds(img.Bounds()))
  //g.Draw(v, img)

  g := gift.New( gift.Grayscale() )
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)

  bounds := img.Bounds()
  floatData := make([][]float32, bounds.Max.Y-bounds.Min.Y)
  for i := range floatData {
    floatData[i] = make([]float32, bounds.Max.X-bounds.Min.X)
  }
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      p1r,p1g,p1b,_ := dst.At(x,y).RGBA()
      g1 := 0.2125*float64(p1r) + 0.7154*float64(p1g) + 0.0721*float64(p1b)
      g2 := float64(m[x][y])

      floatData[x][y] = float32((g1-g2) * (g1-g2))
    }
  }
  return m,floatData
}


func variance( img image.Image, disk int ) (image.Image, image.Image) {
  m := mean(img, disk) // gets a grayscale copy

  // create a grayscale version of the original
  g := gift.New( gift.Grayscale() )
  v := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(v, img)

  bounds := img.Bounds()
  dst := image.NewGray(bounds)
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      p1r,p1g,p1b,_ := v.At(x,y).RGBA()
      p2r,p2g,p2b,_ := m.At(x,y).RGBA()
      g1 := 0.2125*float64(p1r) + 0.7154*float64(p1g) + 0.0721*float64(p1b)
      g2 := 0.2125*float64(p2r) + 0.7154*float64(p2g) + 0.0721*float64(p2b)

      dst.Set(x,y,color.Gray16{ uint16(math.Sqrt((g1-g2) * (g1-g2))) })
      //fmt.Println("value: ", math.Sqrt((g1-g2) * (g1-g2)))
    }
  }
  return m,dst
}

func subMean( img image.Image, disk int ) (image.Image) {
  m := mean(img, disk)

  bounds := img.Bounds()

  floatData := make([][]float32, bounds.Max.Y-bounds.Min.Y)
  for i := range floatData {
    floatData[i] = make([]float32, bounds.Max.X-bounds.Min.X)
  }

  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      p1r,p1g,p1b,_ := img.At(x,y).RGBA()
      p2r,p2g,p2b,_ := m.At(x,y).RGBA()
      g1 := 0.2125*float64(p1r) + 0.7154*float64(p1g) + 0.0721*float64(p1b)
      g2 := 0.2125*float64(p2r) + 0.7154*float64(p2g) + 0.0721*float64(p2b)
      floatData[x][y] = float32(g1-g2)
    }
  }

  min := floatData[0][0]
  max := floatData[0][0]
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      if min > floatData[x][y] {
        min = floatData[x][y]
      }
      if max < floatData[x][y] {
        max = floatData[x][y]
      }
    }
  }

  g := gift.New( gift.Grayscale() )
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)

  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      dst.Set(x,y,color.Gray16{ uint16( (floatData[x][y] - min )/(max-min) * 65535) })
    }
  }
  return dst
}

// take an image and revert its rgb channels (255-channel)
func invert( img image.Image ) (image.Image) {

  g := gift.New( gift.Grayscale(),
    gift.Invert() )
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)

  return dst
}

// try to focus on a specific scale by using a mexican hat filter
func focus( img image.Image, s float32 ) (image.Image) {

  onethird := s/3.0
  small := blur(img,s-onethird)
  large := blur(img,s+onethird)

  bounds := img.Bounds()
  floatData := make([][]float32, bounds.Max.Y-bounds.Min.Y)
  for i := range floatData {
    floatData[i] = make([]float32, bounds.Max.X-bounds.Min.X)
  }

  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      p1r,p1g,p1b,_ := small.At(x,y).RGBA()
      p2r,p2g,p2b,_ := large.At(x,y).RGBA()
      g1 := 0.2125*float64(p1r) + 0.7154*float64(p1g) + 0.0721*float64(p1b)
      g2 := 0.2125*float64(p2r) + 0.7154*float64(p2g) + 0.0721*float64(p2b)
      floatData[x][y] = float32(g1-g2)
    }
  }

  min := floatData[0][0]
  max := floatData[0][0]
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      if min > floatData[x][y] {
        min = floatData[x][y]
      }
      if max < floatData[x][y] {
        max = floatData[x][y]
      }
    }
  }

  g := gift.New( gift.Grayscale() )
  dst := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(dst, img)

  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      dst.Set(x,y,color.Gray16{ uint16( (floatData[x][y] - min )/(max-min) * 65535) })
    }
  }
  return dst
}

// subtract mean and divide by variance
// in order for this to work we have to use floating point artihmetic
func whitening( img image.Image, disk int ) (image.Image, image.Image, image.Image) {

  mean, vari := varianceF(img, disk)

  bounds := img.Bounds()
  floatData := make([][]float32, bounds.Max.Y-bounds.Min.Y)
  for i := range floatData {
    floatData[i] = make([]float32, bounds.Max.X-bounds.Min.X)
  }
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      g1 := float64(mean[x][y])
      g2 := float64(vari[x][y])
      p3r,p3g,p3b,_ := img.At(x,y).RGBA()
      g3 := 0.2125*float64(p3r) + 0.7154*float64(p3g) + 0.0721*float64(p3b)

      if math.Abs(g2) < 1e-6 { // close to zero
        floatData[x][y] = float32(0.0)
      } else {
        floatData[x][y] = float32((g3-g1)*1.0/math.Sqrt(g2))
        fmt.Println("floatData: ", g3, ",", g1,",", g2, ".....", (g3-g1)/g2)
      }
    }
  }
  // now convert all three results back into images
  // scale everything to min max of 16bit gray

  meanImage := image.NewGray(bounds)
  // find min/max
  min := mean[0][0]
  max := mean[0][0]
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      if min > mean[x][y] {
        min = mean[x][y]
      }
      if max < mean[x][y] {
        max = mean[x][y]
      }
    }
  }
  fmt.Println(" mean: ", min, max)
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      meanImage.Set(x,y,color.Gray16{ uint16( (mean[x][y]-min)/(max-min) * 65535 )})
    }
  }


  min = vari[0][0]
  max = vari[0][0]
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      if min > vari[x][y] {
        min = vari[x][y]
      }
      if max < vari[x][y] {
        max = vari[x][y]
      }
    }
  }
  fmt.Println(" vari: ", min, max)
  variImage := image.NewGray(bounds)
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      variImage.Set(x,y,color.Gray16{ uint16( (vari[x][y]-min)/(max-min) * 65535 )})
    }
  }


  min = floatData[0][0]
  max = floatData[0][0]
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      if min > floatData[x][y] {
        min = floatData[x][y]
      }
      if max < floatData[x][y] {
        max = floatData[x][y]
      }
    }
  }
  //min = -.01
  //max = .01
  fmt.Println(" white: ", min, max)
  whitenedImage := image.NewGray(bounds)
  for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for  x := bounds.Min.X; x < bounds.Max.X; x++ {
      val :=  (floatData[x][y]-min)/(max-min) * 65535
      if val > 65535 {
        val = 65535
      }
      if val < 0 {
        val = 0
      }
      whitenedImage.Set(x,y,color.Gray16{ uint16( val )})
      //fmt.Println("val: ", floatData[x][y], "-", min, "-",max, "-",(floatData[x][y]-min)/(max-min)* 65535)
    }
  }

  return meanImage,variImage,whitenedImage
}

func histogram( img image.Image, bins int ) []int {
   mi,_,_ := getMin(img)
   ma,_,_ := getMax(img)

   //fmt.Printf("max and min are : %v %v\n", mi, ma)
   bounds := img.Bounds()
   var h []int
   h = make([]int, bins, bins)

   for  y := bounds.Min.Y; y < bounds.Max.Y; y++ {
     for  x := bounds.Min.X; x < bounds.Max.X; x++ {
       r,g,b,_ := img.At(x,y).RGBA()
       a := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
       idx := int(math.Floor((a-mi)/(ma-mi)*(float64(bins)-1.0)))
       if idx < 0 {
        idx = 0
       }
       if idx >= bins {
        idx = bins-1
       }
       h[idx]++
     }
   }
   return h
}

// calculate quantiles of the cummulative distribution of the histogram
// returns for each quantile the intensity
func quantiles(img image.Image, quant []float64 ) []float64 {
  mi,_,_ := getMin(img)
  ma,_,_ := getMax(img)
  h := histogram(img, 512)

  cumhist := make([]float64, 512)
  cumhist[0] = float64(h[0])
  for i := 1; i < len(h); i++ {
    cumhist[i] = float64(h[i]) + cumhist[i-1]
  }
  // scale by last entry == 1
  for i := 0; i < len(h); i++ {
    cumhist[i] = cumhist[i]/cumhist[len(h)-1];
  }
  erg := make([]float64, len(quant))
  for i := 0; i < len(quant); i++ {
    for j := 0; j < len(cumhist); j++ {
       if cumhist[j] > quant[i] {
          erg[i] = mi + ((ma - mi)/float64(len(h)) * float64(j-1))
          break;
       }
    }
  }

  return erg
}

func Otzu( img image.Image ) float64 {
  mi,_,_ := getMin(img)
  ma,_,_ := getMax(img)
  h := histogram(img, 256)

  bounds := img.Bounds()
  total := float64(bounds.Max.X - bounds.Min.X) * float64(bounds.Max.Y - bounds.Min.Y)

  sum := 0.0
  for i := 1; i < 256; i++ {
    //fmt.Printf("%d, ", h[i])
    sum = sum + float64(i) * float64(h[i])
  }
  var sumB = 0.0
  var wB = 0.0
  var wF = 0.0
  var mB = 0.0
  var mF = 0.0
  var max = 0.0
  var between = 0.0
  var threshold1 = 0.0
  var threshold2 = 0.0
  for i := 0; i < 256; i++ {
    wB = wB + float64(h[i])
    if wB == 0 {
      continue
    }
    wF = total - wB
    if wF == 0 {
      break;
    }
    sumB = sumB + float64(i) * float64(h[i])
    mB = sumB / wB
    mF = (sum - sumB) / wF
    between = wB * wF * (mB - mF) * (mB - mF)
    if between >= max {
      threshold1 = float64(i);
      if between > max {
        threshold2 = float64(i)          
      }
      max = between
    }
  }

  return (float64(threshold1 + threshold2) / 2.0)/256*(ma-mi)+mi
}

// create a color from HSV values (H is in 0...360, S and V are 0 to 1)
func Hsv(H, S, V float64) color.Color {
    Hp := H/60.0
    C := V*S
    X := C*(1.0-math.Abs(math.Mod(Hp, 2.0)-1.0))

    m := V-C;
    r, g, b := 0.0, 0.0, 0.0

    switch {
    case 0.0 <= Hp && Hp < 1.0: r = C; g = X
    case 1.0 <= Hp && Hp < 2.0: r = X; g = C
    case 2.0 <= Hp && Hp < 3.0: g = C; b = X
    case 3.0 <= Hp && Hp < 4.0: g = X; b = C
    case 4.0 <= Hp && Hp < 5.0: r = X; b = C
    case 5.0 <= Hp && Hp < 6.0: r = C; b = X
    }

    return color.RGBA{uint8(255*(m+r)), uint8(255*(m+g)), uint8(255*(m+b)), 255}
}

//
// Segmentation by region growing
//
// This function will use a series of tresholds and perform a 
// threshold based segmentation (everything that is brighter than the threshold).
// For each found area it will calculate how large the area is (number of pixel)
// and how compact the area is (aspect ratio of smallest to longest eigenvalue from PCA).
// Only regions that agree with the segmentation parameters are kept. The operation is repeated with
// the next threshold level until the maximum threshold is reached.
// The selection of thresholds is done based on percentile of pixel above that
// threshold. That should ensure that many regions close together in luminance get more
// finely segmented thresholds (makes the procedure faster as well).
//
func segment1( img image.Image, lowSizeThreshold int, highSizeThreshold int, aspectRatioThreshold float64 ) (image.Image) {

  // calculate n quantiles of the cummulative distribution
  quants := []float64{}
  for i := 0; i < 200; i++ {
     quants = append( quants, 0.05 + (0.99-0.05)/200.0*float64(i) )
  }

  thresholds := quantiles(img, quants)

  // set a uniform color for the output segmentation field
  g := gift.ColorFunc(
                    func(r0, g0, b0, a0 float32) (r,g,b,a float32) {
                     r = 0
                     g = 0
                     b = 0
                     a = a0
                     return
                    },
  )
  seg := image.NewRGBA(g.Bounds(img.Bounds()))
  g.Draw(seg, img, nil);

  bounds := img.Bounds()

  // lets remember if we have visited a location before
  vis := make([][]uint8, bounds.Max.Y-bounds.Min.Y)
  randGenerator := rand.New(rand.NewSource(99))
  currentLabel  := 0

  for _,threshold := range thresholds {

    // reset our visit buffer (we need to visit everything again)
    for i := range vis {
      vis[i] = make([]uint8, bounds.Max.X-bounds.Min.X)
      for  j := 0; j < bounds.Max.X-bounds.Min.X; j++ {
        vis[i][j] = 0 // nothing visited yet
      }
    }

    // we have to work off a queue of border pixel to fill our image
    queue := make([][]int,1)
    queue[0] = []int{ 0, 0 }
    vis[0][0] = 1
    for len(queue) > 0 {
      // take out the first element
      x := queue[0][0]
      y := queue[0][1]
      queue = queue[1:] // keep all the other elements

      r,g,b,_ := seg.At(x,y).RGBA()        
      a := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)        
      if a > 0 { // don't do anything if that pixel has been visited before
        continue
      }
      r,g,b,_ = img.At(x,y).RGBA()
      avg := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
      // fmt.Printf("avg: %v threshold %v\n", avg, threshold)

      if avg > threshold {
        // now we need to do something because this pixel is not yet part of a segment, but it should belong to one
       
        queue2 := make([][]int,1)
        queue2[0] = []int{ x, y }
        vis[0][0] = 1;
        currentSegment := make([][]int,1)
        currentSegment[0] = []int{ x, y }

        for len(queue2) > 0 {
          x2 := queue2[0][0]
          y2 := queue2[0][1]
          queue2 = queue2[1:] // keep all the other elements
          r,g,b,_ = seg.At(x2,y2).RGBA()
          a := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)      
          if a > 0 { // don't do anything if that pixel has a non-zero value already
            //fmt.Printf("found a pixel with label %v at %v %v\n", a, x2, y2)
            continue
          }
          r,g,b,_ := img.At(x2,y2).RGBA()
          avg := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)      
          if avg > threshold {
             currentSegment = append(currentSegment, []int{ x2, y2 })

             // check all the neighbors
             if x2-1 > bounds.Min.X && vis[x2-1][y2] == 0 {
                queue2 = append(queue2, []int{ x2-1,y2 })       
                queue  = append(queue, []int{ x2-1,y2 })       
                vis[x2-1][y2] = 1
             }
             if x2+1 < bounds.Max.X && vis[x2+1][y2] == 0 {
               queue2 = append(queue2, []int{ x2+1,y2 })
               queue = append(queue, []int{ x2+1,y2 })
               vis[x2+1][y2] = 1
             }
             if y2-1 > bounds.Min.Y && vis[x2][y2-1] == 0 {
               queue2 = append(queue2, []int{ x2,y2-1 })
               queue = append(queue, []int{ x2,y2-1 })
               vis[x2][y2-1] = 1
             }
             if y2+1 < bounds.Max.Y && vis[x2][y2+1] == 0 {
               queue2 = append(queue2, []int{ x2,y2+1 })
               queue = append(queue, []int{ x2,y2+1 })
               vis[x2][y2+1] = 1
             }
          }
        }
        // if the size of the region is too large, don't put it in the label
        if len(currentSegment) < highSizeThreshold && len(currentSegment) > lowSizeThreshold {
            // start with creating a new label id
            col := Hsv(randGenerator.Float64()*360, 0.8, 0.5)

            // set back to background
            var centerOfMassX float64 = 0.0
            var centerOfMassY float64 = 0.0

            for j := range currentSegment {
               centerOfMassX = centerOfMassX + float64(currentSegment[j][0])
               centerOfMassY = centerOfMassY + float64(currentSegment[j][1])
            }

            // create a region marker (should be center and size)
            centerOfMassX = centerOfMassX / float64(len(currentSegment))
            centerOfMassY = centerOfMassY / float64(len(currentSegment))

            currentSegmentZeroMean := make([][]int,len(currentSegment))
            for k:= range currentSegmentZeroMean {
              currentSegmentZeroMean[k] = []int{ currentSegment[k][0] - int(centerOfMassX), currentSegment[k][1] - int(centerOfMassY) }
            }

            covar := make([][]float64, 2)
            covar[0] = make([]float64, 2)
            covar[1] = make([]float64, 2)
            for k := 0; k < 2; k++ {
              for l := 0; l < 2; l++ {
                covar[k][l] = 0.0
                 for j := range currentSegmentZeroMean {
                    covar[k][l] = covar[k][l] + float64(currentSegmentZeroMean[j][l]) * float64(currentSegmentZeroMean[j][k])
                 }
                 covar[k][l] = covar[k][l] / float64(len(currentSegment))
               }
            }
            //fmt.Printf(" %v %v %v %v\n", covar[0][0], covar[0][1], covar[1][0], covar[1][1])
            T := covar[0][0] + covar[1][1]
            D := covar[0][0]*covar[1][1] - covar[1][0]*covar[0][1]
            L1 := T/2.0 + math.Sqrt((T*T) / 4.0 - D)
            L2 := T/2.0 - math.Sqrt((T*T) / 4.0 - D)

            // and save the segment as result
            // first read the old locations
            if aspectRatioThreshold < 0 {
              if L1/L2 > -aspectRatioThreshold {
                fmt.Printf("i: %d, x: %v, y: %v, s: %v, a: %.4f\n", currentLabel, int(math.Floor(centerOfMassX + .5)), int(math.Floor(centerOfMassY + .5)), len(currentSegment),
                  L1/L2)
                currentLabel = currentLabel + 1
                for j := range currentSegment {
                   seg.Set(currentSegment[j][0], currentSegment[j][1], col)
                }
              }
            } else {
              if L1/L2 <= aspectRatioThreshold {
                fmt.Printf("i: %d, x: %v, y: %v, s: %v, a: %.4f\n", currentLabel, int(math.Floor(centerOfMassX + .5)), int(math.Floor(centerOfMassY + .5)), len(currentSegment),
                  L1/L2)
                currentLabel = currentLabel + 1
                for j := range currentSegment {
                   seg.Set(currentSegment[j][0], currentSegment[j][1], col)
                }
              }
            }

        }
      }

      // add all the neighbors (if they have not been visited yet)
      if x-1 > bounds.Min.X && vis[x-1][y] == 0 {
        queue = append(queue, []int{ x-1,y })
        vis[x-1][y] = 1
      }
      if x+1 < bounds.Max.X  && vis[x+1][y] == 0 {
        queue = append(queue, []int{ x+1,y })
        vis[x+1][y] = 1
      }
      if y-1 > bounds.Min.Y  && vis[x][y-1] == 0 {
        queue = append(queue, []int{ x,y-1 })
        vis[x][y-1] = 1
      }
      if y+1 < bounds.Max.Y && vis[x][y+1] == 0 {
        queue = append(queue, []int{ x,y+1 })
        vis[x][y+1] = 1
      }
    }
  }

  return(seg)
}

// what about center-surround color cells? What about detecting larger regions of interest
// we could do a color segmentation to get smooth muscle from background from cells





