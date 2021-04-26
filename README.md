# go-coloring

## Description

We will be implementing 4 algorithms that do a O(Δ+1) reduction on a k-coloring, where the graph has max degree Δ. 
One of these algorithms is non-parallel, while three others are distributed and must be run in parallel workloads. 
We feel that Golang is especially suited to this type of task since goroutines essentially function as lightweight, flexible, threads. 
We will test the success of these algorithms using an iterative algorithm to verify that a graph is correctly colored. 
We will also benchmark and visualize the differences in runtime between the different algorithms, as well as output the correct graph coloring given by each algorithm.

## Rubric

| Section (Person) | Description |  Points |
|----------|-------------|------:|
| Testing (Tyler, Thomas) | The harness will assert the validity of the color-reduction using visualization (see below) for smaller graphs and the Verifier for larger graphs. We will know that the actual algorithm has been implemented based on comparisons of established runtimes and our empirical runtime trends. | 10 |
| Visualization (Michael, Thomas) | An informative visualization of both the graph outputs themselves and runtime analysis trends of the various algorithms. This will be done using go-echarts. |   10 |
| Algorithms to Implement | 4 of the 5 following algorithms will be implemented in this project. We are planning on doing the four specified below, but in the case that insurmountable difficulties arise in the implementation of one of the algorithms, we will look up one other distributed algorithm and implement that in its place. <br/> - Naive color-reduction (not parallel) to O(Δ+1) colors **Tyler** <br/> - Kuhn-Wattenhofer color-reduction (parallel) to O(Δ+1) colors **Michael** <br/> - Linial’s Algorithm with Kuhn-Wattenhofer color-reduction applied to its result O(Δ+1) colors **Thomas** <br/> - Cole-Vishkin Algorithm to O(Δ+1) colors **Tyler** <br/> - One other color reduction algorithm to O(Δ+1) colors | 40 (10 each) |
| Extensions (All) | As the project progresses, we will likely implement either: slowed iteration to display color-reduction in real time, another language’s implementation for runtime comparison, or more algorithms. | 5 |
| Total | Everything | 65 |

## Resources

- [Graph coloring and isSafe()](https://www.geeksforgeeks.org/m-coloring-problem-backtracking-5/)
- [Additional graph coloring](https://www.geeksforgeeks.org/graph-coloring-set-2-greedy-algorithm/)
- [Core color-reduction examples](https://stanford.edu/~rezab/classes/cme323/S16/projects_reports/bae.pdf)
- [Additional color-reduction examples](https://www.cs.bgu.ac.il/~elkinm/book.pdf)
- [Golang visualizations](https://golangexample.com/a-graceful-data-visualizing-library-for-golang)
- [Distributed Largest-First](https://ieeexplore-ieee-org.proxy.lib.duke.edu/document/1115204)
- [Cole-Vishkin Reduction](https://www.zhengqunkoo.com:8443/zhengqunkoo/site/src/commit/ebbab6e24911a02c97b380f2e39f06d9c3e83770/worker.js)
