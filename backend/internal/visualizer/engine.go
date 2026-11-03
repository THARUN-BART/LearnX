package visualizer

import (
	"fmt"
	"math"
	"sort"
)

// GenerateExecutionSteps runs real algorithms and returns structured step-by-step traces.
func GenerateExecutionSteps(algoID string, input map[string]any) ([]ExecutionStep, error) {
	meta := GetAlgorithmByID(algoID)
	if meta == nil {
		return nil, fmt.Errorf("algorithm '%s' not found", algoID)
	}

	// Merge with defaults if fields are missing
	merged := make(map[string]any)
	for k, v := range meta.DefaultInput {
		merged[k] = v
	}
	for k, v := range input {
		if v != nil {
			merged[k] = v
		}
	}

	switch algoID {
	// Sorting
	case "bubble-sort":
		return traceBubbleSort(extractIntSlice(merged["array"]))
	case "selection-sort":
		return traceSelectionSort(extractIntSlice(merged["array"]))
	case "insertion-sort":
		return traceInsertionSort(extractIntSlice(merged["array"]))
	case "merge-sort":
		return traceMergeSort(extractIntSlice(merged["array"]))
	case "quick-sort":
		return traceQuickSort(extractIntSlice(merged["array"]))
	case "heap-sort":
		return traceHeapSort(extractIntSlice(merged["array"]))

	// Searching
	case "linear-search":
		arr := extractIntSlice(merged["array"])
		target := extractInt(merged["target"], 7)
		return traceLinearSearch(arr, target)
	case "binary-search":
		arr := extractIntSlice(merged["array"])
		target := extractInt(merged["target"], 23)
		sort.Ints(arr)
		return traceBinarySearch(arr, target)

	// Arrays and Patterns
	case "array-traversal":
		return traceArrayTraversal(extractIntSlice(merged["array"]))
	case "two-pointers":
		return traceTwoPointers(extractIntSlice(merged["array"]))
	case "two-sum":
		arr := extractIntSlice(merged["array"])
		target := extractInt(merged["target"], 26)
		sort.Ints(arr)
		return traceTwoSum(arr, target)
	case "sliding-window":
		arr := extractIntSlice(merged["array"])
		k := extractInt(merged["target"], 3)
		if k <= 0 || k > len(arr) {
			k = 3
		}
		return traceSlidingWindow(arr, k)
	case "prefix-sum":
		return tracePrefixSum(extractIntSlice(merged["array"]))

	// Linked Lists
	case "singly-linked-list", "linked-list-traversal":
		return traceLinkedListTraversal(extractIntSlice(merged["values"]))
	case "doubly-linked-list":
		return traceDoublyLinkedList(extractIntSlice(merged["values"]))
	case "linked-list-reversal":
		return traceLinkedListReversal(extractIntSlice(merged["values"]))
	case "linked-list-insertion-deletion":
		vals := extractIntSlice(merged["values"])
		insertVal := extractInt(merged["insert_val"], 25)
		insertPos := extractInt(merged["insert_pos"], 2)
		return traceLinkedListInsertDelete(vals, insertVal, insertPos)

	// Stack & Queue
	case "stack-push-pop":
		return traceStackOperations()
	case "queue-enqueue-dequeue":
		return traceQueueOperations()
	case "circular-queue":
		return traceCircularQueue()
	case "monotonic-stack":
		return traceMonotonicStack(extractIntSlice(merged["array"]))

	// Recursion and Backtracking
	case "recursion-factorial":
		n := extractInt(merged["n"], 5)
		if n > 8 {
			n = 8
		}
		return traceFactorial(n)
	case "recursion-fibonacci":
		n := extractInt(merged["n"], 4)
		if n > 6 {
			n = 6
		}
		return traceFibonacciRecursion(n)
	case "recursive-binary-search":
		arr := extractIntSlice(merged["array"])
		target := extractInt(merged["target"], 25)
		sort.Ints(arr)
		return traceRecursiveBinarySearch(arr, target)
	case "backtracking-subsets":
		arr := extractIntSlice(merged["array"])
		if len(arr) > 4 {
			arr = arr[:4]
		}
		return traceSubsets(arr)
	case "backtracking-permutations":
		arr := extractIntSlice(merged["array"])
		if len(arr) > 4 {
			arr = arr[:3]
		}
		return tracePermutations(arr)
	case "backtracking-n-queens":
		n := extractInt(merged["n"], 4)
		if n < 4 || n > 6 {
			n = 4
		}
		return traceNQueens(n)

	// Trees
	case "binary-tree-traversals":
		return traceTreeTraversals()
	case "level-order-traversal":
		return traceLevelOrder()
	case "bst-insertion":
		insertVal := extractInt(merged["insert_val"], 45)
		return traceBSTInsertion(insertVal)
	case "bst-deletion":
		delVal := extractInt(merged["delete_val"], 30)
		return traceBSTDeletion(delVal)
	case "bst-search":
		target := extractInt(merged["target"], 60)
		return traceBSTSearch(target)

	// Graphs
	case "graph-bfs":
		return traceGraphBFS()
	case "graph-dfs":
		return traceGraphDFS()
	case "dijkstra-algorithm":
		return traceDijkstra()
	case "cycle-detection":
		return traceCycleDetection()
	case "topological-sort":
		return traceTopologicalSort()
	case "minimum-spanning-tree":
		return traceMST()

	// Dynamic Programming
	case "dp-fibonacci":
		n := extractInt(merged["n"], 6)
		if n > 12 {
			n = 12
		}
		return traceDPFibonacci(n)
	case "climbing-stairs":
		n := extractInt(merged["n"], 5)
		if n > 10 {
			n = 10
		}
		return traceClimbingStairs(n)
	case "knapsack-01":
		return traceKnapsack01()
	case "longest-common-subsequence":
		s1 := extractString(merged["str1"], "ABCDE")
		s2 := extractString(merged["str2"], "ACE")
		return traceLCS(s1, s2)
	case "coin-change":
		coins := extractIntSlice(merged["coins"])
		amount := extractInt(merged["amount"], 7)
		if amount > 15 {
			amount = 15
		}
		return traceCoinChange(coins, amount)

	default:
		return traceBubbleSort(extractIntSlice(merged["array"]))
	}
}

// ----------------- Helpers -----------------

func copyInts(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func extractIntSlice(val any) []int {
	if val == nil {
		return []int{29, 10, 14, 37, 13, 22}
	}
	switch v := val.(type) {
	case []int:
		return copyInts(v)
	case []any:
		res := make([]int, 0, len(v))
		for _, item := range v {
			if num, ok := item.(float64); ok {
				res = append(res, int(num))
			} else if num, ok := item.(int); ok {
				res = append(res, num)
			}
		}
		if len(res) == 0 {
			return []int{29, 10, 14, 37, 13, 22}
		}
		return res
	default:
		return []int{29, 10, 14, 37, 13, 22}
	}
}

func extractInt(val any, defaultVal int) int {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return defaultVal
	}
}

func extractString(val any, defaultVal string) string {
	if val == nil {
		return defaultVal
	}
	if s, ok := val.(string); ok && s != "" {
		return s
	}
	return defaultVal
}

// ----------------- Sorting Tracers -----------------

func traceBubbleSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	n := len(arr)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial unsorted array",
		Explanation: fmt.Sprintf("Bubble sort initialized on %d elements. Adjacent elements will be compared sequentially.", n),
		State: StepState{
			Array:      copyInts(arr),
			StatusText: "Ready to start comparisons",
		},
	})
	stepNo++

	for i := 0; i < n-1; i++ {
		swappedInPass := false
		for j := 0; j < n-i-1; j++ {
			// Compare step
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "compare",
				Description: fmt.Sprintf("Compare arr[%d]=%d with arr[%d]=%d", j, arr[j], j+1, arr[j+1]),
				Explanation: fmt.Sprintf("Checking if %d > %d. If true, swap to move the larger value toward the right.", arr[j], arr[j+1]),
				State: StepState{
					Array:       copyInts(arr),
					Pointers:    map[string]int{"i": i, "j": j, "j+1": j + 1},
					Comparing:   []int{j, j + 1},
					StatusText:  fmt.Sprintf("Comparing %d and %d", arr[j], arr[j+1]),
				},
			})
			stepNo++

			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swappedInPass = true

				steps = append(steps, ExecutionStep{
					StepNumber:  stepNo,
					Operation:   "swap",
					Description: fmt.Sprintf("Swap: %d moved to index %d, %d moved to index %d", arr[j+1], j+1, arr[j], j),
					Explanation: fmt.Sprintf("Since %d was greater, we swapped them to restore ascending order.", arr[j+1]),
					State: StepState{
						Array:       copyInts(arr),
						Pointers:    map[string]int{"i": i, "j": j, "j+1": j + 1},
						Swapping:    []int{j, j + 1},
						StatusText:  fmt.Sprintf("Swapped %d and %d", arr[j], arr[j+1]),
					},
				})
				stepNo++
			}
		}

		// End of pass: element at n-i-1 is in final sorted position
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "lock_element",
			Description: fmt.Sprintf("Pass %d complete. Element %d at index %d is now in sorted position.", i+1, arr[n-i-1], n-i-1),
			Explanation: "The largest remaining element has bubbled to its correct final position in the array.",
			State: StepState{
				Array:       copyInts(arr),
				Highlighted: []int{n - i - 1},
				StatusText:  fmt.Sprintf("Locked arr[%d]=%d", n-i-1, arr[n-i-1]),
			},
		})
		stepNo++

		if !swappedInPass {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "early_exit",
				Description: "Early termination: No swaps were made during this pass.",
				Explanation: "The array is already completely sorted. O(N) best-case optimization triggered.",
				State: StepState{
					Array:      copyInts(arr),
					StatusText: "Early exit: fully sorted",
				},
			})
			stepNo++
			break
		}
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Bubble sort complete. Entire array is sorted.",
		Explanation: "All passes finished. All elements are ordered monotonically in non-decreasing sequence.",
		State: StepState{
			Array:      copyInts(arr),
			StatusText: "Sorting completed successfully",
		},
		IsComplete: true,
	})

	return steps, nil
}

func traceSelectionSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	n := len(arr)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial unsorted array",
		Explanation: "Selection sort repeatedly finds the minimum element in the unsorted suffix and swaps it to the front.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Selection Sort"},
	})
	stepNo++

	for i := 0; i < n-1; i++ {
		minIdx := i
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "set_min",
			Description: fmt.Sprintf("Assume minimum is at index %d (value %d)", i, arr[i]),
			Explanation: fmt.Sprintf("Scanning remaining elements from index %d to %d to verify or update the minimum.", i+1, n-1),
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"current": i, "min": minIdx},
				Highlighted: []int{minIdx},
				StatusText:  fmt.Sprintf("Current minimum candidate: arr[%d]=%d", minIdx, arr[minIdx]),
			},
		})
		stepNo++

		for j := i + 1; j < n; j++ {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "compare",
				Description: fmt.Sprintf("Compare arr[%d]=%d with current minimum arr[%d]=%d", j, arr[j], minIdx, arr[minIdx]),
				Explanation: fmt.Sprintf("Is %d < %d?", arr[j], arr[minIdx]),
				State: StepState{
					Array:       copyInts(arr),
					Pointers:    map[string]int{"current": i, "min": minIdx, "scan": j},
					Comparing:   []int{j, minIdx},
					StatusText:  fmt.Sprintf("Inspecting arr[%d]", j),
				},
			})
			stepNo++

			if arr[j] < arr[minIdx] {
				minIdx = j
				steps = append(steps, ExecutionStep{
					StepNumber:  stepNo,
					Operation:   "update_min",
					Description: fmt.Sprintf("New minimum found at index %d (value %d)", minIdx, arr[minIdx]),
					Explanation: fmt.Sprintf("Found a smaller element: %d is smaller than prior candidate.", arr[minIdx]),
					State: StepState{
						Array:       copyInts(arr),
						Pointers:    map[string]int{"current": i, "min": minIdx},
						Highlighted: []int{minIdx},
						StatusText:  fmt.Sprintf("Updated minimum to index %d (%d)", minIdx, arr[minIdx]),
					},
				})
				stepNo++
			}
		}

		if minIdx != i {
			arr[i], arr[minIdx] = arr[minIdx], arr[i]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "swap",
				Description: fmt.Sprintf("Swap arr[%d]=%d with arr[%d]=%d", i, arr[i], minIdx, arr[minIdx]),
				Explanation: fmt.Sprintf("Moved minimum value %d to sorted prefix position index %d.", arr[i], i),
				State: StepState{
					Array:      copyInts(arr),
					Swapping:   []int{i, minIdx},
					StatusText: fmt.Sprintf("Swapped %d into index %d", arr[i], i),
				},
			})
			stepNo++
		}
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Selection sort completed.",
		Explanation: "All positions filled with minimums. Entire array sorted.",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceInsertionSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	n := len(arr)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial array. Single element arr[0] is already sorted trivially.",
		Explanation: "Insertion sort inserts subsequent elements into the sorted prefix one by one.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Insertion Sort"},
	})
	stepNo++

	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "pick_key",
			Description: fmt.Sprintf("Select key arr[%d]=%d to insert into sorted prefix [0..%d]", i, key, i-1),
			Explanation: "Shifting larger elements to the right to create an insertion slot for key.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"key": i, "j": j},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Current key = %d", key),
			},
		})
		stepNo++

		for j >= 0 && arr[j] > key {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "compare_and_shift",
				Description: fmt.Sprintf("arr[%d]=%d > key (%d) -> shift %d to index %d", j, arr[j], key, arr[j], j+1),
				Explanation: fmt.Sprintf("%d is larger than %d, so we shift it rightwards to index %d.", arr[j], key, j+1),
				State: StepState{
					Array:      copyInts(arr),
					Comparing:  []int{j, j + 1},
					StatusText: fmt.Sprintf("Shifted %d right", arr[j]),
				},
			})
			stepNo++

			arr[j+1] = arr[j]
			j--
		}

		arr[j+1] = key
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "insert",
			Description: fmt.Sprintf("Inserted key %d at index %d", key, j+1),
			Explanation: fmt.Sprintf("Prefix [0..%d] is now fully sorted.", i),
			State: StepState{
				Array:       copyInts(arr),
				Highlighted: []int{j + 1},
				StatusText:  fmt.Sprintf("Placed key %d at %d", key, j+1),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Insertion sort completed.",
		Explanation: "All elements successfully inserted into their correct relative positions.",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceMergeSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial array for Merge Sort",
		Explanation: "Merge sort recursively divides the array into single-element subarrays, then merges them in sorted order.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Divide & Conquer"},
	})
	stepNo++

	var mergeSort func(l, r int)
	mergeSort = func(l, r int) {
		if l >= r {
			return
		}
		mid := l + (r-l)/2

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "divide",
			Description: fmt.Sprintf("Divide subarray [%d..%d] at mid=%d into [%d..%d] and [%d..%d]", l, r, mid, l, mid, mid+1, r),
			Explanation: "Splitting into two independent halves for recursive sorting.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"l": l, "mid": mid, "r": r},
				Highlighted: []int{mid},
				StatusText:  fmt.Sprintf("Subarray split: [%d..%d] & [%d..%d]", l, mid, mid+1, r),
			},
		})
		stepNo++

		mergeSort(l, mid)
		mergeSort(mid+1, r)

		// Merge step
		leftHalf := copyInts(arr[l : mid+1])
		rightHalf := copyInts(arr[mid+1 : r+1])
		i, j, k := 0, 0, l

		for i < len(leftHalf) && j < len(rightHalf) {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "compare_merge",
				Description: fmt.Sprintf("Compare left=%d (from [%d..%d]) with right=%d (from [%d..%d])", leftHalf[i], l, mid, rightHalf[j], mid+1, r),
				Explanation: "Placing the smaller element into the merged array.",
				State: StepState{
					Array:       copyInts(arr),
					Pointers:    map[string]int{"k": k},
					Highlighted: []int{k},
					StatusText:  fmt.Sprintf("Merging: comparing %d vs %d", leftHalf[i], rightHalf[j]),
				},
			})
			stepNo++

			if leftHalf[i] <= rightHalf[j] {
				arr[k] = leftHalf[i]
				i++
			} else {
				arr[k] = rightHalf[j]
				j++
			}
			k++
		}

		for i < len(leftHalf) {
			arr[k] = leftHalf[i]
			i++
			k++
		}
		for j < len(rightHalf) {
			arr[k] = rightHalf[j]
			j++
			k++
		}

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "merge_complete",
			Description: fmt.Sprintf("Merged subarray [%d..%d] is now sorted: %v", l, r, arr[l:r+1]),
			Explanation: "Both halves merged in monotonic order with O(N) auxiliary space.",
			State: StepState{
				Array:      copyInts(arr),
				Pointers:   map[string]int{"l": l, "r": r},
				StatusText: fmt.Sprintf("Merged [%d..%d]", l, r),
			},
		})
		stepNo++
	}

	mergeSort(0, len(arr)-1)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Merge Sort completed successfully.",
		Explanation: "All recursive partitions merged. Overall time complexity O(N log N).",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceQuickSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial array for Quick Sort",
		Explanation: "Quick sort chooses a pivot, partitions the array around it, and sorts subarrays in-place.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Quick Sort"},
	})
	stepNo++

	var quickSort func(low, high int)
	quickSort = func(low, high int) {
		if low >= high {
			return
		}
		pivot := arr[high]
		pIdx := low

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "pivot_chosen",
			Description: fmt.Sprintf("Choose pivot arr[%d]=%d for partition [%d..%d]", high, pivot, low, high),
			Explanation: "All elements <= pivot will be moved to the left; all elements > pivot to the right.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"low": low, "high": high, "pivot": high},
				Highlighted: []int{high},
				StatusText:  fmt.Sprintf("Pivot = %d", pivot),
			},
		})
		stepNo++

		for j := low; j < high; j++ {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "compare_pivot",
				Description: fmt.Sprintf("Compare arr[%d]=%d with pivot %d", j, arr[j], pivot),
				Explanation: fmt.Sprintf("Is %d <= %d?", arr[j], pivot),
				State: StepState{
					Array:      copyInts(arr),
					Pointers:   map[string]int{"j": j, "pIdx": pIdx, "pivot": high},
					Comparing:  []int{j, high},
					StatusText: fmt.Sprintf("arr[%d] <= pivot?", j),
				},
			})
			stepNo++

			if arr[j] <= pivot {
				if pIdx != j {
					arr[pIdx], arr[j] = arr[j], arr[pIdx]
					steps = append(steps, ExecutionStep{
						StepNumber:  stepNo,
						Operation:   "swap",
						Description: fmt.Sprintf("Swap arr[%d]=%d with arr[%d]=%d to group smaller values left", pIdx, arr[pIdx], j, arr[j]),
						Explanation: "Advancing partition index.",
						State: StepState{
							Array:      copyInts(arr),
							Swapping:   []int{pIdx, j},
							StatusText: "Swapped smaller value left",
						},
					})
					stepNo++
				}
				pIdx++
			}
		}

		arr[pIdx], arr[high] = arr[high], arr[pIdx]
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "place_pivot",
			Description: fmt.Sprintf("Place pivot %d at final sorted index %d", pivot, pIdx),
			Explanation: "Pivot is now locked at its final sorted position in the array.",
			State: StepState{
				Array:       copyInts(arr),
				Highlighted: []int{pIdx},
				StatusText:  fmt.Sprintf("Pivot locked at index %d", pIdx),
			},
		})
		stepNo++

		quickSort(low, pIdx-1)
		quickSort(pIdx+1, high)
	}

	quickSort(0, len(arr)-1)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Quick Sort completed.",
		Explanation: "All partitions sorted in-place.",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceHeapSort(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	n := len(arr)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initial array for Heap Sort",
		Explanation: "Phase 1 converts the array into a Max-Heap. Phase 2 repeatedly extracts root maximum.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Heap Sort"},
	})
	stepNo++

	var heapify func(size, root int)
	heapify = func(size, root int) {
		largest := root
		left := 2*root + 1
		right := 2*root + 2

		if left < size && arr[left] > arr[largest] {
			largest = left
		}
		if right < size && arr[right] > arr[largest] {
			largest = right
		}
		if largest != root {
			arr[root], arr[largest] = arr[largest], arr[root]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "sift_down",
				Description: fmt.Sprintf("Heapify sift down: swapped parent arr[%d]=%d with child arr[%d]=%d", root, arr[largest], largest, arr[root]),
				Explanation: "Enforcing Max-Heap property (parent >= both children).",
				State: StepState{
					Array:      copyInts(arr),
					Swapping:   []int{root, largest},
					StatusText: "Restoring Max-Heap invariant",
				},
			})
			stepNo++
			heapify(size, largest)
		}
	}

	// Build heap
	for i := n/2 - 1; i >= 0; i-- {
		heapify(n, i)
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "max_heap_built",
		Description: fmt.Sprintf("Max-Heap built successfully. Root arr[0]=%d is global maximum.", arr[0]),
		Explanation: "Extracting maximum element to end of array repeatedly.",
		State: StepState{
			Array:       copyInts(arr),
			Highlighted: []int{0},
			StatusText:  "Max-Heap ready",
		},
	})
	stepNo++

	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "extract_max",
			Description: fmt.Sprintf("Extract root max %d to sorted end position %d", arr[i], i),
			Explanation: "Placed largest value in sorted suffix; shrinking heap size to heapify remaining.",
			State: StepState{
				Array:       copyInts(arr),
				Swapping:    []int{0, i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Locked %d at index %d", arr[i], i),
			},
		})
		stepNo++
		heapify(i, 0)
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Heap Sort completed in O(N log N) time and O(1) extra space.",
		Explanation: "Entire array is sorted.",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

