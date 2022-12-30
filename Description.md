Implement a distributed system consisting of a set of nodes that provides a Distributed Hash Table (DHT) function. The system is such that: 
1. it provides two operations called *put* and *get*
2. operation *put* has two parameters (it has two parameters of type *int*), a key and a value, and returns a boolean (it returns a value of type *bool*)
   1. operation *put* updates the hash table at the requested key entry with the given value
   2. operation *put* returns a confirmation of whether it succeeded or not
3. operation *get* has one parameter (a key) and it returns an integer (a value)
   1. operation *get* accesses the entry specified by the given key
   2. operation *get* returns the value associated to the given key
4. **(Property 1 - Repeatable Read)** if a call *put(key,val)* is made, with no subsequent calls to put, then a new call *get(key)* to that node will return *val*. 
5. **(Property 2 - Initialization to zero)** if no call to *put* has been made for a given key, any *get* call returns 0 (zero). 
6. **(Property 3 - Consistency)** given a non-empty sequence of successful (returning **true**) put operations performed on the hash table, any subsequent get operation at any node should return the same value
7. **(Property 4 - Liveness)** every call to *put* and *get* returns a value. If a call *put(k,v)* returns **false**, then some call in a repeated sequence of *put(k,v)* will return **true** (eventually, *put* and *get* succeed).
8. **(Property 5 - Reliability)** your system can tolerate a crash-failure of 1 node.

Partial submissions are accepted, e.g., a not fully working implementation.  In this situation, a pseudocode answer that describes the main algorithm and ideas is acceptable, but working *Go* code gives extra credits. If you hand in only pseudo code, this must be included in the report.txt file (or report.pdf).

Supplement your working implementation with a log.txt file that shows a correct execution (only if you are submitting a working implementation).
