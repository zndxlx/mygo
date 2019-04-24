package util

/*For edge case weights like { 5, 1, 1 } we now produce { a, a, b, a, c, a, a }
sequence instead of { c, b, a, a, a, a, a } produced previously.

Algorithm is as follows: on each peer selection we increase current_weight
of each eligible peer by its weight, select peer with greatest current_weight
and reduce its current_weight by total number of weight points distributed
among peers.

In case of { 5, 1, 1 } weights this gives the following sequence of
current_weight's:

a  b  c
0  0  0  (initial state)

5  1  1  (a selected)
-2  1  1

3  2  2  (a selected)
-4  2  2

1  3  3  (b selected)
1 -4  3

6 -3  4  (a selected)
-1 -3  4

4 -2  5  (c selected)
4 -2 -2

9 -1 -1  (a selected)
2 -1 -1

7  0  0  (a selected)
0  0  0*/

type smoothWeightEntry struct {
    item          interface{}
    weight        int
    currentWeight int
}

type SmoothRoundRobinAlg struct {
    items []*smoothWeightEntry
    n     int
}

func (sRRA *SmoothRoundRobinAlg) Add(item interface{}, weight int) {
    weightEntryItem := &smoothWeightEntry{item: item, weight: weight, currentWeight: 0}
    sRRA.items = append(sRRA.items, weightEntryItem)
    sRRA.n++
}

func (sRRA *SmoothRoundRobinAlg) Next() (item interface{}) {

    weightEntryItem := smoothRoundRobinAlg(sRRA.items)
    if weightEntryItem == nil {
        item = nil
    } else {
        item = weightEntryItem.item
    }
    return
}

func (sRRA *SmoothRoundRobinAlg) Remove(item interface{}) {
    for i := 0; i < sRRA.n; i++ {
        if sRRA.items[i].item == item {
            sRRA.items = append(sRRA.items[:i], sRRA.items[i+1:]...)
            sRRA.n--
            break
        }
    }
}

func smoothRoundRobinAlg(items []*smoothWeightEntry) (bestItem *smoothWeightEntry) {
    if items == nil || len(items) == 0 {
        return nil
    }

    total := 0
    for _, item := range items {
        if item == nil {
            continue
        }

        item.currentWeight += item.weight
        total += item.currentWeight

        if bestItem == nil || item.currentWeight > bestItem.currentWeight {
            bestItem = item
        }
    }

    bestItem.currentWeight -= total
    return
}
