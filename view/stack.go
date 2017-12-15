package view

type Stack struct {
	top  *element
	size int
}

type element struct {
	value string // All types satisfy the empty interface, so we can store anything here.
	next  *element
}

// Return the stack's length
func (s *Stack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *Stack) Push(value string) {
	s.top = &element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *Stack) Read() (value string) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return ""
}
