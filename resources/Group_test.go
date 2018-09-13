package resources

import (
	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	groupXML = "xml/group.xml"
)

var _ = ginkgo.Describe("Group", func() {
	var (
		doc   *etree.Document
		group *Group
		err   error
	)

	ginkgo.Describe("getters", func() {
		ginkgo.BeforeEach(func() {
			// create user with data
			doc = etree.NewDocument()
			err = doc.ReadFromFile(groupXML)
			group = &Group{Resource: Resource{XMLData: doc.Root()}}
		})

		ginkgo.It("should find all group attributes", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach
			gomega.Expect(group).ShouldNot(gomega.BeNil())

			gomega.Expect(group.ID()).To(gomega.Equal(102))
			gomega.Expect(group.Users()).To(gomega.HaveLen(1))
			gomega.Expect(group.Admins()).To(gomega.HaveLen(3))
		})
	})

	ginkgo.Describe("create group", func() {
		ginkgo.It("should create group", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

			group = CreateGroup(222)

			gomega.Expect(group.ID()).To(gomega.Equal(222))
		})

		ginkgo.Context("when group doesn't have given attribute", func() {
			ginkgo.BeforeEach(func() {
				group = CreateGroup(222)
				err = nil
			})

			ginkgo.It("should return that group doesn't have name", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = group.Name()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that group doesn't have users", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				var users []int
				users, err = group.Users()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(users).ShouldNot(gomega.BeNil())
				gomega.Expect(users).To(gomega.HaveLen(0))
			})

			ginkgo.It("should return that group doesn't have admins", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				var admins []int
				admins, err = group.Admins()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(admins).ShouldNot(gomega.BeNil())
				gomega.Expect(admins).To(gomega.HaveLen(0))
			})
		})
	})
})
