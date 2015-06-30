package main

import "fmt"
import "os"
import "strconv"
import "path"
import "io"
import "log"
import "math"
import "net/http"
import "github.com/codegangsta/cli"

import (
  "image"
  _  "image/color"
  //"image/jpeg"
  "image/png"
  )

// example image
// https://137.110.172.64/GoogleMaps2/18585/Result/t/q/s/tqsrrrqt.jpg

func p( t string) {
  fmt.Println("> ", t)
}

func main() {

     app := cli.NewApp()
     app.Name    = "talk"
     app.Usage   = "2D image segmentation\n\n" +
                   "   Program to describe an image of many small objects. Each object is filtered by\n" +
                   "   size and aspect ratio. Objects that fulfill the criteria of the filter are printed\n" +
                   "   on standard-out as well highlighted in an output image."
     app.Version = "0.0.1"
     app.Author  = "Hauke Bartsch"
     app.Email   = "HaukeBartsch@gmail.com"
     app.Flags = []cli.Flag {
       cli.BoolFlag {
         Name:  "verbose",
         Usage: "Generate verbose output with intermediate files.",
       },
     }

     app.Commands = []cli.Command {
       {
         Name: "segment1",
         ShortName: "s1",
         Usage: "Detect dark regions in the input image with:\n\t      ./talk segment1 tqsrrrqt.jpg 10 28 3",
         Description: "Compute a label field, save it as a color image and a list all found segments on stdout.\n\n" +
                      "   This command requires four arguments, the file name, minimum number of pixel of valid\n" +
                      "   objects, maximum number of pixel of valid objects and the maximum aspect ratio allowed\n" +
                      "   for a valid object. If the aspect ratio is negative it specifies the minimum allowed\n" +
                      "   aspect ratio. The compactness value prefers more circular objects for smaller values.",
         Flags: []cli.Flag{
           cli.IntFlag {
             Name: "meansize",
             Value: 13,
             Usage: "Size of the region from which the local mean intensity is calculated",
           },
           cli.Float64Flag {
             Name: "compactness",
             Value: -2.0, // default value that prevents compactness calculation
             Usage: "Filter by compactness [1..0] defined by P^2/(4 pi A) where A is area and P is perimeter",
           },
           cli.Float64Flag {
             Name: "focussize",
             Value: 1.8,
             Usage: "Size in pixel that we focus on, structures larger and smaller are blurred",
           },
           cli.BoolFlag {
             Name: "notinvert",
             Usage: "Do not invert the image (default is to invert, detect dark objects)",
           },
         },
         Action: func(c *cli.Context) {
           if len(c.Args()) < 4 {
             fmt.Printf("  Error: need an image name, a minimum (10) and maximum (28) size in pixel\n" +
                        "  and a maximum allowed aspect ration (3)\n\n")
           } else {
             verbose     := c.GlobalBool("verbose")
             invert_flag := !c.Bool("notinvert") 
             if (verbose) {
               p("verbose on")
               p("run segment1")
             }

             // call segment
             var file *os.File
             if _, err := os.Stat(c.Args()[0]); err == nil { // read using direct io
               file, err = os.Open(c.Args()[0])
               if err != nil {
                log.Fatal(err)
               }
             } else { // this part only works if https has a valid non-self-signed certificate
                if verbose {
                  p("try to download file")
                }
                out, _ := os.Create(".download")
                defer out.Close()
                resp, err := http.Get(c.Args()[0])
                if err != nil {
                  log.Fatal(err)
                }
                defer resp.Body.Close()
                _, err = io.Copy(out, resp.Body)
                if err != nil {
                  log.Fatal(err)
                }
                file, err = os.Open(".download")
                if err != nil {
                  log.Fatal(err)
                }
             }

             defer file.Close()
             img, _, err := image.Decode(file)
             if err != nil {
               log.Fatal(err)
             }
             // how big is the image?
             var size image.Rectangle
             size = img.Bounds()
             p(fmt.Sprintf("image size is %d by %d pixel", size.Max.X-size.Min.X, size.Max.Y-size.Min.Y))
             mmin,_,_ := getMin(img)
             mmax,_,_ := getMax(img)
             p(fmt.Sprintf("min/max: %f %f", mmin, mmax))
             d, f  := path.Split(c.Args()[0])

             lST,_ := strconv.ParseInt(c.Args()[1], 0, 32)
             hST,_ := strconv.ParseInt(c.Args()[2], 0, 32)
             meansizevalue := c.Int("meansize")
             if !c.IsSet("meansize") {
                // we can do better by getting a good mean size (twice radius) based on size range
                meansizevalue = int(math.Floor(math.Sqrt( (float64(hST))/3.1415927 ) * 2.0 + 0.5))*4
             }

             meanoff := subMean(img, meansizevalue)
             if verbose {
                p(fmt.Sprintf("mean size used: %d", meansizevalue))
                // save the meanoff image
                fn    := path.Join(d, f[0:len(f)-len(path.Ext(f))] + "_001_meanoff.png")
                out, err := os.Create(fn)
                if err != nil {
                  fmt.Println(err)
                  os.Exit(1)
                }
                p(fmt.Sprintf("write out mean removed image %s", fn))
                png.Encode(out, meanoff)
                out.Close()
             }

             focussizevalue := c.Float64("focussize")
             if !c.IsSet("focussize") {
                focussizevalue = ( math.Sqrt( (float64(lST))/3.1415927 ) * 2.0 ) / 2.0
             }

             focusI := focus(meanoff,float32(focussizevalue))
             if invert_flag {
               p(fmt.Sprintf("invert the image before segmentation"))
               focusI = invert(focusI)
             }
             if verbose {
                p(fmt.Sprintf("focus size used: %g", float32(focussizevalue)))
                // save the focussed image
                fn    := path.Join(d, f[0:len(f)-len(path.Ext(f))] + "_002_focus.png")
                out, err := os.Create(fn)
                if err != nil {
                  fmt.Println(err)
                  os.Exit(1)
                }
                p(fmt.Sprintf("write out the focused and inverted image %s", fn))
                png.Encode(out, focusI)
                out.Close()
             }

             // try some segmentation (size threshold from 10 to 28, aspect ration smaller than 3)
             ar,_  := strconv.ParseFloat(c.Args()[3], 64)
             if verbose {
                if ar < 0 {
                  p(fmt.Sprintf("segment1 with size thresholds %d .. %d, and minimum aspect ratio of %g", lST, hST, ar))
                } else {
                  p(fmt.Sprintf("segment1 with size thresholds %d .. %d, and maximum aspect ratio of %g", lST, hST, ar))
                }
             }
             compactness := c.Float64("compactness")
             if c.IsSet("compactness") && verbose {
                if compactness > 0 {
                   p(fmt.Sprintf("compactness filter, values have to be smaller than %g to be compact enough", compactness))
                } else {
                   p(fmt.Sprintf("compactness filter, values have to be larger than %g to be less compact", -compactness))
                }
             }
             seg   := segment1(focusI, int(lST), int(hST), ar, compactness)
             fn    := path.Join(d, f[0:len(f)-len(path.Ext(f))] + "_seg.png")
             out, err := os.Create(fn)
             if err != nil {
                fmt.Println(err)
                os.Exit(1)
             }
             if verbose {
                p(fmt.Sprintf("write out the found segmentation %s", fn))
             }
             png.Encode(out, seg)
             out.Close()
           }
         },
       },
     }
     app.Run(os.Args)
}
