package sdm630

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"log"
	"math"
	"os"
	"time"
)

// See http://bg-etech.de/download/manual/SDM630Register.pdf
const (
	OpCodeL1Voltage = 0x0000
	OpCodeL2Voltage = 0x0002
	OpCodeL3Voltage = 0x0004
	OpCodeL1Current = 0x0006
	OpCodeL2Current = 0x0008
	OpCodeL3Current = 0x000A
	//TODO: This is broken.
	OpCodeL1Consumption = 0x000C
	OpCodeL2Consumption = 0x000E
	OpCodeL3Consumption = 0x0010
	OpCodeL1Production  = 0x000C
	OpCodeL2Production  = 0x000E
	OpCodeL3Production  = 0x0010
	OpCodeL1PowerFactor = 0x001e
	OpCodeL2PowerFactor = 0x0020
	OpCodeL3PowerFactor = 0x0022

	MaxRetryCount = 3
)

type QueryEngine struct {
	client     modbus.Client
	interval   int
	handler    modbus.RTUClientHandler
	datastream ReadingChannel
}

func NewQueryEngine(
	rtuDevice string,
	interval int,
	verbose bool,
	channel ReadingChannel,
) *QueryEngine {
	// Modbus RTU/ASCII
	mbhandler := modbus.NewRTUClientHandler(rtuDevice)
	mbhandler.BaudRate = 9600
	mbhandler.DataBits = 8
	mbhandler.Parity = "N"
	mbhandler.StopBits = 1
	mbhandler.SlaveId = 1
	mbhandler.Timeout = 1000 * time.Millisecond
	if verbose {
		mbhandler.Logger = log.New(os.Stdout, "RTUClientHandler: ", log.LstdFlags)
		log.Printf("Connecting to RTU via %s\r\n", rtuDevice)
	}

	err := mbhandler.Connect()
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	mbclient := modbus.NewClient(mbhandler)

	return &QueryEngine{client: mbclient, interval: interval,
		handler: *mbhandler, datastream: channel}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float32,
	err error) {
	results, err := q.client.ReadInputRegisters(opcode, 2)
	if err == nil {
		retval = RtuToFloat32(results)
	}
	return retval, err
}

func (q *QueryEngine) queryOrFail(opcode uint16) (retval float32) {
	var err error
	tryCnt := 0
	for tryCnt = 0; tryCnt < MaxRetryCount; tryCnt++ {
		retval, err = q.retrieveOpCode(opcode)
		if err != nil {
			log.Printf("Closing broken handler, reconnecting attempt %d\r\n", tryCnt)
			// Note: Just close the handler here. If a new handler is manually
			// created it will create a resource leak (file descriptors). Just
			// close the handler, the modbus library will recreate one as
			// needed.
			q.handler.Close()
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			break
		}
	}
	if tryCnt == MaxRetryCount {
		log.Fatal("Cannot query the sensor, reached maximum retry count.")
	}
	return retval
}

func (q *QueryEngine) Produce() {
	for {
		q.datastream <- Readings{
			Timestamp: time.Now(),
			Voltage: VoltageReadings{
				L1: q.queryOrFail(OpCodeL1Voltage),
				L2: q.queryOrFail(OpCodeL2Voltage),
				L3: q.queryOrFail(OpCodeL3Voltage),
			},
			Current: CurrentReadings{
				L1: q.queryOrFail(OpCodeL1Current),
				L2: q.queryOrFail(OpCodeL2Current),
				L3: q.queryOrFail(OpCodeL3Current),
			},
			Consumption: PowerReadings{
				L1: q.queryOrFail(OpCodeL1Consumption),
				L2: q.queryOrFail(OpCodeL2Consumption),
				L3: q.queryOrFail(OpCodeL3Consumption),
			},
			Production: PowerReadings{
				L1: q.queryOrFail(OpCodeL1Production),
				L2: q.queryOrFail(OpCodeL2Production),
				L3: q.queryOrFail(OpCodeL3Production),
			},
			Cosphi: CosphiReadings{
				L1: q.queryOrFail(OpCodeL1PowerFactor),
				L2: q.queryOrFail(OpCodeL2PowerFactor),
				L3: q.queryOrFail(OpCodeL3PowerFactor),
			},
		}
		time.Sleep(time.Duration(q.interval) * time.Second)
	}
	q.handler.Close()
}

func RtuToFloat32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
