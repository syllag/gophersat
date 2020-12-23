package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/crillab/gophersat/solver"

	"math/rand"
)

import "C"

var solvers map[int32]*solver.Solver

func init() {
	log.Println("gophersat dynlib is loaded")
	solvers = make(map[int32]*solver.Solver)
}

//export createSolver
func createSolver(nbVars int) int32 {
	id := rand.Int31()
	for _, ok := solvers[id]; ok; _, ok = solvers[id] {
		id = rand.Int31()
	}

	cnf := [][]int{}
	pb := solver.ParseSliceNb(cnf, nbVars)

	solvers[id] = solver.New(pb)

	return id
}

//export deleteSolver
// -1: unknown ID number
//  0: OK
func deleteSolver(solverID int32) int {
	_, ok := solvers[solverID]
	if ok {
		delete(solvers, solverID)
		return 0
	}
	return -1 // unknown ID number
}

//export addClause
// -1: unknown solver ID
//  0: OK
func addClause(solverID int32, clause *int32, size int32) int {
	s, ok := solvers[solverID]

	if !ok {
		return -1 // unknown ID number
	}

	// magic hack from https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	c := (*[1<<32 - 1]int32)(unsafe.Pointer(clause))[:size:size]
	lits := solver.IntsToLits(c...)
	s.AppendClause(solver.NewClause(lits))

	return 0
}

//export solve
// -1: unknown solver ID
//  1: SAT
//  2: UNSAT
func solve(solverID int32) int {
	s, ok := solvers[solverID]
	if !ok {
		return -1
	}
	return int(s.Solve())
}

//export nbVars
// -1: unknown solver ID
//  0: OK
func nbVars(solverID int32) int {
	s, ok := solvers[solverID]
	if !ok {
		return -1
	}
	return s.NbVars()
}

//export model
// -1: unknown solver ID
//  0: OK
func model(solverID int32, buffSize int, buff *bool) int {
	s, ok := solvers[solverID]
	if !ok {
		return -1
	}

	m := s.Model()
	b := (*[1<<32 - 1]bool)(unsafe.Pointer(buff))[:buffSize:buffSize]

	for i, e := range m {
		b[i] = e
	}

	return 0
}

//export assume
// -1: unknown solver ID
//  1: SAT
//  2: UNSAT
func assume(solverID int32, lit *int32, nbLits int32) int {
	s, ok := solvers[solverID]
	if !ok {
		return -1
	}

	c := (*[1<<32 - 1]int32)(unsafe.Pointer(lit))[:nbLits:nbLits]
	lits := solver.IntsToLits(c...)

	return int(s.Assume(lits))
}

func main() {
	id := createSolver(12)
	fmt.Printf("id: %d\nsolvers: %v\n", id, solvers)

	fmt.Printf("status: %s\n", solvers[id].Solve()) // SAT
	solvers[id].AppendClause(solver.NewClause(nil))
	fmt.Printf("status: %s\n", solvers[id].Solve()) // UNSAT
}
