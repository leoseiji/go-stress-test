package cmd

import (
	"log"
	"net/http"
	"sync"
)

func RunStressTest(url string, requests int, concurrency int) {
	log.Printf("Run stress test started: URL[%s], requests[%d], concurrency[%d]", url, requests, concurrency)

	report := &Report{}
	report.StartExecution()

	wg := &sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Try to acquire a token.

		go func() {
			callEndpoint(url, wg, report)
			<-semaphore // Release the token after finishing.
		}()
	}
	wg.Wait()
	close(semaphore) // Closing the semaphore channel is not necessary but can be done for cleanup.

	report.EndExecution()
	report.Show()

	log.Println("Run stress test finished")
}

func callEndpoint(url string, wg *sync.WaitGroup, report *Report) {
	defer wg.Done()

	report.TotalRequests.Add(1)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Printf("error creating request. err:%s \n", err.Error())
		report.TotalUndefinedError.Add(1)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error calling endpoint. err:%s \n", err.Error())
		if resp != nil {
			report.IncrementStatusCount(resp.StatusCode)
		} else {
			report.TotalUndefinedError.Add(1)
		}
		return
	}
	defer resp.Body.Close()

	report.IncrementStatusCount(resp.StatusCode)
}
