// .
// └── organization/
//     └── boiler/
//         ├── input/
//         │   ├── temperature
//         │   └── pressure
//         └── output/
//             ├── temperature
//             └── pressure

package main

type Boiler struct {
	Input  Pipe
	Output Pipe
}

type Pipe struct {
	Temperature float32
	Pressure    float32
}

// Hotel=building/Boiler=boiler/Input=pipe/Temperature=float32
