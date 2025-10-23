package main

import (
	"fmt"
	"sync"
	"time"
)

// typedef struct __node_t {
// 2 int value;
// 3 struct __node_t *next;
// 4 } node_t;
// Node represents a single queue element.
type Node struct {
	value int
	next  *Node
}

// typedef struct __queue_t {
// 7 node_t *head;
// 8 node_t *tail;
// 9 pthread_mutex_t head_lock, tail_lock;
// 10 } queue_t;
// Queue represents a thread-safe queue with separate head/tail locks.
type Queue struct {
	head      *Node
	tail      *Node
	headLock  sync.Mutex
	tailLock  sync.Mutex
}

// void Queue_Init(queue_t *q) {
// 13 node_t *tmp = malloc(sizeof(node_t));
// 14 tmp->next = NULL;
// 15 q->head = q->tail = tmp;
// 16 pthread_mutex_init(&q->head_lock, NULL);
// 17 pthread_mutex_init(&q->tail_lock, NULL);
// 18 }
// QueueInit initializes the queue with a dummy node.
func (q *Queue) Init() {
	tmp := &Node{}
	q.head = tmp
	q.tail = tmp
	// sync.Mutex doesn't require explicit initialization
}

// void Queue_Enqueue(queue_t *q, int value) {
// 21 node_t *tmp = malloc(sizeof(node_t));
// 22 assert(tmp != NULL);
// 23 tmp->value = value;
// 24 tmp->next = NULL;
// 25
// 26 pthread_mutex_lock(&q->tail_lock);
// 27 q->tail->next = tmp;
// 28 q->tail = tmp;
// 29 pthread_mutex_unlock(&q->tail_lock);
// 30 }
// Enqueue adds a new value to the tail of the queue.
func (q *Queue) Enqueue(value int) {
	tmp := &Node{value: value, next: nil}

	q.tailLock.Lock()
	q.tail.next = tmp
	q.tail = tmp
	q.tailLock.Unlock()
}

// int Queue_Dequeue(queue_t *q, int *value) {
// 33 pthread_mutex_lock(&q->head_lock);
// 34 node_t *tmp = q->head;
// 35 node_t *new_head = tmp->next;
// 36 if (new_head == NULL) {
// 37 pthread_mutex_unlock(&q->head_lock);
// 38 return -1; // queue was empty
// 39 }
// 40 *value = new_head->value;
// 41 q->head = new_head;
// 42 pthread_mutex_unlock(&q->head_lock);
// 43 free(tmp);
// 44 return 0;
// 45 }
// Dequeue removes a value from the head of the queue.
// Returns 0 on success, -1 if the queue was empty.
func (q *Queue) Dequeue(value *int) int {
	q.headLock.Lock()
	tmp := q.head
	newHead := tmp.next
	if newHead == nil {
		q.headLock.Unlock()
		return -1 // queue empty
	}
	*value = newHead.value
	q.head = newHead
	q.headLock.Unlock()
	// tmp is automatically garbage collected, no explicit free needed
	return 0
}

// Example usage
func main() {
	// var q Queue
	// q.Init()

	// q.Enqueue(10)
	// q.Enqueue(20)
	// q.Enqueue(30)

	// var val int
	// if q.Dequeue(&val) == 0 {
	// 	fmt.Println("Dequeued:", val)
	// }
	// if q.Dequeue(&val) == 0 {
	// 	fmt.Println("Dequeued:", val)
	// }
	// if q.Dequeue(&val) == 0 {
	// 	fmt.Println("Dequeued:", val)
	// }
	// if q.Dequeue(&val) == -1 {
	// 	fmt.Println("Queue empty")
	// }

	start := time.Now()
	defer func() {
		fmt.Printf("Execution Time: %s\n", time.Since(start))
	}()
	
	// initialize queue
	var q Queue
	q.Init()
	
	// number of goroutines
	numProducers := 10
	numConsumers := 10
	totalOps := 10000
	var wg sync.WaitGroup
	
	// Enqueueing producers
	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < totalOps; j++ {
				q.Enqueue(id*totalOps + j)
			}
		}(i)
	}

	// Dequeueing consumers
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var val int
			for {
				time.Sleep(time.Microsecond)
				fmt.Printf("Producer %d dequeuing count %d\n", id, val)
				if q.Dequeue(&val) == -1 {
					if val >= numProducers*totalOps-1 {
						break
					}
					continue
				}
			}
		}(i)
	}
	wg.Wait()

	fmt.Println("All operations completed.")
	fmt.Printf("Execution Time: %s\n", time.Since(start))

}