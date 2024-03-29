PNG = require('pngjs').PNG

const width = 1280
const height = 720

const imageBuffer = new Uint8Array(width * height * 4)

const spheres = []
for (let i = 0; i < 5000; i++) {
  spheres.push({
    x: Math.random() * 70 - 35,
    y: Math.random() * 40 - 20,
    z: Math.random() * 40 - 20 + 40,
  })
}

const sphereRadius = 5
const sphereColour = { r: 1, g: 0, b: 0 }
const camPos = { x: 0, y: 0, z: -50 }

const lightPos = { x: 300, y: -40, z: -200 }

const now = Date.now()
for (let y = 0; y < height; y++) {
  for (let x = 0; x < width; x++) {
    const ix = x - width / 2
    const iy = y - height / 2

    let ray = { x: ix / height, y: iy / height, z: 1 }
    ray = normalize(ray)

    // cast a ray
    rayDir = { x: ix / height, y: iy / height, z: 1 }
    rayDir = normalize(rayDir)

    let hitSphere = -1
    let minT = 10000

    for (let i = 0; i < spheres.length; i++) {
      const sphere = spheres[i]

      const oc = { x: camPos.x - sphere.x, y: camPos.y - sphere.y, z: camPos.z - sphere.z }

      const a = rayDir.x * rayDir.x + rayDir.y * rayDir.y + rayDir.z * rayDir.z
      const b = 2 * (oc.x * rayDir.x + oc.y * rayDir.y + oc.z * rayDir.z)
      const c = oc.x * oc.x + oc.y * oc.y + oc.z * oc.z - sphereRadius * sphereRadius

      const disc = b * b - 4 * a * c

      if (disc > 0) {
        const t1 = (-b - Math.sqrt(disc)) / (2 * a)
        const t2 = (-b + Math.sqrt(disc)) / (2 * a)
        if (t1 < 0 && t2 < 0) {
          continue
        }

        const t = Math.min(t1, t2)

        if (t < minT) {
          minT = t
          hitSphere = i
        }
      }
    }

    if (hitSphere !== -1) {
      const sphere = spheres[hitSphere]

      const hitPoint = {
        x: camPos.x + rayDir.x * minT,
        y: camPos.y + rayDir.y * minT,
        z: camPos.z + rayDir.z * minT,
      }
      const normal = normalize({
        x: hitPoint.x - sphere.x,
        y: hitPoint.y - sphere.y,
        z: hitPoint.z - sphere.z,
      })

      const lightDir = normalize({
        x: lightPos.x - hitPoint.x,
        y: lightPos.y - hitPoint.y,
        z: lightPos.z - hitPoint.z,
      })

      const lambert = Math.max(0, normal.x * lightDir.x + normal.y * lightDir.y + normal.z * lightDir.z)

      // specular
      const reflected = {
        x: 2 * normal.x * lambert - lightDir.x,
        y: 2 * normal.y * lambert - lightDir.y,
        z: 2 * normal.z * lambert - lightDir.z,
      }
      const viewDir = normalize({
        x: camPos.x - hitPoint.x,
        y: camPos.y - hitPoint.y,
        z: camPos.z - hitPoint.z,
      })

      const specular = Math.pow(
        Math.max(0, reflected.x * viewDir.x + reflected.y * viewDir.y + reflected.z * viewDir.z),
        30
      )

      outColour = {
        r: (sphereColour.r * lambert + specular) * 255,
        g: (sphereColour.g * lambert + specular) * 255,
        b: (sphereColour.b * lambert + specular) * 255,
      }

      setPixel(x, y, Math.min(outColour.r, 255), Math.min(outColour.g, 255), Math.min(outColour.b, 255), 255)
    } else {
      setPixel(x, y, 0, 0, 0, 255)
    }
  }
}

console.log(`Time taken: ${Date.now() - now}ms`)

// Save as PNG
const fs = require('fs')
const png = new PNG({
  width: width,
  height: height,
})
png.data = imageBuffer
png.pack().pipe(fs.createWriteStream('output-node.png'))

//
// ============================================
//

function setPixel(x, y, r, g, b, a) {
  const index = (y * width + x) * 4
  imageBuffer[index] = r
  imageBuffer[index + 1] = g
  imageBuffer[index + 2] = b
  imageBuffer[index + 3] = a
}

function normalize(v) {
  const mag = Math.sqrt(v.x * v.x + v.y * v.y + v.z * v.z)
  return { x: v.x / mag, y: v.y / mag, z: v.z / mag }
}
