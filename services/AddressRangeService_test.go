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
	vnArAdd        = "records/virtualNetwork/addressRange/addToExisting"
	vnArUpdate     = "records/virtualNetwork/addressRange/update"
	vnArAddWrong   = "records/virtualNetwork/addressRange/addWrong"
	vnArAddNo      = "records/virtualNetwork/addressRange/addNo"
	vnArAddWrongAR = "records/virtualNetwork/addressRange/addWrongAR"

	vnArDelete        = "records/virtualNetwork/addressRange/delete"
	vnArDeleteWrong   = "records/virtualNetwork/addressRange/deleteWrong"
	vnArDeleteNo      = "records/virtualNetwork/addressRange/deleteNo"
	vnArDeleteWrongAR = "records/virtualNetwork/addressRange/deleteWrongAR"
)

var _ = ginkgo.Describe("Address Range Service", func() {
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

	ginkgo.Describe("Manage address range in virtual network", func() {
		var vn *resources.VirtualNetwork
		var ar *resources.AddressRange
		var retAR *resources.AddressRange
		var oneVn *resources.VirtualNetwork

		ginkgo.Context("when add address range", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArAdd
			})

			ginkgo.It("should create new address range in virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				size := 5
				ip := net.ParseIP("10.0.0.10")
				ar = &resources.AddressRange{
					Type: "IP4",
					IP:   &ip,
					Size: &size,
				}

				ar, err = client.AddressRangeService.Add(context.TODO(), *vn, *ar)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ar).NotTo(gomega.BeNil())

				// check whether reservation really exists in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), vnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				containIP := false
				for _, ar := range addressRanges {
					if ar.IP.Equal(ip) {
						containIP = true
						break
					}
				}
				gomega.Expect(containIP).Should(gomega.BeTrue())
			})
		})

		ginkgo.Context("when update address range", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArUpdate
			})

			ginkgo.It("should update address range in virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				size := 3
				id := 7
				ar = &resources.AddressRange{
					ID:   &id,
					Size: &size,
				}

				retAR, err = client.AddressRangeService.Update(context.TODO(), *vn, *ar)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(retAR).ShouldNot(gomega.BeNil())

				// check whether reservation really exists in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), vnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				contain := false
				for _, ar := range addressRanges {
					if *ar.ID == id {
						if *ar.Size == size {
							contain = true
							break
						}
					}
				}
				gomega.Expect(contain).Should(gomega.BeTrue())
			})
		})

		ginkgo.Context("when vnet has wrong ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArAddWrong
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 55
				vn = resources.CreateVirtualNetworkWithID(vnID)

				size := 5
				ip := net.ParseIP("10.0.0.10")
				ar = &resources.AddressRange{
					Type: "IP4",
					IP:   &ip,
					Size: &size,
				}

				ar, err = client.AddressRangeService.Add(context.TODO(), *vn, *ar)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(ar).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when vnet has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArAddNo
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				size := 5
				ip := net.ParseIP("10.0.0.10")
				ar = &resources.AddressRange{
					Type: "IP4",
					IP:   &ip,
					Size: &size,
				}

				ar, err = client.AddressRangeService.Add(context.TODO(), resources.VirtualNetwork{}, *ar)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(ar).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when missing element (type) in Address Range", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArAddWrongAR
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				size := 5
				ip := net.ParseIP("10.0.0.10")
				ar = &resources.AddressRange{
					IP:   &ip,
					Size: &size,
				}

				ar, err = client.AddressRangeService.Add(context.TODO(), *vn, *ar)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(ar).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when delete an address range from vnet", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArDelete
			})

			ginkgo.It("should delete the given address range from virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				arID := 2
				ar = &resources.AddressRange{ID: &arID}

				err = client.AddressRangeService.Delete(context.TODO(), *vn, *ar)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether reservation really exists in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), vnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				containIP := false
				for _, ar := range addressRanges {
					if *ar.ID == arID {
						containIP = true
						break
					}
				}
				gomega.Expect(containIP).Should(gomega.BeFalse())
			})
		})

		ginkgo.Context("delete when vnet has wrong ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArDeleteWrong
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 45
				vn = resources.CreateVirtualNetworkWithID(vnID)

				arID := 3
				ar = &resources.AddressRange{ID: &arID}

				err = client.AddressRangeService.Delete(context.TODO(), *vn, *ar)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("delete when vnet has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArDeleteNo
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				arID := 3
				ar = &resources.AddressRange{ID: &arID}

				err = client.AddressRangeService.Delete(context.TODO(), resources.VirtualNetwork{}, *ar)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("delete when AR has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = vnArDeleteWrongAR
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vnID := 3
				vn = resources.CreateVirtualNetworkWithID(vnID)

				err = client.AddressRangeService.Delete(context.TODO(), *vn, resources.AddressRange{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

	})
})
