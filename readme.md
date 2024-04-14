# NanoRay - A Network Distributed Path Tracer

Words here

- Written in Go
- Uses Protobuf and gRPC
- Distributed architecture: controller + worker(s)
- Designed to parallel process, over CPU cores on each worker machine and across multiple workers
- Frontend uses HTMX
- YAML based description language
- Path tracing code based heavily on https://raytracing.github.io/books/RayTracingInOneWeekend.html

Image output at 2560x1080 with 300 samples per pixel, rendered in ~3mins on two laptop grade machines.

![screen](./examples/2024-04-14_17_52_31.png)
