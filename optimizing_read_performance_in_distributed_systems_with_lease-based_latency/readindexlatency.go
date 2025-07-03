import (
"fmt"
"math/rand"
"time"
)
func simulateReadIndexLatency(clusterSize int) float64 {
baseLatency := 1.0
latency := baseLatency + float64(clusterSize)
latency += (rand.Float64() - 0.5)
if latency < 0 {
latency = 0
}
return latency
}
func simulateReadOperation(clusterSize int) float64 {
return simulateReadIndexLatency(clusterSize)
}
func collectMetrics(clusterSize int, numReads int) (float64, int) {
totalLatency := 0.0
successfulReads := 0
for i := 0; i < numReads; i++ {
latency := simulateReadOperation(clusterSize)
totalLatency += latency
successfulReads++
}
return totalLatency, successfulReads
}
func displayMetrics(clusterSizes []int, numReads int) {
fmt.Println("Cluster Size (Nodes) | Successful Reads | Average ReadIndex Latency (ms)")
for _, size := range clusterSizes {
totalLatency, successfulReads := collectMetrics(size, numReads)
averageLatency := totalLatency / float64(successfulReads)
fmt.Printf("%d\t\t\t%d\t\t\t%.2f\n", size, successfulReads, averageLatency)
time.Sleep(200 * time.Millisecond)
}
}
func main() {
rand.Seed(time.Now().UnixNano())
clusterSizes := []int{3, 5, 7, 9, 11}
numReads := 100
displayMetrics(clusterSizes, numReads)
}