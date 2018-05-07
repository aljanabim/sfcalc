# vCalc - Go vector Calculus Calculator

### Go package to calculate the following:
* Gradient and laplacian of scalar field
* Divergence, rotation and vector laplacian of vector field<br>

A scalar field is a three-dimensional mathematical function describing some property of the space it occupies.
By entering a scalar field as a string, a point in the 3D-space and an approximation accuracy value, you can use the functions grad, div or rot to calculate the gradient, divergence or rotation respectivly. 

### Installation
Once you have [installed Go](https://golang.org/doc/install), run this command
to install the `sfcalc` package:

    go get github.com/aljanabim/sfcalc

### Usage


### Documentation
#### type scalarField
	type scalarField {
    	// contains the expression, point, coordinate system and precision
	}

#### func NewScalarField
	func NewScalarField(e string) scalarField
New creates a new scalar field with given expression

#### func NewVectorField
	func NewVectorField(e1,e2,e3 string) vectorField
New creates a new scalar field with given expression

#### func (scalarField) Grad
	func (s scalarField) Grad(c []float64) []float64
Grad calculates gradient of scalar field at given coordinates

#### func (scalarField) Laplacian
	func (s scalarField) Grad(c []float64) float64
Laplacian calculates laplacian of scalar field at given coordinates

#### type vectorField
	type vectorField {
		// contains the expression of each coordiante, point, coordinate system and precision
	}

#### func (vectorField) Div
	func (v vectorField) Div(c []float64) float64
Div calculates divergence of vector field at given coordinates

#### func (vectorField) Rot
	func (v vectorField) Rot(c []float64) []float64
Rot calculates rotation/curl of vector field at given coordinates

#### func (vectorField) VectorLaplacian
	func (v vectorField) Rot(c []float64) []float64
VectorLaplacian calculates the vector laplacian of vector field at given coordinates

### Exampel


Mustafa Al-Janabi