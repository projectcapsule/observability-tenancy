package e2e_test

import (
	"context"
	"encoding/json"

	"k8s.io/client-go/kubernetes"
)

// tenantGet queries a monitoring-system service through the API server service proxy, as the given tenant.
func tenantGet(service, path string, params map[string]string, tenant string, out any) error {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	req := cs.CoreV1().RESTClient().Get().
		AbsPath("/api/v1/namespaces/monitoring-system/services/"+service+"/proxy"+path).
		SetHeader("X-Scope-OrgID", tenant)
	for k, v := range params {
		req = req.Param(k, v)
	}

	raw, err := req.DoRaw(context.TODO())
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, out)
}

// Namespaces of the podinfo series (scraped by Prometheus, written through cortex-proxy) visible to a Mimir tenant.
func mimirNamespaces(tenant string) func() ([]string, error) {
	return func() ([]string, error) {
		var resp struct {
			Data struct {
				Result []struct {
					Metric map[string]string `json:"metric"`
				} `json:"result"`
			} `json:"data"`
		}
		if err := tenantGet("mimir:http", "/prometheus/api/v1/query",
			map[string]string{"query": "count by (namespace) (http_requests_total)"}, tenant, &resp); err != nil {
			return nil, err
		}

		namespaces := []string{}
		for _, r := range resp.Data.Result {
			namespaces = append(namespaces, r.Metric["namespace"])
		}

		return namespaces, nil
	}
}

// Namespaces of the log streams (collected by Alloy, written through loki-proxy) visible to a Loki tenant.
func lokiNamespaces(tenant string) func() ([]string, error) {
	return func() ([]string, error) {
		var resp struct {
			Data []string `json:"data"`
		}
		if err := tenantGet("loki:http-metrics", "/loki/api/v1/label/namespace/values",
			map[string]string{"since": "1h"}, tenant, &resp); err != nil {
			return nil, err
		}

		return append([]string{}, resp.Data...), nil
	}
}
