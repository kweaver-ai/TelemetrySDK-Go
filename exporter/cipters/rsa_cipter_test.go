package cipters

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRsaEncryptBase64(t *testing.T) {
	Convey("Test RsaEncryptBase64", t, func() {
		var t = "anyrobot"
		encryData, err := RsaEncryptBase64(t)
		So(err, ShouldBeNil)

		tt, err := RsaDecryptBase64(encryData)
		So(err, ShouldBeNil)
		So(tt, ShouldEqual, t)
	})
}

func TestRsaEncryptBase64AR(t *testing.T) {
	u, _ := RsaEncryptBase64("anyrobot")
	p, _ := RsaEncryptBase64("eisoo.com123")
	t.Log(u)
	t.Log(p)
	//anyrobot
	//eisoo.com123
}