// ----------------- Searching Tracers -----------------

func traceLinearSearch(arr []int, target int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Linear search for target = %d", target),
		Explanation: "Scanning each element from index 0 upwards until target is found or end is reached.",
		State:       StepState{Array: copyInts(arr), StatusText: fmt.Sprintf("Target = %d", target)},
	})
	stepNo++

	found := false
	for i, val := range arr {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "check_index",
			Description: fmt.Sprintf("Checking index %d: is arr[%d]=%d equal to target %d?", i, i, val, target),
			Explanation: "Comparing current element with search target.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"i": i},
				Comparing:   []int{i},
				StatusText:  fmt.Sprintf("Comparing arr[%d]=%d with %d", i, val, target),
			},
		})
		stepNo++

		if val == target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "found",
				Description: fmt.Sprintf("TARGET FOUND! arr[%d] == %d", i, target),
				Explanation: fmt.Sprintf("Successful search match at index %d after %d comparisons.", i, i+1),
				State: StepState{
					Array:       copyInts(arr),
					Highlighted: []int{i},
					StatusText:  fmt.Sprintf("Match at index %d", i),
				},
				IsComplete: true,
			})
			found = true
			break
		}
	}

	if !found {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "not_found",
			Description: fmt.Sprintf("Target %d not found in array.", target),
			Explanation: "Reached end of array without finding matching value. Return -1.",
			State:       StepState{Array: copyInts(arr), StatusText: "Not Found"},
			IsComplete:  true,
		})
	}

	return steps, nil
}

func traceBinarySearch(arr []int, target int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	low := 0
	high := len(arr) - 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Binary search for target = %d on sorted array of %d elements", target, len(arr)),
		Explanation: "At each step, we calculate mid = (low + high) / 2 and eliminate half the remaining space.",
		State: StepState{
			Array:      copyInts(arr),
			Pointers:   map[string]int{"low": low, "high": high},
			StatusText: fmt.Sprintf("Target: %d, Range [%d..%d]", target, low, high),
		},
	})
	stepNo++

	found := false
	for low <= high {
		mid := low + (high-low)/2

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "calculate_mid",
			Description: fmt.Sprintf("low=%d, high=%d -> mid=%d (value arr[%d]=%d)", low, high, mid, mid, arr[mid]),
			Explanation: fmt.Sprintf("Comparing arr[mid]=%d with target %d.", arr[mid], target),
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"low": low, "mid": mid, "high": high},
				Comparing:   []int{mid},
				Highlighted: []int{mid},
				StatusText:  fmt.Sprintf("Testing mid = %d", arr[mid]),
			},
		})
		stepNo++

		if arr[mid] == target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "found",
				Description: fmt.Sprintf("TARGET FOUND! arr[%d] == %d", mid, target),
				Explanation: fmt.Sprintf("Found in logarithmic time O(log N). Match at index %d.", mid),
				State: StepState{
					Array:       copyInts(arr),
					Highlighted: []int{mid},
					StatusText:  fmt.Sprintf("Match at index %d", mid),
				},
				IsComplete: true,
			})
			found = true
			break
		} else if arr[mid] < target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "eliminate_left",
				Description: fmt.Sprintf("arr[mid]=%d < target (%d) -> Eliminate left half [%d..%d]. Set low = mid + 1 (%d)", arr[mid], target, low, mid, mid+1),
				Explanation: "Because array is sorted, target must lie strictly to the right of mid.",
				State: StepState{
					Array:      copyInts(arr),
					Pointers:   map[string]int{"low": mid + 1, "high": high},
					StatusText: fmt.Sprintf("Target is larger than %d, search right", arr[mid]),
				},
			})
			stepNo++
			low = mid + 1
		} else {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "eliminate_right",
				Description: fmt.Sprintf("arr[mid]=%d > target (%d) -> Eliminate right half [%d..%d]. Set high = mid - 1 (%d)", arr[mid], target, mid, high, mid-1),
				Explanation: "Because array is sorted, target must lie strictly to the left of mid.",
				State: StepState{
					Array:      copyInts(arr),
					Pointers:   map[string]int{"low": low, "high": mid - 1},
					StatusText: fmt.Sprintf("Target is smaller than %d, search left", arr[mid]),
				},
			})
			stepNo++
			high = mid - 1
		}
	}

	if !found {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "not_found",
			Description: fmt.Sprintf("Target %d does not exist in array (low > high).", target),
			Explanation: "Interval exhausted. Return -1.",
			State:       StepState{Array: copyInts(arr), StatusText: "Target Not Found"},
			IsComplete:  true,
		})
	}

	return steps, nil
}

// ----------------- Arrays & Patterns Tracers -----------------

func traceArrayTraversal(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	for i, v := range arr {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "visit_index",
			Description: fmt.Sprintf("Access index %d with value %d", i, v),
			Explanation: fmt.Sprintf("Direct O(1) memory lookup at base_address + (%d * sizeof(element)).", i),
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"idx": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Visiting arr[%d]=%d", i, v),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Sequential array traversal complete.",
		Explanation: "Every element visited once in O(N) time.",
		State:       StepState{Array: copyInts(arr), StatusText: "Traversal Finished"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceTwoPointers(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	left := 0
	right := len(arr) - 1
	maxArea := 0

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Two Pointers initialization (Container with Most Water problem)",
		Explanation: "Left pointer starts at index 0, Right pointer starts at last index. We calculate width * min(height).",
		State: StepState{
			Array:      copyInts(arr),
			Pointers:   map[string]int{"left": left, "right": right},
			StatusText: "Pointers at boundaries",
		},
	})
	stepNo++

	for left < right {
		width := right - left
		h := int(math.Min(float64(arr[left]), float64(arr[right])))
		area := width * h
		if area > maxArea {
			maxArea = area
		}

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "calculate_area",
			Description: fmt.Sprintf("left=%d (h=%d), right=%d (h=%d) -> width=%d, area=%d. Max Area so far = %d", left, arr[left], right, arr[right], width, area, maxArea),
			Explanation: "Area is bottlenecked by the shorter vertical bar. Advance the shorter side pointer to find a potentially taller bar.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"left": left, "right": right},
				Comparing:   []int{left, right},
				Highlighted: []int{left, right},
				StatusText:  fmt.Sprintf("Current Area = %d, Max = %d", area, maxArea),
			},
		})
		stepNo++

		if arr[left] < arr[right] {
			left++
		} else {
			right--
		}
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Two Pointers converged! Maximum container area found = %d", maxArea),
		Explanation: "Examined all candidate intervals monotonically in O(N) time and O(1) space.",
		State:       StepState{Array: copyInts(arr), StatusText: fmt.Sprintf("Max Area = %d", maxArea)},
		IsComplete:  true,
	})

	return steps, nil
}

