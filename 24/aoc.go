package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type hailstone struct {
	//k          float64
	//n          float64
	//v          [3]int
	//pos        [3]int
	x, y, z    int
	vx, vy, vz int
}

//const (
//	x = iota
//	y
//	z
//)

func loadHailstones(lines []string) []hailstone {
	hailstones := make([]hailstone, 0)

	for _, l := range lines {
		// Found a much better alternative to loading the line
		var x, y, z, vx, vy, vz int
		fmt.Sscanf(l, "%d, %d, %d @ %d, %d, %d", &x, &y, &z, &vx, &vy, &vz)
		h := hailstone{x, y, z, vx, vy, vz}
		//h := hailstone{x, y, 0, vx, vy, 0}

		//h := hailstone{pos: [3]int{0, 0, 0}, v: [3]int{0, 0, 0}}
		//s := strings.Split(l, " @ ")
		//posSplit := strings.Split(s[0], ", ")
		//speedSplit := strings.Split(s[1], ", ")
		//
		//for i, vs := range posSplit {
		//	v, _ := strconv.Atoi(vs)
		//	h.pos[i] = v
		//}
		//for i, vs := range speedSplit {
		//	v, _ := strconv.Atoi(vs)
		//	h.v[i] = v
		//}
		//
		//// Determine value at x=0
		//vx := h.v[x]
		//vy := h.v[y]
		//if vx > 0 {
		//	vx *= -1
		//	vy *= -1
		//}
		//factor := float64(h.pos[x]) / float64(vx)
		//n := float64(h.pos[y]) + -1*factor*float64(vy)
		//// k = (y2-y1) / (x2-x1)
		//k := (float64(h.pos[y]) - n) / float64(h.pos[x]) // x1 is 0
		//h.k = k
		//h.n = n
		hailstones = append(hailstones, h)
	}

	return hailstones
}

// solution for intersection using vector form... as it makes sense to do
func intersect2d(h1, h2 hailstone) (float64, float64, int, int, bool) {
	// D = 1 / (v1_x * (-1) * v2_y - ((-1) * v2_x * v1_y))
	// [t1, t2] = D * [(p2_x - p1_x) * (-1) * v2_y + (p2_y - p1_y) * v2_x, (p2_x-p1_x)*(-1)*v1_y + (p2_y-p1_y)*v1_x]
	// t1 = D * ((p2_x - p1_x) * (-1) * v2_y + (p2_y - p1_y) * v2_x)
	// t2 = D * ((p2_x-p1_x)*(-1)*v1_y + (p2_y-p1_y)*v1_x)

	denominatorCheck := float64(h1.vx*-1*h2.vy - (-1 * h2.vx * h1.vy))
	if denominatorCheck == 0.0 {
		return 0, 0, 0, 0, false
	}

	denominator := 1 / denominatorCheck
	t1 := denominator * float64((h2.x-h1.x)*-1*h2.vy+(h2.y-h1.y)*h2.vx)
	t2 := denominator * float64((h2.x-h1.x)*-1*h1.vy+(h2.y-h1.y)*h1.vx)

	// TODO: Should this be checked still in here?
	if t1 < 0 || t2 < 0 {
		// Intersect in the past
		//fmt.Println("Intersect in the past", "t1", t1, "t2", t2)
		return 0, 0, 0, 0, false
	}
	// [x, y] = [v_x, v_y] * t + [pos_x, pos_y]
	x := float64(h1.vx)*t1 + float64(h1.x)
	y := float64(h1.vy)*t1 + float64(h1.y)

	return x, y, int(t1), int(t2), true
}

// Solution for part 1 using y = kx + n form
//func doHailstonesIntersect(h1, h2 hailstone, minValue, maxValue float64) bool {
//	if h1.k == h2.k {
//		return false
//	}
//
//	xCross := (h2.n - h1.n) / (h1.k - h2.k)
//	yCross := h1.k*xCross + h1.n
//
//	// Check that they do not cross in the past
//	if xCross > float64(h1.pos[x]) && h1.v[x] < 0 ||
//		xCross < float64(h1.pos[x]) && h1.v[x] > 0 ||
//		xCross > float64(h2.pos[x]) && h2.v[x] < 0 ||
//		xCross < float64(h2.pos[x]) && h2.v[x] > 0 {
//		return false
//	}
//
//	if xCross >= minValue && yCross >= minValue && xCross <= maxValue && yCross <= maxValue {
//		//fmt.Println("Intersection", h1, "and", h2, "at", "x", xCross, "y", yCross)
//		return true
//	}
//	return false
//}

