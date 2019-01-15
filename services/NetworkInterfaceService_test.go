package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	virtualMachineNICAttach               = "records/virtualMachine/nic/attach"
	virtualMachineNICAttachWrongID        = "records/virtualMachine/nic/attachWrongID"
	virtualMachineNICAttachNoID           = "records/virtualMachine/nic/attachNoID"
	virtualMachineNICAttachWrongBlueprint = "records/virtualMachine/nic/attachWrongBlueprint"

	virtualMachineNICDetach           = "records/virtualMachine/nic/detach"
	virtualMachineNICDetachWrongID    = "records/virtualMachine/nic/detachWrongID"
	virtualMachineNICDetachNoID       = "records/virtualMachine/nic/detachNoID"
	virtualMachineNICDetachWrongNICID = "records/virtualMachine/nic/detachWrongNICID"
	virtualMachineNICDetachNoNICID    = "records/virtualMachine/nic/detachNoNICID"
)

var _ = ginkgo.Describe("Network Interface Service", func() {
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

	ginkgo.Describe("attach NIC", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			nic       *resources.NIC
			networkID int
		)

		ginkgo.Context("when virtual machine exists and NIC is correct", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICAttach

				virtualMachineID = 134
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				networkID = 3
				nic = &resources.NIC{NetworkID: networkID}
			})

			ginkgo.It("should attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Attach(context.TODO(), *virtualMachine, *nic)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether NIC was really attached in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var nics []*resources.NIC
				nics, err = oneVirtualMachine.NICs()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(nics).To(gomega.HaveLen(1))
				gomega.Expect(nics[0].NetworkID).To(gomega.Equal(networkID))
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICAttachWrongID

				virtualMachineID = 100
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				networkID = 3
				nic = &resources.NIC{NetworkID: networkID}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Attach(context.TODO(), *virtualMachine, *nic)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICAttachNoID

				networkID = 3
				nic = &resources.NIC{NetworkID: networkID}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Attach(context.TODO(), resources.VirtualMachine{}, *nic)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine exists but NIC is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICAttachWrongBlueprint

				virtualMachineID = 88
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Attach(context.TODO(), *virtualMachine, resources.NIC{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("detach NIC", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			nic               *resources.NIC
			virtualMachineID  int
			nicID             int
		)

		ginkgo.Context("when virtual machine exists and has given NIC", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICDetach

				virtualMachineID = 137
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				nicID = 0
				nic = &resources.NIC{NicID: nicID}
			})

			ginkgo.It("should attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Detach(context.TODO(), *virtualMachine, *nic)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether NIC was really detached in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				var nics []*resources.NIC
				nics, err = oneVirtualMachine.NICs()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(nics).To(gomega.BeEmpty())
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICDetachWrongID

				virtualMachineID = 100
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				nicID = 0
				nic = &resources.NIC{NicID: nicID}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Detach(context.TODO(), *virtualMachine, *nic)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICDetachNoID

				nicID = 0
				nic = &resources.NIC{NicID: nicID}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Detach(context.TODO(), resources.VirtualMachine{}, *nic)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine exists but NIC is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICDetachWrongNICID

				virtualMachineID = 88
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				nicID = 0
				nic = &resources.NIC{NicID: nicID}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Detach(context.TODO(), *virtualMachine, *nic)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine exists but NIC is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineNICDetachNoNICID

				virtualMachineID = 88
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("shouldn't attach NIC to a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.NetworkInterfaceService.Detach(context.TODO(), *virtualMachine, resources.NIC{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})
})
