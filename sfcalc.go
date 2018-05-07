package main

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

// Returns a new scalar field
func NewScalarField(expression string, point []float64, h float64) scalarField {
	v := scalarField{}
	v.expression = expression
	v.precision = h
	if len(point) < 3 || len(point) > 3 {
		panic("Too many or too few points coordinates given")
	}
	v.point = point
	v.setCoordinateSystem()
	return v
}

// Take a mathematical expression as string and return the coordinat system used
// given the variables used in the expression
// Returns string with name of the coordinate system used
func (v *scalarField) setCoordinateSystem() {
	if regexp.MustCompile(`x[^p]`).MatchString(v.expression) &&
		regexp.MustCompile(`y`).MatchString(v.expression) &&
		regexp.MustCompile(`z`).MatchString(v.expression) {
		v.coordsys = "cartesian"
	} else if regexp.MustCompile(`[^q]r`).MatchString(v.expression) &&
		regexp.MustCompile(`phi`).MatchString(v.expression) &&
		regexp.MustCompile(`z`).MatchString(v.expression) {
		v.coordsys = "cylinder"
	} else if regexp.MustCompile(`[^q]r`).MatchString(v.expression) &&
		regexp.MustCompile(`theta`).MatchString(v.expression) &&
		regexp.MustCompile(`phi`).MatchString(v.expression) {
		v.coordsys = "spherical"
	} else {
		panic("Insufficient coordinate name given, the following coordinate names are allowed together: (x,y,z) for cartesian coordinates, (r,phi,z) for cylinder coordinates and (r,theta,phi) for spherical coordinates")
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
	return res //expressionParser(_1, _2, _3, "-32sin(3x^2)^4", coordsys) - expressionParser(_1, _2, _3, "3y*z", coordsys) + expressionParser(_1, _2, _3, "3x", coordsys)
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
	case "cartesian":
		COORD1 = "x"
		COORD2 = "y"
		COORD3 = "z"

	case "cylinder":
		COORD1 = "r"
		COORD2 = "phi"
		COORD3 = "z"

	case "spherical":
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
func (v scalarField) grad() []float64 {
	h := v.precision
	switch v.coordsys {
	case "cartesian":
		x := v.point[0]
		y := v.point[1]
		z := v.point[2]

		return []float64{
			(fn(x+h, y, z, v.expression, v.coordsys) - fn(x-h, y, z, v.expression, v.coordsys)) / (2 * h),
			(fn(x, y+h, z, v.expression, v.coordsys) - fn(x, y-h, z, v.expression, v.coordsys)) / (2 * h),
			(fn(x, y, z+h, v.expression, v.coordsys) - fn(x, y, z-h, v.expression, v.coordsys)) / (2 * h)}

	case "cylinder":
		r := v.point[0]
		phi := v.point[1]
		z := v.point[2]

		return []float64{
			(fn(r+h, phi, z, v.expression, v.coordsys) - fn(r-h, phi, z, v.expression, v.coordsys)) / (2 * h),
			(fn(r, phi+h, z, v.expression, v.coordsys) - fn(r, phi-h, z, v.expression, v.coordsys)) / (2 * h * r),
			(fn(r, phi, z+h, v.expression, v.coordsys) - fn(r, phi, z-h, v.expression, v.coordsys)) / (2 * h)}

	case "spherical":
		r := v.point[0]
		theta := v.point[1]
		phi := v.point[2]

		return []float64{
			(fn(r+h, theta, phi, v.expression, v.coordsys) - fn(r-h, theta, phi, v.expression, v.coordsys)) / (2 * h),
			(fn(r, theta+h, phi, v.expression, v.coordsys) - fn(r, theta-h, phi, v.expression, v.coordsys)) / (2 * h * r),
			(fn(r, theta, phi+h, v.expression, v.coordsys) - fn(r, theta, phi-h, v.expression, v.coordsys)) / (2 * h * r * math.Sin(theta))}
	}
	return []float64{}
}

// Calculates the divergence of scalarField
// Returns float64 containg the calculated divergence

// Calculates the rotation/curl of scalarField
// Returns float64 containg the calculated rotation

func main() {
	// v := NewScalarField("5+7^3+3r+2phi-cos(z)", []float64{3, 4, 5})
	// v := NewScalarField("3y*z+3*x/5cos(y)+z", []float64{3.002, -4, 5}, 0.0001)
	//fmt.Println(v)          32sin(3x^2)^4

	v := NewScalarField("-3sin(2r^3)^5+phi*theta^2", []float64{1, 1, 1}, 0.0001)
	fmt.Println(v.grad())‘
	fmt.Println(v.coordsys)
	// h := 0.00001
	// F := scalarField{"3x+2y+cos(z)", []float64{2, 2, 2}, "cartesian"}
	// fmt.Println(F.grad(h))

	//parseExpression(v.point[0], v.point[1], v.point[2], v.expression, v.coordsys)
}
