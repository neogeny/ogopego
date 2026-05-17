//go:build ignore
// +build ignore

// Initialize an empty stack (slice)
var stack []int

// PUSH: Append an element to the end
stack = append(stack, 42)

// POP: Get the last element and slice it off
if len(stack) > 0 {
    // 1. Get the top element
    top := stack[len(stack)-1]

    // 2. Remove it by reslicing
    stack = stack[:len(stack)-1]
}
