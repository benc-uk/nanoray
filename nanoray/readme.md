# NanoRay CLI

This is a standalone command line version of the renderer, it essentially acts as a single local worker.
It will slice the image into a number of jobs equal to the number of CPU cores, and render each in parallel

```
NanoRay - A path based ray tracer
  -aspect float
        Aspect ratio of the output image (default 1.7)
  -depth int
        Maximum ray recursion depth (default 5)
  -file string
        Scene file to render, in YAML format
  -output string
        Rendered output PNG file name (default "render.png")
  -samples int
        Samples per pixel, higher values give better quality but slower rendering (default 20)
  -width int
        Width of the output image (default 800)
```