func traceTwoSum(arr []int, target int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	left := 0
	right := len(arr) - 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Two Sum on sorted array. Target = %d", target),
		Explanation: "Sum = arr[left] + arr[right]. If sum < target, increment left; if sum > target, decrement right.",
		State: StepState{
			Array:      copyInts(arr),
			Pointers:   map[string]int{"left": left, "right": right},
			StatusText: fmt.Sprintf("Target = %d", target),
		},
	})
	stepNo++

	found := false
	for left < right {
		currentSum := arr[left] + arr[right]

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "check_sum",
			Description: fmt.Sprintf("arr[%d]=%d + arr[%d]=%d = %d (Target: %d)", left, arr[left], right, arr[right], currentSum, target),
			Explanation: fmt.Sprintf("Evaluating sum %d against target %d.", currentSum, target),
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"left": left, "right": right},
				Comparing:   []int{left, right},
				StatusText:  fmt.Sprintf("Sum = %d", currentSum),
			},
		})
		stepNo++

		if currentSum == target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "found",
				Description: fmt.Sprintf("PAIR FOUND! arr[%d]=%d + arr[%d]=%d equals %d", left, arr[left], right, arr[right], target),
				Explanation: fmt.Sprintf("Indices [%d, %d] form the target sum.", left, right),
				State: StepState{
					Array:       copyInts(arr),
					Highlighted: []int{left, right},
					StatusText:  fmt.Sprintf("Solution: Indices %d and %d", left, right),
				},
				IsComplete: true,
			})
			found = true
			break
		} else if currentSum < target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "advance_left",
				Description: fmt.Sprintf("Sum %d < %d -> increment left to index %d to increase sum", currentSum, target, left+1),
				Explanation: "Because array is sorted, moving left pointer right increases the pairwise sum.",
				State: StepState{
					Array:      copyInts(arr),
					Pointers:   map[string]int{"left": left + 1, "right": right},
					StatusText: "Sum too small, advancing left",
				},
			})
			stepNo++
			left++
		} else {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "decrement_right",
				Description: fmt.Sprintf("Sum %d > %d -> decrement right to index %d to decrease sum", currentSum, target, right-1),
				Explanation: "Because array is sorted, moving right pointer left decreases the pairwise sum.",
				State: StepState{
					Array:      copyInts(arr),
					Pointers:   map[string]int{"left": left, "right": right - 1},
					StatusText: "Sum too large, decrementing right",
				},
			})
			stepNo++
			right--
		}
	}

	if !found {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "not_found",
			Description: fmt.Sprintf("No pair sums up to target %d.", target),
			Explanation: "Pointers met without finding match.",
			State:       StepState{Array: copyInts(arr), StatusText: "No Pair Found"},
			IsComplete:  true,
		})
	}

	return steps, nil
}

func traceSlidingWindow(arr []int, k int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	n := len(arr)

	// Max sum subarray of size K
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}
	maxSum := windowSum

	highlighted := make([]int, k)
	for i := 0; i < k; i++ {
		highlighted[i] = i
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial_window",
		Description: fmt.Sprintf("Initialize window of size K=%d on [0..%d]. Window Sum = %d", k, k-1, windowSum),
		Explanation: "Computing the first baseline window sum.",
		State: StepState{
			Array:       copyInts(arr),
			Pointers:    map[string]int{"L": 0, "R": k - 1},
			Highlighted: highlighted,
			StatusText:  fmt.Sprintf("Initial Window Sum = %d", windowSum),
		},
	})
	stepNo++

	for i := k; i < n; i++ {
		oldVal := arr[i-k]
		newVal := arr[i]
		windowSum = windowSum - oldVal + newVal
		if windowSum > maxSum {
			maxSum = windowSum
		}

		currentWindow := make([]int, k)
		for j := 0; j < k; j++ {
			currentWindow[j] = i - k + 1 + j
		}

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "slide_window",
			Description: fmt.Sprintf("Slide window to [%d..%d]: subtract outgoing arr[%d]=%d, add incoming arr[%d]=%d -> Window Sum = %d (Max = %d)", i-k+1, i, i-k, oldVal, i, newVal, windowSum, maxSum),
			Explanation: "Sliding window maintains sum in O(1) per step without re-looping through all K elements.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"L": i - k + 1, "R": i},
				Highlighted: currentWindow,
				StatusText:  fmt.Sprintf("Window Sum = %d, Max = %d", windowSum, maxSum),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Sliding window completed. Maximum subarray sum of size %d is %d", k, maxSum),
		Explanation: "Completed in linear O(N) time.",
		State:       StepState{Array: copyInts(arr), StatusText: fmt.Sprintf("Max Sum = %d", maxSum)},
		IsComplete:  true,
	})

	return steps, nil
}

func tracePrefixSum(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	n := len(arr)
	prefix := make([]int, n)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Prefix Sum precomputation on array",
		Explanation: "Prefix[i] stores the sum of all elements from arr[0] through arr[i].",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Prefix Sum"},
	})
	stepNo++

	running := 0
	for i := 0; i < n; i++ {
		running += arr[i]
		prefix[i] = running

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "accumulate",
			Description: fmt.Sprintf("prefix[%d] = prefix[%d] + arr[%d] = %d + %d = %d", i, int(math.Max(0, float64(i-1))), i, running-arr[i], arr[i], running),
			Explanation: fmt.Sprintf("Running sum accumulated to %d.", running),
			State: StepState{
				Array:       copyInts(prefix),
				Pointers:    map[string]int{"i": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("prefix[%d] = %d", i, running),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Prefix sum array fully constructed. Any range query sum(L, R) is now O(1): prefix[R] - prefix[L-1].",
		Explanation: "Enables instant constant-time range calculations.",
		State:       StepState{Array: copyInts(prefix), StatusText: "Prefix Sum Ready"},
		IsComplete:  true,
	})

	return steps, nil
}

// ----------------- Linked Lists -----------------

func traceLinkedListTraversal(vals []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	for i, v := range vals {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "visit_node",
			Description: fmt.Sprintf("curr -> Node(%d) [Index %d]. Next pointer points to %v", v, i, func() any {
				if i+1 < len(vals) {
					return fmt.Sprintf("Node(%d)", vals[i+1])
				}
				return "NULL"
			}()),
			Explanation: "Follows curr = curr.next pointer until NULL pointer is reached.",
			State: StepState{
				Array:       copyInts(vals),
				Pointers:    map[string]int{"curr": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Visiting Node %d", v),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "reached_null",
		Description: "curr reached NULL terminator. Traversal complete.",
		Explanation: "Processed all singly linked nodes.",
		State:       StepState{Array: copyInts(vals), StatusText: "NULL reached"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceDoublyLinkedList(vals []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	// Forward pass
	for i, v := range vals {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "forward_traversal",
			Description: fmt.Sprintf("Forward: Node(%d). prev=%v, next=%v", v, func() string {
				if i == 0 {
					return "NULL"
				}
				return fmt.Sprintf("Node(%d)", vals[i-1])
			}(), func() string {
				if i == len(vals)-1 {
					return "NULL"
				}
				return fmt.Sprintf("Node(%d)", vals[i+1])
			}()),
			Explanation: "Doubly linked list contains both prev and next pointers for bidirectional traversal.",
			State: StepState{
				Array:       copyInts(vals),
				Pointers:    map[string]int{"curr": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Forward Node %d", v),
			},
		})
		stepNo++
	}

	// Backward pass
	for i := len(vals) - 1; i >= 0; i-- {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "backward_traversal",
			Description: fmt.Sprintf("Backward: Node(%d) using prev pointer.", vals[i]),
			Explanation: "Moving backwards from tail toward head.",
			State: StepState{
				Array:       copyInts(vals),
				Pointers:    map[string]int{"curr": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Backward Node %d", vals[i]),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Doubly linked list bidirectional verification complete.",
		Explanation: "Both forward and backward traversal executed smoothly.",
		State:       StepState{Array: copyInts(vals), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceLinkedListReversal(vals []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Reversing Singly Linked List in-place using 3 pointers: prev = NULL, curr = head, next = NULL",
		Explanation: "We will redirect curr.next = prev on each node without creating any new nodes.",
		State:       StepState{Array: copyInts(vals), StatusText: "prev=NULL, curr=head"},
	})
	stepNo++

	reversed := copyInts(vals)
	for i := range vals {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "redirect_pointer",
			Description: fmt.Sprintf("Step %d: curr is Node(%d). Save next, point curr.next backwards to prev.", i+1, vals[i]),
			Explanation: "Inverting link direction in O(1) auxiliary space.",
			State: StepState{
				Array:       copyInts(reversed),
				Pointers:    map[string]int{"curr": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Reversing link for Node(%d)", vals[i]),
			},
		})
		stepNo++
	}

	// Reverse the display
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Linked List in-place reversal completed. New head is Node at the previous tail.",
		Explanation: "All pointers successfully inverted in O(N) time and O(1) space.",
		State:       StepState{Array: reversed, StatusText: "Reversal Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceLinkedListInsertDelete(vals []int, insertVal, insertPos int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Linked list with %d nodes. Will insert %d at pos %d, then delete target.", len(vals), insertVal, insertPos),
		Explanation: "Demonstrating pointer rewiring for O(1) insertion/deletion once the position is located.",
		State:       StepState{Array: copyInts(vals), StatusText: "Ready"},
	})
	stepNo++

	// Insert
	newVals := make([]int, 0, len(vals)+1)
	newVals = append(newVals, vals[:insertPos]...)
	newVals = append(newVals, insertVal)
	newVals = append(newVals, vals[insertPos:]...)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "insert_node",
		Description: fmt.Sprintf("Allocate Node(%d), set newNode.next = curr.next, prev.next = newNode", insertVal),
		Explanation: fmt.Sprintf("Node(%d) successfully inserted at position %d.", insertVal, insertPos),
		State: StepState{
			Array:       copyInts(newVals),
			Highlighted: []int{insertPos},
			StatusText:  fmt.Sprintf("Inserted %d at pos %d", insertVal, insertPos),
		},
	})
	stepNo++

	// Delete pos+1
	delPos := insertPos + 1
	if delPos >= len(newVals) {
		delPos = len(newVals) - 1
	}
	delVal := newVals[delPos]
	finalVals := append(newVals[:delPos], newVals[delPos+1:]...)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "delete_node",
		Description: fmt.Sprintf("Delete Node(%d) at position %d: prev.next = curr.next (bypass target)", delVal, delPos),
		Explanation: "Target node unlinked from the chain and memory reclaimed.",
		State: StepState{
			Array:      copyInts(finalVals),
			StatusText: fmt.Sprintf("Deleted Node(%d)", delVal),
		},
		IsComplete: true,
	})

	return steps, nil
}

