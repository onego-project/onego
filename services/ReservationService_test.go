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
	vnReservationAddToExisting = "records/virtualNetwork/reservation/addToExisting"
	vnReservationAddNew        = "records/virtualNetwork/reservation/addNew"
	vnReserveToExisting        = "records/virtualNetwork/reservation/toExisting"
	vnReserveAsNew             = "records/virtualNetwork/reservation/asNew"
	vnReserveWrongID           = "records/virtualNetwork/reservation/wrongID"
	vnReserveEmptyVN           = "records/virtualNetwork/reservation/emptyVN"
	vnReserveEmptyReservation  = "records/virtualNetwork/reservation/emptyReservation"
	vnReservationFree          = "records/virtualNetwork/reservation/free"
)

var _ = ginkgo.Describe("Reservation Service", func() {
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

	ginkgo.Describe("reserve network addresses", func() {
		var vn *resources.VirtualNetwork
		var reservation *resources.Reservation
		var resVn *resources.VirtualNetwork
		var resVnID int
		var oneVn *resources.VirtualNetwork

		ginkgo.Context("when add to existing reservation", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReservationAddToExisting
			})

			ginkgo.It("should create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(3)

				size := 1
				arID := 0
				ip := net.ParseIP("10.0.0.14")
				nID := 7
				reservation = &resources.Reservation{Size: &size, AddressRangeID: &arID, IP: &ip, VirtualNetworkID: &nID}

				resVn, err = client.ReservationService.Reserve(context.TODO(), *vn, *reservation)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resVn).NotTo(gomega.BeNil())

				// check whether reservation really exists in OpenNebula
				resVnID, err = resVn.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), resVnID)
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

		ginkgo.Context("when reserve as new", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReservationAddNew
			})

			ginkgo.It("should create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(3)

				size := 1
				reservation = &resources.Reservation{Size: &size, Name: "newReservation"}

				resVn, err = client.ReservationService.Reserve(context.TODO(), *vn, *reservation)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether reservation really exists in OpenNebula
				resVnID, err = resVn.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), resVnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVn.ID()).Should(gomega.Equal(resVnID))
			})
		})

		ginkgo.Context("when reserve to the existing reservation", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReserveToExisting
			})

			ginkgo.It("should add reservation to existing virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(3)

				resVn, err = client.ReservationService.ReserveToExisting(context.TODO(), *vn, 1, 1, 7)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether reservation really exists in OpenNebula
				resVnID, err = resVn.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), resVnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(addressRanges).ShouldNot(gomega.BeEmpty())

				ip := net.ParseIP("10.0.0.10")
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

		ginkgo.Context("when reserve as new virtual network", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReserveAsNew
			})

			ginkgo.It("should create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(3)

				resVn, err = client.ReservationService.ReserveAsNew(context.TODO(), *vn, 1, "newVN")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether reservation really exists in OpenNebula
				resVnID, err = resVn.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), resVnID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVn.ID()).Should(gomega.Equal(resVnID))
			})
		})

		ginkgo.Context("when v. network does not exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReserveWrongID
			})

			ginkgo.It("should not create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(553)

				resVn, err = client.ReservationService.ReserveAsNew(context.TODO(), *vn, 1, "newVN")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when v. network is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReserveEmptyVN
			})

			ginkgo.It("should not create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				resVn, err = client.ReservationService.ReserveAsNew(context.TODO(), resources.VirtualNetwork{}, 1, "newVN")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when reservation is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReserveEmptyReservation
			})

			ginkgo.It("should not create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(3)

				resVn, err = client.ReservationService.Reserve(context.TODO(), *vn, resources.Reservation{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("free network address", func() {
		var vn *resources.VirtualNetwork
		var oneVn *resources.VirtualNetwork

		ginkgo.Context("when free reservation", func() {
			ginkgo.BeforeEach(func() {
				recName = vnReservationFree
			})

			ginkgo.It("should delete reservation", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vn = resources.CreateVirtualNetworkWithID(10)

				id := 1

				err = client.ReservationService.Free(context.TODO(), *vn, resources.AddressRange{ID: &id})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether reservation was really removed in OpenNebula
				oneVn, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), 10)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var addressRanges []*resources.AddressRange

				addressRanges, err = oneVn.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				containIP := false
				for _, ar := range addressRanges {
					if ar.ID == &id {
						containIP = true
						break
					}
				}
				gomega.Expect(containIP).ShouldNot(gomega.BeTrue())
			})
		})
	})
})
