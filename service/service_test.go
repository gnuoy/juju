// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package service_test

import (
	"runtime"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/service"
	"github.com/juju/juju/service/common"
	"github.com/juju/juju/service/upstart"
	"github.com/juju/juju/service/windows"
)

type serviceSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&serviceSuite{})

func (*serviceSuite) TestDiscoverService(c *gc.C) {
	name := "a-service"
	conf := common.Conf{
		Desc:      "some service",
		ExecStart: "<do something>",
	}
	svc, err := service.DiscoverService(name, conf)
	c.Assert(err, jc.ErrorIsNil)

	switch runtime.GOOS {
	case "linux":
		c.Check(svc, gc.FitsTypeOf, &upstart.Service{})
		conf.InitDir = "/etc/init"
	case "windows":
		c.Check(svc, gc.FitsTypeOf, &windows.Service{})
	default:
		c.Errorf("unrecognized os %q", runtime.GOOS)
	}
	c.Check(svc.Name(), gc.Equals, "a-service")
	c.Check(svc.Conf(), jc.DeepEquals, conf)
}

func (*serviceSuite) TestListServicesCommand(c *gc.C) {
	cmd := service.ListServicesCommand()

	c.Check(cmd, gc.Equals, ""+
		`if [[ "$(cat /proc/1/cmdline)" == "/sbin/init" ]]; then `+
		`sudo initctl list | awk '{print $1}' | sort | uniq`+"\n"+
		`else exit 1`+"\n"+
		`fi`)
}