You have to implement the following distributed system: A set of nodes tries to implement a distributed Increment function. We want the following properties to hold:

1. Every call to Increment returns a value. Value of Increment is always a natural number.
2. First call to Increment returns 0.
3. Monotonicity: For every node, if it has subsequent calls to Increment, with values Inc_1 and Inc_2 respectively, we have that Inc_1 < Inc_2
4. Uniqueness: No call to Increment should ever return a previously returned value.
5. Liveness: every call to Increment returns a value
6. Reliability: your system can tolerate a crash-failure of 1 node.
Partial submissions are accepted, e.g., a not fully working implementation. In this situation, a pseudocode answer that describes the main algorithm and ideas is acceptable, but working Golang code gives full points.
