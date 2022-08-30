package main

import "os"

func writeDummyTargets(filename string, count int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(`{"status": "success", "data": { "activeTargets": [`); err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		if _, err := f.WriteString(`
{
	"discoveredLabels": {
		"__address__": "123m%[1]d",
		"__meta_filepath": "/etc/prometheus/targets/targets.json",
		"__metrics_path__": "/probe",
		"__param_module": "http_200",
		"__scheme__": "http",
		"__scrape_interval__": "60s",
		"__scrape_timeout__": "20s",
		"location_id": "123",
		"location_name": "pop123",
		"instance": "123m%[1]d",
		"job": "mock-job-name",
		"node_status": "v",
		"node_type": "server",
		"region": "global"
	},
	"labels": {
		"location_id": "123",
		"location_name": "pop123",
		"instance": "123m%[1]d",
		"job": "mock-job-name",
		"node_status": "v",
		"node_type": "server",
		"region": "global",
		"target": "https://scrape-target:9090/probe"
	},
	"scrapePool": "mock-job-name",
	"scrapeUrl": "http://123m%[1]d:9090/probe?module=http_200&target=https%3A%2F%2Fscrape-target%3A9090%2Fprobe",
	"globalUrl": "http://123m%[1]d:9090/probe?module=http_200&target=https%3A%2F%2Fscrape-target%3A9090%2Fprobe",
	"lastError": "",
	"lastScrape": "2022-08-23T12:39:41.121731854Z",
	"lastScrapeDuration": 0.003068728,
	"health": "up",
	"scrapeInterval": "60s",
	"scrapeTimeout": "20s"
}`); err != nil {
			return err
		}
		if i < count-1 {
			if _, err := f.WriteString(","); err != nil {
				return err
			}
		}
	}

	if _, err := f.WriteString(`]}}`); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := writeDummyTargets("targets.json", 200000); err != nil {
		panic(err)
	}
}
