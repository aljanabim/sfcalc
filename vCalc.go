// vcalc by Mustafa Al-Janabi, v0.5.4

package vclac

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type scalarField struct {
	expression string
	point      []float64
	coordsys   string
	precision  float64
}

type vectorField struct {
	expressionCoord1 string
	expressionCoord2 string
	expressionCoord3 string
	point            []float64
	coordsys         string
	precision        float64
}

// Returns a new scalar field
func NewScalarField(expression string, coordsys string) scalarField {
	s := scalarField{}
	s.expression = expression
	s.coordsys = coordsys
	s.precision = 0.0001
	checkCoords(expression, coordsys)
	return s
}

// Returns a new scalar field
func NewVectorField(e1, e2, e3, coordsys string) vectorField {
	v := vectorField{}
	v.expressionCoord1 = e1
	v.expressionCoord2 = e2
	v.expressionCoord3 = e3
	v.coordsys = coordsys
	v.precision = 0.0001
	checkCoords(e1+"+"+e2+"+"+e3, coordsys)
	return v
}

// Checks if user has used right coordinate names, panics if not
func checkCoords(expression string, coordsys string) {
	switch coordsys {
	case "car":
		if regexp.MustCompile(`[^q]r`).MatchString(expression) ||
			regexp.MustCompile(`theta`).MatchString(expression) ||
			regexp.MustCompile(`phi`).MatchString(expression) {
			panic("Insufficient coordinate names given, the following coordinate names are allowed together: (x,y,z) for cartesian coordinates, (r,phi,z) for cylinder coordinates and (r,theta,phi) for spherical coordinates")
		}
	case "cyl":
		if regexp.MustCompile(`[^q]r`).MatchString(expression) ||
			regexp.MustCompile(`y`).MatchString(expression) ||
			regexp.MustCompile(`theta`).MatchString(expression) {
			panic("Insufficient coordinate names given, the following coordinate names are allowed together: (x,y,z) for cartesian coordinates, (r,phi,z) for cylinder coordinates and (r,theta,phi) for spherical coordinates")
		}
	case "sph":
		if regexp.MustCompile(`x[^p]`).MatchString(expression) ||
			regexp.MustCompile(`y`).MatchString(expression) ||
			regexp.MustCompile(`z`).MatchString(expression) {
			panic("Insufficient coordinate names given, the following coordinate names are allowed together: (x,y,z) for cartesian coordinates, (r,phi,z) for cylinder coordinates and (r,theta,phi) for spherical coordinates")
		}
	}
}

// Returns the calculation of the expression given the points _1, _2, _3 in coordinate system
// Effectivley this function sorts out the addition and subtraction of the terms in the expression
func fn(_1, _2, _3 float64, expression string, coordsys string) float64 {
	var res float64

	// Split "+" and "-" and calculate sub parts containg "*" and "/"
	partsToCombine := strings.FieldsFunc(expression, fnSplitHelper)
	OPRList := listOPR(expression, len(partsToCombine)+1)
	// firstPart := partsToCombine[0]

	for i, v := range partsToCombine {
		// Get the sign of the other terms
		switch OPRList[i] {
		case "+":
			res += calculateTerm(_1, _2, _3, v, coordsys)
		case "-":
			res -= calculateTerm(_1, _2, _3, v, coordsys)
		default:
			res += calculateTerm(_1, _2, _3, v, coordsys)
		}

	}
	return res
}

// Takes type rune and returns true if rune has - or + operator
// Taken from https://stackoverflow.com/questions/39862613/how-to-split-multiple-delimiter-in-golang?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
func fnSplitHelper(r rune) bool {
	return r == '-' || r == '+'
}

// Returns a calucaltion of the given term in expression at the points _1, _2, _3 given the coordinate system coordsys
// This function effectively manages the calculation of multiplication and division in each term in expression
// It's the used in fn which preforms the addition and subtraction
func calculateTerm(_1, _2, _3 float64, expression string, coordsys string) float64 {
	// var res float64
	var arg float64 = 1
	var term float64 = 1
	// Regular expressions
	var OPR string
	var COEFFUNC string
	var FUNC string
	var EXPFUNC string
	var COEFCOORD string
	var COORD string
	var EXPCOORD string

	submatches := mathParser(expression)

	for i, match := range submatches {
		OPR = match[1]
		COEFFUNC = match[3]
		FUNC = match[4]
		COEFCOORD = match[5]
		COORD = match[6]
		EXPCOORD = match[7]
		EXPFUNC = match[8]

		arg = getCOEF(COEFCOORD) * math.Pow(getCOORD(_1, _2, _3, COORD, coordsys), getCOEF(EXPCOORD))

		if i == 0 && (OPR == "*" || OPR == "/") {
			term = 1
		}

		// Calculate the term given the operator infront of it
		switch OPR {
		case "":
			term = getCOEF(COEFFUNC) * math.Pow(getFUNC(FUNC, arg), getCOEF(EXPFUNC))
		case "*":
			term *= getCOEF(COEFFUNC) * math.Pow(getFUNC(FUNC, arg), getCOEF(EXPFUNC))
		case "/":
			term /= getCOEF(COEFFUNC) * math.Pow(getFUNC(FUNC, arg), getCOEF(EXPFUNC))
		}
	}
	// fmt.Println(term)
	return term
}

