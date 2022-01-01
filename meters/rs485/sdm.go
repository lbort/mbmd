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
	ops := Opcodes{
		VoltageL1:     0x0000, // 220, 230
		VoltageL2:     0x0002,
		VoltageL3:     0x0004,
		CurrentL1:     0x0006, // 220, 230
		CurrentL2:     0x0008,
		CurrentL3:     0x000A,
		PowerL1:       0x000C, //      230
		PowerL2:       0x000E,
		PowerL3:       0x0010,

		//			   0x0012, // apparent power L1
		//			   0x0014, // apparent power L2
		//			   0x0016, // apparent power L3
		//             0x0018, // reactive power L1
		//			   0x001A, // reactive power L2
		//			   0x001C, // reactive power L3
		//			   0x001E, // power factor L1 (already below!)
		//			   0x0020, // power factor L2 (already below!)
		//			   0x0022, // power factor L3 (already below!)
		//			   0x0024, // phase angle L1
		//			   0x0026, // phase angle L2
		//			   0x0028, // phase angle L3
		//			   0x002A, // Average Line to neutral volts
		//			   0x002E, // Average Line current
		//			   0x0030, // Sum of Line currents

		Power:         0x0034,
		ApparentPower: 0x0038,
		ReactivePower: 0x003C,

		//			   0x003E, // total power factor (already below!)
		//			   0x0042, // total phase angle
		//			   0x0046, // freq. (already below!)
		//			   0x0048, // import kwh since last reset (already below!)
		//			   0x004A, // export kwh since last reset
		//			   0x004C, // import kVARh since last reset
		//			   0x004E, // export kVARh since last reset
		//			   0x0050, // varh since last reset
		//			   0x0052, // aH since last reset

		ImportPower:   0x0054,

		//			   0x0056, // maximum total power demand kW
		//			   0x0064, // total power demand kVA 
		//			   0x0066, // maximum total power demand kVA 
		//			   0x0068, // neutral current demand
		//			   0x006A, // maximum neutral current demand
		//			   0x00C8, // Line 1 to line 2 volts
		//			   0x00CA, // Line 2 to line 3 volts
		//			   0x00CC, // Line 3 to line 1 volts
		//			   0x00CE, // Average Line to line volt

		//			   0x00E0, // Neutral current
		//			   0x00EA, // L1-N volt. THD
		//			   0x00EC, // L2-N volt. THD
		//			   0x00EE, // L3-N volt. THD
		//			   0x00F0, // L1 current THD
		//			   0x00F2, // L2 current THD
		//			   0x00F4, // L3 current THD
		//			   0x00F8, // average line to neutral thd voltage
		//			   0x00FA, // average line current thd
		//			   0x00FE, // total system power factor

		//			   0x0102, //
		//			   0x0104, //
		//			   0x0106, //
		//			   0x0108, //
		//			   0x010A, //
		//			   0x010C, //

		//			   0x014E, //
		//			   0x0150, //
		//			   0x0152, //
		//			   0x0154, //
		//			   0x0156, //
		//			   0x0158, //
		
		ImportL1:      0x015a,
		ImportL2:      0x015c,
		ImportL3:      0x015e,
		Import:        0x0048, // 220, 230 // (taken from above!)
		ExportL1:      0x0160,
		ExportL2:      0x0162,
		ExportL3:      0x0164,
		Export:        0x004a, // 220, 230 // (taken from above!)
		SumL1:         0x0166,
		SumL2:         0x0168,
		SumL3:         0x016a,
		//			   0x016c, //
		//			   0x016e, //
		//			   0x0170, //
		//			   0x0172, //
		//			   0x0174, //
		//			   0x0176, //
		//			   0x0178, //
		//			   0x017A, //
		//			   0x017C, //

		Sum:           0x0156, // 220 // (taken from above!)


		// out of order! 
		CosphiL1:      0x001e, //      230
		CosphiL2:      0x0020,
		CosphiL3:      0x0022,
		Cosphi:        0x003e,
		THDL1:         0x00ea, // voltage
		THDL2:         0x00ec, // voltage
		THDL3:         0x00ee, // voltage
		THD:           0x00F8, // voltage
		Frequency:     0x0046, //      230
		//L1THDCurrent: 0x00F0, // current
		//L2THDCurrent: 0x00F2, // current
		//L3THDCurrent: 0x00F4, // current
		//AvgTHDCurrent: 0x00Fa, // current
		//ApparentImportPower: 0x0064,
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
