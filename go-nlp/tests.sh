#!/bin/bash

failed=0

for pkg in $(go list ./...); do
    echo "Testing: $pkg"

    if ! go test "$pkg"; then
        echo "❌ FAILED: $pkg"
        (( failed+=1 ))
    else
        echo "✅ PASSED: $pkg"
    fi

    echo
done

if [ $failed -eq 0 ]; then
    echo "✅ Todos los tests pasaron."
else
    echo "❌ Algunos tests fallaron."
    echo $failed
    exit 1
fi