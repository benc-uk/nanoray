use image::RgbImage;

#[derive(Copy, Clone)]
struct Point3D {
    x: f64,
    y: f64,
    z: f64,
}

fn main() {
    // image width and height
    let width = 1280;
    let height = 1720;
    let width_f = width as f64;
    let height_f = height as f64;

    let now = std::time::Instant::now();

    // array of pixels each of which is an array of 3 integers
    //let mut image = vec![vec![vec![0; 3]; width]; height];
    let mut img = RgbImage::new(width, height);

    let cam_pos = Point3D {
        x: 0.0,
        y: 0.0,
        z: -20.0,
    };

    let mut sphere_positions = vec![
        Point3D {
            x: 0.0,
            y: 0.0,
            z: 0.0
        };
        1000
    ];

    // Random sphere positions
    for i in 0..1000 {
        sphere_positions[i] = Point3D {
            x: rand::random::<f64>() * 70.0 - 35.0,
            y: rand::random::<f64>() * 40.0 - 20.0,
            z: 40.0,
        };
    }

    let sphere_radius = 5.0;
    let light_dir = Point3D {
        x: 0.0,
        y: 1.0,
        z: 0.0,
    };

    // Render
    for y in 0..height {
        for x in 0..width {
            //println!("rendering pixel: {}, {}", x, y);

            let ix = x as f64 - width_f / 2.0;
            let iy = y as f64 - height_f / 2.0;

            let ray_dir = Point3D {
                x: ix / height_f,
                y: iy / height_f,
                z: 1.0,
            };

            let mut hit = false;
            for i in 0..1000 {
                let sphere = sphere_positions[i];

                let oc = Point3D {
                    x: cam_pos.x - sphere.x,
                    y: cam_pos.y - sphere.y,
                    z: cam_pos.z - sphere.z,
                };

                let a = ray_dir.x * ray_dir.x + ray_dir.y * ray_dir.y + ray_dir.z * ray_dir.z;
                let b = 2.0 * (oc.x * ray_dir.x + oc.y * ray_dir.y + oc.z * ray_dir.z);
                let c = oc.x * oc.x + oc.y * oc.y + oc.z * oc.z - sphere_radius * sphere_radius;

                let discriminant = b * b - 4.0 * a * c;

                if discriminant >= 0.0 {
                    let t1 = (-b - discriminant.sqrt()) / (2.0 * a);
                    let t2 = (-b + discriminant.sqrt()) / (2.0 * a);

                    if t1 < 0.0 && t2 < 0.0 {
                        continue;
                    }

                    let t = if t1 < t2 { t1 } else { t2 };

                    hit = true;
                    break;
                }
            }

            if hit {
                img.put_pixel(x, y, image::Rgb([255, 255, 255]));
            } else {
                img.put_pixel(x, y, image::Rgb([0, 0, 0]));
            }
        }
    }

    println!("Time: {}ms", now.elapsed().as_millis());

    img.save("output-rust.png").unwrap();
}
