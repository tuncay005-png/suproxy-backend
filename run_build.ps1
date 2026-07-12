go test ./test/integration -c -o integration.test.exe 2>&1 | Tee-Object -FilePath fresh_compile_errors.txt
