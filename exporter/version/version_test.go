package version

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	convey.Convey("TestVersion", t, func() {
		convey.So(TelemetrySDKVersion, convey.ShouldEqual, TelemetrySDKVersion)
	})
}
