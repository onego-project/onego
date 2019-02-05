package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var (
	vnLeaseHold        = "records/virtualNetwork/lease/hold"
	vnLeaseHoldWrongVN = "records/virtualNetwork/lease/holdWrongVN"

	vnLeaseRelease        = "records/virtualNetwork/lease/release"
	vnLeaseReleaseWrongIP = "records/virtualNetwork/lease/releaseWrongIP"
)

var _ = ginkgo.Describe("Lease Service", func() {
	var (
		recName string
		rec     *recorder.Recorder
		client  *onego.Client
		err     error
	)

	ginkgo.JustBeforeEach(func() {
		// Start recorder
		rec, err = recorder.New(recName)
		if err != nil {
			return
		}

		rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			var b bytes.Buffer
			if _, err = b.ReadFrom(r.Body); err != nil {
				return false
			}
			r.Body = ioutil.NopCloser(&b)
			return cassette.DefaultMatcher(r, i) && (b.String() == "" || b.String() == i.Body)
		})

		// Create an HTTP client and inject our transport
		clientHTTP := &http.Client{
			Transport: rec, // Inject as transport!
		}

		// create onego client
		client = onego.CreateClient(endpoint, token, clientHTTP)
		if client == nil {
			err = errors.ErrNoClient
			return
		}
	})

	ginkgo.AfterEach(func() {
		rec.Stop()
	})

	ginkgo.Describe("Manage lease in virtual network", func() {
		var vn *resources.VirtualNetwork
		var oneVn *resources.VirtualNetwork

		ginkgo.Context("when hold lease", func() {
			ginkgo.BeforeEach(func() {
				recName = vnLeaseHold
			})

			ginkgo.It("should hold the lease in virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				ip := net.ParseIP("10.0.0.11")
				ip2 := net.ParseIP("10.0.0.12")
				ips := make([]net.IP, 2)
				ips[0] = ip
				ips[1] = ip2

				err = client.LeaseService.Hold(context.TODO(), *vn, ips)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether lease is really hold in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), vnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				containIP := false
				for _, ar := range addressRanges {
					for _, lease := range ar.Leases {
						if lease.IP.Equal(ip) {
							containIP = true
							break
						}
					}
				}
				gomega.Expect(containIP).Should(gomega.BeTrue())
			})
		})

		ginkgo.Context("when vnet has wrong ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnLeaseHoldWrongVN
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 55
				vn = resources.CreateVirtualNetworkWithID(vnID)

				err = client.LeaseService.Hold(context.TODO(), *vn, make([]net.IP, 2))
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vnet has no ID", func() {
			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.LeaseService.Hold(context.TODO(), resources.VirtualNetwork{}, make([]net.IP, 2))
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when IP array is empty", func() {
			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				err = client.LeaseService.Hold(context.TODO(), *vn, make([]net.IP, 2))
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when release a lease from vnet", func() {
			ginkgo.BeforeEach(func() {
				recName = vnLeaseRelease
			})

			ginkgo.It("should release a lease in virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				ip := net.ParseIP("10.0.0.12")
				ips := make([]net.IP, 1)
				ips[0] = ip

				err = client.LeaseService.Release(context.TODO(), *vn, ips)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether lease was really released in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), vnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				containIP := false
				for _, ar := range addressRanges {
					for _, lease := range ar.Leases {
						if lease.IP.Equal(ip) {
							containIP = true
							break
						}
					}
				}
				gomega.Expect(containIP).Should(gomega.BeFalse())
			})
		})

		ginkgo.Context("when IP is not hold", func() {
			ginkgo.BeforeEach(func() {
				recName = vnLeaseReleaseWrongIP
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				ip := net.ParseIP("10.0.0.15")
				ips := make([]net.IP, 1)
				ips[0] = ip

				err = client.LeaseService.Hold(context.TODO(), *vn, ips)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})
})
