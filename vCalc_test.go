package vcalc

import (
	"math"
	"reflect"
	"testing"
)

func TestNewScalarField(t *testing.T) {
	var tests = []struct {
		expression string
		coordsys   string
	}{
		{"-3sin(2r^3)^5+phi*theta^2", "sph"},
		{"72+3x^2+5cos(y^2)-3z", "car"},
	}
	for _, v := range tests {
		if exp := NewScalarField(v.expression, v.coordsys); exp.coordsys != v.coordsys || exp.expression != v.expression {
			t.Error("Test failed: {", v.expression, v.coordsys, " } inputted, expected {", v, "} and got {", exp, "}")
		}
	}
}

func TestNewVectorField(t *testing.T) {
	var tests = []struct {
		expression1 string
		expression2 string
		expression3 string
		coordsys    string
	}{
		{"3r^2", "sqrt(1-theta^2)-5phi+3", "theta^3", "sph"},
		{"72+3x^2", "cos(y)", "z", "car"},
	}
	for _, v := range tests {
		if exp := NewVectorField(v.expression1, v.expression2, v.expression3, v.coordsys); exp.coordsys != v.coordsys || exp.expressionCoord1 != v.expression1 || exp.expressionCoord2 != v.expression2 || exp.expressionCoord3 != v.expression3 {
			t.Error("Test failed: {", v.expression1, v.expression2, v.expression3, v.coordsys, " } inputted, expected {", v, "} and got {", exp, "}")
		}
	}
}