func intersect3d(a, b hailstone) (float64, float64, float64, bool) {
	var xa, ya, za = float64(a.x), float64(a.y), float64(a.z)
	var xb, yb, _ = float64(b.x), float64(b.y), float64(b.z)

	var vxa, vya, vza = float64(a.vx), float64(a.vy), float64(a.vz)
	var vxb, vyb, _ = float64(b.vx), float64(b.vy), float64(b.vz)

	// xa + vxa * tx = xb + vxb * ty
	// ya + vya * tx = yb + vyb * ty
	// za + vza * tx = zb + vzb * ty

	// vxa * tx - vxb * ty = xb - xa
	// tx = (xb - xa + vxb * ty) / vxa
	// ya + vya * (xb - xa + vxb * ty) / vxa = yb + vyb * ty
	// ya + (vya/vxa)*(xb-xa) + vya*vxb*ty/vxa = yb + vyb * ty
	// ya + (vya/vxa)*(xb-xa) - yb =  ty * (vyb - vya*vxb/vxa)

	if xb-xa == 0 {
		fmt.Println("xb-xa == 0")
		return 0, 0, 0, false
	}
	if vxa == 0 {
		fmt.Println("vxa == 0")
		return 0, 0, 0, false
	}
	if vyb-vya*vxb/vxa == 0 {
		fmt.Println("vyb - vya*vxb/vxa == 0")
		return 0, 0, 0, false
	}

	ty := (ya - yb + (vya/vxa)*(xb-xa)) / (vyb - vya*vxb/vxa)
	tx := (xb - xa + vxb*ty) / vxa

	if tx < 0 || ty < 0 {
		return 0, 0, 0, false
	}
	//fmt.Printf("tx: %f, ty: %f, cross (at x=%f, y=%f)\n", tx, ty, xa+vxa*tx, ya+vya*tx)

	return xa + vxa*tx, ya + vya*tx, za + vza*tx, true
}

func PartOne(lines []string, minValue, maxValue float64) int {
	hailstones := loadHailstones(lines)
	//for _, h := range hailstones {
	//	fmt.Println(h)
	//}

	count := 0
	for i := range hailstones {
		h1 := hailstones[i]
		for j := i + 1; j < len(hailstones); j++ {
			x, y, t1, t2, ok := intersect2d(h1, hailstones[j])
			if ok && x >= minValue && y >= minValue && x <= maxValue && y <= maxValue && t1 >= 0 && t2 >= 0 {
				count++
			}
		}
	}
	return count
}

func findXY(data []hailstone, n int) (int, int, int, int) {
	for vx := -n; vx < n; vx++ {
		for vy := -n; vy < n; vy++ {
			d0 := data[0]
			d0.vx -= vx
			d0.vy -= vy
			d1 := data[1]
			d1.vx -= vx
			d1.vy -= vy
			x, y, _, _, ok := intersect2d(d0, d1)
			ix, iy := math.Round(x), math.Round(y)
			if ok {
				var found = true
				//for i := 2; i < len(data); i++ {
				for i := 2; i < 100; i++ {
					d := data[i]
					d.vx -= vx
					d.vy -= vy
					x0, y0, _, _, ok0 := intersect2d(d, d0)
					ix0, iy0 := math.Round(x0), math.Round(y0)
					//x1, y1, _, ok1 := intersect2d(d, d1)
					if !ok0 || ix0 != ix || iy0 != iy {
						//if !ok1 || x1 != x || y1 != y {
						found = false
						break
						//}
					}
				}
				if !found {
					continue
				}

				return int(ix), int(iy), vx, vy
			}
		}
	}
	return 0, 0, 0, 0
}

func findZ(data []hailstone, x, y, vx, vy int, n int) (int, int) {
	for vz := -n; vz < n; vz++ {
		d0 := data[0]
		d0.vx -= vx
		d0.vy -= vy
		d0.vz -= vz

		d1 := data[1]
		d1.vx -= vx
		d1.vy -= vy
		d1.vz -= vz

		x2, y2, z2, ok := intersect3d(d0, d1)
		ix2, iy2, iz2 := math.Round(x2), math.Round(y2), math.Round(z2)
		if ok {
			var found = true
			//for i := 2; i < len(data); i++ {
			for i := 2; i < 100; i++ {
				d := data[i]
				d.vx -= vx
				d.vy -= vy
				d.vz -= vz
				xx, yy, zz, ok := intersect3d(d, d0)
				ixx, iyy, izz := math.Round(xx), math.Round(yy), math.Round(zz)
				if !ok || ixx != ix2 || iyy != iy2 || izz != iz2 {
					found = false
					break
				}
			}
			if !found {
				continue
			}
			return int(iz2), vz
		}
	}
	return 0, 0
}

// Part two is copied from THE INTERNETS, but it doesn't even work on my input...
func PartTwo(lines []string) int {
	hailstones := loadHailstones(lines)

	n := 1000
	x, y, vx, vy := findXY(hailstones, n)
	fmt.Printf("x: %d, y: %d, vx: %d, vy: %d\n", x, y, vx, vy)

	z, vz := findZ(hailstones, x, y, vx, vy, n)
	fmt.Printf("z: %d, vz: %d\n", z, vz)

	return x + y + z
}

func LoadLines(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")
	return lines[:len(lines)-1], nil
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines, 200000000000000, 400000000000000))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