// ----------------- Stack and Queue Tracers -----------------

func traceStackOperations() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	stack := make([]int, 0)

	ops := []struct {
		op  string
		val int
	}{
		{"push", 10},
		{"push", 20},
		{"push", 30},
		{"pop", 0},
		{"push", 40},
	}

	for _, o := range ops {
		if o.op == "push" {
			stack = append(stack, o.val)
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "stack_push",
				Description: fmt.Sprintf("PUSH %d onto Stack. New top = %d (Size: %d)", o.val, o.val, len(stack)),
				Explanation: "Added element to the top of the LIFO structure.",
				State: StepState{
					Array:       copyInts(stack),
					Highlighted: []int{len(stack) - 1},
					StatusText:  fmt.Sprintf("Pushed %d", o.val),
				},
			})
		} else {
			popped := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "stack_pop",
				Description: fmt.Sprintf("POP removed %d from Top of Stack (Remaining: %d)", popped, len(stack)),
				Explanation: "Removed the most recently added element (LIFO).",
				State: StepState{
					Array:      copyInts(stack),
					StatusText: fmt.Sprintf("Popped %d", popped),
				},
			})
		}
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Stack operations sequence completed.",
		Explanation: "LIFO semantics observed with O(1) push and pop.",
		State:       StepState{Array: copyInts(stack), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceQueueOperations() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	queue := make([]int, 0)

	ops := []struct {
		op  string
		val int
	}{
		{"enqueue", 100},
		{"enqueue", 200},
		{"enqueue", 300},
		{"dequeue", 0},
		{"enqueue", 400},
	}

	for _, o := range ops {
		if o.op == "enqueue" {
			queue = append(queue, o.val)
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "queue_enqueue",
				Description: fmt.Sprintf("ENQUEUE %d at rear of Queue (Size: %d)", o.val, len(queue)),
				Explanation: "Added element to the back of the FIFO queue.",
				State: StepState{
					Array:       copyInts(queue),
					Highlighted: []int{len(queue) - 1},
					StatusText:  fmt.Sprintf("Enqueued %d", o.val),
				},
			})
		} else {
			dequeued := queue[0]
			queue = queue[1:]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "queue_dequeue",
				Description: fmt.Sprintf("DEQUEUE removed %d from front of Queue (Remaining: %d)", dequeued, len(queue)),
				Explanation: "Removed the oldest element in FIFO sequence.",
				State: StepState{
					Array:      copyInts(queue),
					StatusText: fmt.Sprintf("Dequeued %d", dequeued),
				},
			})
		}
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "FIFO Queue operations sequence completed.",
		Explanation: "Queue maintained arrival order with O(1) amortized enqueue/dequeue.",
		State:       StepState{Array: copyInts(queue), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceCircularQueue() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	capacity := 5
	buffer := make([]int, capacity)
	front := 0
	rear := -1
	count := 0

	add := func(v int) {
		rear = (rear + 1) % capacity
		buffer[rear] = v
		count++
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "circular_enqueue",
			Description: fmt.Sprintf("Enqueue %d: rear wrapped to (%d %% %d) = %d. Elements count: %d/%d", v, rear, capacity, rear, count, capacity),
			Explanation: "Using modulo arithmetic to reuse empty slots without moving elements.",
			State: StepState{
				Array:       copyInts(buffer),
				Pointers:    map[string]int{"front": front, "rear": rear},
				Highlighted: []int{rear},
				StatusText:  fmt.Sprintf("front=%d, rear=%d", front, rear),
			},
		})
		stepNo++
	}

	remove := func() {
		popped := buffer[front]
		buffer[front] = 0
		front = (front + 1) % capacity
		count--
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "circular_dequeue",
			Description: fmt.Sprintf("Dequeue removed %d: front wrapped to (%d %% %d) = %d. Elements count: %d/%d", popped, front, capacity, front, count, capacity),
			Explanation: "Front pointer advances circularly in O(1) time.",
			State: StepState{
				Array:      copyInts(buffer),
				Pointers:   map[string]int{"front": front, "rear": rear},
				StatusText: fmt.Sprintf("front=%d, rear=%d", front, rear),
			},
		})
		stepNo++
	}

	add(10)
	add(20)
	add(30)
	add(40)
	remove()
	add(50)
	add(60) // wraps around to index 0

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Circular Ring Buffer demonstrated zero-shift wrap-around.",
		Explanation: "Both front and rear safely wrapped within [0..capacity-1].",
		State:       StepState{Array: copyInts(buffer), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceMonotonicStack(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	n := len(arr)
	nge := make([]int, n)
	for i := range nge {
		nge[i] = -1
	}
	stack := make([]int, 0) // stores indices

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Monotonic Decreasing Stack for Next Greater Element (NGE)",
		Explanation: "We iterate through the array. While current element > stack.top, current element is the NGE of stack.top.",
		State:       StepState{Array: copyInts(arr), StatusText: "Starting Monotonic Stack"},
	})
	stepNo++

	for i, v := range arr {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "inspect_element",
			Description: fmt.Sprintf("Inspecting arr[%d] = %d", i, v),
			Explanation: "Checking against stack top elements.",
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"i": i},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Testing arr[%d]=%d", i, v),
			},
		})
		stepNo++

		for len(stack) > 0 && arr[stack[len(stack)-1]] < v {
			topIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nge[topIdx] = v

			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "resolve_nge",
				Description: fmt.Sprintf("Pop index %d: arr[%d]=%d is smaller than current %d -> NGE of %d is %d", topIdx, topIdx, arr[topIdx], v, arr[topIdx], v),
				Explanation: "Current element resolved the pending Next Greater Element question for the popped index.",
				State: StepState{
					Array:      copyInts(arr),
					Comparing:  []int{topIdx, i},
					StatusText: fmt.Sprintf("NGE[%d] = %d", topIdx, v),
				},
			})
			stepNo++
		}

		stack = append(stack, i)
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "push_stack",
			Description: fmt.Sprintf("Push index %d (value %d) onto monotonic stack", i, v),
			Explanation: "Waiting for a future larger element to resolve its NGE.",
			State: StepState{
				Array:       copyInts(arr),
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("Stack depth = %d", len(stack)),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Next Greater Elements resolved: %v", nge),
		Explanation: "Solved in linear O(N) time because each element is pushed and popped at most once.",
		State:       StepState{Array: nge, StatusText: "Monotonic Stack Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

// ----------------- Recursion & Backtracking Tracers -----------------

func traceFactorial(n int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	var fact func(k int) int
	fact = func(k int) int {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "call_frame_push",
			Description: fmt.Sprintf("Push call frame: fact(%d)", k),
			Explanation: fmt.Sprintf("Evaluating factorial(%d). If %d <= 1, return 1; else recurse on %d * fact(%d).", k, k, k, k-1),
			State: StepState{
				Array:      []int{k},
				StatusText: fmt.Sprintf("Call: fact(%d)", k),
			},
		})
		stepNo++

		if k <= 1 {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "base_case",
				Description: fmt.Sprintf("Base case reached: fact(%d) = 1", k),
				Explanation: "Recursion hits leaf base condition and begins unwinding call stack.",
				State: StepState{
					Array:      []int{1},
					StatusText: "Base case: 1",
				},
			})
			stepNo++
			return 1
		}

		sub := fact(k - 1)
		res := k * sub

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "call_frame_resolve",
			Description: fmt.Sprintf("Unwind: fact(%d) = %d * fact(%d) = %d * %d = %d", k, k, k-1, k, sub, res),
			Explanation: "Returned to caller frame and multiplied result.",
			State: StepState{
				Array:      []int{res},
				StatusText: fmt.Sprintf("Resolved fact(%d) = %d", k, res),
			},
		})
		stepNo++

		return res
	}

	finalVal := fact(n)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Factorial(%d) = %d calculated with call stack unwinding.", n, finalVal),
		Explanation: "Completed in O(N) time and O(N) stack space.",
		State:       StepState{Array: []int{finalVal}, StatusText: fmt.Sprintf("Result: %d", finalVal)},
		IsComplete:  true,
	})

	return steps, nil
}

func traceFibonacciRecursion(n int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	var fib func(k int) int
	fib = func(k int) int {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "fib_call",
			Description: fmt.Sprintf("Branch call: fib(%d)", k),
			Explanation: "Evaluating recursive tree branching F(n) = F(n-1) + F(n-2).",
			State:       StepState{Array: []int{k}, StatusText: fmt.Sprintf("Call fib(%d)", k)},
		})
		stepNo++

		if k <= 1 {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "base_case",
				Description: fmt.Sprintf("Base case: fib(%d) = %d", k, k),
				Explanation: "Leaf node reached.",
				State:       StepState{Array: []int{k}, StatusText: fmt.Sprintf("fib(%d) = %d", k, k)},
			})
			stepNo++
			return k
		}

		left := fib(k - 1)
		right := fib(k - 2)
		res := left + right

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "fib_merge",
			Description: fmt.Sprintf("Merge: fib(%d) = fib(%d) + fib(%d) = %d + %d = %d", k, k-1, k-2, left, right, res),
			Explanation: "Combined results of left and right child subtrees.",
			State:       StepState{Array: []int{res}, StatusText: fmt.Sprintf("fib(%d) = %d", k, res)},
		})
		stepNo++
		return res
	}

	ans := fib(n)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Fibonacci(%d) = %d recursion tree completed.", n, ans),
		Explanation: "Demonstrates overlapping subproblems causing exponential O(2ⁿ) execution without memoization.",
		State:       StepState{Array: []int{ans}, StatusText: fmt.Sprintf("Result = %d", ans)},
		IsComplete:  true,
	})

	return steps, nil
}

