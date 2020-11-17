package main

var domainInt int
var domainIntStop int

func main() {
	args := getArgparser("Test")
	mode := args.get("", "mode", "lock", "can be lock, chan, goroutine")
	threads := args.getInt("", "threads", "20", "threads count")
	testCount := args.getInt("", "testCount", "2000", "测试的数量")
	toCheckDomain := args.getBool("", "toCheckDomain", "true", "是否依赖外部资源")
	args.parseArgs()

	domainInt = 4874088
	domainIntStop = domainInt + testCount

	startTime := now()
	go func() {
		for {
			if domainInt >= domainIntStop {
				print(fmtTimeDuration(now() - startTime))
				exit(0)
			}
			sleep(0.01)
		}
	}()

	// mode lock
	if mode == "lock" {
		lock := getLock()
		for range rangeInt(threads) {
			go func() {
				for {
					lock.acquire()

					domain := numToBHex(domainInt, 36) + ".ws"
					domainInt++

					lock.release()

					if toCheckDomain {
						try(func() {
							gethostbyname(domain)
						})
					}

				}
			}()
		}
		select {}
		// mode chan
	} else if mode == "chan" {
		c := make(chan string, 100)

		go func() {
			for {
				domain := numToBHex(domainInt, 36) + ".ws"
				domainInt++
				c <- domain
			}
		}()

		for range rangeInt(threads) {
			go func() {
				for {
					if toCheckDomain {
						try(func() {
							gethostbyname(<-c)
						})
					} else {
						<-c
					}
				}
			}()
		}
		select {}
		// mode goroutine
	} else if mode == "goroutine" {
		runningThread := 0
		for {
			for runningThread < threads {
				domain := numToBHex(domainInt, 36)
				domainInt++
				if toCheckDomain {
					go func(domain string) {
						try(func() {
							gethostbyname(domain)
						})
					}(domain)
				}
				sleep(0.01)
			}
			sleep(0.01)
		}
	}
}
