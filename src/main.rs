use image::RgbImage;

#[derive(Copy, Clone, Debug)]
struct Point3D {
    x: f64,
    y: f64,
    z: f64,
}

const MAX_SPHERES: usize = 5000;
const SPHERE_RAD: f64 = 5.0;

fn main() {
    // image width and height
    let width = 1280;
    let height = 720;
    let width_f = width as f64;
    let height_f = height as f64;

    let mut img = RgbImage::new(width, height);

    let cam_pos = Point3D {
        x: 0.0,
        y: 0.0,
        z: -50.0,
    };

    let mut sphere_positions = vec![
        Point3D {
            x: 0.0,
            y: 0.0,
            z: 0.0
        };
        MAX_SPHERES
    ];

    // Random sphere positions
    for i in 0..MAX_SPHERES {
        sphere_positions[i] = Point3D {
            x: rand::random::<f64>() * 70.0 - 35.0,
            y: rand::random::<f64>() * 40.0 - 20.0,
            z: (rand::random::<f64>() * 40.0 - 20.0) + 40.0,
        };
    }

    let mut global_light_dir = Point3D {
        x: 500.0,
        y: -50.0,
        z: -100.0,
    };
    let global_light_len = (global_light_dir.x * global_light_dir.x
        + global_light_dir.y * global_light_dir.y
        + global_light_dir.z * global_light_dir.z)
        .sqrt();
    global_light_dir = Point3D {
        x: global_light_dir.x / global_light_len,
        y: global_light_dir.y / global_light_len,
        z: global_light_dir.z / global_light_len,
    };
    let sphere_colour = [1.0, 0.0, 0.0];

    let now = std::time::Instant::now();

    // Render
    for y in 0..height {
        for x in 0..width {
            let ix = x as f64 - width_f / 2.0;
            let iy = y as f64 - height_f / 2.0;

            let ray_dir = Point3D {
                x: ix / height_f,
                y: iy / height_f,
                z: 1.0,
            };
            // normalize ray_dir
            let ray_len =
                (ray_dir.x * ray_dir.x + ray_dir.y * ray_dir.y + ray_dir.z * ray_dir.z).sqrt();
            let ray_dir = Point3D {
                x: ray_dir.x / ray_len,
                y: ray_dir.y / ray_len,
                z: ray_dir.z / ray_len,
            };

            let mut min_t = f64::MAX;
            let mut hit_sphere: i32 = -1;

            for i in 0..MAX_SPHERES {
                let sphere = sphere_positions[i];
                let t = calc_t(sphere, cam_pos, ray_dir);
                if t >= 0.0 && t < min_t {
                    min_t = t;
                    hit_sphere = i as i32;
                }
            }

            if hit_sphere >= 0 {
                let sphere = sphere_positions[hit_sphere as usize];

                let hit_pos = Point3D {
                    x: cam_pos.x + ray_dir.x * min_t,
                    y: cam_pos.y + ray_dir.y * min_t,
                    z: cam_pos.z + ray_dir.z * min_t,
                };

                let normal = Point3D {
                    x: hit_pos.x - sphere.x,
                    y: hit_pos.y - sphere.y,
                    z: hit_pos.z - sphere.z,
                };

                let normal_len =
                    (normal.x * normal.x + normal.y * normal.y + normal.z * normal.z).sqrt();
                let normal = Point3D {
                    x: normal.x / normal_len,
                    y: normal.y / normal_len,
                    z: normal.z / normal_len,
                };

                let mut dot = normal.x * global_light_dir.x
                    + normal.y * global_light_dir.y
                    + normal.z * global_light_dir.z;
                if dot < 0.0 {
                    dot = 0.0;
                }

                let reflected_dir = Point3D {
                    x: normal.x * 2.0 * dot - global_light_dir.x,
                    y: normal.y * 2.0 * dot - global_light_dir.y,
                    z: normal.z * 2.0 * dot - global_light_dir.z,
                };

                let view_dir = Point3D {
                    x: -ray_dir.x,
                    y: -ray_dir.y,
                    z: -ray_dir.z,
                };

                // dot product of view_dir and reflected_dir
                let mut specular = reflected_dir.x * view_dir.x
                    + reflected_dir.y * view_dir.y
                    + reflected_dir.z * view_dir.z;

                // max of 0 and specular
                specular = if specular > 0.0 { specular } else { 0.0 };
                specular = specular.powf(12.0);

                let out_r = (sphere_colour[0] * dot) + specular;
                let out_g = (sphere_colour[1] * dot) + specular;
                let out_b = (sphere_colour[2] * dot) + specular;

                img.put_pixel(
                    x,
                    y,
                    image::Rgb([
                        (out_r * 255.0) as u8,
                        (out_g * 255.0) as u8,
                        (out_b * 255.0) as u8,
                    ]),
                );
                continue;
            }

            img.put_pixel(x, y, image::Rgb([0, 0, 0]));
        }
    }

    println!("Time taken: {}ms", now.elapsed().as_millis());

    img.save("output-rust.png").unwrap();
}

fn calc_t(sphere: Point3D, cam_pos: Point3D, ray_dir: Point3D) -> f64 {
    let oc = Point3D {
        x: cam_pos.x - sphere.x,
        y: cam_pos.y - sphere.y,
        z: cam_pos.z - sphere.z,
    };

    let mut t = -1.0;
    let a = ray_dir.x * ray_dir.x + ray_dir.y * ray_dir.y + ray_dir.z * ray_dir.z;
    let b = 2.0 * (oc.x * ray_dir.x + oc.y * ray_dir.y + oc.z * ray_dir.z);
    let c = oc.x * oc.x + oc.y * oc.y + oc.z * oc.z - SPHERE_RAD * SPHERE_RAD;

    let discriminant = b * b - 4.0 * a * c;

    if discriminant >= 0.0 {
        let t1 = (-b - discriminant.sqrt()) / (2.0 * a);
        let t2 = (-b + discriminant.sqrt()) / (2.0 * a);

        if t1 < 0.0 && t2 < 0.0 {
            t = -1.0;
        } else {
            t = if t1 < t2 { t1 } else { t2 };
        }
    }

    return t;
}
