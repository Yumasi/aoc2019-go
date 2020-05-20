package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// Machine opcodes
const (
	Add         = 1
	Mul         = 2
	In          = 3
	Out         = 4
	JumpIfTrue  = 5
	JumpIfFalse = 6
	LessThan    = 7
	Equals      = 8
	Halt        = 99
)

// Machine states
const (
	Running = iota
	Stopped = iota
)

// Parameter modes
const (
	PositionMode  = 0
	ImmediateMode = 1
)

// IntMachine is an int machine
type IntMachine struct {
	state  int
	pc     int
	memory map[int]int
}

func newMemory(reader io.Reader) (res *IntMachine, err error) {
	var input []byte

	if input, err = ioutil.ReadAll(reader); err != nil {
		return
	}

	inputStr := string(input)
	inputStr = strings.TrimSpace(inputStr)
	res = &IntMachine{
		Stopped,
		0,
		map[int]int{},
	}

	tokens := strings.Split(inputStr, ",")

	for i, val := range tokens {
		var tmp int
		if tmp, err = strconv.Atoi(val); err != nil {
			return
		}

		res.memory[i] = tmp
	}

	return
}

// Exec runs the IntMachine
func (m *IntMachine) Exec() error {
	m.state = Running

	for m.state == Running {
		instr := m.getInstr()
		args := instr[1:]

		switch instr[0] {
		case Add:
			m.execAdd(args)
		case Mul:
			m.execMul(args)
		case In:
			m.execIn(args)
		case Out:
			m.execOut(args)
		case JumpIfTrue:
			m.execJumpIfTrue(args)
		case JumpIfFalse:
			m.execJumpIfFalse(args)
		case LessThan:
			m.execLessThan(args)
		case Equals:
			m.execEquals(args)
		case Halt:
			log.Printf("Halting.\nmem: %v\nPos zero: %d", m.memory, m.memory[0])
			m.state = Stopped
		default:
			return fmt.Errorf("invalid opcode: %v @%d", instr, m.pc)
		}
	}

	return nil
}

func (m *IntMachine) execAdd(args []int) {
	firstArg := args[0]
	secondArg := args[1]
	resultPos := args[2]

	m.memory[resultPos] = m.memory[firstArg] + m.memory[secondArg]

	m.pc += 4
}

func (m *IntMachine) execMul(args []int) {
	firstArg := args[0]
	secondArg := args[1]
	resultPos := args[2]

	m.memory[resultPos] = m.memory[firstArg] * m.memory[secondArg]

	m.pc += 4
}

func (m *IntMachine) execIn(args []int) {
	resultPos := args[0]

	var tmp int
	if _, err := fmt.Scanln(&tmp); err != nil {
		log.Fatalf("could not read standard input: %v", err)
	}

	m.memory[resultPos] = tmp

	m.pc += 2
}

func (m *IntMachine) execOut(args []int) {
	srcPos := args[0]
	if _, err := fmt.Println(m.memory[srcPos]); err != nil {
		log.Fatalf("could not print to standard output: %v", err)
	}

	m.pc += 2
}

func (m *IntMachine) execJumpIfTrue(args []int) {
	if m.memory[args[0]] != 0 {
		m.pc = m.memory[args[1]]
	} else {
		m.pc += 3
	}

}

func (m *IntMachine) execJumpIfFalse(args []int) {
	if m.memory[args[0]] == 0 {
		m.pc = m.memory[args[1]]
	} else {
		m.pc += 3
	}
}

func (m *IntMachine) execLessThan(args []int) {
	if m.memory[args[0]] < m.memory[args[1]] {
		m.memory[args[2]] = 1
	} else {
		m.memory[args[2]] = 0
	}

	m.pc += 4
}

func (m *IntMachine) execEquals(args []int) {
	if m.memory[args[0]] == m.memory[args[1]] {
		m.memory[args[2]] = 1
	} else {
		m.memory[args[2]] = 0
	}

	m.pc += 4
}

func (m *IntMachine) getInstr() []int {
	inst := m.memory[m.pc]
	res := make([]int, 4)

	res[0] = getOpcode(inst)
	for i := 0; i < 3; i++ {
		res[i+1] = m.getAddrFromMode(m.pc+i+1, getMode(inst, i))
	}

	return res
}

func (m *IntMachine) getAddrFromMode(arg, mode int) (res int) {
	switch mode {
	case PositionMode:
		res = m.memory[arg]
	case ImmediateMode:
		res = arg
	default:
		log.Fatal("invalid parameter mode")
	}

	return
}

func getOpcode(val int) int {
	return val % 100
}

func getMode(val, modeNb int) int {
	val = val / 100

	for i := 0; i < modeNb; i++ {
		val /= 10
	}

	return val % 10
}
