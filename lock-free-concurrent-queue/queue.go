package main

import (
	"fmt"
	"sync/atomic"
)

// structure node t fvalue: data type, next: pointer t
// node_t represents a queue node.
type node_t struct {
	value int
	next  atomic.Pointer[node_t]
}

// structure queue t fHead: pointer t, Tail: pointer tg
// queue_t represents the concurrent lock-free queue.
type queue_t struct {
	Head atomic.Pointer[node_t]
	Tail atomic.Pointer[node_t]
}

// initialize(Q: pointer to queue t)
// node = new node() # Allocate a free node
// node–>next.ptr = NULL # Make it the only node in the linked list
// Q–>Head = Q–>Tail = node
// initialize sets up a dummy node so that Head and Tail point to it.
func initialize(Q *queue_t) {
	node := &node_t{}
	Q.Head.Store(node)
	Q.Tail.Store(node)
}

// enqueue(Q: pointer to queue t, value: data type)
// E1: node = new node() # Allocate a new node from the free list
// E2: node–>value = value # Copy enqueued value into node
// E3: node–>next.ptr = NULL # Set next pointer of node to NULL
// E4: loop # Keep trying until Enqueue is done
// E5: tail = Q–>Tail # Read Tail.ptr and Tail.count together
// E6: next = tail.ptr–>next # Read next ptr and count fields together
// E7: if tail == Q–>Tail # Are tail and next consistent?
// E8: if next.ptr == NULL # Was Tail pointing to the last node?
// E9: if CAS(&tail.ptr–>next, next, <node, next.count+1>) # Try to link node at the end of the linked list
// E10: break # Enqueue is done. Exit loop
// E11: endif
// E12: else # Tail was not pointing to the last node
// E13: CAS(&Q–>Tail, tail, <next.ptr, tail.count+1>) # Try to swing Tail to the next node
// E14: endif
// E15: endif
// E16: endloop
// E17: CAS(&Q–>Tail, tail, <node, tail.count+1>)
// enqueue adds a new node with the given value to the queue.
func enqueue(Q *queue_t, value int) {
	node := &node_t{value: value} // Create new node and copy value
	for {                         // Keep trying until Enqueue is done
		tail := Q.Tail.Load()      // Read Tail
		next := tail.next.Load()   // Read next
		if tail == Q.Tail.Load() { // Are tail and next consistent?
			if next == nil { // Was Tail pointing to the last node?
				if tail.next.CompareAndSwap(nil, node) { // Try to link node at the end of the linked list
					break // Enqueue is done. Exit loop
				}
			} else { // Tail was not pointing to the last node
				Q.Tail.CompareAndSwap(tail, next) // Try to swing Tail to the next node
			}
		}
	}
	Q.Tail.CompareAndSwap(Q.Tail.Load(), node) // Swing Tail to the inserted node
}

// dequeue(Q: pointer to queue t, pvalue: pointer to data type): boolean
// D1: loop # Keep trying until Dequeue is done
// D2: head = Q–>Head # Read Head
// D3: tail = Q–>Tail # Read Tail
// D4: next = head–>next # Read Head.ptr–>next
// D5: if head == Q–>Head # Are head, tail, and next consistent?
// D6: if head.ptr == tail.ptr # Is queue empty or Tail falling behind?
// D7: if next.ptr == NULL # Is queue empty?
// D8: return FALSE # Queue is empty, couldn’t dequeue
// D9: endif
// D10: CAS(&Q–>Tail, tail, <next.ptr, tail.count+1>) # Tail is falling behind. Try to advance it
// D11: else # No need to deal with Tail
// # Read value before CAS, otherwise another dequeue might free the next node
// D12: *pvalue = next.ptr–>value
// D13: if CAS(&Q–>Head, head, <next.ptr, head.count+1>) # Try to swing Head to the next node
// D14: break # Dequeue is done. Exit loop
// D15: endif
// D16: endif
// D17: endif
// D18: endloop
// D19: free(head.ptr) # It is safe now to free the old dummy node
// D20: return TRUE # Queue was not empty, dequeue succeeded
// dequeue removes one element from the queue.
// Returns true if successful, false if queue was empty.
func dequeue(Q *queue_t, pvalue *int) bool {
	for { // Keep trying until Dequeue is done
		head := Q.Head.Load()      // Read Head
		tail := Q.Tail.Load()      // Read Tail
		next := head.next.Load()   // Read Head.next
		if head == Q.Head.Load() { // Are head, tail, and next consistent?
			if head == tail { // Is queue empty or Tail falling behind?
				if next == nil { // Is queue empty?
					return false // Queue is empty, couldn’t dequeue
				}
				Q.Tail.CompareAndSwap(tail, next) // Tail is falling behind. Try to advance it
			} else { // No need to deal with Tail
				*pvalue = next.value                   // Read value before CAS, otherwise another dequeue might free the next node
				if Q.Head.CompareAndSwap(head, next) { // Try to swing Head to the next node
					break // Dequeue is done. Exit loop
				}
			}
		}
	}
	return true // Queue was not empty, dequeue succeeded
}

// Example usage
func main() {
	
}
