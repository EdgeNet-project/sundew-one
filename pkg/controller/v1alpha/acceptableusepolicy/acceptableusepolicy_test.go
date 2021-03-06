package acceptableusepolicy

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	apps_v1alpha "github.com/EdgeNet-project/edgenet/pkg/apis/apps/v1alpha"
	"github.com/EdgeNet-project/edgenet/pkg/generated/clientset/versioned"
	edgenettestclient "github.com/EdgeNet-project/edgenet/pkg/generated/clientset/versioned/fake"
	"github.com/EdgeNet-project/edgenet/pkg/util"
	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

// The main structure of test group
type TestGroup struct {
	authorityObj  apps_v1alpha.Authority
	AUPObj        apps_v1alpha.AcceptableUsePolicy
	client        kubernetes.Interface
	edgenetClient versioned.Interface
	handler       Handler
}

func TestMain(m *testing.M) {
	flag.String("dir", "../../../..", "Override the directory.")
	flag.String("smtp-path", "../../../../configs/smtp_test.yaml", "Set SMTP path.")
	flag.Parse()

	log.SetOutput(ioutil.Discard)
	logrus.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

// Init syncs the test group
func (g *TestGroup) Init() {
	authorityObj := apps_v1alpha.Authority{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Authority",
			APIVersion: "apps.edgenet.io/v1alpha",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "edgenet",
		},
		Spec: apps_v1alpha.AuthoritySpec{
			FullName:  "EdgeNet",
			ShortName: "EdgeNet",
			URL:       "https://www.edge-net.org",
			Address: apps_v1alpha.Address{
				City:    "Paris - NY - CA",
				Country: "France - US",
				Street:  "4 place Jussieu, boite 169",
				ZIP:     "75005",
			},
			Contact: apps_v1alpha.Contact{
				Email:     "joe.public@edge-net.org",
				FirstName: "Joe",
				LastName:  "Public",
				Phone:     "+33NUMBER",
				Username:  "joepublic",
			},
			Enabled: true,
		},
	}
	AUPObj := apps_v1alpha.AcceptableUsePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AcceptableUsePolicy",
			APIVersion: "apps.edgenet.io/v1alpha",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "joepublic",
			Namespace: "authority-edgenet",
		},
		Spec: apps_v1alpha.AcceptableUsePolicySpec{
			Accepted: false,
		},
	}
	g.authorityObj = authorityObj
	g.AUPObj = AUPObj
	g.client = testclient.NewSimpleClientset()
	g.edgenetClient = edgenettestclient.NewSimpleClientset()
	// authorityHandler := authority.Handler{}
	// authorityHandler.Init(g.client, g.edgenetClient)
	// Create Authority
	g.edgenetClient.AppsV1alpha().Authorities().Create(context.TODO(), g.authorityObj.DeepCopy(), metav1.CreateOptions{})
	namespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: g.AUPObj.GetNamespace()}}
	namespaceLabels := map[string]string{"owner": "authority", "owner-name": g.authorityObj.GetName(), "authority-name": g.authorityObj.GetName()}
	namespace.SetLabels(namespaceLabels)
	g.client.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{}) // Invoke ObjectCreated to create namespace
	// Create a user as admin on authority
	user := apps_v1alpha.User{}
	user.SetName(strings.ToLower(g.authorityObj.Spec.Contact.Username))
	user.Spec.Email = g.authorityObj.Spec.Contact.Email
	user.Spec.FirstName = g.authorityObj.Spec.Contact.FirstName
	user.Spec.LastName = g.authorityObj.Spec.Contact.LastName
	user.Spec.Active = true
	user.Status.AUP = false
	user.Status.Type = "admin"
	g.edgenetClient.AppsV1alpha().Users(fmt.Sprintf("authority-%s", g.authorityObj.GetName())).Create(context.TODO(), user.DeepCopy(), metav1.CreateOptions{})
	// authorityHandler.ObjectCreated(g.authorityObj.DeepCopy())
}

func TestHandlerInit(t *testing.T) {
	// Sync the test group
	g := TestGroup{}
	g.Init()
	// Initialize the handler
	g.handler.Init(g.client, g.edgenetClient)
	util.Equals(t, g.client, g.handler.clientset)
	util.Equals(t, g.edgenetClient, g.handler.edgenetClientset)
}

