# vCalc - Vector calculus calculator
by Mustafa Al-Janabi

## Go package to calculate the following:
* Gradient and laplacian of scalar field
* Divergence and rotation of vector field


## Installation
Once you have [installed Go](https://golang.org/doc/install), run this command
to install the `vcalc` package:

    go get github.com/aljanabim/vcalc

## Usage
The package is used by defining a new scalar- or a vector field. Definig a field requires a mathematical expression as well as a coordinate system. The gradient or laplacian are then calcualted for scalar fields by entering a point in the assigned coordinates system to the functions Grad and Laplacian respectivly. Using the methods Div and Rot with a point in the assigned coordinate system the divergence and rotation are calculated for a vector field.

#### How to define a scalar field?
To define a new scalar field you use the function NewScalarField. You enter the expression of your scalar field and the coordinate system as arguments.
<pre><code>NewScalarField{<b>EXPRESSION</b>, <b>"COORDINATE SYSTEM"</b>}</code></pre>
#### How to define a vector field?
To define a new vector field you use the function NewVectorField. You enter the expression of each axis of the vector field as well as the coordinate system as inputs. 
<pre><code>NewVectorField{<b>EXPRESSION</b>, <b>EXPRESSION</b>, <b>EXPRESSION</b>, <b>"COORDINATE SYSTEM"</b>}</code></pre>
#### How to write an __EXPRESSION__?
You enter an expression as a string. The expression string should follow the following pattern
<pre><code><b>"[OPR]a[FUNC](b[COORD]^c)^d"</b></code></pre>
Where

| Expression parts | Possible string values |
| :-------------: | :------ |
| [OPR]     | "", "+", "-", "*" or "/"  |
| [FUNC]     | "", "sin", "cos", "tan", "exp" or "sqrt" |
| [COORD] | "", "x", "y" or "z" for cartesian coordinates<br> "", "r", "phi" or "z" for cylinder coordinates<br> "", "r", "theta" or "phi" for spherical coordinates    |
| a,b,c and d| arbitrary positive integers as strings |

_Note that every part is optional and can be left empty to cover diverce mathematical expressions, however_ __every__ _coordinate has to be present in the expression._

For example 
```
"72+3x^2+5cos(y^2*z^2)-3z"
```
is valid. While 
```
"72+3x^2+5cos(0)-3z"
```
is invalid, because the y-coordinate is missing.

#### How to write a __COORDINATE SYSTEM__?
You enter a coordinate system as a string. The string should be
* "car" for cartesian coordinates
* "cyl" for cylinder coordinates
* "sph" for spherical coordinates

#### How to Calculate gradient, laplacian, divergence and rotation
To calculate gradient and laplacian you use the methods Grad and Laplacian on a scalar field at a specific point in the 3-dimensional space.

To calculate divergence and rotation you use the methods Div and Rot on a vector field.

#### Examples
```go
	s := NewScalarField("3^5-7x^2-y+3cos(z^2)^2", "car")
	fmt.Println(s.Grad([]float64{4, 2, 7.2}))
	// Prints 
	// [3,3,3]
```

```go
	s := NewScalarField("-3sin(2r^3)^5+phi*z^2", "cyl")
	fmt.Println(s.Laplacian([]float64{-1, 2, 0}))
	// Prints 
	// [3,3,3]
```

```go
	v := NewVectorField("x^2+cos(7y)", "y^2", "3z^2", "car")
	fmt.Println(v.Div([]float64{-1, 2.76, 0}))
	// Prints 
	// 22
```

```go
	v := NewVectorField("3theta*5^2-r^2",
			    "sqrt(theta*phi)",
			    "3theta*7r+2exp(3phi^2)",
			    "sph")
	fmt.Println(v.Rot([]float64{-11, 3.14, 2}))
	// Prints 
	// [3,3,3]
```
	


## Documentation

type scalarField
* [func NewScalarField(e, c string) scalarField](#type
* scalarField)
* func (s scalarField) Grad(c []float64) []float64
* func (s scalarField) Grad(c []float64) float64

type vectorField
* func NewVectorField(e1,e2,e3, c string) vectorField
* func (v vectorField) Div(c []float64) float64
* func (v vectorField) Rot(c []float64) []float64

#### type scalarField
	type scalarField {
    	// contains the expression, point, coordinate system and precision
	}

#### func NewScalarField
	func NewScalarField(e, c string) scalarField
New creates a new scalar field with given expression and coordinate system

#### func (scalarField) Grad
	func (s scalarField) Grad(c []float64) []float64
Grad calculates gradient of scalar field at given coordinates

#### func (scalarField) Laplacian
	func (s scalarField) Laplacian(c []float64) float64
Laplacian calculates laplacian of scalar field at given coordinates

#### type vectorField
	type vectorField {
		// contains the expression of each coordiante, point, coordinate system and precision
	}

#### func NewVectorField
	func NewVectorField(e1,e2,e3, c string) vectorField
New creates a new scalar field with given expression and coordinate system

#### func (vectorField) Div
	func (v vectorField) Div(c []float64) float64
Div calculates divergence of vector field at given coordinates

#### func (vectorField) Rot
	func (v vectorField) Rot(c []float64) []float64
Rot calculates rotation/curl of vector field at given coordinates

## Roadmap
* The package has yet to support "pi" and floats in the expression.

Mustafa Al-Janabi