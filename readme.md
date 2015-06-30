talk - Segmentation of 2D images
==================================

Talk is a command line utility to compute regions of interest in a 2D image. It is useful if many small
objects can be detected by a size range a maximum aspect ratio and a maximum compactness value. 
As an output the list of found objects is printed together with a label image that highlights
each of the found objects.

You can download a compiled version of the program, or build it yourself from source using golang.

* Linux
   wget https://github.com/HaukeBartsch/talk/raw/master/binary/Linux64/talk
* MacOS
   wget https://github.com/HaukeBartsch/talk/raw/master/binary/MacOS/talk

<img title="input image" src="tqsrrrqt.jpg"></img>
<img title="image after removing the local mean intensities" src="tqsrrrqt_001_meanoff.png"></img>
<img title="image after mexican hat filter focused on size range" src="tqsrrrqt_002_focus.png"></img>
<img title="output image with segmented regions (in color)" src="tqsrrrqt_seg.png"></img>

Program Help:

```
NAME:
   talk - 2D image segmentation

   Program to describe an image of many small objects. Each object is filtered by
   size and aspect ratio. Objects that fulfill the criteria of the filter are printed
   on standard-out as well highlighted in an output image.

USAGE:
   talk [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
  Hauke Bartsch - <HaukeBartsch@gmail.com>

COMMANDS:
   segment1, s1	Detect dark regions in input image with:
		      ./talk segment1 tqsrrrqt.jpg 10 28 3
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --verbose		Generate verbose output with intermediate files.
   --help, -h		show help
   --version, -v	print the version
```

Command line examples
----------------------

The segment1 command will normalize the image and sharpen it using a mexican hat filter. Region growing is executed at different threshold levels and for each threshold value a region growing will extract objects that are filtered by min/max size and aspect ratio.

```
./talk segment1 --help
NAME:
   segment1 - Detect dark regions in the input image with:
         ./talk segment1 tqsrrrqt.jpg 10 28 3

USAGE:
   command segment1 [command options] [arguments...]

DESCRIPTION:
   Compute a label field, save it as a color image and a list all found segments on stdout.

   This command requires four arguments, the file name, minimum number of pixel of valid
   objects, maximum number of pixel of valid objects and the maximum aspect ratio allowed
   for a valid object. If the aspect ratio is negative it specifies the minimum allowed
   aspect ratio. The compactness value prefers more circular objects for smaller values.

OPTIONS:
   --meansize "13"   Size of the region from which the local mean intensity is calculated
   --compactness "-2"   Filter by compactness [1..0] defined by P^2/(4 pi A) where A is area and P is perimeter
   --focussize "1.8" Size in pixel that we focus on, structures larger and smaller are blurred
   --notinvert    Do not invert the image (default is to invert, detect dark objects)
```

In order to run the program call talk with the name of the image, the size range of the objects (in pixel) and the 
maximum allowed aspect ratio of the objects.

```
./talk --verbose segment1 tqsrrrqt.jpg 10 28 3
>  verbose on
>  run segment1
>  image size is 1920 by 1440 pixel
>  min/max: 4280.566300 61311.770400
>  mean size used: 16
>  write out mean removed image HiResHE_001_meanoff.png
>  invert the image before segmentation
>  focus size used: 2
>  write out the focused and inverted image HiResHE_002_focus.png
>  segment1 with size thresholds 25 .. 200, and maximum aspect ratio of 4
i: 0, x: 817, y: 1346, s: 70, a: 1.1661, c: 0.768491
i: 1, x: 667, y: 1288, s: 95, a: 1.3367, c: 0.804989
i: 2, x: 874, y: 1254, s: 81, a: 1.1454, c: 0.826230
...
i: 6592, x: 1613, y: 338, s: 30, a: 1.2899, c: 0.679061
i: 6593, x: 46, y: 169, s: 29, a: 3.6812, c: 0.889072
i: 6594, x: 1799, y: 808, s: 30, a: 3.4860, c: 0.859437
>  write out the found segmentation HiResHE_seg.png
```

The returned lines have the following structure for each detected region of interest:
  * i: index of region
  * x: x-coordinate of center of mass
  * y: y-coordinate of center of mass
  * s: size in pixel for region
  * a: aspect ratio for region
  * c: compactness value for region of interest


Here an example on how to perform a segmentation using more of the available options:

```
./talk --verbose segment1 HiResHE.jpg 25 200 4 --meansize 16 --focussize 2 --compactness 1.1
```

In this case the 'verbose' option will print out information about detected and estimated parameters
of the algorithm. We look for objects that have an area of between 30 and 150 pixel. To remove uneven illumination the image intensities will be corrected using a gliding mean subtraction filter of 16x16 pixel. A mexican hat sharpening filter focused on a pixel size of 2 is used afterwards. The resulting image is used for detection
of contiguous regions with a maximum aspect ratio of 1:4 and a maximum compactness value of 1.1 which filters out objects that are not compact enough.
