# tests/stdlibs

This directory contains test-specific standard libraries. These are only
available when testing gno code in `_test.gno` and `_filetest.gno` files.
Re-declarations of functions already existing override the definitions of the
normal stdlibs directory.

Adding imports that don't exist in the corresponding normal stdlib is undefined behavior