func TestFn(t *testing.T) {
	var tests = []struct {
		point      []float64
		expression string
		coordsys   string
		exp        float64
	}{
		{[]float64{3, 3, 1}, "4x*y-3+z", "car", 34},
		{[]float64{2, 4, 0}, "r^3-cos(phi)/theta", "sph", 7.75},
		{[]float64{2, 4, 0}, "", "car", 0},
		{[]float64{2, 4, 0}, "", "cyl", 0},
		{[]float64{2, 4, 0}, "", "sph", 0},
		{[]float64{0, 0, 0}, "r+phi+z", "cyl", 0},
	}
	for _, v := range tests {
		if exp := fn(v.point[0], v.point[1], v.point[2], v.expression, v.coordsys); exp != v.exp {
			t.Error("Test failed: {", v.point, v.expression, v.coordsys, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestCalculateTerm(t *testing.T) {
	var tests = []struct {
		point      []float64
		expression string
		coordsys   string
		exp        float64
	}{
		{[]float64{3, 3, 1}, "4x*y*z", "car", 36},
		{[]float64{2, 4, 0}, "r^3*cos(phi)/theta", "sph", 2},
		{[]float64{2, 4, 0}, "", "car", 0},
		{[]float64{2, 4, 0}, "", "cyl", 0},
		{[]float64{2, 4, 0}, "", "sph", 0},
		{[]float64{0, 1, 0}, "r/phi*z", "cyl", 0},
	}
	for _, v := range tests {
		if exp := calculateTerm(v.point[0], v.point[1], v.point[2], v.expression, v.coordsys); exp != v.exp {
			t.Error("Test failed: {", v.point, v.expression, v.coordsys, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestListOPR(t *testing.T) {
	var tests = []struct {
		expression string
		exp        []string
	}{
		{"", []string{""}},
		{"*+/-", []string{"", "+", "-"}},
		{"+", []string{"+"}},
		{"3+35-43/54+cos(3x^3)^3", []string{"", "+", "-", "+"}},
		{"+++--**//+", []string{"+", "+", "+", "-", "-", "+"}},
		{"-3cos(x)-65y+theta/phi", []string{"-", "-", "+"}},
	}
	for _, v := range tests {
		if exp := listOPR(v.expression); reflect.DeepEqual(exp, v.exp) == false {
			t.Error("Test failed: {", v.expression, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestMathParser(t *testing.T) {
	var tests = []struct {
		expression string
		exp        [][]string
	}{
		{"", [][]string{[]string{"", "", "", "", "", "", "", ""}}},
		{"-3sin(2r^3)^5+phi*theta^2",
			[][]string{
				[]string{"-3sin(2r^3)^5", "-", "3", "sin", "2", "r", "3", "5"},
				[]string{"+phi", "+", "", "", "", "phi", "", ""},
				[]string{"*theta^2", "*", "", "", "", "theta", "2", ""},
			},
		},
		{"72+3x^2+5cos(y^2)-3z",
			[][]string{
				[]string{"72", "", "", "", "72", "", "", ""},
				[]string{"+3x^2", "+", "", "", "3", "x", "2", ""},
				[]string{"+5cos(y^2)", "+", "5", "cos", "", "y", "2", ""},
				[]string{"-3z", "-", "", "", "3", "z", "", ""},
			},
		},
	}

	for _, v := range tests {
		if exp := mathParser(v.expression); reflect.DeepEqual(exp, v.exp) == false {
			t.Error("Test failed: {", v.expression, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestGetCOEF(t *testing.T) {
	var tests = []struct {
		in  string
		exp float64
	}{
		{"", 1},
		{"1", 1},
		{"-4", -4},
		{"3.43", 3.43},
		{"0", 0},
	}
	for _, v := range tests {
		if exp := getCOEF(v.in); exp != v.exp {
			t.Error("Test failed: {", v.in, "} inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestGetCoord(t *testing.T) {
	var tests = []struct {
		_1, _2, _3 float64
		COORD      string
		coordsys   string
		exp        float64
	}{
		{3.72, -4, 0, "x", "car", 3.72},
		{3.72, -4, 0, "y", "car", -4},
		{3.72, -4, 0, "z", "car", 0},
		{342, -423423, 0.4452, "r", "cyl", 342},
		{342, -423423, 0.4452, "phi", "cyl", -423423},
		{342, -423423, 0.4452, "z", "cyl", 0.4452},
		{3, -0.4, 74, "r", "sph", 3},
		{3, -0.4, 74, "theta", "sph", -0.4},
		{3, -0.4, 74, "phi", "sph", 74},
	}
	for _, v := range tests {
		if exp := getCOORD(v._1, v._2, v._3, v.COORD, v.coordsys); exp != v.exp {
			t.Error("Test failed: {", v._1, v._2, v._3, v.COORD+" "+v.coordsys+" } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestGetFUNC(t *testing.T) {
	var tests = []struct {
		FUNC string
		arg  float64
		exp  float64
	}{
		{"sin", 0, 0},
		{"sin", math.Pi / 2, 1},
		{"cos", 0, 1},
		{"exp", 0, 1},
		{"tan", math.Pi / 4, 1},
		{"", -3, -3},
		{"", 4.53, 4.53},
	}
	for _, v := range tests {
		if exp := getFUNC(v.FUNC, v.arg); exp != v.exp {
			t.Error("Test failed: {", v.FUNC, v.arg, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestGrad(t *testing.T) {
	var tests = []struct {
		point []float64
		s     scalarField
		exp   []float64
	}{
		{[]float64{1, 3.14, math.Pi}, NewScalarField("-3sin(2r^3)^5+phi*theta^2", "sph"), []float64{25.604287369116463, 19.72920186458893, 6190.677138721322}},
		{[]float64{-1, -1, -1}, NewScalarField("-3sin(2r^3)^5+phi*z^2", "cyl"), []float64{25.604287369124233, -0.9999999999998899, 1.9999999999992246}},
		{[]float64{0, 0, 0}, NewScalarField("72+3x^2+5cos(y^2)-3z", "car"), []float64{0, 0, -2.999999999957481}},
	}
	for _, v := range tests {
		if exp := v.s.Grad(v.point); reflect.DeepEqual(exp, v.exp) == false {
			t.Error("Test failed: {", v.point, v.s, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestDiv(t *testing.T) {
	var tests = []struct {
		point []float64
		s     vectorField
		exp   float64
	}{
		{[]float64{1, 3.14, math.Pi}, NewVectorField("3r^2", "5cos(theta^3*phi)", "sqrt(1-theta^2)-5phi+3", "sph"), -11774.075610130705},
		{[]float64{-1, -1, -1}, NewVectorField("3r^2", "5cos(phi^3*z)", "sqrt(1-phi^2)-z+3", "cyl"), 2.622064867998322},
		{[]float64{0, 0, 0}, NewVectorField("x^2+cos(7y)", "y^2", "3z^2", "car"), 0},
	}
	for _, v := range tests {
		if exp := v.s.Div(v.point); exp != v.exp {
			t.Error("Test failed: {", v.point, v.s, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}

func TestRot(t *testing.T) {
	var tests = []struct {
		point []float64
		s     vectorField
		exp   []float64
	}{
		{[]float64{1, 3.14, math.Pi}, NewVectorField("3r^2", "5cos(theta^3*phi)", "sqrt(1-theta^2)-5phi+3", "sph"), []float64{-11464.210893801464, 3.484117402575223e+07, 14.097523877137295}},
		{[]float64{-1, -1, -1}, NewVectorField("3r^2", "5cos(phi^3*z)", "sqrt(1-phi^2)-z+3", "cyl"), []float64{-4.701511529343616, 0, 2.701511529340699}},
		{[]float64{0, 0, 0}, NewVectorField("x^2+cos(7y)", "y^2", "3z^2", "car"), []float64{0, 0, 0}},
	}
	for _, v := range tests {
		if exp := v.s.Rot(v.point); reflect.DeepEqual(exp, v.exp) == false {
			t.Error("Test failed: {", v.point, v.s, " } inputted, expected {", v.exp, "} and got {", exp, "}")
		}
	}
}
