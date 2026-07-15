@echo off
echo Compiling integration tests...
go test -c ./test/integration > test_compile_output.txt 2>&1
echo Done. Check test_compile_output.txt for results.
type test_compile_output.txt
