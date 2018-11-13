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

var (
	userGenerateToken           = "records/user/token/generateToken"
	userGenerateInfiniteToken   = "records/user/token/generateInfiniteToken"
	userGenerateTokenWrongUser  = "records/user/token/generateTokenWrongUser"
	userGenerateTokenWrongGroup = "records/user/token/generateTokenWrongGroup"
	userGenerateTokenEmptyGroup = "records/user/token/generateTokenEmptyGroup"

	userRevokeToken = "records/user/token/revokeToken"
)

var _ = ginkgo.Describe("Token Service", func() {
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

	ginkgo.Describe("Manage login token for the user", func() {
		var token string
		var group *resources.Group

		ginkgo.Context("when generate token", func() {
			ginkgo.BeforeEach(func() {
				recName = userGenerateToken
			})

			ginkgo.It("should generate a login token for a given user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				group = resources.CreateGroupWithID(120)

				token, err = client.TokenService.GenerateToken(context.TODO(), "Frantisek", 10, *group)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(token).NotTo(gomega.Equal(""))
			})
		})

		ginkgo.Context("when generate infinite token", func() {
			ginkgo.BeforeEach(func() {
				recName = userGenerateInfiniteToken
			})

			ginkgo.It("should generate an infinite login token for a given user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				group = resources.CreateGroupWithID(120)

				token, err = client.TokenService.GenerateInfiniteToken(context.TODO(), "Frantisek", *group)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(token).NotTo(gomega.Equal(""))
			})
		})

		ginkgo.Context("when the user with the given name does not exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userGenerateTokenWrongUser
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				token, err = client.TokenService.GenerateToken(context.TODO(), "Kamil", 10, *group)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(token).To(gomega.Equal(""))
			})
		})

		ginkgo.Context("when the group does not exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userGenerateTokenWrongGroup
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				group = resources.CreateGroupWithID(1008)

				token, err = client.TokenService.GenerateToken(context.TODO(), "Frantisek", 10, *group)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(token).To(gomega.Equal(""))
			})
		})

		ginkgo.Context("when the group is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userGenerateTokenEmptyGroup
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				token, err = client.TokenService.GenerateToken(context.TODO(), "Frantisek", 10,
					resources.Group{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(token).To(gomega.Equal(""))
			})
		})

		ginkgo.Context("when revoke token", func() {
			ginkgo.BeforeEach(func() {
				recName = userRevokeToken
			})

			ginkgo.It("should revoke a login token for a given user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				group = resources.CreateGroupWithID(120)

				err = client.TokenService.RevokeToken(context.TODO(), "Frantisek")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})
	})
})
