package gocouchbase

import "time"
import "fmt"
import "github.com/couchbase/gocouchbaseio"

type Cluster struct {
	manager *ClusterManager

	spec              connSpec
	connectionTimeout time.Duration
}

func Connect(connSpecStr string) (*Cluster, error) {
	spec := parseConnSpec(connSpecStr)
	if spec.Scheme == "" {
		spec.Scheme = "http"
	}
	if spec.Scheme != "couchbase" && spec.Scheme != "couchbases" && spec.Scheme != "http" {
		panic("Unsupported Scheme!")
	}
	cluster := &Cluster{
		spec:              spec,
		connectionTimeout: 10000 * time.Millisecond,
	}
	return cluster, nil
}

func (c *Cluster) OpenBucket(bucket, password string) (*Bucket, error) {
	var memdHosts []string
	var httpHosts []string
	isHttpHosts := c.spec.Scheme == "http"
	isSslHosts := c.spec.Scheme == "couchbases"
	for _, specHost := range c.spec.Hosts {
		if specHost.Port == 0 {
			if !isHttpHosts {
				if !isSslHosts {
					specHost.Port = 11210
				} else {
					specHost.Port = 11207
				}
			} else {
				panic("HTTP configuration not yet supported")
				//specHost.Port = 8091
			}
		}
		memdHosts = append(memdHosts, fmt.Sprintf("%s:%d", specHost.Host, specHost.Port))
	}

	authFn := func(srv *gocouchbaseio.MemdServer) error {
		fmt.Printf("Want to auth for %s\n", srv.Address())
		return nil
	}
	cli, err := gocouchbaseio.CreateAgent(memdHosts, httpHosts, isSslHosts, authFn)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		client: cli,
	}, nil
}

func (c *Cluster) Manager(username, password string) *ClusterManager {
	if c.manager == nil {
		c.manager = &ClusterManager{}
	}
	return c.manager
}