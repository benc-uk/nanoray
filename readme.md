# NanoRay - A Network Distributed Path Tracer

Every few years I have to write a new ray tracer, this is my latest attempt in 2024

- Written in Go
- Uses [Protobuf](https://protobuf.dev/) and gRPC
- Distributed architecture: controller + worker(s)
- Designed to parallel process, over CPU cores on each worker machine and across multiple workers
- Frontend uses [HTMX](https://htmx.org/)
- YAML based [scene description language](./schemas/scene.json)
- Path tracing code based heavily on https://raytracing.github.io/books/RayTracingInOneWeekend.html

Caustics and lights

![screen](./examples/renders/2024-04-20_16_14_53.jpg)

Christmas!

![screen](./examples/renders/2024-04-14_17_52_31.png)

Example with depth of field

![screen](./examples/renders/focal-01.png)
