# 51Degrees Go API
# ----------------
# Makefile copies required sources to go/src directory and generates the swig 
# wrapper. This is to facilitate the building of Swig generated wrappers and 
# underlying C/C++ sources. The root /src/ directory contains other C/C++ 
# sources not pertaining to the Go implementation. 
# If building against the src directory in the repository root is attempted, 
# `go build` will fail as it will attempt to compile everything in the
# directory, including other wrappers files.

           
all: 
	swig -go -cgo -c++ -intgosize 64 -o 51Degrees_go.cxx 51Degrees.i
