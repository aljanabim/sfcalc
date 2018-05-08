# vCalc - Go vector Calculus Calculator
by Mustafa Al-Janabi

### Go package to calculate the following:
* Gradient and laplacian of scalar field
* Divergence and rotation of vector field<br>


### Installation
Once you have [installed Go](https://golang.org/doc/install), run this command
to install the `vcalc` package:

    go get github.com/aljanabim/vcalc

### Usage
The package is used by defining a new scalar- or a vector field. Definig a field requires a mathematical expression as well as a coordinate system. The gradient or laplacian are then calcualted for scalar fields by entering a point in the assigned coordinates system to the functions Grad and Laplacian respectivly. Using the methods Div and Rot with a point in the assigned coordinate system the divergence and rotation are calculated for a vector field.

#### How to define a scalar field


#### How to calculate gradient and laplacian

#### How to define a vector field

#### How to calculate divergence and rotation

### Documentation
#### type scalarField
	type scalarField {
    	// contains the expression, point, coordinate system and precision
	}

#### func NewScalarField
	func NewScalarField(e, c string) scalarField
New creates a new scalar field with given expression and coordinate system

#### func NewVectorField
	func NewVectorField(e1,e2,e3, c string) vectorField
New creates a new scalar field with given expression and coordinate system

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


Mustafa Al-Janabi