# eece5811-hw4: ABA problem + lock free concurrent queues

## Q1:

The ABA problem occurs in concurrent program when a value at location L is read as A by thread T_1, T_1 gets preempted, another thread T_2 runs and changes the value at L from A to B then back again to A, T_2 gets preempted and T_1 resumes, T_1 sees the value A at L_i and assumes that it has not changed from A when in fact it has changed twice. This can lead to incorrect behavior in concurrent programs, especially when using lock-free data structures that rely on atomic operations like compare-and-swap (CAS). For example, in a queue implementation that uses CAS, if a thread reads the head pointer as A, then gets preempted while another thread dequeues an element (changing the head pointer to B) and then enqueues a new element (changing the head pointer back to A), when the first thread resumes, it may incorrectly assume that the head pointer is still valid and proceed with its operation, leading to data corruption or lost updates. This is the exact problem we saw in the lock-free queue implementation discussed in class when the algorithm used CAS to update pointers without additional checks.

A common approach of solving the ABA problem is to use versioning or tagging. This involves associating a version number or tag with the value at location L, where each time the value at location L is changed, the version number ot tag is incremented. When a thread reads the value, it also reads the version number or tag so when it performs a CAS operation, it can check both the value and the version number or tag to ensure that neither has changed since it was last read. If either has changed, the CAS operation fails, and the thread can retry its operation.

### Sources:
- Wikipedia article on ABA problem: https://en.wikipedia.org/wiki/ABA_problem
    - Credible source: Reviewed sources and citations. (Actually found the publication by Bjarne Stroustrup, et al. through this article)
- https://www.stroustrup.com/isorc2010.pdf
    - Credible source: Bjarne Stroustrup is the creator of C++ and this is a publication from a reputable conference.
- https://spcl.inf.ethz.ch/Teaching/2020-pp/lectures/PP-l21-ConcurrencyTheory.pdf
    - Credible source: Lecture notes from ETH Zurich, a reputable university. Sponsored by industry leaders like Google and Microsoft.

## Q2:

Current implementation of both my lock and lock-free queues do not work. I did not have enough time to debug them. Will fix and evaluate performance on my own time.
 
### Libraries Used:
- `"fmt"`: For input and output operations.
- `"sync"`: For synchronization primitives
- `"time"`: For benchmarking and adding delays