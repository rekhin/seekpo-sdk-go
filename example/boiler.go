// .
// └── Organization/
//     └── HotWaterSupply/
//         ├── input/
//         │   ├── temperature
//         │   └── pressure
//         └── output/
//             ├── temperature
//             └── pressure

// # Seekpo=Organization
// ## Hotel=Building
// ### Boiler=HotWaterSupply
// #### Direct=Pipe
// ##### Inpup
// ###### Temperature
// ###### Pressure
// ###### VolumeFlowRate
// ###### MassFlowRate
// ##### Output
// ###### Enthalpy
// ###### Density
// ###### MassFlow
// #### Reverse=Pipe
// ##### Inpup
// ###### Temperature
// ###### Pressure
// ###### VolumeFlowRate
// ###### MassFlowRate
// ##### Output
// ###### Enthalpy
// ###### Density
// ###### MassFlow

// .
// └── Seekpo=Organization
//     └── Hotel=Building
//         └── Boiler=HotWaterSupply
//             ├── Direct=Pipe
//             │   ├── Temperature=floa64
//             │   ├── Pressure=floa64
//             │   ├── VolumeFlowRate=floa64
//             │   ├── MassFlowRate=floa64
//             │   ├── Enthalpy=floa64
//             │   ├── Density=floa64
//             │   └── MassFlow=floa64
//             └── Reverse=Pipe
//                 ├── Temperature=floa64
//                 ├── Pressure=floa64
//                 ├── VolumeFlowRate=floa64
//                 ├── MassFlowRate=floa64
//                 ├── Enthalpy=floa64
//                 ├── Density=floa64
//                 └── MassFlow=floa64

package main

type HotWaterSupply struct {
	Direct  Pipe
	Reverse Pipe
}

type Pipe struct {
	Inpup  Input
	Output Output
}

type Input struct {
	Temperature    float64
	Pressure       float64
	VolumeFlowRate float64
	MassFlowRate   float64
}

type Output struct {
	Enthalpy float64
	Density  float64
	MassFlow float64
}

// Hotel=building/Boiler=boiler/Input=pipe/Temperature=float32