func TestCreate(t *testing.T) {
	g := TestGroup{}
	g.Init()
	g.handler.Init(g.client, g.edgenetClient)

	regular := g.AUPObj
	regular.SetUID("regular")
	accepted := g.AUPObj
	accepted.SetUID("accepted")
	accepted.Spec.Accepted = true
	recreation := g.AUPObj
	recreation.SetUID("recreation")
	recreation.Spec.Accepted = true
	recreation.Status.Expires = &metav1.Time{
		Time: time.Now().Add(1000 * time.Hour),
	}
	recreationExpired := g.AUPObj
	recreationExpired.SetUID("recreationExpired")
	recreationExpired.Spec.Accepted = true
	recreationExpired.Status.Expires = &metav1.Time{
		Time: time.Now().Add(-1000 * time.Hour),
	}
	t.Run("regular", func(t *testing.T) {
		g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(regular.GetNamespace()).Create(context.TODO(), regular.DeepCopy(), metav1.CreateOptions{})
		defer g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(regular.GetNamespace()).Delete(context.TODO(), regular.GetName(), metav1.DeleteOptions{})
		g.handler.ObjectCreated(regular.DeepCopy())
		AUP, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(regular.GetNamespace()).Get(context.TODO(), regular.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		util.Equals(t, success, AUP.Status.State)
		t.Run("user status", func(t *testing.T) {
			user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, false, user.Status.AUP)
		})
	})
	t.Run("accepted already", func(t *testing.T) {
		user, err := g.edgenetClient.AppsV1alpha().Users(g.AUPObj.GetNamespace()).Get(context.TODO(), g.AUPObj.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		user.Status.AUP = false
		g.edgenetClient.AppsV1alpha().Users(user.GetNamespace()).Update(context.TODO(), user, metav1.UpdateOptions{})

		g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(accepted.GetNamespace()).Create(context.TODO(), accepted.DeepCopy(), metav1.CreateOptions{})
		defer g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(accepted.GetNamespace()).Delete(context.TODO(), accepted.GetName(), metav1.DeleteOptions{})
		g.handler.ObjectCreated(accepted.DeepCopy())
		AUP, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(accepted.GetNamespace()).Get(context.TODO(), accepted.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		util.Equals(t, success, AUP.Status.State)
		t.Run("user status", func(t *testing.T) {
			user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, true, user.Status.AUP)
		})
	})
	t.Run("recreation", func(t *testing.T) {
		user, err := g.edgenetClient.AppsV1alpha().Users(g.AUPObj.GetNamespace()).Get(context.TODO(), g.AUPObj.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		user.Status.AUP = false
		g.edgenetClient.AppsV1alpha().Users(user.GetNamespace()).Update(context.TODO(), user, metav1.UpdateOptions{})

		g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreation.GetNamespace()).Create(context.TODO(), recreation.DeepCopy(), metav1.CreateOptions{})
		defer g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreation.GetNamespace()).Delete(context.TODO(), recreation.GetName(), metav1.DeleteOptions{})
		g.handler.ObjectCreated(recreation.DeepCopy())
		AUP, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreation.GetNamespace()).Get(context.TODO(), recreation.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		util.Equals(t, "", AUP.Status.State)
		t.Run("user status", func(t *testing.T) {
			user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, true, user.Status.AUP)
		})
	})
	t.Run("recreation of expired one", func(t *testing.T) {
		user, err := g.edgenetClient.AppsV1alpha().Users(g.AUPObj.GetNamespace()).Get(context.TODO(), g.AUPObj.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		user.Status.AUP = false
		g.edgenetClient.AppsV1alpha().Users(user.GetNamespace()).Update(context.TODO(), user, metav1.UpdateOptions{})

		g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreationExpired.GetNamespace()).Create(context.TODO(), recreationExpired.DeepCopy(), metav1.CreateOptions{})
		defer g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreationExpired.GetNamespace()).Delete(context.TODO(), recreationExpired.GetName(), metav1.DeleteOptions{})
		g.handler.ObjectCreated(recreationExpired.DeepCopy())
		AUP, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(recreationExpired.GetNamespace()).Get(context.TODO(), recreationExpired.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		util.Equals(t, failure, AUP.Status.State)
		t.Run("user status", func(t *testing.T) {
			user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, false, user.Status.AUP)
		})
	})
}

func TestAccept(t *testing.T) {
	g := TestGroup{}
	g.Init()
	g.handler.Init(g.client, g.edgenetClient)
	// Create AUP to update later
	g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(g.AUPObj.GetNamespace()).Create(context.TODO(), g.AUPObj.DeepCopy(), metav1.CreateOptions{})
	// Invoke ObjectCreated func to create a AUP
	g.handler.ObjectCreated(g.AUPObj.DeepCopy())
	AUP, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(g.AUPObj.GetNamespace()).Get(context.TODO(), g.AUPObj.GetName(), metav1.GetOptions{})
	util.OK(t, err)
	// Update of AUP status
	// Building field parameter
	var field fields
	field.accepted = true
	AUP.Spec.Accepted = true
	g.handler.ObjectUpdated(AUP.DeepCopy(), field)
	time.Sleep(time.Millisecond * 100)

	AUP, err = g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(g.AUPObj.GetNamespace()).Get(context.TODO(), g.AUPObj.GetName(), metav1.GetOptions{})
	t.Run("update", func(t *testing.T) {
		util.OK(t, err)
		util.Equals(t, success, AUP.Status.State)
	})
	t.Run("set expiry date", func(t *testing.T) {
		expected := metav1.Time{
			Time: time.Now().Add(4382 * time.Hour),
		}
		util.Equals(t, expected.Day(), AUP.Status.Expires.Day())
		util.Equals(t, expected.Month(), AUP.Status.Expires.Month())
		util.Equals(t, expected.Year(), AUP.Status.Expires.Year())
	})
	t.Run("user status", func(t *testing.T) {
		user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
		util.OK(t, err)
		util.Equals(t, true, user.Status.AUP)
	})
	t.Run("timeout", func(t *testing.T) {
		go g.handler.runApprovalTimeout(AUP)
		AUP.Status.Expires = &metav1.Time{
			Time: time.Now().Add(10 * time.Millisecond),
		}
		_, err := g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(AUP.GetNamespace()).Update(context.TODO(), AUP.DeepCopy(), metav1.UpdateOptions{})
		util.OK(t, err)
		time.Sleep(100 * time.Millisecond)
		t.Run("expired", func(t *testing.T) {
			AUP, err = g.edgenetClient.AppsV1alpha().AcceptableUsePolicies(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, false, AUP.Spec.Accepted)
		})
		t.Run("user status", func(t *testing.T) {
			user, err := g.edgenetClient.AppsV1alpha().Users(AUP.GetNamespace()).Get(context.TODO(), AUP.GetName(), metav1.GetOptions{})
			util.OK(t, err)
			util.Equals(t, false, user.Status.AUP)
		})
	})
}