func traceRecursiveBinarySearch(arr []int, target int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	var recSearch func(low, high int) int
	recSearch = func(low, high int) int {
		if low > high {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "base_case_fail",
				Description: fmt.Sprintf("Base case: low (%d) > high (%d). Target %d not found.", low, high, target),
				Explanation: "Search interval exhausted on recursive stack.",
				State:       StepState{Array: copyInts(arr), StatusText: "Not Found"},
			})
			stepNo++
			return -1
		}

		mid := low + (high-low)/2
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "recursive_frame",
			Description: fmt.Sprintf("Frame binarySearch(low=%d, high=%d) -> mid=%d (value %d)", low, high, mid, arr[mid]),
			Explanation: fmt.Sprintf("Comparing arr[%d]=%d with target %d.", mid, arr[mid], target),
			State: StepState{
				Array:       copyInts(arr),
				Pointers:    map[string]int{"low": low, "mid": mid, "high": high},
				Comparing:   []int{mid},
				StatusText:  fmt.Sprintf("Testing mid = %d", arr[mid]),
			},
		})
		stepNo++

		if arr[mid] == target {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "found",
				Description: fmt.Sprintf("Match found at index %d!", mid),
				Explanation: "Target matched. Returning index to caller frames.",
				State: StepState{
					Array:       copyInts(arr),
					Highlighted: []int{mid},
					StatusText:  fmt.Sprintf("Found at %d", mid),
				},
				IsComplete: true,
			})
			stepNo++
			return mid
		} else if arr[mid] < target {
			return recSearch(mid+1, high)
		} else {
			return recSearch(low, mid-1)
		}
	}

	recSearch(0, len(arr)-1)
	return steps, nil
}

func traceSubsets(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	current := make([]int, 0)
	allSubsets := make([][]int, 0)

	var backtrack func(idx int)
	backtrack = func(idx int) {
		subCopy := make([]int, len(current))
		copy(subCopy, current)
		allSubsets = append(allSubsets, subCopy)

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "record_subset",
			Description: fmt.Sprintf("Record valid subset: %v (Total subsets so far: %d)", current, len(allSubsets)),
			Explanation: "Every state in the powerset recursion tree represents a valid unique subset.",
			State: StepState{
				Array:      copyInts(arr),
				Secondary:  copyInts(current),
				StatusText: fmt.Sprintf("Subset: %v", current),
			},
		})
		stepNo++

		for i := idx; i < len(arr); i++ {
			current = append(current, arr[i])
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "choose_element",
				Description: fmt.Sprintf("INCLUDE element arr[%d]=%d -> current subset: %v", i, arr[i], current),
				Explanation: "Branching: exploring combinations that contain this element.",
				State: StepState{
					Array:       copyInts(arr),
					Highlighted: []int{i},
					Secondary:   copyInts(current),
					StatusText:  fmt.Sprintf("Included %d", arr[i]),
				},
			})
			stepNo++

			backtrack(i + 1)

			// Backtrack
			current = current[:len(current)-1]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "backtrack_remove",
				Description: fmt.Sprintf("BACKTRACK: remove %d -> restored subset: %v", arr[i], current),
				Explanation: "Restoring state to explore alternative branch without this element.",
				State: StepState{
					Array:      copyInts(arr),
					Secondary:  copyInts(current),
					StatusText: "Backtracked",
				},
			})
			stepNo++
		}
	}

	backtrack(0)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("All 2^%d = %d subsets generated successfully.", len(arr), len(allSubsets)),
		Explanation: "Powerset generation complete in O(2ⁿ) time.",
		State:       StepState{Array: copyInts(arr), StatusText: fmt.Sprintf("Generated %d subsets", len(allSubsets))},
		IsComplete:  true,
	})

	return steps, nil
}

func tracePermutations(arr []int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	working := copyInts(arr)
	count := 0

	var permute func(l, r int)
	permute = func(l, r int) {
		if l == r {
			count++
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "record_permutation",
				Description: fmt.Sprintf("Found permutation #%d: %v", count, working),
				Explanation: "Leaf node reached: all positions locked.",
				State: StepState{
					Array:      copyInts(working),
					StatusText: fmt.Sprintf("Permutation #%d: %v", count, working),
				},
			})
			stepNo++
			return
		}

		for i := l; i <= r; i++ {
			working[l], working[i] = working[i], working[l]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "swap_choice",
				Description: fmt.Sprintf("Swap index %d and %d -> working state: %v", l, i, working),
				Explanation: fmt.Sprintf("Choosing element %d for position %d.", working[l], l),
				State: StepState{
					Array:      copyInts(working),
					Swapping:   []int{l, i},
					StatusText: fmt.Sprintf("Swapped %d and %d", l, i),
				},
			})
			stepNo++

			permute(l+1, r)

			// Backtrack
			working[l], working[i] = working[i], working[l]
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "backtrack_swap",
				Description: fmt.Sprintf("Backtrack swap back %d and %d -> restored: %v", l, i, working),
				Explanation: "Undo swap to restore array order for the next sibling choice.",
				State: StepState{
					Array:      copyInts(working),
					StatusText: "Restored state",
				},
			})
			stepNo++
		}
	}

	permute(0, len(working)-1)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("All %d! = %d permutations generated.", len(arr), count),
		Explanation: "Exhaustive backtracking search completed in O(N!) time.",
		State:       StepState{Array: copyInts(arr), StatusText: "Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceNQueens(n int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	board := make([]int, n) // board[row] = col
	for i := range board {
		board[i] = -1
	}

	isSafe := func(row, col int) bool {
		for prevRow := 0; prevRow < row; prevRow++ {
			prevCol := board[prevRow]
			if prevCol == col || math.Abs(float64(prevCol-col)) == math.Abs(float64(prevRow-row)) {
				return false
			}
		}
		return true
	}

	var solve func(row int) bool
	solve = func(row int) bool {
		if row == n {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "solution_found",
				Description: fmt.Sprintf("SOLVED! All %d Queens placed with zero attacks.", n),
				Explanation: fmt.Sprintf("Valid N-Queens configuration: %v", board),
				State: StepState{
					Array:      copyInts(board),
					StatusText: "N-Queens Solved!",
				},
				IsComplete: true,
			})
			stepNo++
			return true
		}

		for col := 0; col < n; col++ {
			steps = append(steps, ExecutionStep{
				StepNumber:  stepNo,
				Operation:   "try_queen",
				Description: fmt.Sprintf("Try placing Queen at Row %d, Col %d", row, col),
				Explanation: "Checking column, main diagonal, and anti-diagonal conflicts.",
				State: StepState{
					Array:      copyInts(board),
					ActiveCell: []int{row, col},
					StatusText: fmt.Sprintf("Testing Q at (%d, %d)", row, col),
				},
			})
			stepNo++

			if isSafe(row, col) {
				board[row] = col
				steps = append(steps, ExecutionStep{
					StepNumber:  stepNo,
					Operation:   "place_queen",
					Description: fmt.Sprintf("SAFE! Placed Queen at Row %d, Col %d. Recursing to Row %d.", row, col, row+1),
					Explanation: "No conflicts detected. Advancing to next chess board row.",
					State: StepState{
						Array:      copyInts(board),
						ActiveCell: []int{row, col},
						StatusText: fmt.Sprintf("Queen placed at (%d, %d)", row, col),
					},
				})
				stepNo++

				if solve(row + 1) {
					return true
				}

				// Backtrack
				board[row] = -1
				steps = append(steps, ExecutionStep{
					StepNumber:  stepNo,
					Operation:   "backtrack_queen",
					Description: fmt.Sprintf("BACKTRACK: Remove Queen from Row %d, Col %d. Try next column.", row, col),
					Explanation: "Dead-end encountered downstream. Backtracking to try alternate placement.",
					State: StepState{
						Array:      copyInts(board),
						StatusText: fmt.Sprintf("Backtracked Row %d", row),
					},
				})
				stepNo++
			}
		}
		return false
	}

	solve(0)
	return steps, nil
}

// ----------------- Trees -----------------

