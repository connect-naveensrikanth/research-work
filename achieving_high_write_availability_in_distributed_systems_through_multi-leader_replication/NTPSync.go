package main

import (
	"fmt"
	"log"
	"time"
	"github.com/beevik/ntp"
)

func fetchNTPTime(server string) (time.Time, error) {
	ntpTime, err := ntp.Time(server)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch time from NTP server %s: %v", server, err)
	}
	return ntpTime, nil
}
func compareTimes(ntpTime, localTime time.Time) {
	offset := ntpTime.Sub(localTime)
	fmt.Printf("Time Offset (NTP - Local): %v\n", offset)
	syncedTime := localTime.Add(offset)
	fmt.Printf("Adjusted Time (Local + Offset): %v\n", syncedTime)
	if offset > 0 {
		fmt.Println("Local system is behind the NTP server.")
	} else if offset < 0 {
		fmt.Println("Local system is ahead of the NTP server.")
	} else {
		fmt.Println("Local system time is in sync with the NTP server.")
	}
}
func printFormattedTimes(ntpTime, localTime time.Time) {
	currentTimeStr := ntpTime.Format(time.RFC3339)
	fmt.Println("NTP Time (RFC3339):", currentTimeStr)
	localTimeStr := localTime.Format(time.RFC3339)
	fmt.Println("Local Time (RFC3339):", localTimeStr)
}
func printUnixTimestamps(ntpTime, localTime time.Time) {
	fmt.Printf("Current timestamp from NTP server: %d\n", ntpTime.Unix())
	fmt.Printf("Current timestamp from local system: %d\n", localTime.Unix())
}
func displayServerTimeDetails(server string) {
	ntpTime, err := fetchNTPTime(server)
	if err != nil {
		log.Fatal(err)
	}
	localTime := time.Now()
	fmt.Printf("NTP Time from %s: %v\n", server, ntpTime)
	fmt.Printf("Local Time: %v\n", localTime)
	compareTimes(ntpTime, localTime)
	printFormattedTimes(ntpTime, localTime)
	printUnixTimestamps(ntpTime, localTime)
}
func main() {
	server := "time.google.com"
	displayServerTimeDetails(server)
	server2 := "pool.ntp.org"
	displayServerTimeDetails(server2)
	ntpTime, _ := fetchNTPTime(server)
	localTime := time.Now()
	timestampDiff := ntpTime.Sub(localTime)
	fmt.Printf("Timestamp difference between NTP and Local time: %v\n", timestampDiff)
}
