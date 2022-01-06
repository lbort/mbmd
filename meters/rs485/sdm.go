package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewSDMProducer)
}

const (
	METERTYPE_SDM = "SDM"
)

type SDMProducer struct {
	Opcodes
}

func NewSDMProducer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM630.
	 * See http://bg-etech.de/download/manual/SDM630Register.pdf
	 * This is to a large extent a superset of all SDM devices, however there are
	 * subtle differences (see 220, 230). Some opcodes might not work on some devices.
	 */
	ops := Opcodes{								// explanation from SDM630 datasheet
		VoltageL1:     0x0000, // 220, 230		// Phase 1 line to neutral volts
		VoltageL2:     0x0002,					// Phase 2 line to neutral volts
		VoltageL3:     0x0004,					// Phase 3 line to neutral volts
		CurrentL1:     0x0006, // 220, 230		// Phase 1 current
		CurrentL2:     0x0008,					// Phase 2 current
		CurrentL3:     0x000A,					// Phase 3 current
		PowerL1:       0x000C, //      230		// Phase 1 active power
		PowerL2:       0x000E,					// Phase 2 active power
		PowerL3:       0x0010,					// Phase 3 active power

		ApparentPowerL1: 0x0012, 				// apparent power L1
		ApparentPowerL2: 0x0014, 				// apparent power L2
		ApparentPowerL3: 0x0016, 				// apparent power L3
		ReactivePowerL1: 0x0018, 				// reactive power L1 (positive = capacitive?)
		ReactivePowerL2: 0x001A, 				// reactive power L2 (positive = capacitive?)
		ReactivePowerL3: 0x001C, 				// reactive power L3 (positive = capacitive?)

		//			   0x001E, 					// power factor L1 (already defined below!)
		//			   0x0020, 					// power factor L2 (already defined below!)
		//			   0x0022, 					// power factor L3 (already defined below!)
		// PhaseAngleL1:  0x0024, 					// phase angle L1
		// PhaseAngleL2:  0x0026, 					// phase angle L2
		// PhaseAngleL3:  0x0028, 					// phase angle L3
		Voltage:	   0x002A, 					// Average Line to neutral volts (todo: is this a good idea?)
		//			   0x002E, 					// Average Line current
		//			   0x0030, 					// Sum of Line currents

		Power:         0x0034,					// Total system (active) power
		ApparentPower: 0x0038,					// Total system volt-amps (apparent power)
		ReactivePower: 0x003C,					// Total system VAr (positive = capacitive?)

		//			   0x003E, 					// total power factor (already defined below!)
		PhaseAngle:	   0x0042, 					// total phase angle
		//			   0x0046, 					// freq. (already defined below!)
		//			   0x0048, 					// import kwh since last reset (already defined below!)
		//			   0x004A, 					// export kwh since last reset (already defined below!)
		//			   0x004C, 					// import kVARh since last reset
		//			   0x004E, 					// export kVARh since last reset
		//			   0x0050, 					// varh since last reset
		//			   0x0052, 					// aH since last reset

		ImportPower:   0x0054,					// total system power demand (is this averaged over set interval?)

		//			   0x0056, 					// maximum total power demand W
		//			   0x0064, 					// total power demand VA
		//			   0x0066, 					// maximum total power demand VA
		//			   0x0068, 					// neutral current demand
		//			   0x006A, 					// maximum neutral current demand
		//			   0x00C8, 					// Line 1 to line 2 volts
		//			   0x00CA, 					// Line 2 to line 3 volts
		//			   0x00CC, 					// Line 3 to line 1 volts
		//			   0x00CE, 					// Average line to line voltage

		Current:	   0x00E0, 					// Neutral current (bad idea! not really intuitive that "total" = neutral...)
		//			   0x00EA, 					// L1-N voltage THD (already defined below!)
		//			   0x00EC, 					// L2-N voltage THD (already defined below!)
		//			   0x00EE, 					// L3-N voltage THD (already defined below!)
		THDiL1:		   0x00F0, 					// L1 current THD
		THDiL2:		   0x00F2, 					// L2 current THD
		THDiL3:		   0x00F4, 					// L3 current THD
		//			   0x00F8, 					// average line to neutral thd voltage (already defined below!)
		THDi:		   0x00FA, 					// average line current thd
		//			   0x00FE, 					// negative total system power factor (same as 0x003E, but sign inverted?? but here it is degrees?)

		//			   0x0102, 					// Phase 1 current demand
		//			   0x0104, 					// Phase 2 current demand
		//			   0x0106, 					// Phase 3 current demand
		//			   0x0108, 					// Maximum phase 1 current demand
		//			   0x010A, 					// Maximum phase 2 current demand
		//			   0x010C, 					// Maximum phase 3 current demand

		//			   0x014E, 					// Line 1 to line 2 voltage THD
		//			   0x0150, 					// Line 2 to line 3 voltage THD
		//			   0x0152, 					// Line 3 to line 1 voltage THD
		//			   0x0154, 					// Average Line to line voltage THD
		//			   0x0156, 					// total kWh (already defined below!)
		ReactiveSum:   0x0158, 					// total kvarh
		
		ImportL1:      0x015A,					// L1 import kWh
		ImportL2:      0x015C,					// L2 import kWh
		ImportL3:      0x015E,					// L3 import kWh
		Import:        0x0048, // 220, 230 	// (out of order: taken from above!)
		ExportL1:      0x0160,					// L1 export kWh
		ExportL2:      0x0162,					// L2 export kWh
		ExportL3:      0x0164,					// L3 export kWh
		Export:        0x004a, // 220, 230 // (out of order: taken from above!)
		SumL1:         0x0166,					// L1 total kWh
		SumL2:         0x0168,					// L2 total kWh
		SumL3:         0x016a,					// L3 total kWh
		ReactiveImportL1: 0x016c, 				// L1 import kVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveImportL2: 0x016e, 				// L2 import KVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveImportL3: 0x0170, 				// L3 import KVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveExportL1: 0x0172, 				// L1 export KVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveExportL2: 0x0174, 				// L2 export KVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveExportL3: 0x0176, 				// L3 export KVARh (import and export seem to imply "capacitive" and "inductive". sign not clear yet)
		ReactiveSumL1: 0x0178, 					// L1 total KVARh
		ReactiveSumL2: 0x017A, 					// L2 total KVARh
		ReactiveSumL3: 0x017C, 					// L3 total KVARh

		Sum:           0x0156, // 220 // (out of order: taken from above!)


		// out of order! 
		CosphiL1:      0x001e, //      230 // (taken from above!)
		CosphiL2:      0x0020,			   // (taken from above!)
		CosphiL3:      0x0022,			   // (taken from above!)
		Cosphi:        0x003e,			   // (taken from above!)
		THDL1:         0x00ea, // voltage  // (taken from above!)
		THDL2:         0x00ec, // voltage  // (taken from above!)
		THDL3:         0x00ee, // voltage  // (taken from above!)
		THD:           0x00F8, // voltage  // (taken from above!)
		Frequency:     0x0046, //      230 // (taken from above!)
		//L1THDCurrent: 0x00F0, // current  (these ones do not compile!)
		//L2THDCurrent: 0x00F2, // current  (these ones do not compile!)
		//L3THDCurrent: 0x00F4, // current  (these ones do not compile!)
		//AvgTHDCurrent: 0x00Fa, // current (these ones do not compile!)
		//ApparentImportPower: 0x0064,      (these ones do not compile!)
	}
	return &SDMProducer{Opcodes: ops}
}

func (p *SDMProducer) Type() string {
	return METERTYPE_SDM
}

func (p *SDMProducer) Description() string {
	return "Eastron SDM630"
}

func (p *SDMProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDMProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDMProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