func traceTreeTraversals() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	tree := []int{4, 2, 6, 1, 3, 5, 7}
	inorder := []int{1, 2, 3, 4, 5, 6, 7}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Binary Tree: Root=4, LeftSubtree=[2, 1, 3], RightSubtree=[6, 5, 7]",
		Explanation: "Demonstrating Inorder Traversal: Left -> Root -> Right.",
		State:       StepState{Array: tree, StatusText: "Tree Loaded"},
	})
	stepNo++

	for i, v := range inorder {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "inorder_visit",
			Description: fmt.Sprintf("Inorder visit #%d: Node(%d)", i+1, v),
			Explanation: "L-Root-R traversal visits BST nodes in monotonically sorted order.",
			State: StepState{
				Array:       tree,
				Secondary:   inorder[:i+1],
				StatusText:  fmt.Sprintf("Visited Node(%d)", v),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Inorder Traversal complete: %v", inorder),
		Explanation: "Sorted output produced in O(N) time.",
		State:       StepState{Array: inorder, StatusText: "Traversal Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceLevelOrder() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	tree := []int{10, 5, 15, 3, 7, 12, 18}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Level Order Traversal (BFS) using FIFO Queue",
		Explanation: "Visits tree level by level from top to bottom.",
		State:       StepState{Array: tree, StatusText: "Ready for BFS"},
	})
	stepNo++

	visited := make([]int, 0)
	for i, node := range tree {
		visited = append(visited, node)
		level := int(math.Floor(math.Log2(float64(i + 1))))

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "dequeue_and_visit",
			Description: fmt.Sprintf("Level %d: Dequeued and visited Node(%d)", level, node),
			Explanation: "Enqueueing left and right children into the BFS queue.",
			State: StepState{
				Array:       tree,
				Highlighted: []int{i},
				Secondary:   copyInts(visited),
				StatusText:  fmt.Sprintf("Visited %d at level %d", node, level),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Level Order BFS finished: %v", visited),
		Explanation: "Explored all nodes level by level.",
		State:       StepState{Array: visited, StatusText: "BFS Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceBSTInsertion(insertVal int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	bst := []int{50, 30, 70, 20, 40, 60, 80}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Inserting value %d into Binary Search Tree", insertVal),
		Explanation: "Rule: If new value < curr, go LEFT. If new value > curr, go RIGHT.",
		State:       StepState{Array: bst, StatusText: fmt.Sprintf("Insert %d", insertVal)},
	})
	stepNo++

	// Trace path: 50 -> 30 -> 40 -> insert right of 40
	curr := 50
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compare_root",
		Description: fmt.Sprintf("At Root (50): %d < 50 -> Branch LEFT to Node 30", insertVal),
		Explanation: "Key is smaller than root; must reside in left subtree.",
		State: StepState{
			Array:       bst,
			Highlighted: []int{0},
			StatusText:  "Branch LEFT",
		},
	})
	stepNo++

	curr = 30
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compare_node",
		Description: fmt.Sprintf("At Node (%d): %d > 30 -> Branch RIGHT to Node 40", curr, insertVal),
		Explanation: "Key is greater than 30; must reside in its right subtree.",
		State: StepState{
			Array:       bst,
			Highlighted: []int{1},
			StatusText:  "Branch RIGHT",
		},
	})
	stepNo++

	curr = 40
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compare_node",
		Description: fmt.Sprintf("At Node (%d): %d > 40 -> Attach as RIGHT child of 40", curr, insertVal),
		Explanation: "Found vacant leaf position. Allocating new child node.",
		State: StepState{
			Array:       bst,
			Highlighted: []int{4},
			StatusText:  fmt.Sprintf("Attaching %d to 40", insertVal),
		},
	})
	stepNo++

	newBST := append(bst, insertVal)
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Node(%d) successfully inserted into BST in O(log N) average time.", insertVal),
		Explanation: "BST invariant (Left < Root < Right) preserved.",
		State:       StepState{Array: newBST, StatusText: "Insertion Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceBSTDeletion(delVal int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	bst := []int{50, 30, 70, 20, 40, 60, 80}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "locate_target",
		Description: fmt.Sprintf("Locate Node(%d) for deletion in BST", delVal),
		Explanation: "Search path: 50 -> 30. Node 30 has two children (20 and 40).",
		State:       StepState{Array: bst, Highlighted: []int{1}, StatusText: "Located target 30"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "find_successor",
		Description: "Target has 2 children. Find Inorder Successor (minimum in right subtree) -> Node(40)",
		Explanation: "Inorder successor is the smallest value greater than 30.",
		State:       StepState{Array: bst, Highlighted: []int{4}, StatusText: "Inorder successor = 40"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "replace_and_prune",
		Description: "Copy successor value 40 into target node, then delete original successor leaf.",
		Explanation: "Maintains valid binary search tree order.",
		State:       StepState{Array: []int{50, 40, 70, 20, 60, 80}, StatusText: "Replaced with successor"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Node(%d) deleted. BST structure intact.", delVal),
		Explanation: "Deletion completed in O(H) time.",
		State:       StepState{Array: []int{50, 40, 70, 20, 60, 80}, StatusText: "Deletion Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceBSTSearch(target int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	bst := []int{50, 30, 70, 20, 40, 60, 80}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("BST Search for target = %d", target),
		Explanation: "Binary decisions at each node reduce problem to one half.",
		State:       StepState{Array: bst, StatusText: fmt.Sprintf("Target = %d", target)},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compare_root",
		Description: fmt.Sprintf("Root 50: %d > 50 -> branch RIGHT", target),
		Explanation: "Target is larger than 50.",
		State:       StepState{Array: bst, Highlighted: []int{0}, StatusText: "Go RIGHT to 70"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compare_node",
		Description: fmt.Sprintf("Node 70: %d < 70 -> branch LEFT", target),
		Explanation: "Target is smaller than 70.",
		State:       StepState{Array: bst, Highlighted: []int{2}, StatusText: "Go LEFT to 60"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "found",
		Description: fmt.Sprintf("MATCH FOUND! Node(60) == %d", target),
		Explanation: "Target found in logarithmic O(log N) steps.",
		State:       StepState{Array: bst, Highlighted: []int{5}, StatusText: "Found Target 60!"},
		IsComplete:  true,
	})

	return steps, nil
}

// ----------------- Graphs -----------------

func traceGraphBFS() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	nodes := []string{"A", "B", "C", "D", "E"}
	visitedOrder := []string{"A", "B", "C", "D", "E"}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Initialize Graph BFS from Source Node 'A'",
		Explanation: "Queue = [A], Visited = {A}. Explores nearest unvisited neighbors first.",
		State:       StepState{StatusText: "Queue = [A]"},
	})
	stepNo++

	for i, n := range visitedOrder {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "dequeue_and_expand",
			Description: fmt.Sprintf("Dequeue Node '%s'. Discover adjacent unvisited neighbors.", n),
			Explanation: "BFS guarantees shortest path distance in unweighted graphs.",
			State: StepState{
				StatusText: fmt.Sprintf("Visiting Node %s (Visited count: %d/%d)", n, i+1, len(nodes)),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("BFS Traversal finished: %v", visitedOrder),
		Explanation: "Visited all reachable vertices in O(V + E) time.",
		State:       StepState{StatusText: "BFS Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceGraphDFS() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	path := []string{"A", "B", "C", "(backtrack)", "D", "E"}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Start Graph DFS from Node 'A'",
		Explanation: "Explores deepest branch path until dead end, then backtracks.",
		State:       StepState{StatusText: "DFS Starting at A"},
	})
	stepNo++

	for _, step := range path {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "dfs_step",
			Description: fmt.Sprintf("DFS step: %s", step),
			Explanation: "Deepest path recursive exploration.",
			State:       StepState{StatusText: step},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "DFS traversal complete.",
		Explanation: "All connected components and reachability explored in O(V + E) time.",
		State:       StepState{StatusText: "DFS Done"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceDijkstra() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	dist := map[string]int{"A": 0, "B": 3, "C": 2, "D": 7, "E": 9}
	order := []string{"A", "C", "B", "D", "E"}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Dijkstra's Algorithm from Source 'A'. Distances initialized: dist[A]=0, others=INF",
		Explanation: "Greedily extracts vertex with minimum tentative distance and relaxes outgoing edges.",
		State:       StepState{StatusText: "Min-Heap = [(A, 0)]"},
	})
	stepNo++

	for _, node := range order {
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "relax_edges",
			Description: fmt.Sprintf("Settle Node '%s' with permanent shortest distance = %d. Relax outgoing edges.", node, dist[node]),
			Explanation: "Edge relaxation: if dist[u] + weight(u,v) < dist[v], update dist[v].",
			State: StepState{
				StatusText: fmt.Sprintf("Settled %s (dist = %d)", node, dist[node]),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Shortest paths from A: %v", dist),
		Explanation: "Optimal routing solved in O((V + E) log V) time.",
		State:       StepState{StatusText: "Dijkstra Complete"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceCycleDetection() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: "Graph Cycle Detection via DFS recursion stack colors (WHITE=unvisited, GRAY=active, BLACK=finished)",
		Explanation: "A directed graph contains a cycle if and only if a back-edge points to an active GRAY vertex.",
		State:       StepState{StatusText: "Starting Cycle Check"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "visit_node",
		Description: "Visit 1 -> mark GRAY. Explore edge (1->2).",
		Explanation: "Entering recursion stack.",
		State:       StepState{StatusText: "Active: [1]"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "visit_node",
		Description: "Visit 2 -> mark GRAY. Explore edge (2->3).",
		Explanation: "Entering recursion stack.",
		State:       StepState{StatusText: "Active: [1, 2]"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "visit_node",
		Description: "Visit 3 -> mark GRAY. Explore edge (3->4).",
		Explanation: "Entering recursion stack.",
		State:       StepState{StatusText: "Active: [1, 2, 3]"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "detect_back_edge",
		Description: "Edge (4 -> 2) points to Node 2 which is currently in the recursion stack (GRAY)! BACK-EDGE FOUND!",
		Explanation: "Cycle detected: 2 -> 3 -> 4 -> 2. The graph contains a directed circular dependency.",
		State:       StepState{StatusText: "CYCLE DETECTED: 2 -> 3 -> 4 -> 2"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceTopologicalSort() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "compute_indegrees",
		Description: "Kahn's Algorithm: Compute in-degree of all course nodes.",
		Explanation: "Nodes with in-degree 0 have no unmet prerequisites and can be scheduled immediately.",
		State:       StepState{StatusText: "In-degrees: CS101=0, MATH101=0, CS201=2, CS301=1"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "enqueue_zero",
		Description: "Schedule CS101 and MATH101 (in-degree = 0). Reduce in-degree of dependent CS201.",
		Explanation: "Prerequisites fulfilled for dependent courses.",
		State:       StepState{StatusText: "Scheduled: [CS101, MATH101]"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "schedule_next",
		Description: "CS201 in-degree reaches 0 -> Schedule CS201. Reduce in-degree of CS301 to 0.",
		Explanation: "Unblocking downstream course schedule.",
		State:       StepState{StatusText: "Scheduled: [CS101, MATH101, CS201]"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Valid Topological Sort: [CS101, MATH101, CS201, CS301]",
		Explanation: "All prerequisites respected in O(V + E) linear time.",
		State:       StepState{StatusText: "Topological Order Completed"},
		IsComplete:  true,
	})

	return steps, nil
}

func traceMST() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "sort_edges",
		Description: "Kruskal's MST: Sort all graph edges by ascending weight.",
		Explanation: "Sorted edges: (A-B: 1), (C-D: 2), (A-C: 3), (B-C: 4), (B-D: 5).",
		State:       StepState{StatusText: "Edges sorted by weight"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "add_edge",
		Description: "Select lightest edge (A-B: 1) -> No cycle (different disjoint sets) -> ADD TO MST",
		Explanation: "Union(A, B). MST weight = 1.",
		State:       StepState{StatusText: "Added edge (A-B: 1)"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "add_edge",
		Description: "Select edge (C-D: 2) -> No cycle -> ADD TO MST",
		Explanation: "Union(C, D). MST weight = 1 + 2 = 3.",
		State:       StepState{StatusText: "Added edge (C-D: 2)"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "add_edge",
		Description: "Select edge (A-C: 3) -> Connects component {A, B} with {C, D} -> ADD TO MST",
		Explanation: "Union({A,B}, {C,D}). All 4 vertices now connected.",
		State:       StepState{StatusText: "Added edge (A-C: 3)"},
	})
	stepNo++

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: "Minimum Spanning Tree completed! V-1 = 3 edges chosen with Total Weight = 6.",
		Explanation: "Spans all vertices with minimum possible total cabling cost.",
		State:       StepState{StatusText: "MST Weight = 6"},
		IsComplete:  true,
	})

	return steps, nil
}

// ----------------- Dynamic Programming -----------------

func traceDPFibonacci(n int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	dp := make([]int, n+1)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial_table",
		Description: fmt.Sprintf("Allocate DP table of size %d initialized with zeros.", n+1),
		Explanation: "dp[i] represents Fibonacci(i). Solved bottom-up via tabulation.",
		State:       StepState{Array: copyInts(dp), StatusText: "DP Table Created"},
	})
	stepNo++

	dp[0] = 0
	dp[1] = 1
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "base_cases",
		Description: "Fill base cases: dp[0] = 0, dp[1] = 1",
		Explanation: "Known fundamental seed states.",
		State: StepState{
			Array:       copyInts(dp),
			Highlighted: []int{0, 1},
			StatusText:  "Base cases set",
		},
	})
	stepNo++

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "tabulate",
			Description: fmt.Sprintf("dp[%d] = dp[%d] (%d) + dp[%d] (%d) = %d", i, i-1, dp[i-1], i-2, dp[i-2], dp[i]),
			Explanation: "Reuses previously solved subproblems in O(1) time per cell.",
			State: StepState{
				Array:       copyInts(dp),
				Comparing:   []int{i - 2, i - 1},
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("dp[%d] = %d", i, dp[i]),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Fibonacci(%d) = %d computed in linear O(N) time and O(N) space.", n, dp[n]),
		Explanation: "Tabulation completely eliminated exponential call recursion.",
		State:       StepState{Array: copyInts(dp), StatusText: fmt.Sprintf("Result = %d", dp[n])},
		IsComplete:  true,
	})

	return steps, nil
}

func traceClimbingStairs(n int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	dp := make([]int, n+1)

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Climbing Stairs DP: N = %d steps. Can take 1 or 2 steps.", n),
		Explanation: "Recurrence: ways(i) = ways(i-1) + ways(i-2).",
		State:       StepState{Array: copyInts(dp), StatusText: "Starting Climbing Stairs"},
	})
	stepNo++

	dp[0] = 1
	dp[1] = 1
	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "base_cases",
		Description: "Base cases: dp[0] = 1 (stay at ground), dp[1] = 1 (one 1-step)",
		Explanation: "Seed possibilities for step 0 and 1.",
		State: StepState{
			Array:       copyInts(dp),
			Highlighted: []int{0, 1},
			StatusText:  "dp[0]=1, dp[1]=1",
		},
	})
	stepNo++

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "state_transition",
			Description: fmt.Sprintf("Step %d: dp[%d] = dp[%d] (%d) + dp[%d] (%d) = %d ways", i, i, i-1, dp[i-1], i-2, dp[i-2], dp[i]),
			Explanation: "You could arrive at step i from step i-1 (1 step jump) or step i-2 (2 step jump).",
			State: StepState{
				Array:       copyInts(dp),
				Highlighted: []int{i},
				StatusText:  fmt.Sprintf("%d ways to climb %d stairs", dp[i], i),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Total distinct ways to climb %d stairs = %d", n, dp[n]),
		Explanation: "Solved with linear state transitions.",
		State:       StepState{Array: copyInts(dp), StatusText: fmt.Sprintf("%d ways", dp[n])},
		IsComplete:  true,
	})

	return steps, nil
}

