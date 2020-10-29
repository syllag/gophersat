package main

import (
	"fmt"
	"unsafe"

	"github.com/crillab/gophersat/solver"

	"math/rand"
)

import "C"

var solvers map[int32]*solver.Solver

func init() {
	fmt.Println("gophersat dynlib is loaded")
	solvers = make(map[int32]*solver.Solver)
}

//export createSolver
func createSolver(nbVars int) int32 {
	id := rand.Int31()
	for _, ok := solvers[id]; ok; _, ok = solvers[id] {
		id = rand.Int31()
	}

	// NewProblem avec un nombre de variables ?
	// var pb solver.Problem
	// pb.NbVars = nbVars

	cnf := [][]int{}
	pb := solver.ParseSliceNb(cnf, nbVars)

	solvers[id] = solver.New(pb)

	return id
}

//export addClause
// 1: unknown ID number
// 0: OK
func addClause(solverID int32, clause *int32, size int32) int {
	s, ok := solvers[solverID]

	if !ok {
		return 1 // unknown ID number
	}

	// hack from https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	c := (*[1<<32 - 1]int32)(unsafe.Pointer(clause))[:size:size]
	lits := solver.IntsToLits(c...)
	s.AppendClause(solver.NewClause(lits))

	return 0
}

//export solve
// 0: UNSAT
// 1: SAT
func solve(solverID int32) int32 {
	return int32(solvers[solverID].Solve())
}

//func model(solverID int) []bool

//func (s *Solver) Model() []bool

//export sum
func sum(liste *int32, len int32) int32 {
	ptrSize := unsafe.Sizeof(liste)
	fmt.Println("ptrSize:", ptrSize)

	// from https://github.com/golang/go/wiki/cgo
	slice := (*[1<<32 - 1]int32)(unsafe.Pointer(liste))[:len:len]
	var s int32
	for i, val := range slice {
		fmt.Printf("%d: %d\n", i, val)
		s += val
	}
	return s
}

func main() {
	id := createSolver(12)
	fmt.Printf("id: %d\nsolvers: %v\n", id, solvers)

	fmt.Printf("status: %s\n", solvers[id].Solve()) // SAT
	solvers[id].AppendClause(solver.NewClause(nil))
	fmt.Printf("status: %s\n", solvers[id].Solve()) // UNSAT
}
