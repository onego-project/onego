package resources

import (
	"time"

	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	userXML = "xml/user.xml"
)

var _ = ginkgo.Describe("User", func() {
	var (
		doc  *etree.Document
		user *User
		err  error
	)

	ginkgo.Describe("getters", func() {
		var loginTokens []LoginToken

		ginkgo.BeforeEach(func() {
			// create user with data
			doc = etree.NewDocument()
			err = doc.ReadFromFile(userXML)
			user = CreateUserFromXML(doc.Root())
		})

		ginkgo.It("should find all user attributes", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach
			gomega.Expect(user).ShouldNot(gomega.BeNil())

			gomega.Expect(user.ID()).To(gomega.Equal(22))
			gomega.Expect(user.Name()).To(gomega.Equal("Karol"))
			gomega.Expect(user.Password()).To(gomega.Equal("6adfb183a4a2c94a2f92dab5ade762a47889a5a1"))
			gomega.Expect(user.AuthDriver()).To(gomega.Equal("non-existing-driver"))
			gomega.Expect(user.MainGroup()).ShouldNot(gomega.BeNil())
			gomega.Expect(user.Groups()).To(gomega.HaveLen(1))
			gomega.Expect(user.Enabled()).To(gomega.BeTrue())

			loginTokens, err = user.LoginTokens()
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(loginTokens).To(gomega.HaveLen(1))
			gomega.Expect(loginTokens[0].EGID).To(gomega.Equal(-1))
			expTime := time.Unix(int64(1497632279), 0)
			gomega.Expect(loginTokens[0].ExpirationTime).To(gomega.Equal(&expTime))
			gomega.Expect(loginTokens[0].Token).To(gomega.Equal("325efc70635143dd1b33a8d73d0a1aa159515dd2"))
		})
	})

	ginkgo.Describe("create user", func() {
		ginkgo.It("should create user", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

			user = CreateUserWithID(42)

			gomega.Expect(user.ID()).To(gomega.Equal(42))
		})

		ginkgo.Context("when user doesn't have given attribute", func() {
			ginkgo.BeforeEach(func() {
				user = CreateUserWithID(42)
				err = nil
			})

			ginkgo.It("should return that user doesn't have name", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.Name()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that user doesn't have password", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.Password()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that user doesn't have AuthDriver", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.AuthDriver()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that user doesn't have main group", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.MainGroup()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that user doesn't have groups", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.Groups()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return that user doesn't have Enabled", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				_, err = user.Enabled()
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})

			ginkgo.It("should return empty array of LoginTokens and no error", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				var loginTokens []LoginToken
				loginTokens, err = user.LoginTokens()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(loginTokens).Should(gomega.HaveLen(0))
			})
		})
	})
})