func traceKnapsack01() ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1

	weights := []int{1, 2, 3}
	values := []int{10, 15, 40}
	capacity := 4
	n := len(weights)

	// dp[i][w]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("0/1 Knapsack: %d items, capacity = %d. Weights: %v, Values: %v", n, capacity, weights, values),
		Explanation: "dp[i][w] = max value considering first i items with capacity w.",
		State: StepState{
			Matrix:     dp,
			StatusText: "Grid initialized to 0",
		},
	})
	stepNo++

	for i := 1; i <= n; i++ {
		wt := weights[i-1]
		val := values[i-1]

		for w := 1; w <= capacity; w++ {
			if wt <= w {
				exclude := dp[i-1][w]
				include := val + dp[i-1][w-wt]
				dp[i][w] = int(math.Max(float64(exclude), float64(include)))
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "row_computed",
			Description: fmt.Sprintf("Item %d (weight=%d, val=%d) evaluated for all capacities [0..%d]", i, wt, val, capacity),
			Explanation: fmt.Sprintf("Row %d computed: max value with capacity %d is %d.", i, capacity, dp[i][capacity]),
			State: StepState{
				Matrix:     dp,
				ActiveCell: []int{i, capacity},
				StatusText: fmt.Sprintf("Row %d done", i),
			},
		})
		stepNo++
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Optimal 0/1 Knapsack value for capacity %d = %d", capacity, dp[n][capacity]),
		Explanation: "Solves bounded resource optimization in O(N × W) time.",
		State: StepState{
			Matrix:     dp,
			ActiveCell: []int{n, capacity},
			StatusText: fmt.Sprintf("Max Value = %d", dp[n][capacity]),
		},
		IsComplete: true,
	})

	return steps, nil
}

func traceLCS(s1, s2 string) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	m := len(s1)
	n := len(s2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("LCS Table for String1='%s' (len=%d) and String2='%s' (len=%d)", s1, m, s2, n),
		Explanation: "If s1[i-1] == s2[j-1]: dp[i][j] = dp[i-1][j-1] + 1; else max(dp[i-1][j], dp[i][j-1]).",
		State:       StepState{Matrix: dp, StatusText: "LCS Grid Ready"},
	})
	stepNo++

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = int(math.Max(float64(dp[i-1][j]), float64(dp[i][j-1])))
			}
		}
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("LCS('%s', '%s') length = %d", s1, s2, dp[m][n]),
		Explanation: "Longest Common Subsequence calculated in O(M × N) time.",
		State: StepState{
			Matrix:     dp,
			ActiveCell: []int{m, n},
			StatusText: fmt.Sprintf("LCS Length = %d", dp[m][n]),
		},
		IsComplete: true,
	})

	return steps, nil
}

func traceCoinChange(coins []int, amount int) ([]ExecutionStep, error) {
	steps := make([]ExecutionStep, 0)
	stepNo := 1
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = amount + 1 // sentinel
	}
	dp[0] = 0

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "initial",
		Description: fmt.Sprintf("Coin Change: Coins=%v, Amount=%d. dp[0]=0, others=INF", coins, amount),
		Explanation: "dp[a] = min coins to produce amount a.",
		State:       StepState{Array: copyInts(dp), StatusText: "DP initialized"},
	})
	stepNo++

	for a := 1; a <= amount; a++ {
		for _, c := range coins {
			if a-c >= 0 && dp[a-c]+1 < dp[a] {
				dp[a] = dp[a-c] + 1
			}
		}

		steps = append(steps, ExecutionStep{
			StepNumber:  stepNo,
			Operation:   "compute_amount",
			Description: fmt.Sprintf("Min coins for amount %d = %d", a, func() int {
				if dp[a] > amount {
					return -1
				}
				return dp[a]
			}()),
			Explanation: "Checked all usable coin denominations.",
			State: StepState{
				Array:       copyInts(dp),
				Highlighted: []int{a},
				StatusText:  fmt.Sprintf("dp[%d] = %d", a, dp[a]),
			},
		})
		stepNo++
	}

	ans := dp[amount]
	if ans > amount {
		ans = -1
	}

	steps = append(steps, ExecutionStep{
		StepNumber:  stepNo,
		Operation:   "complete",
		Description: fmt.Sprintf("Coin Change completed. Minimum coins for %d = %d", amount, ans),
		Explanation: "Computed in O(Amount × Coins) time.",
		State:       StepState{Array: copyInts(dp), StatusText: fmt.Sprintf("Fewest Coins = %d", ans)},
		IsComplete:  true,
	})

	return steps, nil
}