// Returns list of the addition and subtraction operation in the order they appear
func listOPR(expression string, length int) []string {
	var OPRList []string

	submatches := mathParser(expression)

	OPRList = append(OPRList, submatches[0][1])
	for i, match := range submatches {
		if i > 0 && (match[1] == "+" || match[1] == "-") {
			OPRList = append(OPRList, match[1])
		}
	}
	return OPRList
}

// Returns a list of strings with submatches of the parsed mathematical expression
func mathParser(expression string) [][]string {
	OPRERATIONS := `\+|\-|\*|\/`
	FUNCTIONS := `sin|cos|exp|sqrt|tan`
	COORDINATES := `x|y|z|r|phi|theta`
	REGEXP := `\s?(?P<OPR>` + OPRERATIONS + `)?\s?((?P<COEFFUNC>\d+)?(?P<FUNC>` + FUNCTIONS + `))?\(?(?P<COEFCOORD>\d+)?(?P<COORD>` + COORDINATES + `)?\^?(?P<EXPCOORD>\d+)?\)?\^?(?P<EXPFUNC>\d+)?`
	re := regexp.MustCompile(REGEXP)
	submatches := re.FindAllStringSubmatch(expression, -1)
	return submatches
}

// Takes a string containing a coefficient
// If string is empty it returns 1, thus the coefficient has no effect in expression
// else convert to float64 and return it
func getCOEF(COEF string) float64 {
	if COEF == "" {
		return 1
	} else {
		COEFF, _ := strconv.ParseFloat(COEF, 64)
		return COEFF
	}
}

// Takes a the coordinates seperatly as float64, the cordinate and coordinatesystem as strings
// If COORD is empty it returns 1, thus the cordinate has no effect in expression
// else check which coordinate it is given coordsys and return the given value of the coordinate
func getCOORD(_1, _2, _3 float64, COORD string, coordsys string) float64 {
	var COORD1 string
	var COORD2 string
	var COORD3 string

	switch coordsys { // Get the right name of coordinates given coordinate system
	case "car":
		COORD1 = "x"
		COORD2 = "y"
		COORD3 = "z"

	case "cyl":
		COORD1 = "r"
		COORD2 = "phi"
		COORD3 = "z"

	case "sph":
		COORD1 = "r"
		COORD2 = "theta"
		COORD3 = "phi"
	}

	switch COORD {
	case COORD1: // Use _1 as coordinate value
		return _1
	case COORD2: // Use _2 as coordinate value
		return _2
	case COORD3: // Use _3 as coordinate value
		return _3
	default:
		return 1
	}
}

// Takes a mathematical function as string and its arguments as float64
// Returns calculated value of the actual function given the arguments
// The functions are defined in expressionParser as FUNCTIONS
func getFUNC(FUNC string, arg float64) float64 {
	switch FUNC {
	case "sin":
		return math.Sin(arg)
	case "cos":
		return math.Cos(arg)
	case "exp":
		return math.Exp(arg)
	case "sqrt":
		return math.Sqrt(arg)
	case "tan":
		return math.Tan(arg)
	default:
		return arg
	}
}

// Calculates the gradient of scalarField
// Returns a slice of float64 containg the calculated gradient
func (s scalarField) Grad(c []float64) []float64 {
	if len(c) < 3 || len(c) > 3 {
		panic("Too many or too few points coordinates given")
	}
	h := s.precision
	switch s.coordsys {
	case "car":
		x := c[0]
		y := c[1]
		z := c[2]

		return []float64{
			(fn(x+h, y, z, s.expression, s.coordsys) - fn(x-h, y, z, s.expression, s.coordsys)) / (2 * h),
			(fn(x, y+h, z, s.expression, s.coordsys) - fn(x, y-h, z, s.expression, s.coordsys)) / (2 * h),
			(fn(x, y, z+h, s.expression, s.coordsys) - fn(x, y, z-h, s.expression, s.coordsys)) / (2 * h)}

	case "cyl":
		r := c[0]
		phi := c[1]
		z := c[2]

		return []float64{
			(fn(r+h, phi, z, s.expression, s.coordsys) - fn(r-h, phi, z, s.expression, s.coordsys)) / (2 * h),
			(fn(r, phi+h, z, s.expression, s.coordsys) - fn(r, phi-h, z, s.expression, s.coordsys)) / (2 * h * r),
			(fn(r, phi, z+h, s.expression, s.coordsys) - fn(r, phi, z-h, s.expression, s.coordsys)) / (2 * h)}

	case "sph":
		r := c[0]
		theta := c[1]
		phi := c[2]

		return []float64{
			(fn(r+h, theta, phi, s.expression, s.coordsys) - fn(r-h, theta, phi, s.expression, s.coordsys)) / (2 * h),
			(fn(r, theta+h, phi, s.expression, s.coordsys) - fn(r, theta-h, phi, s.expression, s.coordsys)) / (2 * h * r),
			(fn(r, theta, phi+h, s.expression, s.coordsys) - fn(r, theta, phi-h, s.expression, s.coordsys)) / (2 * h * r * math.Sin(theta))}
	}
	return []float64{}
}

// Calculates the divergence of scalarField
// Returns float64 containg the calculated divergence

// Calculates the rotation/curl of scalarField
// Returns float64 containg the calculated rotation

// func main() {
// 	s := NewScalarField("-3sin(2r^3)^5+phi*theta^2", "sph")
// 	fmt.Println(s.grad([]float64{1, 1, 1}))
// 	fmt.Println(s.coordsys)

// 	v := NewVectorField("3x^2", "5cos(y^3*z)", "sqrt(1-y^2)-5z+3", "car")
// 	fmt.Println(v.coordsys)

// }
